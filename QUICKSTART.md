# Quick Start Guide

Get up and running with weatherornot in 5 minutes!

## Step 1: Get an API Key

1. Go to [OpenWeatherMap](https://openweathermap.org/api)
2. Sign up for a free account
3. Navigate to your API keys section
4. Copy your API key

## Step 2: Install

### Option A: Using go install
```bash
go install github.com/jc/weatherornot/cmd/weatherornot@latest
```

### Option B: Download Binary
Download from the [Releases](https://github.com/jc/weatherornot/releases) page

### Option C: Build from Source
```bash
git clone https://github.com/jc/weatherornot.git
cd weatherornot
go build -o weatherornot cmd/weatherornot/main.go
sudo mv weatherornot /usr/local/bin/
```

## Step 3: Initialize Configuration

```bash
weatherornot config init
```

You'll be prompted for:
- Your OpenWeatherMap API key
- A default location (optional)

## Step 4: Get Weather!

```bash
# Use your default location
weatherornot

# Or specify a location
weatherornot 90210
weatherornot "San Francisco,CA"
weatherornot "40.7128,-74.0060"
```

## Customize Your Experience

### Try Different Display Modes

```bash
# Widget mode (default) - Beautiful boxed panels
weatherornot --mode widget

# Neofetch mode - ASCII art style
weatherornot --mode neofetch
```

### Change Units

```bash
# Metric (Celsius, m/s)
weatherornot --units metric

# Imperial (Fahrenheit, mph)
weatherornot --units imperial
```

### Add Favorite Locations

```bash
weatherornot config favorite add home "Seattle,WA"
weatherornot config favorite add work "San Francisco,CA"

# Use favorites
weatherornot -f home
```

## Common Commands

```bash
# Show help
weatherornot --help

# View configuration
weatherornot config show

# List favorites
weatherornot config favorite list

# Disable colors
weatherornot --no-color

# Show more hourly forecasts
weatherornot --hours 24
```

## Troubleshooting

### "No API key configured"
Run `weatherornot config init` or manually set: `weatherornot config set api_key YOUR_KEY`

### "Location not found"
Try different formats:
- ZIP: `90210`
- City,State: `"Los Angeles,CA"`
- Coordinates: `"34.0522,-118.2437"`

### Can't find config file
Check path: `weatherornot config path`
Default: `~/.config/weatherornot/weatherornot.toml`

## Next Steps

- Read the full [README.md](README.md) for all features
- Explore all CLI flags with `weatherornot --help`
- Customize your config file directly at `~/.config/weatherornot/weatherornot.toml`

Happy weather checking! 🌤️

