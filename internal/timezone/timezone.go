package timezone

import (
	"time"
)

var timezoneCountryMap = map[string]string{
	"Africa/Abidjan":       "Ivory Coast",
	"Africa/Accra":         "Ghana",
	"Africa/Algiers":       "Algeria",
	"Africa/Cairo":         "Egypt",
	"Africa/Casablanca":    "Morocco",
	"Africa/Johannesburg":  "South Africa",
	"Africa/Lagos":         "Nigeria",
	"Africa/Nairobi":       "Kenya",
	"America/Anchorage":    "United States",
	"America/Bogota":       "Colombia",
	"America/Buenos_Aires": "Argentina",
	"America/Chicago":      "United States",
	"America/Denver":       "United States",
	"America/Halifax":      "Canada",
	"America/Lima":         "Peru",
	"America/Los_Angeles":  "United States",
	"America/Mexico_City":  "Mexico",
	"America/New_York":     "United States",
	"America/Phoenix":      "United States",
	"America/Santiago":     "Chile",
	"America/Sao_Paulo":    "Brazil",
	"America/Toronto":      "Canada",
	"America/Vancouver":    "Canada",
	"Asia/Baghdad":         "Iraq",
	"Asia/Bangkok":         "Thailand",
	"Asia/Colombo":         "Sri Lanka",
	"Asia/Dhaka":           "Bangladesh",
	"Asia/Dubai":           "United Arab Emirates",
	"Asia/Ho_Chi_Minh":     "Vietnam",
	"Asia/Hong_Kong":       "Hong Kong",
	"Asia/Jakarta":         "Indonesia",
	"Asia/Jerusalem":       "Israel",
	"Asia/Karachi":         "Pakistan",
	"Asia/Kolkata":         "India",
	"Asia/Kuala_Lumpur":    "Malaysia",
	"Asia/Manila":          "Philippines",
	"Asia/Riyadh":          "Saudi Arabia",
	"Asia/Seoul":           "South Korea",
	"Asia/Shanghai":        "China",
	"Asia/Singapore":       "Singapore",
	"Asia/Taipei":          "Taiwan",
	"Asia/Tehran":          "Iran",
	"Asia/Tokyo":           "Japan",
	"Atlantic/Reykjavik":   "Iceland",
	"Australia/Adelaide":   "Australia",
	"Australia/Brisbane":   "Australia",
	"Australia/Darwin":     "Australia",
	"Australia/Hobart":     "Australia",
	"Australia/Melbourne":  "Australia",
	"Australia/Perth":      "Australia",
	"Australia/Sydney":     "Australia",
	"Europe/Amsterdam":     "Netherlands",
	"Europe/Athens":        "Greece",
	"Europe/Belgrade":      "Serbia",
	"Europe/Berlin":        "Germany",
	"Europe/Brussels":      "Belgium",
	"Europe/Bucharest":     "Romania",
	"Europe/Budapest":      "Hungary",
	"Europe/Copenhagen":    "Denmark",
	"Europe/Dublin":        "Ireland",
	"Europe/Helsinki":      "Finland",
	"Europe/Istanbul":      "Turkey",
	"Europe/Kyiv":          "Ukraine",
	"Europe/Lisbon":        "Portugal",
	"Europe/London":        "United Kingdom",
	"Europe/Madrid":        "Spain",
	"Europe/Moscow":        "Russia",
	"Europe/Oslo":          "Norway",
	"Europe/Paris":         "France",
	"Europe/Prague":        "Czech Republic",
	"Europe/Riga":          "Latvia",
	"Europe/Rome":          "Italy",
	"Europe/Sofia":         "Bulgaria",
	"Europe/Stockholm":     "Sweden",
	"Europe/Tallinn":       "Estonia",
	"Europe/Vienna":        "Austria",
	"Europe/Vilnius":       "Lithuania",
	"Europe/Warsaw":        "Poland",
	"Europe/Zurich":        "Switzerland",
	"Pacific/Auckland":     "New Zealand",
	"Pacific/Fiji":         "Fiji",
	"Pacific/Guam":         "Guam",
	"Pacific/Honolulu":     "United States",
	"Pacific/Samoa":        "Samoa",
	"UTC":                  "UTC",
}

