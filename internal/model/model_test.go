package model

import "testing"

func TestNew(t *testing.T) {
	m := New()

	if m.timezones == nil {
		t.Error("New() should initialize timezones as empty slice, not nil")
	}
	if len(m.timezones) != 0 {
		t.Errorf("New() should have empty timezones, got %d", len(m.timezones))
	}
	if m.selected != 0 {
		t.Errorf("New() selected = %d, want 0", m.selected)
	}
	if m.view != ListView {
		t.Errorf("New() view = %v, want ListView", m.view)
	}
	if m.done != false {
		t.Error("New() done should be false")
	}
	if m.themeIndex != 0 {
		t.Errorf("New() themeIndex = %d, want 0", m.themeIndex)
	}
	if m.ticker == nil {
		t.Error("New() should initialize ticker")
	}
}

func TestModelSetTimezones(t *testing.T) {
	m := New()

	testTimezones := []string{"America/New_York", "Europe/London", "Asia/Tokyo"}
	m.SetTimezones(testTimezones)

	if len(m.timezones) != len(testTimezones) {
		t.Errorf("SetTimezones() len = %d, want %d", len(m.timezones), len(testTimezones))
	}

	for i, tz := range testTimezones {
		if m.timezones[i] != tz {
			t.Errorf("timezones[%d] = %q, want %q", i, m.timezones[i], tz)
		}
	}
}

func TestModelSetTheme(t *testing.T) {
	m := New()

	m.SetTheme("nord")
	if m.themeIndex != 1 {
		t.Errorf("SetTheme(nord) themeIndex = %d, want 1", m.themeIndex)
	}

	m.SetTheme("paper")
	if m.themeIndex != 7 {
		t.Errorf("SetTheme(paper) themeIndex = %d, want 7", m.themeIndex)
	}

	m.SetTheme("invalid")
	if m.themeIndex != 0 {
		t.Errorf("SetTheme(invalid) themeIndex = %d, want 0 (default)", m.themeIndex)
	}
}

func TestModelGetThemeName(t *testing.T) {
	m := New()

	m.SetTheme("dracula")
	if m.GetThemeName() != "dracula" {
		t.Errorf("GetThemeName() = %q, want dracula", m.GetThemeName())
	}

	m.SetTheme("nord")
	if m.GetThemeName() != "nord" {
		t.Errorf("GetThemeName() = %q, want nord", m.GetThemeName())
	}
}

func TestThemeCyclingForward(t *testing.T) {
	m := New()
	m.SetTheme("dracula")

	m.themeIndex = 0
	m.themeIndex = (m.themeIndex + 1) % 8
	if m.themeIndex != 1 {
		t.Errorf("After +1: themeIndex = %d, want 1", m.themeIndex)
	}

	m.themeIndex = 7
	m.themeIndex = (m.themeIndex + 1) % 8
	if m.themeIndex != 0 {
		t.Errorf("After wrap: themeIndex = %d, want 0", m.themeIndex)
	}
}

func TestThemeCyclingBackward(t *testing.T) {
	m := New()
	m.SetTheme("dracula")

	m.themeIndex = 0
	m.themeIndex = (m.themeIndex - 1 + 8) % 8
	if m.themeIndex != 7 {
		t.Errorf("After -1 from 0: themeIndex = %d, want 7", m.themeIndex)
	}

	m.themeIndex = 7
	m.themeIndex = (m.themeIndex - 1 + 8) % 8
	if m.themeIndex != 6 {
		t.Errorf("After -1 from 7: themeIndex = %d, want 6", m.themeIndex)
	}
}

func TestSplitDisplayName(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		wantTimezone string
		wantCountry  string
	}{
		{"with country", "America/New_York (United States)", "America/New_York", "United States"},
		{"without country", "UTC", "UTC", ""},
		{"empty", "", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tz, country := splitDisplayName(tt.input)
			if tz != tt.wantTimezone {
				t.Errorf("splitDisplayName(%q) timezone = %q, want %q", tt.input, tz, tt.wantTimezone)
			}
			if country != tt.wantCountry {
				t.Errorf("splitDisplayName(%q) country = %q, want %q", tt.input, country, tt.wantCountry)
			}
		})
	}
}

func TestContains(t *testing.T) {
	slice := []string{"a", "b", "c"}

	if !contains(slice, "a") {
		t.Error("contains(slice, 'a') = false, want true")
	}
	if !contains(slice, "b") {
		t.Error("contains(slice, 'b') = false, want true")
	}
	if contains(slice, "d") {
		t.Error("contains(slice, 'd') = true, want false")
	}
	if contains(nil, "a") {
		t.Error("contains(nil, 'a') = true, want false")
	}
}

func TestMin(t *testing.T) {
	if min(1, 2) != 1 {
		t.Error("min(1, 2) = 2, want 1")
	}
	if min(2, 1) != 1 {
		t.Error("min(2, 1) = 2, want 1")
	}
	if min(1, 1) != 1 {
		t.Error("min(1, 1) = 1, want 1")
	}
}

func TestMax(t *testing.T) {
	if max(1, 2) != 2 {
		t.Error("max(1, 2) = 1, want 2")
	}
	if max(2, 1) != 2 {
		t.Error("max(2, 1) = 1, want 2")
	}
	if max(1, 1) != 1 {
		t.Error("max(1, 1) = 1, want 1")
	}
}

func TestJoinLines(t *testing.T) {
	result := joinLines([]string{"a", "b", "c"})
	expected := "a\nb\nc"
	if result != expected {
		t.Errorf("joinLines() = %q, want %q", result, expected)
	}

	result = joinLines([]string{"single"})
	if result != "single" {
		t.Errorf("joinLines(single) = %q, want %q", result, "single")
	}

	result = joinLines([]string{})
	if result != "" {
		t.Errorf("joinLines(empty) = %q, want %q", result, "")
	}
}
