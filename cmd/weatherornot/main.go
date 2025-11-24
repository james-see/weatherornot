package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"github.com/james-see/weatherornot/internal/api"
	"github.com/james-see/weatherornot/internal/config"
	"github.com/james-see/weatherornot/internal/display"
	"github.com/james-see/weatherornot/internal/location"
)

var (
	// Global flags
	cfgFile      string
	displayMode  string
	units        string
	favorite     string
	noColor      bool
	showGraph    bool
	hours        int
	days         int
	showLocation bool
)

var rootCmd = &cobra.Command{
	Use:   "weatherornot [location]",
	Short: "A beautiful weather CLI tool",
	Long: `weatherornot is a command-line weather forecast tool that displays current conditions,
hourly forecasts, and daily forecasts in a beautiful terminal interface.

Location can be specified as:
  - ZIP code: 10001 or 10001,US
  - City: "San Francisco" or "San Francisco,CA" or "San Francisco,CA,US"
  - Coordinates: "37.7749,-122.4194"
  - Favorite: Use -f or --favorite flag

If no location is provided, uses default_location from config.`,
	Example: `  weatherornot 90210
  weatherornot "New York,NY"
  weatherornot "40.7128,-74.0060"
  weatherornot -f home
  weatherornot --mode neofetch "London,GB"`,
	Args: cobra.MaximumNArgs(1),
	RunE: runWeather,
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&displayMode, "mode", "m", "", "Display mode: widget or neofetch (default from config)")
	rootCmd.PersistentFlags().StringVarP(&units, "units", "u", "", "Units: metric, imperial, or standard (default from config)")
	rootCmd.PersistentFlags().StringVarP(&favorite, "favorite", "f", "", "Use a favorite location from config")
	rootCmd.PersistentFlags().BoolVar(&noColor, "no-color", false, "Disable colored output")
	rootCmd.PersistentFlags().BoolVar(&showGraph, "graph", true, "Show temperature graph")
	rootCmd.PersistentFlags().IntVar(&hours, "hours", 12, "Number of hourly forecasts to show")
	rootCmd.PersistentFlags().IntVar(&days, "days", 5, "Number of daily forecasts to show")
	rootCmd.PersistentFlags().BoolVar(&showLocation, "show-location", true, "Show location name in output")

	// Config subcommand
	rootCmd.AddCommand(configCmd)
}

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage configuration",
	Long:  "View and manage weatherornot configuration file",
}

func init() {
	configCmd.AddCommand(configInitCmd)
	configCmd.AddCommand(configSetCmd)
	configCmd.AddCommand(configShowCmd)
	configCmd.AddCommand(configPathCmd)
	configCmd.AddCommand(configFavoriteCmd)
}

var configInitCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize configuration file",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := config.DefaultConfig()
		
		// Prompt for API key
		fmt.Print("Enter your OpenWeatherMap API key: ")
		fmt.Scanln(&cfg.APIKey)
		
		// Prompt for default location
		fmt.Print("Enter default location (e.g., 10001 or San Francisco,CA): ")
		var location string
		fmt.Scanln(&location)
		cfg.DefaultLocation = location

		if err := config.Create(cfg); err != nil {
			return fmt.Errorf("failed to create config: %w", err)
		}

		path, _ := config.GetConfigPath()
		fmt.Printf("Configuration file created at: %s\n", path)
		return nil
	},
}

var configSetCmd = &cobra.Command{
	Use:   "set <key> <value>",
	Short: "Set a configuration value",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return err
		}

		key := args[0]
		value := args[1]

		switch key {
		case "api_key":
			cfg.APIKey = value
		case "default_location":
			cfg.DefaultLocation = value
		case "units":
			if value != "metric" && value != "imperial" && value != "standard" {
				return fmt.Errorf("units must be metric, imperial, or standard")
			}
			cfg.Units = value
		case "display_mode":
			if value != "widget" && value != "neofetch" {
				return fmt.Errorf("display_mode must be widget or neofetch")
			}
			cfg.DisplayMode = value
		case "show_colors":
			boolVal, err := strconv.ParseBool(value)
			if err != nil {
				return fmt.Errorf("show_colors must be true or false")
			}
			cfg.ShowColors = boolVal
		default:
			return fmt.Errorf("unknown config key: %s", key)
		}

		if err := config.Save(cfg); err != nil {
			return err
		}

		fmt.Printf("Set %s = %s\n", key, value)
		return nil
	},
}

var configShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show current configuration",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return err
		}

		fmt.Printf("API Key:          %s\n", maskAPIKey(cfg.APIKey))
		fmt.Printf("Provider:         %s\n", cfg.Provider)
		fmt.Printf("Default Location: %s\n", cfg.DefaultLocation)
		fmt.Printf("Units:            %s\n", cfg.Units)
		fmt.Printf("Display Mode:     %s\n", cfg.DisplayMode)
		fmt.Printf("Show Colors:      %t\n", cfg.ShowColors)
		
		if len(cfg.Favorites) > 0 {
			fmt.Println("\nFavorites:")
			for name, loc := range cfg.Favorites {
				fmt.Printf("  %s: %s\n", name, loc)
			}
		}

		return nil
	},
}

