package location

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// LocationType represents the type of location input
type LocationType int

const (
	TypeUnknown LocationType = iota
	TypeZip
	TypeCity
	TypeCoords
)

// ParsedLocation represents a parsed location input
type ParsedLocation struct {
	Type        LocationType
	Zip         string
	CountryCode string
	City        string
	State       string
	Country     string
	Latitude    float64
	Longitude   float64
}

// Parse parses a location string and determines its type
func Parse(input string) (*ParsedLocation, error) {
	if input == "" {
		return nil, fmt.Errorf("location input cannot be empty")
	}

	input = strings.TrimSpace(input)

	// Try to parse as coordinates (lat,lng)
	if loc, ok := parseCoords(input); ok {
		return loc, nil
	}

	// Try to parse as zip code
	if loc, ok := parseZip(input); ok {
		return loc, nil
	}

	// Otherwise, treat as city/state/country
	return parseCity(input), nil
}

// parseCoords tries to parse input as coordinates
func parseCoords(input string) (*ParsedLocation, bool) {
	// Match patterns like: "40.7128,-74.0060" or "40.7128, -74.0060"
	coordRegex := regexp.MustCompile(`^(-?\d+\.?\d*)\s*,\s*(-?\d+\.?\d*)$`)
	matches := coordRegex.FindStringSubmatch(input)
	
	if len(matches) != 3 {
		return nil, false
	}

	lat, err1 := strconv.ParseFloat(matches[1], 64)
	lng, err2 := strconv.ParseFloat(matches[2], 64)

	if err1 != nil || err2 != nil {
		return nil, false
	}

	// Validate coordinate ranges
	if lat < -90 || lat > 90 || lng < -180 || lng > 180 {
		return nil, false
	}

	return &ParsedLocation{
		Type:      TypeCoords,
		Latitude:  lat,
		Longitude: lng,
	}, true
}

// parseZip tries to parse input as a zip code
func parseZip(input string) (*ParsedLocation, bool) {
	// US ZIP: 5 digits or 5+4 format
	usZipRegex := regexp.MustCompile(`^(\d{5})(-\d{4})?$`)
	if matches := usZipRegex.FindStringSubmatch(input); matches != nil {
		return &ParsedLocation{
			Type:        TypeZip,
			Zip:         matches[1],
			CountryCode: "US",
		}, true
	}

	// ZIP with country code: "12345,US" or "SW1A 1AA,GB"
	// Must contain at least one digit to be considered a zip code
	zipWithCountryRegex := regexp.MustCompile(`^([A-Z0-9]+(?:[\s-][A-Z0-9]+)?),\s*([A-Z]{2})$`)
	upperInput := strings.ToUpper(input)
	if matches := zipWithCountryRegex.FindStringSubmatch(upperInput); matches != nil {
		// Additional check: must contain at least one digit to be a zip code
		firstPart := matches[1]
		if !regexp.MustCompile(`\d`).MatchString(firstPart) {
			// No digits found - likely a city name, not a zip
			return nil, false
		}
		return &ParsedLocation{
			Type:        TypeZip,
			Zip:         strings.TrimSpace(matches[1]),
			CountryCode: matches[2],
		}, true
	}

	// Just digits, assume US ZIP
	if matched, _ := regexp.MatchString(`^\d{5}$`, input); matched {
		return &ParsedLocation{
			Type:        TypeZip,
			Zip:         input,
			CountryCode: "US",
		}, true
	}

	return nil, false
}

// parseCity parses input as city/state/country
func parseCity(input string) *ParsedLocation {
	parts := strings.Split(input, ",")
	
	loc := &ParsedLocation{
		Type: TypeCity,
	}

	switch len(parts) {
	case 1:
		// Just city name
		loc.City = strings.TrimSpace(parts[0])
	case 2:
		// City, State or City, Country
		loc.City = strings.TrimSpace(parts[0])
		second := strings.TrimSpace(parts[1])
		
		// Treat 2-letter codes as state (US convention: "City,ST")
		// For international use, specify as: "City,State,CC" or just "City" for well-known cities
		loc.State = second
	case 3:
		// City, State, Country
		loc.City = strings.TrimSpace(parts[0])
		loc.State = strings.TrimSpace(parts[1])
		loc.Country = strings.TrimSpace(strings.ToUpper(parts[2]))
	default:
		// More than 3 parts, just take first as city
		loc.City = strings.TrimSpace(parts[0])
	}

	return loc
}

// isAllUpperOrLower checks if a string is all uppercase or all lowercase
func isAllUpperOrLower(s string) bool {
	return s == strings.ToUpper(s) || s == strings.ToLower(s)
}

// String returns a string representation of the parsed location
func (l *ParsedLocation) String() string {
	switch l.Type {
	case TypeZip:
		return fmt.Sprintf("ZIP: %s, %s", l.Zip, l.CountryCode)
	case TypeCoords:
		return fmt.Sprintf("Coordinates: %.4f, %.4f", l.Latitude, l.Longitude)
	case TypeCity:
		parts := []string{l.City}
		if l.State != "" {
			parts = append(parts, l.State)
		}
		if l.Country != "" {
			parts = append(parts, l.Country)
		}
		return fmt.Sprintf("City: %s", strings.Join(parts, ", "))
	default:
		return "Unknown location type"
	}
}

