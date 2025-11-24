package location

import (
	"testing"
)

func TestParseCoordinates(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		expectLat float64
		expectLon float64
		expectErr bool
	}{
		{
			name:      "valid coordinates with spaces",
			input:     "40.7128, -74.0060",
			expectLat: 40.7128,
			expectLon: -74.0060,
			expectErr: false,
		},
		{
			name:      "valid coordinates without spaces",
			input:     "37.7749,-122.4194",
			expectLat: 37.7749,
			expectLon: -122.4194,
			expectErr: false,
		},
		{
			name:      "negative coordinates",
			input:     "-33.8688,151.2093",
			expectLat: -33.8688,
			expectLon: 151.2093,
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			loc, err := Parse(tt.input)
			if (err != nil) != tt.expectErr {
				t.Errorf("Parse() error = %v, expectErr %v", err, tt.expectErr)
				return
			}
			if err == nil {
				if loc.Type != TypeCoords {
					t.Errorf("Expected type TypeCoords, got %v", loc.Type)
				}
				if loc.Latitude != tt.expectLat {
					t.Errorf("Expected latitude %f, got %f", tt.expectLat, loc.Latitude)
				}
				if loc.Longitude != tt.expectLon {
					t.Errorf("Expected longitude %f, got %f", tt.expectLon, loc.Longitude)
				}
			}
		})
	}
}

func TestParseZipCode(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expectZip   string
		expectCountry string
	}{
		{
			name:        "5-digit US ZIP",
			input:       "90210",
			expectZip:   "90210",
			expectCountry: "US",
		},
		{
			name:        "ZIP with country code",
			input:       "10001,US",
			expectZip:   "10001",
			expectCountry: "US",
		},
		{
			name:        "UK postcode",
			input:       "SW1A 1AA,GB",
			expectZip:   "SW1A 1AA",
			expectCountry: "GB",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			loc, err := Parse(tt.input)
			if err != nil {
				t.Errorf("Parse() error = %v", err)
				return
			}
			if loc.Type != TypeZip {
				t.Errorf("Expected type TypeZip, got %v", loc.Type)
			}
			if loc.Zip != tt.expectZip {
				t.Errorf("Expected zip %s, got %s", tt.expectZip, loc.Zip)
			}
			if loc.CountryCode != tt.expectCountry {
				t.Errorf("Expected country %s, got %s", tt.expectCountry, loc.CountryCode)
			}
		})
	}
}

func TestParseCity(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		expectCity    string
		expectState   string
		expectCountry string
	}{
		{
			name:          "city only",
			input:         "London",
			expectCity:    "London",
			expectState:   "",
			expectCountry: "",
		},
		{
			name:          "city and state",
			input:         "San Francisco,CA",
			expectCity:    "San Francisco",
			expectState:   "CA",
			expectCountry: "",
		},
		{
			name:          "city and country code (treated as state)",
			input:         "Paris,FR",
			expectCity:    "Paris",
			expectState:   "FR",
			expectCountry: "",
		},
		{
			name:          "city, state, and country",
			input:         "New York,NY,US",
			expectCity:    "New York",
			expectState:   "NY",
			expectCountry: "US",
		},
		{
			name:          "city with spaces",
			input:         "San Francisco, CA, US",
			expectCity:    "San Francisco",
			expectState:   "CA",
			expectCountry: "US",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			loc, err := Parse(tt.input)
			if err != nil {
				t.Errorf("Parse() error = %v", err)
				return
			}
			if loc.Type != TypeCity {
				t.Errorf("Expected type TypeCity, got %v", loc.Type)
			}
			if loc.City != tt.expectCity {
				t.Errorf("Expected city %s, got %s", tt.expectCity, loc.City)
			}
			if loc.State != tt.expectState {
				t.Errorf("Expected state %s, got %s", tt.expectState, loc.State)
			}
			if loc.Country != tt.expectCountry {
				t.Errorf("Expected country %s, got %s", tt.expectCountry, loc.Country)
			}
		})
	}
}

func TestParseEmptyString(t *testing.T) {
	_, err := Parse("")
	if err == nil {
		t.Error("Expected error for empty string, got nil")
	}
}

func TestLocationString(t *testing.T) {
	tests := []struct {
		name     string
		loc      *ParsedLocation
		expected string
	}{
		{
			name: "coordinates",
			loc: &ParsedLocation{
				Type:      TypeCoords,
				Latitude:  40.7128,
				Longitude: -74.0060,
			},
			expected: "Coordinates: 40.7128, -74.0060",
		},
		{
			name: "zip code",
			loc: &ParsedLocation{
				Type:        TypeZip,
				Zip:         "90210",
				CountryCode: "US",
			},
			expected: "ZIP: 90210, US",
		},
		{
			name: "city",
			loc: &ParsedLocation{
				Type:    TypeCity,
				City:    "San Francisco",
				State:   "CA",
				Country: "US",
			},
			expected: "City: San Francisco, CA, US",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.loc.String()
			if result != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result)
			}
		})
	}
}

