package main

import (
	"context"
	"path/filepath"
	"outline-client/internal/auth"
	"outline-client/internal/vpn"
)

// App struct
type App struct {
	ctx        context.Context
	vpnClient  *vpn.Client
	licManager *auth.LicenseManager
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{
		vpnClient:  vpn.NewClient(),
		licManager: auth.NewLicenseManager(filepath.Join("data", "license.json")),
	}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// ActivateLicense activates a license key
func (a *App) ActivateLicense(key string) map[string]interface{} {
	err := a.licManager.ActivateLicense(key)
	if err != nil {
		return map[string]interface{}{
			"success": false,
			"error":   err.Error(),
		}
	}

	return map[string]interface{}{
		"success": true,
		"info":    a.licManager.GetLicenseInfo(),
	}
}

// GetLicenseInfo returns the current license information
func (a *App) GetLicenseInfo() map[string]interface{} {
	return a.licManager.GetLicenseInfo()
}

// ConnectVPN connects to the VPN server
func (a *App) ConnectVPN(serverIP string, port int, password string) map[string]interface{} {
	if !a.licManager.IsLicenseValid() {
		return map[string]interface{}{
			"success": false,
			"error":   "請先啟用有效的授權",
		}
	}

	config := &vpn.Config{
		ServerIP: serverIP,
		Port:     port,
		Password: password,
		Method:   "aes-256-gcm",
	}

	err := a.vpnClient.Connect(config)
	if err != nil {
		return map[string]interface{}{
			"success": false,
			"error":   err.Error(),
		}
	}

	return map[string]interface{}{
		"success": true,
	}
}

// DisconnectVPN disconnects from the VPN server
func (a *App) DisconnectVPN() map[string]interface{} {
	err := a.vpnClient.Disconnect()
	if err != nil {
		return map[string]interface{}{
			"success": false,
			"error":   err.Error(),
		}
	}

	return map[string]interface{}{
		"success": true,
	}
}

// GetVPNStatus returns the current VPN connection status
func (a *App) GetVPNStatus() map[string]interface{} {
	isConnected := a.vpnClient.IsConnected()
	status := map[string]interface{}{
		"connected": isConnected,
	}

	if isConnected {
		stats := a.vpnClient.GetStats()
		status["stats"] = stats
	}

	return status
}
