package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGatorConfig(t *testing.T) {
	// Mock a config file in a temp folder
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, configFileName)

	mockConfig := `{
		"db_url": "postgresql://localhost:5432/testdb",
		"current_user_name": "testuser"
	}`
	err := os.WriteFile(configPath, []byte(mockConfig), 0644)
	if err != nil {
		t.Fatalf("failed to write mock config to temp config file: %v", err)
	}

	// Temporarily replace the osUserHomeDir return value to the temp mock file
	originalOsUserHomeDir := osUserHomeDir
	defer func() { osUserHomeDir = originalOsUserHomeDir }()
	osUserHomeDir = func() (string, error) {
		return tempDir, nil
	}

	// Call ReadGatorConfig
	cfg, err := ReadGatorConfig()
	if err != nil {
		t.Fatalf("error while reading mock config file: %v", err)
	}

	// Verify the loaded config
	if cfg.DBUrl != "postgresql://localhost:5432/testdb" {
		t.Errorf("wrong DB URL: %v", err)
	}
	if cfg.CurrentUserName != "testuser" {
		t.Errorf("wrong current user name: %v", err)
	}
}
