package vpn

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"sync"
	"sync/atomic"
	"time"
)

// Config represents the VPN connection configuration
type Config struct {
	ServerIP  string `json:"server_ip"`
	Port     int    `json:"port"`
	Password string `json:"password"`
	Method   string `json:"method"`
}

// Client represents an Outline VPN client
type Client struct {
	config     *Config
	isRunning  bool
	mutex      sync.Mutex
	connection net.Conn
	cipher     cipher.AEAD
	startTime  time.Time
	bytesRecv  uint64
	bytesSent  uint64
	stopChan   chan struct{}
}

// NewClient creates a new VPN client
func NewClient() *Client {
	return &Client{
		stopChan: make(chan struct{}),
	}
}

// Connect establishes a connection to the VPN server
func (c *Client) Connect(config *Config) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.isRunning {
		return fmt.Errorf("client is already running")
	}

	// Create AES-GCM cipher
	key := []byte(config.Password)
	if len(key) < 32 {
		// If password length is insufficient, use PKCS7 padding
		newKey := make([]byte, 32)
		copy(newKey, key)
		for i := len(key); i < 32; i++ {
			newKey[i] = byte(32 - len(key))
		}
		key = newKey
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return fmt.Errorf("create cipher failed: %v", err)
	}

	c.cipher, err = cipher.NewGCM(block)
	if err != nil {
		return fmt.Errorf("create GCM failed: %v", err)
	}

	// Connect to server
	addr := fmt.Sprintf("%s:%d", config.ServerIP, config.Port)
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return fmt.Errorf("connect to server failed: %v", err)
	}

	c.connection = conn
	c.config = config
	c.isRunning = true
	c.startTime = time.Now()

	// Start data forwarding
	go c.handleConnection()

	return nil
}

// handleConnection handles the VPN connection
func (c *Client) handleConnection() {
	defer c.Disconnect()

	// Create buffer
	buf := make([]byte, 4096)
	nonce := make([]byte, c.cipher.NonceSize())

	for {
		select {
		case <-c.stopChan:
			return
		default:
			// Read data
			n, err := c.connection.Read(buf)
			if err != nil {
				if err != io.EOF {
					fmt.Printf("read data failed: %v\n", err)
				}
				return
			}

			// Update received bytes
			atomic.AddUint64(&c.bytesRecv, uint64(n))

			// Decrypt data
			_, err = rand.Read(nonce)
			if err != nil {
				fmt.Printf("generate nonce failed: %v\n", err)
				return
			}

			// Encrypt data and send
			ciphertext := c.cipher.Seal(nil, nonce, buf[:n], nil)
			header := make([]byte, c.cipher.NonceSize()+2)
			copy(header, nonce)
			binary.BigEndian.PutUint16(header[c.cipher.NonceSize():], uint16(len(ciphertext)))

			_, err = c.connection.Write(append(header, ciphertext...))
			if err != nil {
				fmt.Printf("send data failed: %v\n", err)
				return
			}

			// Update sent bytes
			atomic.AddUint64(&c.bytesSent, uint64(len(ciphertext)+len(header)))
		}
	}
}

// Disconnect terminates the VPN connection
func (c *Client) Disconnect() error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if !c.isRunning {
		return fmt.Errorf("client is not running")
	}

	close(c.stopChan)

	if c.connection != nil {
		c.connection.Close()
		c.connection = nil
	}

	c.isRunning = false
	return nil
}

// IsConnected returns whether the client is currently connected
func (c *Client) IsConnected() bool {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	return c.isRunning
}

// GetStats returns connection statistics
func (c *Client) GetStats() map[string]interface{} {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	stats := map[string]interface{}{
		"bytesReceived": atomic.LoadUint64(&c.bytesRecv),
		"bytesSent":    atomic.LoadUint64(&c.bytesSent),
		"uptime":       0,
	}

	if c.isRunning {
		stats["uptime"] = int(time.Since(c.startTime).Seconds())
	}

	return stats
}
