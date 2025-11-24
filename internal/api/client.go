package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

const (
	baseURL         = "https://api.openweathermap.org/data/2.5"
	geocodingURL    = "https://api.openweathermap.org/geo/1.0"
)

// Client represents an OpenWeatherMap API client
type Client struct {
	apiKey     string
	httpClient *http.Client
	units      string
}

// NewClient creates a new API client
func NewClient(apiKey, units string) *Client {
	return &Client{
		apiKey: apiKey,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		units: units,
	}
}

// GetWeatherByZip fetches weather data by zip code
func (c *Client) GetWeatherByZip(zip, countryCode string) (*WeatherData, error) {
	// First get current weather
	currentURL := fmt.Sprintf("%s/weather?zip=%s,%s&appid=%s&units=%s", 
		baseURL, zip, countryCode, c.apiKey, c.units)
	
	current, err := c.fetchCurrentWeather(currentURL)
	if err != nil {
		return nil, err
	}

	// Get forecast using coordinates
	forecast, err := c.GetForecast(current.Location.Latitude, current.Location.Longitude)
	if err != nil {
		return nil, err
	}

	return &WeatherData{
		Current:  current.Current,
		Location: current.Location,
		Hourly:   forecast.Hourly,
		Daily:    forecast.Daily,
	}, nil
}

// GetWeatherByCity fetches weather data by city name
func (c *Client) GetWeatherByCity(city, state, country string) (*WeatherData, error) {
	// Build query string
	query := city
	if state != "" {
		query += "," + state
	}
	if country != "" {
		query += "," + country
	}

	currentURL := fmt.Sprintf("%s/weather?q=%s&appid=%s&units=%s", 
		baseURL, url.QueryEscape(query), c.apiKey, c.units)
	
	current, err := c.fetchCurrentWeather(currentURL)
	if err != nil {
		return nil, err
	}

	// Get forecast using coordinates
	forecast, err := c.GetForecast(current.Location.Latitude, current.Location.Longitude)
	if err != nil {
		return nil, err
	}

	return &WeatherData{
		Current:  current.Current,
		Location: current.Location,
		Hourly:   forecast.Hourly,
		Daily:    forecast.Daily,
	}, nil
}

// GetWeatherByCoords fetches weather data by coordinates
func (c *Client) GetWeatherByCoords(lat, lon float64) (*WeatherData, error) {
	currentURL := fmt.Sprintf("%s/weather?lat=%f&lon=%f&appid=%s&units=%s", 
		baseURL, lat, lon, c.apiKey, c.units)
	
	current, err := c.fetchCurrentWeather(currentURL)
	if err != nil {
		return nil, err
	}

	// Get forecast using coordinates
	forecast, err := c.GetForecast(lat, lon)
	if err != nil {
		return nil, err
	}

	return &WeatherData{
		Current:  current.Current,
		Location: current.Location,
		Hourly:   forecast.Hourly,
		Daily:    forecast.Daily,
	}, nil
}

