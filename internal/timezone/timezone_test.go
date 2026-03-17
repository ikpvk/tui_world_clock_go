package timezone

import "testing"

func TestGetDisplayName(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"valid timezone with country", "America/New_York", "America/New_York (United States)"},
		{"valid timezone Europe", "Europe/London", "Europe/London (United Kingdom)"},
		{"valid timezone Asia", "Asia/Tokyo", "Asia/Tokyo (Japan)"},
		{"invalid timezone", "Invalid/Zone", "Invalid/Zone"},
		{"empty string", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetDisplayName(tt.input)
			if result != tt.expected {
				t.Errorf("GetDisplayName(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestFilterTimezones(t *testing.T) {
	tests := []struct {
		name     string
		query    string
		minCount int
	}{
		{"empty query", "", 80},
		{"exact match America", "America/New_York", 1},
		{"partial match Tokyo", "Tokyo", 1},
		{"partial match London", "London", 1},
		{"case insensitive search", "NEW_YORK", 1},
		{"country search", "United States", 5},
		{"country search Japan", "Japan", 1},
		{"no match", "InvalidQuery123", 0},
		{"partial country", "United", 5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FilterTimezones(tt.query)
			if len(result) < tt.minCount {
				t.Errorf("FilterTimezones(%q) returned %d results, expected at least %d", tt.query, len(result), tt.minCount)
			}
		})
	}
}

func TestGetTime(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"valid timezone UTC", "UTC", false},
		{"valid timezone America", "America/New_York", false},
		{"valid timezone Europe", "Europe/London", false},
		{"valid timezone Asia", "Asia/Tokyo", false},
		{"invalid timezone", "Invalid/Zone", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := GetTime(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetTime(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
		})
	}
}

func TestGetCurrentTimeForZone(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		format string
	}{
		{"valid UTC", "UTC", "2006-01-02 15:04:05"},
		{"valid America/New_York", "America/New_York", "2006-01-02 15:04:05"},
		{"invalid returns Invalid", "Invalid/Zone", "Invalid"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetCurrentTimeForZone(tt.input)
			if tt.input == "Invalid/Zone" {
				if result != "Invalid" {
					t.Errorf("GetCurrentTimeForZone(%q) = %q, want %q", tt.input, result, "Invalid")
				}
			} else {
				if len(result) != len(tt.format) {
					t.Errorf("GetCurrentTimeForZone(%q) = %q, unexpected format", tt.input, result)
				}
			}
		})
	}
}

func TestContainsIgnoreCase(t *testing.T) {
	tests := []struct {
		name     string
		s        string
		substr   string
		expected bool
	}{
		{"exact match", "Hello", "Hello", true},
		{"partial match", "Hello World", "World", true},
		{"case insensitive", "Hello", "hello", true},
		{"case insensitive mixed", "HELLO", "hello", true},
		{"no match", "Hello", "XYZ", false},
		{"empty substring", "Hello", "", true},
		{"empty main", "", "Hello", false},
		{"both empty", "", "", true},
		{"substring longer", "Hi", "Hello", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := containsIgnoreCase(tt.s, tt.substr)
			if result != tt.expected {
				t.Errorf("containsIgnoreCase(%q, %q) = %v, want %v", tt.s, tt.substr, result, tt.expected)
			}
		})
	}
}

func TestToLower(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"all uppercase", "HELLO", "hello"},
		{"all lowercase", "hello", "hello"},
		{"mixed case", "HeLLo", "hello"},
		{"empty string", "", ""},
		{"single char upper", "A", "a"},
		{"single char lower", "a", "a"},
		{"numbers unchanged", "Hello123", "hello123"},
		{"special chars unchanged", "Hello!", "hello!"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := toLower(tt.input)
			if result != tt.expected {
				t.Errorf("toLower(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestGetAllTimezones(t *testing.T) {
	result := GetAllTimezones()
	if len(result) == 0 {
		t.Error("GetAllTimezones() should not return empty slice")
	}
	if result[0] != "Africa/Abidjan" {
		t.Errorf("Expected first timezone to be Africa/Abidjan, got %s", result[0])
	}
	if result[len(result)-1] != "UTC" {
		t.Errorf("Expected last timezone to be UTC, got %s", result[len(result)-1])
	}
}
