package auth

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
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
	return &LicenseManager{
		licenseFile: licenseFile,
	}
}

// ValidateLicense checks if a license key is valid
func (lm *LicenseManager) ValidateLicense(key string) error {
	// TODO: Implement actual license validation
	// This should involve:
	// 1. Checking against a license server
	// 2. Validating the signature
	// 3. Checking expiration

	// For now, we'll just do a basic check
	if len(key) < 32 {
		return fmt.Errorf("invalid license key")
	}

	return nil
}

// ActivateLicense activates a license key
func (lm *LicenseManager) ActivateLicense(key string) error {
	if err := lm.ValidateLicense(key); err != nil {
		return err
	}

	license := &License{
		Key:       key,
		IssuedTo:  "User", // This should come from the license server
		IssuedAt:  time.Now(),
		ExpiresAt: time.Now().AddDate(1, 0, 0), // 1 year license
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

	return time.Now().Before(lm.currentLicense.ExpiresAt)
}

// saveLicense saves the current license to disk
func (lm *LicenseManager) saveLicense() error {
	if lm.currentLicense == nil {
		return fmt.Errorf("no license to save")
	}

	data, err := json.Marshal(lm.currentLicense)
	if err != nil {
		return err
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
		return err
	}

	license := &License{}
	if err := json.Unmarshal(data, license); err != nil {
		return err
	}

	lm.currentLicense = license
	return nil
}
