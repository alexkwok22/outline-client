package vpn

import (
	"fmt"
	"net"
	"sync"
	"time"
)

// Config represents the VPN connection configuration
type Config struct {
	ServerIP   string `json:"server_ip"`
	Port      int    `json:"port"`
	Password  string `json:"password"`
	Method    string `json:"method"`
}

// Client represents an Outline VPN client
type Client struct {
	config     *Config
	isRunning  bool
	mutex      sync.Mutex
	connection net.Conn
	startTime  time.Time
	bytesRecv  uint64
	bytesSent  uint64
}

// NewClient creates a new VPN client
func NewClient() *Client {
	return &Client{}
}

// Connect establishes a connection to the VPN server
func (c *Client) Connect(config *Config) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.isRunning {
		return fmt.Errorf("client is already running")
	}

	c.config = config
	// TODO: Implement actual connection logic
	// This will involve:
	// 1. Setting up the shadowsocks connection
	// 2. Configuring system proxy
	// 3. Setting up routing tables

	c.isRunning = true
	c.startTime = time.Now()
	return nil
}

// Disconnect terminates the VPN connection
func (c *Client) Disconnect() error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if !c.isRunning {
		return fmt.Errorf("client is not running")
	}

	if c.connection != nil {
		c.connection.Close()
	}

	// TODO: Implement cleanup
	// 1. Remove system proxy settings
	// 2. Clean up routing tables

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
		"bytesReceived": c.bytesRecv,
		"bytesSent":    c.bytesSent,
		"uptime":       0,
	}

	if c.isRunning {
		stats["uptime"] = int(time.Since(c.startTime).Seconds())
	}

	return stats
}
