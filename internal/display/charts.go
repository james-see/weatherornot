package display

import (
	"fmt"
	"strings"

	"github.com/guptarohit/asciigraph"
	"github.com/jc/weatherornot/internal/api"
)

// ChartDisplay handles temperature trend ASCII graphs
type ChartDisplay struct {
	useColors bool
	units     string
}

// NewChartDisplay creates a new chart display
func NewChartDisplay(useColors bool, units string) *ChartDisplay {
	return &ChartDisplay{
		useColors: useColors,
		units:     units,
	}
}

// RenderHourlyTempChart renders an ASCII chart of hourly temperatures
func (d *ChartDisplay) RenderHourlyTempChart(data *api.WeatherData, hours int) string {
	if len(data.Hourly) == 0 {
		return ""
	}

	maxHours := hours
	if len(data.Hourly) < maxHours {
		maxHours = len(data.Hourly)
	}

	// Extract temperatures
	temps := make([]float64, maxHours)
	for i := 0; i < maxHours; i++ {
		temps[i] = data.Hourly[i].Temperature
	}

	// Create chart
	graph := asciigraph.Plot(temps,
		asciigraph.Height(10),
		asciigraph.Width(60),
		asciigraph.Caption(fmt.Sprintf("Temperature Trend (Next %d Hours)", maxHours)),
	)

	return d.addTempLabels(graph, temps)
}

// RenderDailyTempChart renders an ASCII chart of daily temperatures
func (d *ChartDisplay) RenderDailyTempChart(data *api.WeatherData) string {
	if len(data.Daily) == 0 {
		return ""
	}

	days := len(data.Daily)
	if days > 5 {
		days = 5
	}

	// Extract max temperatures
	maxTemps := make([]float64, days)
	for i := 0; i < days; i++ {
		maxTemps[i] = data.Daily[i].TempMax
	}

	// Create chart
	graph := asciigraph.Plot(maxTemps,
		asciigraph.Height(10),
		asciigraph.Width(60),
		asciigraph.Caption(fmt.Sprintf("Daily Max Temperature (%d Days)", days)),
	)

	return d.addTempLabels(graph, maxTemps)
}

// RenderDailyTempRangeChart renders daily min/max temperature ranges
func (d *ChartDisplay) RenderDailyTempRangeChart(data *api.WeatherData) string {
	if len(data.Daily) == 0 {
		return ""
	}

	var output strings.Builder
	tempUnit := d.getTempUnit()

	output.WriteString("\nDaily Temperature Range:\n")
	output.WriteString(strings.Repeat("─", 50) + "\n")

	days := len(data.Daily)
	if days > 5 {
		days = 5
	}

	for i := 0; i < days; i++ {
		day := data.Daily[i]
		dateStr := day.Date.Format("Mon 01/02")
		
		// Calculate bar lengths
		minTemp := day.TempMin
		maxTemp := day.TempMax
		
		// Normalize to percentage for bar display (assuming reasonable temp range)
		barWidth := 30
		
		output.WriteString(fmt.Sprintf("%s  ", dateStr))
		output.WriteString(fmt.Sprintf("%.0f%s ", minTemp, tempUnit))
		
		// Draw bar
		tempRange := maxTemp - minTemp
		if tempRange > 0 {
			bar := strings.Repeat("█", int((tempRange/20.0)*float64(barWidth)))
			output.WriteString(bar)
		}
		
		output.WriteString(fmt.Sprintf(" %.0f%s\n", maxTemp, tempUnit))
	}

	return output.String()
}

// addTempLabels adds temperature unit labels to the graph
func (d *ChartDisplay) addTempLabels(graph string, temps []float64) string {
	tempUnit := d.getTempUnit()
	
	// Find min and max temps
	min, max := temps[0], temps[0]
	for _, t := range temps {
		if t < min {
			min = t
		}
		if t > max {
			max = t
		}
	}

	var output strings.Builder
	output.WriteString("\n")
	output.WriteString(graph)
	output.WriteString("\n")
	output.WriteString(fmt.Sprintf("Range: %.1f%s - %.1f%s\n", min, tempUnit, max, tempUnit))

	return output.String()
}

// getTempUnit returns the temperature unit symbol
func (d *ChartDisplay) getTempUnit() string {
	switch d.units {
	case "imperial":
		return "°F"
	case "metric":
		return "°C"
	default:
		return "K"
	}
}

// RenderCompactTempGraph renders a compact inline temperature trend
func (d *ChartDisplay) RenderCompactTempGraph(temps []float64) string {
	if len(temps) == 0 {
		return ""
	}

	// Simple sparkline-style graph using block characters
	min, max := temps[0], temps[0]
	for _, t := range temps {
		if t < min {
			min = t
		}
		if t > max {
			max = t
		}
	}

	var output strings.Builder
	blocks := []rune{'▁', '▂', '▃', '▄', '▅', '▆', '▇', '█'}

	for _, temp := range temps {
		// Normalize to 0-7 range
		normalized := 0
		if max > min {
			normalized = int(((temp - min) / (max - min)) * 7)
		}
		if normalized < 0 {
			normalized = 0
		}
		if normalized > 7 {
			normalized = 7
		}
		output.WriteRune(blocks[normalized])
	}

	return output.String()
}

