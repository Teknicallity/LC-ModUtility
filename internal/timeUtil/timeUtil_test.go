package timeUtil

import (
	"testing"
	"time"
)

func timesAreCloseEnough(t1, t2 time.Time, tolerance time.Duration) bool {
	diff := t1.Sub(t2)
	return diff <= tolerance && diff >= -tolerance
}

func TestParseTimeString(t *testing.T) {
	tolerance := time.Second

	testCases := []struct {
		input    string
		expected time.Time
	}{
		{"30 seconds ago", time.Now().Add(-30 * time.Second)},
		{"2 minutes ago", time.Now().Add(-2 * time.Minute)},
		{"An hour ago", time.Now().Add(-1 * time.Hour)},
		{"18 hours ago", time.Now().Add(-18 * time.Hour)},
		{"a month ago", time.Now().AddDate(0, -1, 0)},
		{"5 months ago", time.Now().AddDate(0, -5, 0)},
		{"a year ago", time.Now().AddDate(-1, 0, 0)},
		{"1 year ago", time.Now().AddDate(-1, 0, 0)},
		{"2 years ago", time.Now().AddDate(-2, 0, 0)},
	}

	for _, tc := range testCases {
		result := ParseTimeString(tc.input)

		if !timesAreCloseEnough(result, tc.expected, tolerance) {
			t.Errorf("For input '%s', expected %v, but got %v", tc.input, tc.expected, result)
		}
	}
}
