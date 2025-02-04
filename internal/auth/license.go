package auth

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// License represents a user license
type License struct {
	Key       string    `json:"key"`
	IssuedTo  string    `json:"issued_to"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiresAt time.Time `json:"expires_at"`
	Features  []string  `json:"features"`
}

// LicenseManager handles license validation and storage
type LicenseManager struct {
	currentLicense *License
	licenseFile   string
}

// NewLicenseManager creates a new license manager
func NewLicenseManager(licenseFile string) *LicenseManager {
	manager := &LicenseManager{
		licenseFile: licenseFile,
	}
	
	// Create license directory if it doesn't exist
	if dir := filepath.Dir(licenseFile); dir != "" {
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			log.Fatal(err)
		}
	}
	
	// Try to load existing license
	manager.loadLicense()
	return manager
}

// ValidateLicense checks if a license key is valid
func (lm *LicenseManager) ValidateLicense(key string) error {
	// Remove any whitespace
	key = strings.TrimSpace(key)
	
	// Basic validation
	if len(key) < 32 {
		return fmt.Errorf("授權碼長度不足")
	}

	// Check if the key starts with "00"
	if !strings.HasPrefix(key, "00") {
		return fmt.Errorf("無效的授權碼")
	}

	return nil
}

// ActivateLicense activates a license key
func (lm *LicenseManager) ActivateLicense(key string) error {
	if err := lm.ValidateLicense(key); err != nil {
		return err
	}

	// Create a new license
	license := &License{
		Key:       key,
		IssuedTo:  "User", // In production, this would come from a server
		IssuedAt:  time.Now(),
		ExpiresAt: time.Now().AddDate(0, 1, 0), // 1 month license
		Features:  []string{"basic", "premium"},
	}

	lm.currentLicense = license
	return lm.saveLicense()
}

// GetCurrentLicense returns the current active license
func (lm *LicenseManager) GetCurrentLicense() *License {
	return lm.currentLicense
}

// IsLicenseValid checks if the current license is valid
func (lm *LicenseManager) IsLicenseValid() bool {
	if lm.currentLicense == nil {
		return false
	}

	// Check if license has expired
	if time.Now().After(lm.currentLicense.ExpiresAt) {
		return false
	}

	return true
}

// saveLicense saves the current license to disk
func (lm *LicenseManager) saveLicense() error {
	if lm.currentLicense == nil {
		return fmt.Errorf("沒有授權可以保存")
	}

	data, err := json.MarshalIndent(lm.currentLicense, "", "  ")
	if err != nil {
		return fmt.Errorf("保存授權時發生錯誤: %v", err)
	}

	return ioutil.WriteFile(lm.licenseFile, data, 0600)
}

// loadLicense loads a license from disk
func (lm *LicenseManager) loadLicense() error {
	data, err := ioutil.ReadFile(lm.licenseFile)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("讀取授權檔案時發生錯誤: %v", err)
	}

	license := &License{}
	if err := json.Unmarshal(data, license); err != nil {
		return fmt.Errorf("解析授權檔案時發生錯誤: %v", err)
	}

	lm.currentLicense = license
	return nil
}

// GetLicenseInfo returns formatted license information
func (lm *LicenseManager) GetLicenseInfo() map[string]interface{} {
	if lm.currentLicense == nil {
		return map[string]interface{}{
			"status": "未授權",
		}
	}

	daysLeft := int(time.Until(lm.currentLicense.ExpiresAt).Hours() / 24)
	
	return map[string]interface{}{
		"status":    "已授權",
		"issuedTo":  lm.currentLicense.IssuedTo,
		"expiresIn": fmt.Sprintf("%d 天", daysLeft),
		"features":  lm.currentLicense.Features,
	}
}
