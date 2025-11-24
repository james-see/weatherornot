# weatherornot

A beautiful, feature-rich weather CLI tool for your terminal. Get current weather, hourly forecasts, daily forecasts, and temperature trends with ASCII art and charts.

![License](https://img.shields.io/badge/license-MIT-blue.svg)
![Go Version](https://img.shields.io/badge/go-1.21+-00ADD8.svg)

🌐 **[Visit the Retro Website](https://james-see.github.io/weatherornot/)** 🌐

## Features

- 🌤️ **Multiple Display Modes**: Choose between neofetch-style or widget-style displays
- 📍 **Flexible Location Input**: Supports ZIP codes, city names, and coordinates
- 📊 **Temperature Graphs**: ASCII charts showing temperature trends
- ⭐ **Favorite Locations**: Save and quickly access your favorite locations
- 🎨 **Customizable**: Configure units, colors, and display preferences
- 🌍 **Powered by OpenWeatherMap**: Accurate weather data from a trusted source
- 💻 **Cross-Platform**: Available for Linux, macOS, Windows, and FreeBSD

## Installation

### Via `go install`

```bash
go install github.com/james-see/weatherornot/cmd/weatherornot@latest
```

Make sure `$GOPATH/bin` (or `$GOBIN`) is in your `PATH`.

### Download Pre-built Binaries

Download the latest release for your platform from the [Releases](https://github.com/james-see/weatherornot/releases) page.

#### Linux/macOS
```bash
# Download and extract (replace VERSION, OS, and ARCH with your values)
tar -xzf weatherornot_VERSION_OS_ARCH.tar.gz

# Move to your PATH
sudo mv weatherornot /usr/local/bin/
```

#### Windows
Download the `.zip` file, extract it, and add the executable to your PATH.

### Build from Source

```bash
git clone https://github.com/james-see/weatherornot.git
cd weatherornot
go build -o weatherornot cmd/weatherornot/main.go
```

## Quick Start

1. **Get an API Key**: Sign up at [OpenWeatherMap](https://openweathermap.org/) and get a free API key

2. **Initialize Configuration**:
```bash
weatherornot config init
```

3. **Get Weather**:
```bash
# Use default location from config
weatherornot

# Specify a location
weatherornot 90210
weatherornot "San Francisco,CA"
weatherornot "40.7128,-74.0060"
```

## Usage

### Basic Commands

```bash
# Get weather for a location
weatherornot <location>

# Use a favorite location
weatherornot -f home

# Switch display mode
weatherornot --mode neofetch "London,GB"
weatherornot --mode widget "Paris,FR"

# Change units
weatherornot --units metric "Tokyo,JP"
weatherornot --units imperial "New York,NY"

# Disable colors
weatherornot --no-color

# Show help
weatherornot --help
```

### Location Formats

weatherornot supports multiple location input formats:

- **ZIP Code**: `10001` or `10001,US`
- **City**: `"San Francisco"` or `"San Francisco,CA"` or `"San Francisco,CA,US"`
- **Coordinates**: `"37.7749,-122.4194"`

### Configuration Management

```bash
# Initialize config with prompts
weatherornot config init

# Show current configuration
weatherornot config show

# Set configuration values
weatherornot config set api_key YOUR_API_KEY
weatherornot config set default_location "Seattle,WA"
weatherornot config set units metric
weatherornot config set display_mode widget

# Get config file path
weatherornot config path
```

### Favorite Locations

```bash
# Add a favorite
weatherornot config favorite add home "Seattle,WA"
weatherornot config favorite add work "San Francisco,CA"
weatherornot config favorite add vacation "Honolulu,HI"

# List favorites
weatherornot config favorite list

# Remove a favorite
weatherornot config favorite remove work

# Use a favorite
weatherornot -f home
```

## Configuration File

Configuration is stored at `~/.config/weatherornot/weatherornot.toml`:

```toml
api_key = "your_openweathermap_api_key"
provider = "OpenWeatherMap"
default_location = "90210"
units = "imperial"  # metric, imperial, or standard
display_mode = "widget"  # widget or neofetch
show_colors = true

[favorites]
home = "San Francisco,CA,US"
work = "37.7749,-122.4194"
vacation = "33139"
```

## Display Modes

### Widget Mode (Default)

Beautiful boxed panels with current weather, hourly forecast, and daily forecast.

```bash
weatherornot --mode widget
```

### Neofetch Mode

ASCII art weather icons with clean information display, similar to neofetch.

```bash
weatherornot --mode neofetch
```

## Command-Line Flags

| Flag | Short | Description | Default |
|------|-------|-------------|---------|
| `--mode` | `-m` | Display mode (widget/neofetch) | From config |
| `--units` | `-u` | Units (metric/imperial/standard) | From config |
| `--favorite` | `-f` | Use favorite location | - |
| `--no-color` | - | Disable colored output | false |
| `--graph` | - | Show temperature graph | true |
| `--hours` | - | Number of hourly forecasts | 12 |
| `--days` | - | Number of daily forecasts | 5 |
| `--show-location` | - | Show location name | true |

## Examples

```bash
# Get weather with neofetch style
weatherornot --mode neofetch 10001

# Get weather in metric units
weatherornot --units metric "London,GB"

# Use a favorite location
weatherornot -f home

# Disable temperature graph
weatherornot --graph=false "Tokyo,JP"

# Show only 6 hours of forecast
weatherornot --hours 6 "Paris,FR"
```

## Development

### Prerequisites

- Go 1.21 or higher
- OpenWeatherMap API key

### Build

```bash
go build -o weatherornot cmd/weatherornot/main.go
```

### Run Tests

```bash
go test -v ./...
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgements

- [OpenWeatherMap](https://openweathermap.org/) for weather data
- [stormy](https://github.com/ashish0kumar/stormy) for inspiration
- [wttr.in](https://wttr.in) for ASCII weather icons
- [Cobra](https://github.com/spf13/cobra) for CLI framework
- [Lipgloss](https://github.com/charmbracelet/lipgloss) for beautiful terminal styling
- [asciigraph](https://github.com/guptarohit/asciigraph) for ASCII charts

## Support

If you encounter any issues or have questions, please [open an issue](https://github.com/james-see/weatherornot/issues).

---

Made with ❤️ and Go

