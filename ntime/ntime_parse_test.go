package ntime

import (
	"testing"
	"time"
)

func TestParse(t *testing.T) {
	// 1. empty string
	result := Parse("")
	if !result.IsZero() {
		t.Errorf("Parse(\"\") expected zero time, got %v", result)
	}

	// 2. time.Time type
	now := time.Now()
	result = Parse(now)
	if !result.Equal(now) {
		t.Errorf("Parse(time.Time) expected %v, got %v", now, result)
	}

	// 3. *time.Time non-nil
	ptr := &now
	result = Parse(ptr)
	if !result.Equal(now) {
		t.Errorf("Parse(*time.Time) expected %v, got %v", now, result)
	}

	// 4. *time.Time nil
	var nilPtr *time.Time
	result = Parse(nilPtr)
	if !result.IsZero() {
		t.Errorf("Parse(nil *time.Time) expected zero time, got %v", result)
	}

	// 5. int64 timestamp
	ts := int64(1700000000)
	result = Parse(ts)
	expected := time.Unix(ts, 0)
	if !result.Equal(expected) {
		t.Errorf("Parse(int64) expected %v, got %v", expected, result)
	}

	// 6. int timestamp
	its := int(1700000000)
	result = Parse(its)
	expected = time.Unix(int64(its), 0)
	if !result.Equal(expected) {
		t.Errorf("Parse(int) expected %v, got %v", expected, result)
	}

	// 7. non-string unsupported type (float64)
	result = Parse(3.14)
	if !result.IsZero() {
		t.Errorf("Parse(float64) expected zero time, got %v", result)
	}

	// 8. nil value
	result = Parse(nil)
	if !result.IsZero() {
		t.Errorf("Parse(nil) expected zero time, got %v", result)
	}

	// 9. "now" keyword
	result = Parse("now")
	if result.IsZero() {
		t.Errorf("Parse(\"now\") should not be zero")
	}
	loc, _ := time.LoadLocation(DEFAULT_TIMEZONE)
	if result.Location().String() != loc.String() {
		t.Errorf("Parse(\"now\") expected location %v, got %v", loc, result.Location())
	}

	// 10. "yesterday" keyword
	result = Parse("yesterday")
	if result.IsZero() {
		t.Errorf("Parse(\"yesterday\") should not be zero")
	}

	// 11. "tomorrow" keyword
	result = Parse("tomorrow")
	if result.IsZero() {
		t.Errorf("Parse(\"tomorrow\") should not be zero")
	}

	// 12. valid date string with default timezone
	result = Parse("2024-01-15 12:00:00")
	if result.IsZero() {
		t.Errorf("Parse(\"2024-01-15 12:00:00\") should not be zero")
	}

	// 13. valid date string with custom timezone
	result = Parse("2024-01-15 12:00:00", "UTC")
	if result.IsZero() {
		t.Errorf("Parse with UTC timezone should not be zero")
	}
	utcLoc, _ := time.LoadLocation("UTC")
	if result.Location().String() != utcLoc.String() {
		t.Errorf("Parse with UTC expected location %v, got %v", utcLoc, result.Location())
	}

	// 14. invalid timezone
	result = Parse("2024-01-15 12:00:00", "Invalid/Zone")
	if !result.IsZero() {
		t.Errorf("Parse with invalid timezone expected zero time, got %v", result)
	}

	// 15. invalid date string (unparseable)
	result = Parse("not-a-valid-date")
	if !result.IsZero() {
		t.Errorf("Parse(\"not-a-valid-date\") expected zero time, got %v", result)
	}

	// 16. time.Time with custom timezone (timezone param should be ignored for time.Time input)
	tInUTC := time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC)
	result = Parse(tInUTC, "Asia/Tokyo")
	if !result.Equal(tInUTC) {
		t.Errorf("Parse(time.Time) should return input as-is, expected %v, got %v", tInUTC, result)
	}

	// 17. "now" with custom timezone
	result = Parse("now", "America/New_York")
	if result.IsZero() {
		t.Errorf("Parse(\"now\", \"America/New_York\") should not be zero")
	}
	nyLoc, _ := time.LoadLocation("America/New_York")
	if result.Location().String() != nyLoc.String() {
		t.Errorf("Parse(\"now\") expected location %v, got %v", nyLoc, result.Location())
	}

	// 18. various date formats
	formats := []string{
		"2024-01-15",
		"12:00:00",
		"2024-01-15T12:00:00Z",
		"2024",
	}

	for _, f := range formats {
		result = Parse(f)
		if result.IsZero() {
			t.Errorf("Parse(%q) should not be zero", f)
		}
	}

	// 19. parseTimezone: empty timezone
	_, err := parseTimezone("")
	if err == nil {
		t.Errorf("parseTimezone(\"\") should return error")
	}
}
