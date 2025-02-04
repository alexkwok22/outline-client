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
	app := &App{}
	
	// Initialize license manager with a file in the user's home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		homeDir = "."
	}
	licenseFile := filepath.Join(homeDir, ".outline", "license.json")
	app.licenseManager = auth.NewLicenseManager(licenseFile)
	
	// Initialize VPN client
	app.vpnClient = vpn.NewClient()
	
	return app
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// Connect connects to the VPN server
func (a *App) Connect(serverIP string, port int, password string) error {
	if !a.licenseManager.IsLicenseValid() {
		return fmt.Errorf("請先啟動授權")
	}
	
	config := &vpn.Config{
		ServerIP: serverIP,
		Port:     port,
		Password: password,
		Method:   "aes-256-gcm",
	}
	
	return a.vpnClient.Connect(config)
}

// Disconnect disconnects from the VPN server
func (a *App) Disconnect() error {
	return a.vpnClient.Disconnect()
}

// GetStatus returns the current connection status
func (a *App) GetStatus() bool {
	return a.vpnClient.IsConnected()
}

// GetStats returns the current connection statistics
func (a *App) GetStats() map[string]interface{} {
	return a.vpnClient.GetStats()
}

// ActivateLicense activates a license key
func (a *App) ActivateLicense(key string) error {
	return a.licenseManager.ActivateLicense(key)
}

// GetLicenseStatus returns whether the current license is valid
func (a *App) GetLicenseStatus() bool {
	return a.licenseManager.IsLicenseValid()
}

// GetLicenseInfo returns detailed license information
func (a *App) GetLicenseInfo() map[string]interface{} {
	return a.licenseManager.GetLicenseInfo()
}
