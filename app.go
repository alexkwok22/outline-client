package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"outline-client/internal/auth"
	"outline-client/internal/vpn"
)

// App struct
type App struct {
	ctx            context.Context
	vpnClient      *vpn.Client
	licenseManager *auth.LicenseManager
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	
	// Initialize license manager
	homeDir, _ := os.UserHomeDir()
	licenseFile := filepath.Join(homeDir, ".outline", "license.json")
	a.licenseManager = auth.NewLicenseManager(licenseFile)

	// Initialize VPN client with default config
	config := &vpn.Config{
		Method: "aes-256-gcm",
	}
	a.vpnClient = vpn.NewClient(config)
}

// Connect to VPN server
func (a *App) Connect(serverIP string, port int, password string) error {
	if !a.licenseManager.IsLicenseValid() {
		return fmt.Errorf("no valid license found")
	}

	config := &vpn.Config{
		ServerIP: serverIP,
		Port:     port,
		Password: password,
		Method:   "aes-256-gcm",
	}
	a.vpnClient = vpn.NewClient(config)
	return a.vpnClient.Connect()
}

// Disconnect from VPN server
func (a *App) Disconnect() error {
	return a.vpnClient.Disconnect()
}

// GetStatus returns the current VPN connection status
func (a *App) GetStatus() bool {
	return a.vpnClient.Status()
}

// GetStats returns VPN connection statistics
func (a *App) GetStats() map[string]interface{} {
	return a.vpnClient.GetStats()
}

// ActivateLicense activates a license key
func (a *App) ActivateLicense(key string) error {
	return a.licenseManager.ActivateLicense(key)
}

// GetLicenseStatus returns the current license status
func (a *App) GetLicenseStatus() bool {
	return a.licenseManager.IsLicenseValid()
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}
