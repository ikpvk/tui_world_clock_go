package theme

import "testing"

func TestGetThemeIndex(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected int
	}{
		{"valid dracula", "dracula", 0},
		{"valid nord", "nord", 1},
		{"valid one-dark", "one-dark", 2},
		{"valid gruvbox", "gruvbox", 3},
		{"valid solarized-light", "solarized-light", 4},
		{"valid github-light", "github-light", 5},
		{"valid monokai-light", "monokai-light", 6},
		{"valid paper", "paper", 7},
		{"invalid name", "invalid", 0},
		{"empty name", "", 0},
		{"case sensitive", "DRACULA", 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetThemeIndex(tt.input)
			if result != tt.expected {
				t.Errorf("GetThemeIndex(%q) = %d, want %d", tt.input, result, tt.expected)
			}
		})
	}
}

func TestGetThemeByIndex(t *testing.T) {
	tests := []struct {
		name     string
		index    int
		expected string
	}{
		{"valid index 0", 0, "dracula"},
		{"valid index 1", 1, "nord"},
		{"valid index 7", 7, "paper"},
		{"out of bounds", 100, "dracula"},
		{"negative index", -1, "dracula"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetThemeByIndex(tt.index)
			if result.Name != tt.expected {
				t.Errorf("GetThemeByIndex(%d) = %q, want %q", tt.index, result.Name, tt.expected)
			}
		})
	}
}

func TestThemeDataIntegrity(t *testing.T) {
	for _, th := range Themes {
		t.Run("theme_"+th.Name, func(t *testing.T) {
			if th.Name == "" {
				t.Error("Theme name should not be empty")
			}
			if th.Foreground == "" {
				t.Error("Theme Foreground should not be empty for", th.Name)
			}
			if th.Background == "" {
				t.Error("Theme Background should not be empty for", th.Name)
			}
			if th.Header == "" {
				t.Error("Theme Header should not be empty for", th.Name)
			}
			if th.SelectedBg == "" {
				t.Error("Theme SelectedBg should not be empty for", th.Name)
			}
			if th.SelectedFg == "" {
				t.Error("Theme SelectedFg should not be empty for", th.Name)
			}
			if !isValidHex(th.Foreground) {
				t.Errorf("Theme Foreground %q is not valid hex for %s", th.Foreground, th.Name)
			}
			if !isValidHex(th.Background) {
				t.Errorf("Theme Background %q is not valid hex for %s", th.Background, th.Name)
			}
			if !isValidHex(th.Header) {
				t.Errorf("Theme Header %q is not valid hex for %s", th.Header, th.Name)
			}
			if !isValidHex(th.SelectedBg) {
				t.Errorf("Theme SelectedBg %q is not valid hex for %s", th.SelectedBg, th.Name)
			}
			if !isValidHex(th.SelectedFg) {
				t.Errorf("Theme SelectedFg %q is not valid hex for %s", th.SelectedFg, th.Name)
			}
		})
	}
}

func isValidHex(s string) bool {
	if len(s) != 7 {
		return false
	}
	if s[0] != '#' {
		return false
	}
	for i := 1; i < 7; i++ {
		c := s[i]
		if !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F')) {
			return false
		}
	}
	return true
}

func TestThemeCount(t *testing.T) {
	if len(Themes) != 8 {
		t.Errorf("Expected 8 themes, got %d", len(Themes))
	}
}
