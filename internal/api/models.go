package api

import "time"

// WeatherData represents the complete weather information
type WeatherData struct {
	Current  CurrentWeather
	Hourly   []HourlyForecast
	Daily    []DailyForecast
	Location Location
}

// Location represents geographic location information
type Location struct {
	Name      string
	Country   string
	Latitude  float64
	Longitude float64
	Timezone  string
}

// CurrentWeather represents current weather conditions
type CurrentWeather struct {
	Temperature     float64
	FeelsLike       float64
	Humidity        int
	Pressure        int
	WindSpeed       float64
	WindDegree      int
	Visibility      int
	CloudCover      int
	UVIndex         float64
	Condition       string
	ConditionCode   int
	Icon            string
	Sunrise         time.Time
	Sunset          time.Time
	Time            time.Time
}

// HourlyForecast represents hourly weather forecast
type HourlyForecast struct {
	Time            time.Time
	Temperature     float64
	FeelsLike       float64
	Humidity        int
	WindSpeed       float64
	Condition       string
	ConditionCode   int
	Icon            string
	PrecipChance    int
}

// DailyForecast represents daily weather forecast
type DailyForecast struct {
	Date            time.Time
	TempMax         float64
	TempMin         float64
	Humidity        int
	WindSpeed       float64
	Condition       string
	ConditionCode   int
	Icon            string
	PrecipChance    int
	Sunrise         time.Time
	Sunset          time.Time
}

// OpenWeatherMapResponse represents the response from OpenWeatherMap current weather API
type OpenWeatherMapResponse struct {
	Coord struct {
		Lon float64 `json:"lon"`
		Lat float64 `json:"lat"`
	} `json:"coord"`
	Weather []struct {
		ID          int    `json:"id"`
		Main        string `json:"main"`
		Description string `json:"description"`
		Icon        string `json:"icon"`
	} `json:"weather"`
	Base string `json:"base"`
	Main struct {
		Temp      float64 `json:"temp"`
		FeelsLike float64 `json:"feels_like"`
		TempMin   float64 `json:"temp_min"`
		TempMax   float64 `json:"temp_max"`
		Pressure  int     `json:"pressure"`
		Humidity  int     `json:"humidity"`
	} `json:"main"`
	Visibility int `json:"visibility"`
	Wind       struct {
		Speed float64 `json:"speed"`
		Deg   int     `json:"deg"`
	} `json:"wind"`
	Clouds struct {
		All int `json:"all"`
	} `json:"clouds"`
	Dt  int64 `json:"dt"`
	Sys struct {
		Type    int    `json:"type"`
		ID      int    `json:"id"`
		Country string `json:"country"`
		Sunrise int64  `json:"sunrise"`
		Sunset  int64  `json:"sunset"`
	} `json:"sys"`
	Timezone int    `json:"timezone"`
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Cod      int    `json:"cod"`
}

// OpenWeatherMapForecastResponse represents the response from OpenWeatherMap forecast API
type OpenWeatherMapForecastResponse struct {
	Cod     string `json:"cod"`
	Message int    `json:"message"`
	Cnt     int    `json:"cnt"`
	List    []struct {
		Dt   int64 `json:"dt"`
		Main struct {
			Temp      float64 `json:"temp"`
			FeelsLike float64 `json:"feels_like"`
			TempMin   float64 `json:"temp_min"`
			TempMax   float64 `json:"temp_max"`
			Pressure  int     `json:"pressure"`
			Humidity  int     `json:"humidity"`
		} `json:"main"`
		Weather []struct {
			ID          int    `json:"id"`
			Main        string `json:"main"`
			Description string `json:"description"`
			Icon        string `json:"icon"`
		} `json:"weather"`
		Clouds struct {
			All int `json:"all"`
		} `json:"clouds"`
		Wind struct {
			Speed float64 `json:"speed"`
			Deg   int     `json:"deg"`
		} `json:"wind"`
		Pop      float64 `json:"pop"`
		DtTxt    string  `json:"dt_txt"`
	} `json:"list"`
	City struct {
		ID    int    `json:"id"`
		Name  string `json:"name"`
		Coord struct {
			Lat float64 `json:"lat"`
			Lon float64 `json:"lon"`
		} `json:"coord"`
		Country  string `json:"country"`
		Timezone int    `json:"timezone"`
		Sunrise  int64  `json:"sunrise"`
		Sunset   int64  `json:"sunset"`
	} `json:"city"`
}

// GeocodingResponse represents the response from OpenWeatherMap geocoding API
type GeocodingResponse []struct {
	Name    string            `json:"name"`
	LocalNames map[string]string `json:"local_names,omitempty"`
	Lat     float64           `json:"lat"`
	Lon     float64           `json:"lon"`
	Country string            `json:"country"`
	State   string            `json:"state,omitempty"`
}

