package display

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/jc/weatherornot/internal/api"
)

// WidgetDisplay renders weather data in widget style with boxes
type WidgetDisplay struct {
	useColors bool
	units     string
}

// NewWidgetDisplay creates a new widget-style display
func NewWidgetDisplay(useColors bool, units string) *WidgetDisplay {
	return &WidgetDisplay{
		useColors: useColors,
		units:     units,
	}
}

// Render renders the weather data in widget style
func (d *WidgetDisplay) Render(data *api.WeatherData, showLocation bool, showHourly bool, showDaily bool) string {
	var output strings.Builder

	// Location header
	if showLocation {
		output.WriteString(d.renderLocationHeader(data))
		output.WriteString("\n\n")
	}

	// Current weather
	output.WriteString(d.renderCurrentWeather(data))
	output.WriteString("\n\n")

	// Hourly forecast
	if showHourly && len(data.Hourly) > 0 {
		output.WriteString(d.renderHourlyForecast(data))
		output.WriteString("\n\n")
	}

	// Daily forecast
	if showDaily && len(data.Daily) > 0 {
		output.WriteString(d.renderDailyForecast(data))
	}

	return output.String()
}

// renderLocationHeader renders the location header
func (d *WidgetDisplay) renderLocationHeader(data *api.WeatherData) string {
	location := data.Location.Name
	if data.Location.Country != "" {
		location += ", " + data.Location.Country
	}

	style := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("12")).
		PaddingLeft(2).
		PaddingRight(2)

	if !d.useColors {
		style = lipgloss.NewStyle().Bold(true).PaddingLeft(2).PaddingRight(2)
	}

	return style.Render(location)
}

// renderCurrentWeather renders current weather in a box
func (d *WidgetDisplay) renderCurrentWeather(data *api.WeatherData) string {
	var content strings.Builder

	// Determine if it's night
	now := time.Now()
	isNight := now.Hour() < 6 || now.Hour() > 18

	// Get emoji icon
	emoji := GetSimpleIcon(data.Current.ConditionCode, isNight)
	
	tempUnit := d.getTempUnit()
	windUnit := d.getWindUnit()

	content.WriteString(fmt.Sprintf("%s  %s\n\n", emoji, strings.Title(data.Current.Condition)))
	content.WriteString(fmt.Sprintf("Temperature:  %.1f%s (feels like %.1f%s)\n",
		data.Current.Temperature, tempUnit, data.Current.FeelsLike, tempUnit))
	content.WriteString(fmt.Sprintf("Humidity:     %d%%\n", data.Current.Humidity))
	content.WriteString(fmt.Sprintf("Wind:         %.1f %s\n", data.Current.WindSpeed, windUnit))
	content.WriteString(fmt.Sprintf("Pressure:     %d hPa\n", data.Current.Pressure))
	content.WriteString(fmt.Sprintf("Clouds:       %d%%\n", data.Current.CloudCover))
	content.WriteString(fmt.Sprintf("Visibility:   %.1f km", float64(data.Current.Visibility)/1000))

	return d.renderBox("Current Weather", content.String(), lipgloss.Color("14"))
}

// renderHourlyForecast renders hourly forecast
func (d *WidgetDisplay) renderHourlyForecast(data *api.WeatherData) string {
	var content strings.Builder

	tempUnit := d.getTempUnit()
	
	// Show up to 12 hours
	maxHours := 12
	if len(data.Hourly) < maxHours {
		maxHours = len(data.Hourly)
	}

	for i := 0; i < maxHours; i++ {
		hour := data.Hourly[i]
		timeStr := hour.Time.Format("15:04")
		icon := GetSimpleIcon(hour.ConditionCode, false)
		
		if i > 0 {
			content.WriteString("\n")
		}
		
		content.WriteString(fmt.Sprintf("%s  %s  %.1f%s  %s",
			timeStr, icon, hour.Temperature, tempUnit,
			strings.Title(hour.Condition)))
	}

	return d.renderBox("Hourly Forecast", content.String(), lipgloss.Color("11"))
}

// renderDailyForecast renders daily forecast
func (d *WidgetDisplay) renderDailyForecast(data *api.WeatherData) string {
	var content strings.Builder

	tempUnit := d.getTempUnit()

	for i, day := range data.Daily {
		if i >= 5 { // Limit to 5 days
			break
		}

		dateStr := day.Date.Format("Mon, Jan 02")
		icon := GetSimpleIcon(day.ConditionCode, false)
		
		if i > 0 {
			content.WriteString("\n")
		}
		
		content.WriteString(fmt.Sprintf("%s  %s  %.0f%s / %.0f%s  %s",
			dateStr, icon, day.TempMax, tempUnit, day.TempMin, tempUnit,
			strings.Title(day.Condition)))
	}

	return d.renderBox("5-Day Forecast", content.String(), lipgloss.Color("13"))
}

// renderBox renders content in a bordered box
func (d *WidgetDisplay) renderBox(title, content string, color lipgloss.Color) string {
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(color).
		PaddingLeft(1).
		PaddingRight(1)

	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(color).
		Padding(1, 2)

	if !d.useColors {
		titleStyle = lipgloss.NewStyle().Bold(true).PaddingLeft(1).PaddingRight(1)
		boxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			Padding(1, 2)
	}

	return titleStyle.Render(title) + "\n" + boxStyle.Render(content)
}

// getTempUnit returns the temperature unit symbol
func (d *WidgetDisplay) getTempUnit() string {
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
func (d *WidgetDisplay) getWindUnit() string {
	switch d.units {
	case "imperial":
		return "mph"
	case "metric":
		return "m/s"
	default:
		return "m/s"
	}
}

