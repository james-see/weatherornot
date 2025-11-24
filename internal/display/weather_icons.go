package display

// GetWeatherIcon returns ASCII art for weather conditions based on condition code
func GetWeatherIcon(conditionCode int, isNight bool) []string {
	// OpenWeatherMap condition codes:
	// 2xx: Thunderstorm
	// 3xx: Drizzle
	// 5xx: Rain
	// 6xx: Snow
	// 7xx: Atmosphere (fog, mist, etc.)
	// 800: Clear
	// 80x: Clouds

	switch {
	case conditionCode >= 200 && conditionCode < 300:
		return thunderstorm()
	case conditionCode >= 300 && conditionCode < 400:
		return drizzle()
	case conditionCode >= 500 && conditionCode < 600:
		return rain()
	case conditionCode >= 600 && conditionCode < 700:
		return snow()
	case conditionCode >= 700 && conditionCode < 800:
		return fog()
	case conditionCode == 800:
		if isNight {
			return clearNight()
		}
		return clearDay()
	case conditionCode > 800 && conditionCode < 803:
		return partlyCloudy()
	case conditionCode >= 803:
		return cloudy()
	default:
		return unknown()
	}
}

func clearDay() []string {
	return []string{
		"    \\   /    ",
		"     .-.     ",
		"  ― (   ) ―  ",
		"     `-'     ",
		"    /   \\    ",
	}
}

func clearNight() []string {
	return []string{
		"    .-.      ",
		"   (   )     ",
		"  (  .  )    ",
		"   (___)     ",
		"             ",
	}
}

func partlyCloudy() []string {
	return []string{
		"   \\  /      ",
		" _ /\"\".-.    ",
		"   \\_(   ).  ",
		"   /(___(__) ",
		"             ",
	}
}

func cloudy() []string {
	return []string{
		"             ",
		"     .--.    ",
		"  .-(    ).  ",
		" (___.__)__) ",
		"             ",
	}
}

func rain() []string {
	return []string{
		"     .-.     ",
		"    (   ).   ",
		"   (___(__)  ",
		"  ‚'‚'‚'‚'   ",
		"  ‚'‚'‚'‚'   ",
	}
}

func drizzle() []string {
	return []string{
		"     .-.     ",
		"    (   ).   ",
		"   (___(__)  ",
		"   ‚'‚'‚'    ",
		"   ‚'‚'‚'    ",
	}
}

func thunderstorm() []string {
	return []string{
		"     .-.     ",
		"    (   ).   ",
		"   (___(__)  ",
		"  ⚡'⚡'⚡'   ",
		"  ‚'‚'‚'     ",
	}
}

func snow() []string {
	return []string{
		"     .-.     ",
		"    (   ).   ",
		"   (___(__)  ",
		"   * * * *   ",
		"  * * * *    ",
	}
}

func fog() []string {
	return []string{
		"             ",
		" _ - _ - _ - ",
		"  _ - _ - _  ",
		" _ - _ - _ - ",
		"             ",
	}
}

func unknown() []string {
	return []string{
		"             ",
		"    .-.      ",
		"   (   )     ",
		"    `-'      ",
		"             ",
	}
}

// GetSimpleIcon returns a single character emoji for the weather condition
func GetSimpleIcon(conditionCode int, isNight bool) string {
	switch {
	case conditionCode >= 200 && conditionCode < 300:
		return "⛈️ "
	case conditionCode >= 300 && conditionCode < 400:
		return "🌦️ "
	case conditionCode >= 500 && conditionCode < 600:
		return "🌧️ "
	case conditionCode >= 600 && conditionCode < 700:
		return "❄️ "
	case conditionCode >= 700 && conditionCode < 800:
		return "🌫️ "
	case conditionCode == 800:
		if isNight {
			return "🌙"
		}
		return "☀️ "
	case conditionCode > 800 && conditionCode < 803:
		return "⛅"
	case conditionCode >= 803:
		return "☁️ "
	default:
		return "🌡️ "
	}
}