var configPathCmd = &cobra.Command{
	Use:   "path",
	Short: "Show configuration file path",
	RunE: func(cmd *cobra.Command, args []string) error {
		path, err := config.GetConfigPath()
		if err != nil {
			return err
		}
		fmt.Println(path)
		return nil
	},
}

var configFavoriteCmd = &cobra.Command{
	Use:   "favorite",
	Short: "Manage favorite locations",
}

func init() {
	configFavoriteCmd.AddCommand(favoriteAddCmd)
	configFavoriteCmd.AddCommand(favoriteRemoveCmd)
	configFavoriteCmd.AddCommand(favoriteListCmd)
}

var favoriteAddCmd = &cobra.Command{
	Use:   "add <name> <location>",
	Short: "Add a favorite location",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return err
		}

		name := args[0]
		location := args[1]

		if cfg.Favorites == nil {
			cfg.Favorites = make(map[string]string)
		}

		cfg.Favorites[name] = location

		if err := config.Save(cfg); err != nil {
			return err
		}

		fmt.Printf("Added favorite: %s = %s\n", name, location)
		return nil
	},
}

var favoriteRemoveCmd = &cobra.Command{
	Use:   "remove <name>",
	Short: "Remove a favorite location",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return err
		}

		name := args[0]

		if _, exists := cfg.Favorites[name]; !exists {
			return fmt.Errorf("favorite '%s' not found", name)
		}

		delete(cfg.Favorites, name)

		if err := config.Save(cfg); err != nil {
			return err
		}

		fmt.Printf("Removed favorite: %s\n", name)
		return nil
	},
}

var favoriteListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all favorite locations",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return err
		}

		if len(cfg.Favorites) == 0 {
			fmt.Println("No favorites configured")
			return nil
		}

		fmt.Println("Favorite Locations:")
		for name, location := range cfg.Favorites {
			fmt.Printf("  %s: %s\n", name, location)
		}

		return nil
	},
}

func runWeather(cmd *cobra.Command, args []string) error {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Check if API key is set
	if cfg.APIKey == "" {
		fmt.Println("No API key configured. Please run 'weatherornot config init' or set api_key.")
		return fmt.Errorf("API key not configured")
	}

	// Determine location
	var locationStr string
	if favorite != "" {
		// Use favorite location
		if loc, exists := cfg.Favorites[favorite]; exists {
			locationStr = loc
		} else {
			return fmt.Errorf("favorite '%s' not found", favorite)
		}
	} else if len(args) > 0 {
		// Use command line argument
		locationStr = args[0]
	} else if cfg.DefaultLocation != "" {
		// Use default location from config
		locationStr = cfg.DefaultLocation
	} else {
		return fmt.Errorf("no location specified and no default location configured")
	}

	// Parse location
	loc, err := location.Parse(locationStr)
	if err != nil {
		return fmt.Errorf("failed to parse location: %w", err)
	}

	// Override config with command line flags
	if units != "" {
		cfg.Units = units
	}
	if displayMode != "" {
		cfg.DisplayMode = displayMode
	}
	if noColor {
		cfg.ShowColors = false
	}

	// Create API client
	client := api.NewClient(cfg.APIKey, cfg.Units)

	// Fetch weather data based on location type
	var weatherData *api.WeatherData
	switch loc.Type {
	case location.TypeZip:
		weatherData, err = client.GetWeatherByZip(loc.Zip, loc.CountryCode)
	case location.TypeCity:
		weatherData, err = client.GetWeatherByCity(loc.City, loc.State, loc.Country)
	case location.TypeCoords:
		weatherData, err = client.GetWeatherByCoords(loc.Latitude, loc.Longitude)
	default:
		return fmt.Errorf("unsupported location type")
	}

	if err != nil {
		return fmt.Errorf("failed to fetch weather data: %w", err)
	}

	// Display weather data
	switch strings.ToLower(cfg.DisplayMode) {
	case "neofetch":
		renderer := display.NewNeofetchDisplay(cfg.ShowColors, cfg.Units)
		output := renderer.Render(weatherData, showLocation)
		fmt.Print(output)
		
		// Show graph if requested
		if showGraph {
			chartRenderer := display.NewChartDisplay(cfg.ShowColors, cfg.Units)
			graph := chartRenderer.RenderHourlyTempChart(weatherData, hours)
			fmt.Println("\n" + graph)
		}
	default: // widget
		renderer := display.NewWidgetDisplay(cfg.ShowColors, cfg.Units)
		output := renderer.Render(weatherData, showLocation, true, true)
		fmt.Print(output)
		
		// Show graph if requested
		if showGraph {
			chartRenderer := display.NewChartDisplay(cfg.ShowColors, cfg.Units)
			graph := chartRenderer.RenderHourlyTempChart(weatherData, hours)
			fmt.Println(graph)
		}
	}

	return nil
}

func maskAPIKey(apiKey string) string {
	if len(apiKey) <= 8 {
		return strings.Repeat("*", len(apiKey))
	}
	return apiKey[:4] + strings.Repeat("*", len(apiKey)-8) + apiKey[len(apiKey)-4:]
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