func GetDisplayName(tz string) string {
	country, exists := timezoneCountryMap[tz]
	if exists {
		return tz + " (" + country + ")"
	}
	return tz
}

func GetTime(tz string) (time.Time, error) {
	location, err := time.LoadLocation(tz)
	if err != nil {
		return time.Time{}, err
	}
	return time.Now().In(location), nil
}

func GetCurrentTimeForZone(tz string) string {
	t, err := GetTime(tz)
	if err != nil {
		return "Invalid"
	}
	return t.Format("2006-01-02 15:04:05")
}

func GetAllTimezones() []string {
	return []string{
		"Africa/Abidjan",
		"Africa/Accra",
		"Africa/Algiers",
		"Africa/Cairo",
		"Africa/Casablanca",
		"Africa/Johannesburg",
		"Africa/Lagos",
		"Africa/Nairobi",
		"America/Anchorage",
		"America/Bogota",
		"America/Buenos_Aires",
		"America/Chicago",
		"America/Denver",
		"America/Halifax",
		"America/Lima",
		"America/Los_Angeles",
		"America/Mexico_City",
		"America/New_York",
		"America/Phoenix",
		"America/Santiago",
		"America/Sao_Paulo",
		"America/Toronto",
		"America/Vancouver",
		"Asia/Baghdad",
		"Asia/Bangkok",
		"Asia/Colombo",
		"Asia/Dhaka",
		"Asia/Dubai",
		"Asia/Ho_Chi_Minh",
		"Asia/Hong_Kong",
		"Asia/Jakarta",
		"Asia/Jerusalem",
		"Asia/Karachi",
		"Asia/Kolkata",
		"Asia/Kuala_Lumpur",
		"Asia/Manila",
		"Asia/Riyadh",
		"Asia/Seoul",
		"Asia/Shanghai",
		"Asia/Singapore",
		"Asia/Taipei",
		"Asia/Tehran",
		"Asia/Tokyo",
		"Atlantic/Reykjavik",
		"Australia/Adelaide",
		"Australia/Brisbane",
		"Australia/Darwin",
		"Australia/Hobart",
		"Australia/Melbourne",
		"Australia/Perth",
		"Australia/Sydney",
		"Europe/Amsterdam",
		"Europe/Athens",
		"Europe/Belgrade",
		"Europe/Berlin",
		"Europe/Brussels",
		"Europe/Bucharest",
		"Europe/Budapest",
		"Europe/Copenhagen",
		"Europe/Dublin",
		"Europe/Helsinki",
		"Europe/Istanbul",
		"Europe/Kyiv",
		"Europe/Lisbon",
		"Europe/London",
		"Europe/Madrid",
		"Europe/Moscow",
		"Europe/Oslo",
		"Europe/Paris",
		"Europe/Prague",
		"Europe/Riga",
		"Europe/Rome",
		"Europe/Sofia",
		"Europe/Stockholm",
		"Europe/Tallinn",
		"Europe/Vienna",
		"Europe/Vilnius",
		"Europe/Warsaw",
		"Europe/Zurich",
		"Pacific/Auckland",
		"Pacific/Fiji",
		"Pacific/Guam",
		"Pacific/Honolulu",
		"Pacific/Samoa",
		"UTC",
	}
}

func FilterTimezones(query string) []string {
	all := GetAllTimezones()

	var result []string
	for _, tz := range all {
		displayName := GetDisplayName(tz)
		if containsIgnoreCase(tz, query) || containsIgnoreCase(displayName, query) {
			result = append(result, tz)
		}
	}
	return result
}

func containsIgnoreCase(s, substr string) bool {
	sLower := toLower(s)
	substrLower := toLower(substr)
	for i := 0; i <= len(sLower)-len(substrLower); i++ {
		if sLower[i:i+len(substrLower)] == substrLower {
			return true
		}
	}
	return false
}

func toLower(s string) string {
	result := make([]byte, len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c >= 'A' && c <= 'Z' {
			c += 'a' - 'A'
		}
		result[i] = c
	}
	return string(result)
}
