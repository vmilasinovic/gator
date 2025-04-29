package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

const configFileName = ".gatorconfig.json"

var osUserHomeDir = os.UserHomeDir

type Config struct {
	DBUrl           string        `json:"db_url"`
	CurrentUserName string        `json:"current_user_name"`
	Mux             *sync.RWMutex `json:"-"`
	FilePath        string        `json:"-"`
}

// Read ~/.gatorconfig.json
func ReadGatorConfig() (*Config, error) {
	userHomeDir, err := osUserHomeDir()
	if err != nil {
		return &Config{}, fmt.Errorf("error retrieving home folder location: %w", err)
	}

	gatorConfigLocation := filepath.Join(userHomeDir, configFileName)
	gatorConfig, err := os.ReadFile(gatorConfigLocation)
	if err != nil {
		return &Config{}, fmt.Errorf("error reading gatorconfig file: %w", err)
	}

	var appConfig Config
	if err = json.Unmarshal(gatorConfig, &appConfig); err != nil {
		return &Config{}, fmt.Errorf("error unmarshaling gatorconfig data to appConfig struct: %w", err)
	}

	appConfig.Mux = new(sync.RWMutex)
	appConfig.FilePath = gatorConfigLocation

	return &appConfig, nil
}

func (cfg *Config) WriteToConf(userHomeDir string, data []byte) error {
	cfg.Mux.Lock()
	defer cfg.Mux.Unlock()

	if err := os.WriteFile(userHomeDir, data, 0644); err != nil {
		return fmt.Errorf("error writing to config file: %w", err)
	}

	return nil
}
