package config

// Config represents the application configuration
type Config struct {
	APIKey        string            `mapstructure:"api_key"`
	Provider      string            `mapstructure:"provider"`
	DefaultLocation string          `mapstructure:"default_location"`
	Units         string            `mapstructure:"units"`
	DisplayMode   string            `mapstructure:"display_mode"`
	ShowColors    bool              `mapstructure:"show_colors"`
	Favorites     map[string]string `mapstructure:"favorites"`
}

// DefaultConfig returns a new Config with default values
func DefaultConfig() *Config {
	return &Config{
		APIKey:          "",
		Provider:        "OpenWeatherMap",
		DefaultLocation: "",
		Units:           "imperial",
		DisplayMode:     "widget",
		ShowColors:      true,
		Favorites:       make(map[string]string),
	}
}

