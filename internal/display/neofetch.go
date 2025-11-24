package display

import (
	"fmt"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/james-see/weatherornot/internal/api"
)

// NeofetchDisplay renders weather data in neofetch style
type NeofetchDisplay struct {
	useColors bool
	units     string
}

// NewNeofetchDisplay creates a new neofetch-style display
func NewNeofetchDisplay(useColors bool, units string) *NeofetchDisplay {
	return &NeofetchDisplay{
		useColors: useColors,
		units:     units,
	}
}

// Render renders the weather data in neofetch style
func (d *NeofetchDisplay) Render(data *api.WeatherData, showLocation bool) string {
	var output strings.Builder

	// Determine if it's night time
	now := time.Now()
	isNight := now.Hour() < 6 || now.Hour() > 18

	// Get weather icon
	icon := GetWeatherIcon(data.Current.ConditionCode, isNight)

	// Prepare info lines
	info := d.buildInfoLines(data, showLocation)

	// Combine icon and info side by side
	maxLines := len(icon)
	if len(info) > maxLines {
		maxLines = len(info)
	}

	for i := 0; i < maxLines; i++ {
		var iconLine, infoLine string

		if i < len(icon) {
			iconLine = icon[i]
		} else {
			iconLine = strings.Repeat(" ", 13)
		}

		if i < len(info) {
			infoLine = info[i]
		}

		output.WriteString(iconLine)
		output.WriteString("  ")
		output.WriteString(infoLine)
		output.WriteString("\n")
	}

	return output.String()
}

// buildInfoLines builds the info lines for neofetch display
func (d *NeofetchDisplay) buildInfoLines(data *api.WeatherData, showLocation bool) []string {
	lines := make([]string, 0)

	// Title with location (if enabled)
	if showLocation {
		location := data.Location.Name
		if data.Location.Country != "" {
			location += ", " + data.Location.Country
		}
		lines = append(lines, d.colorize(location, color.FgCyan, true))
		lines = append(lines, d.colorize(strings.Repeat("-", len(location)), color.FgCyan, false))
	}

	// Condition
	lines = append(lines, fmt.Sprintf("%s %s",
		d.colorize("Weather:", color.FgBlue, true),
		strings.Title(data.Current.Condition)))

	// Temperature
	tempUnit := d.getTempUnit()
	lines = append(lines, fmt.Sprintf("%s %s%s",
		d.colorize("Temp:", color.FgBlue, true),
		d.formatTemp(data.Current.Temperature),
		tempUnit))

	// Feels like
	lines = append(lines, fmt.Sprintf("%s %s%s",
		d.colorize("Feels like:", color.FgBlue, true),
		d.formatTemp(data.Current.FeelsLike),
		tempUnit))

	// Humidity
	lines = append(lines, fmt.Sprintf("%s %d%%",
		d.colorize("Humidity:", color.FgBlue, true),
		data.Current.Humidity))

	// Wind
	windUnit := d.getWindUnit()
	lines = append(lines, fmt.Sprintf("%s %.1f %s",
		d.colorize("Wind:", color.FgBlue, true),
		data.Current.WindSpeed,
		windUnit))

	// Pressure
	lines = append(lines, fmt.Sprintf("%s %d hPa",
		d.colorize("Pressure:", color.FgBlue, true),
		data.Current.Pressure))

	// Visibility
	visibilityKm := float64(data.Current.Visibility) / 1000
	lines = append(lines, fmt.Sprintf("%s %.1f km",
		d.colorize("Visibility:", color.FgBlue, true),
		visibilityKm))

	// Cloud cover
	lines = append(lines, fmt.Sprintf("%s %d%%",
		d.colorize("Clouds:", color.FgBlue, true),
		data.Current.CloudCover))

	return lines
}

// colorize applies color to text if colors are enabled
func (d *NeofetchDisplay) colorize(text string, c color.Attribute, bold bool) string {
	if !d.useColors {
		return text
	}

	if bold {
		return color.New(c, color.Bold).Sprint(text)
	}
	return color.New(c).Sprint(text)
}

// getTempUnit returns the temperature unit symbol
func (d *NeofetchDisplay) getTempUnit() string {
	switch d.units {
	case "imperial":
		return "°F"
	case "metric":
		return "°C"
	default:
		return "K"
	}
}

// getWindUnit returns the wind speed unit
func (d *NeofetchDisplay) getWindUnit() string {
	switch d.units {
	case "imperial":
		return "mph"
	case "metric":
		return "m/s"
	default:
		return "m/s"
	}
}

// formatTemp formats temperature value
func (d *NeofetchDisplay) formatTemp(temp float64) string {
	return fmt.Sprintf("%.1f", temp)
}

