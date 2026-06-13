package auth

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"time"
)

// Credentials holds the Twitter auth tokens needed for API access.
type Credentials struct {
	AuthToken string    `json:"auth_token"`
	CT0       string    `json:"ct0"`
	UserID    string    `json:"user_id,omitempty"` // numeric Twitter user ID
	Handle    string    `json:"handle,omitempty"`  // @screen_name (without @)
	SavedAt   time.Time `json:"saved_at"`          // MUST remain time.Time — auth status uses .IsZero() and .Format()
}

// ErrNotAuthenticated is returned when no credentials are stored.
var ErrNotAuthenticated = errors.New("not authenticated — run: twt auth login")

// configDir returns the path to ~/.config/twt/.
func configDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".config", "twt"), nil
}

// credPath returns the full path to the credentials file.
func credPath() (string, error) {
	dir, err := configDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "credentials.json"), nil
}

// Save writes credentials to disk at 0600 permissions.
func Save(creds *Credentials) error {
	dir, err := configDir()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return err
	}

	creds.SavedAt = time.Now()
	data, err := json.MarshalIndent(creds, "", "  ")
	if err != nil {
		return err
	}

	path, err := credPath()
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o600)
}

// Load reads stored credentials. Returns ErrNotAuthenticated if none exist.
func Load() (*Credentials, error) {
	path, err := credPath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return nil, ErrNotAuthenticated
	}
	if err != nil {
		return nil, err
	}

	var creds Credentials
	if err := json.Unmarshal(data, &creds); err != nil {
		return nil, err
	}
	if creds.AuthToken == "" || creds.CT0 == "" {
		return nil, ErrNotAuthenticated
	}
	return &creds, nil
}

// Delete removes the stored credentials file.
func Delete() error {
	path, err := credPath()
	if err != nil {
		return err
	}
	err = os.Remove(path)
	if errors.Is(err, os.ErrNotExist) {
		return nil // already gone
	}
	return err
}

// IsAuthenticated returns true if valid credentials exist.
func IsAuthenticated() bool {
	_, err := Load()
	return err == nil
}
