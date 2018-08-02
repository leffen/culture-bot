package handlers

import (
	"testing"

	calendar "google.golang.org/api/calendar/v3"
)

func TestGetLocationName(t *testing.T) {
	tests := map[string]string{
		"Tạ Quang VN (8)": "Ta Qwan",
		"Edison (8)":      "Edison",
		"Edison8":         "Edison",
		"":                "",
	}

	for loc, expected := range tests {
		actual := getLocationName(&calendar.Event{
			Location: loc,
		})

		if actual != expected {
			t.Errorf("getLocationName(%s) expected %s, got %s", loc, expected, actual)
		}
	}
}