// fetchCurrentWeather fetches current weather from the API
func (c *Client) fetchCurrentWeather(apiURL string) (*WeatherData, error) {
	resp, err := c.httpClient.Get(apiURL)
	if err != nil {
		return nil, fmt.Errorf("error fetching weather data: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
	}

	var owmResp OpenWeatherMapResponse
	if err := json.NewDecoder(resp.Body).Decode(&owmResp); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return c.parseCurrentWeather(&owmResp), nil
}

// GetForecast fetches forecast data
func (c *Client) GetForecast(lat, lon float64) (*WeatherData, error) {
	forecastURL := fmt.Sprintf("%s/forecast?lat=%f&lon=%f&appid=%s&units=%s", 
		baseURL, lat, lon, c.apiKey, c.units)

	resp, err := c.httpClient.Get(forecastURL)
	if err != nil {
		return nil, fmt.Errorf("error fetching forecast data: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
	}

	var forecastResp OpenWeatherMapForecastResponse
	if err := json.NewDecoder(resp.Body).Decode(&forecastResp); err != nil {
		return nil, fmt.Errorf("error decoding forecast response: %w", err)
	}

	return c.parseForecast(&forecastResp), nil
}

// parseCurrentWeather converts OpenWeatherMap response to our WeatherData model
func (c *Client) parseCurrentWeather(resp *OpenWeatherMapResponse) *WeatherData {
	var condition string
	var conditionCode int
	var icon string
	
	if len(resp.Weather) > 0 {
		condition = resp.Weather[0].Description
		conditionCode = resp.Weather[0].ID
		icon = resp.Weather[0].Icon
	}

	return &WeatherData{
		Location: Location{
			Name:      resp.Name,
			Country:   resp.Sys.Country,
			Latitude:  resp.Coord.Lat,
			Longitude: resp.Coord.Lon,
			Timezone:  fmt.Sprintf("UTC%+d", resp.Timezone/3600),
		},
		Current: CurrentWeather{
			Temperature:   resp.Main.Temp,
			FeelsLike:     resp.Main.FeelsLike,
			Humidity:      resp.Main.Humidity,
			Pressure:      resp.Main.Pressure,
			WindSpeed:     resp.Wind.Speed,
			WindDegree:    resp.Wind.Deg,
			Visibility:    resp.Visibility,
			CloudCover:    resp.Clouds.All,
			Condition:     condition,
			ConditionCode: conditionCode,
			Icon:          icon,
			Sunrise:       time.Unix(resp.Sys.Sunrise, 0),
			Sunset:        time.Unix(resp.Sys.Sunset, 0),
			Time:          time.Unix(resp.Dt, 0),
		},
	}
}

// parseForecast converts OpenWeatherMap forecast response to our model
func (c *Client) parseForecast(resp *OpenWeatherMapForecastResponse) *WeatherData {
	data := &WeatherData{
		Hourly: make([]HourlyForecast, 0),
		Daily:  make([]DailyForecast, 0),
	}

	// Process hourly forecasts (up to 40 3-hour intervals)
	for i, item := range resp.List {
		if i >= 16 { // Limit to 48 hours (16 * 3-hour intervals)
			break
		}

		var condition string
		var conditionCode int
		var icon string
		
		if len(item.Weather) > 0 {
			condition = item.Weather[0].Description
			conditionCode = item.Weather[0].ID
			icon = item.Weather[0].Icon
		}

		data.Hourly = append(data.Hourly, HourlyForecast{
			Time:          time.Unix(item.Dt, 0),
			Temperature:   item.Main.Temp,
			FeelsLike:     item.Main.FeelsLike,
			Humidity:      item.Main.Humidity,
			WindSpeed:     item.Wind.Speed,
			Condition:     condition,
			ConditionCode: conditionCode,
			Icon:          icon,
			PrecipChance:  int(item.Pop * 100),
		})
	}

	// Process daily forecasts (aggregate 3-hour data into daily)
	dailyMap := make(map[string]*DailyForecast)
	
	for _, item := range resp.List {
		t := time.Unix(item.Dt, 0)
		dateKey := t.Format("2006-01-02")

		if daily, exists := dailyMap[dateKey]; exists {
			// Update min/max temps
			if item.Main.TempMax > daily.TempMax {
				daily.TempMax = item.Main.TempMax
			}
			if item.Main.TempMin < daily.TempMin {
				daily.TempMin = item.Main.TempMin
			}
			// Average other values
			if item.Pop > float64(daily.PrecipChance)/100 {
				daily.PrecipChance = int(item.Pop * 100)
			}
		} else {
			var condition string
			var conditionCode int
			var icon string
			
			if len(item.Weather) > 0 {
				condition = item.Weather[0].Description
				conditionCode = item.Weather[0].ID
				icon = item.Weather[0].Icon
			}

			dailyMap[dateKey] = &DailyForecast{
				Date:          t,
				TempMax:       item.Main.TempMax,
				TempMin:       item.Main.TempMin,
				Humidity:      item.Main.Humidity,
				WindSpeed:     item.Wind.Speed,
				Condition:     condition,
				ConditionCode: conditionCode,
				Icon:          icon,
				PrecipChance:  int(item.Pop * 100),
			}
		}
	}

	// Convert map to slice and sort by date
	for _, daily := range dailyMap {
		data.Daily = append(data.Daily, *daily)
	}

	// Sort daily forecasts by date
	for i := 0; i < len(data.Daily)-1; i++ {
		for j := i + 1; j < len(data.Daily); j++ {
			if data.Daily[i].Date.After(data.Daily[j].Date) {
				data.Daily[i], data.Daily[j] = data.Daily[j], data.Daily[i]
			}
		}
	}

	// Limit to 5 days
	if len(data.Daily) > 5 {
		data.Daily = data.Daily[:5]
	}

	return data
}

// Geocode converts city name to coordinates
func (c *Client) Geocode(city, state, country string) (float64, float64, error) {
	query := city
	if state != "" {
		query += "," + state
	}
	if country != "" {
		query += "," + country
	}

	geocodeURL := fmt.Sprintf("%s/direct?q=%s&limit=1&appid=%s", 
		geocodingURL, url.QueryEscape(query), c.apiKey)

	resp, err := c.httpClient.Get(geocodeURL)
	if err != nil {
		return 0, 0, fmt.Errorf("error geocoding location: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return 0, 0, fmt.Errorf("geocoding API error (status %d): %s", resp.StatusCode, string(body))
	}

	var geoResp GeocodingResponse
	if err := json.NewDecoder(resp.Body).Decode(&geoResp); err != nil {
		return 0, 0, fmt.Errorf("error decoding geocoding response: %w", err)
	}

	if len(geoResp) == 0 {
		return 0, 0, fmt.Errorf("location not found")
	}

	return geoResp[0].Lat, geoResp[0].Lon, nil
}

