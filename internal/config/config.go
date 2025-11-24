package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

const configFileName = "weatherornot"

// Load loads the configuration from the config file
func Load() (*Config, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("could not get user home directory: %w", err)
	}

	configPath := filepath.Join(homeDir, ".config", configFileName)
	
	viper.SetConfigName(configFileName)
	viper.SetConfigType("toml")
	viper.AddConfigPath(filepath.Join(homeDir, ".config"))
	viper.AddConfigPath(homeDir) // fallback to home directory

	// Set defaults
	cfg := DefaultConfig()
	viper.SetDefault("api_key", cfg.APIKey)
	viper.SetDefault("provider", cfg.Provider)
	viper.SetDefault("default_location", cfg.DefaultLocation)
	viper.SetDefault("units", cfg.Units)
	viper.SetDefault("display_mode", cfg.DisplayMode)
	viper.SetDefault("show_colors", cfg.ShowColors)
	viper.SetDefault("favorites", cfg.Favorites)

	// Try to read config
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found, create it
			return cfg, Create(cfg)
		}
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	// Unmarshal config
	if err := viper.Unmarshal(cfg); err != nil {
		return nil, fmt.Errorf("error unmarshaling config: %w", err)
	}

	// Ensure config directory exists
	if _, err := os.Stat(configPath + ".toml"); os.IsNotExist(err) {
		return cfg, Create(cfg)
	}

	return cfg, nil
}

// Create creates a new config file with default values
func Create(cfg *Config) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("could not get user home directory: %w", err)
	}

	configDir := filepath.Join(homeDir, ".config")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("could not create config directory: %w", err)
	}

	configPath := filepath.Join(configDir, configFileName+".toml")
	
	viper.SetConfigFile(configPath)
	viper.Set("api_key", cfg.APIKey)
	viper.Set("provider", cfg.Provider)
	viper.Set("default_location", cfg.DefaultLocation)
	viper.Set("units", cfg.Units)
	viper.Set("display_mode", cfg.DisplayMode)
	viper.Set("show_colors", cfg.ShowColors)
	viper.Set("favorites", cfg.Favorites)

	if err := viper.WriteConfig(); err != nil {
		return fmt.Errorf("error writing config file: %w", err)
	}

	return nil
}

// Save saves the current configuration to the config file
func Save(cfg *Config) error {
	viper.Set("api_key", cfg.APIKey)
	viper.Set("provider", cfg.Provider)
	viper.Set("default_location", cfg.DefaultLocation)
	viper.Set("units", cfg.Units)
	viper.Set("display_mode", cfg.DisplayMode)
	viper.Set("show_colors", cfg.ShowColors)
	viper.Set("favorites", cfg.Favorites)

	if err := viper.WriteConfig(); err != nil {
		return fmt.Errorf("error saving config file: %w", err)
	}

	return nil
}

// GetConfigPath returns the path to the config file
func GetConfigPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(homeDir, ".config", configFileName+".toml"), nil
}

