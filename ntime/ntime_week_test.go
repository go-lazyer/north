package ntime_test

import (
	"testing"
	"time"

	"github.com/go-lazyer/north/ntime"
)

func TestDayOfWeek(t *testing.T) {
	tests := []struct {
		name string
		date time.Time
		want int
	}{
		{"Monday", time.Date(2026, 5, 18, 0, 0, 0, 0, time.Local), 1},
		{"Tuesday", time.Date(2026, 5, 19, 0, 0, 0, 0, time.Local), 2},
		{"Wednesday", time.Date(2026, 5, 20, 0, 0, 0, 0, time.Local), 3},
		{"Thursday", time.Date(2026, 5, 21, 0, 0, 0, 0, time.Local), 4},
		{"Friday", time.Date(2026, 5, 22, 0, 0, 0, 0, time.Local), 5},
		{"Saturday", time.Date(2026, 5, 23, 0, 0, 0, 0, time.Local), 6},
		{"Sunday", time.Date(2026, 5, 24, 0, 0, 0, 0, time.Local), 7},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ntime.DayOfWeek(tt.date); got != tt.want {
				t.Errorf("DayOfWeek(%v) = %d, want %d", tt.date, got, tt.want)
			}
		})
	}
}

func TestWeekConstants(t *testing.T) {
	if ntime.Sunday != 0 {
		t.Errorf("Sunday = %d, want 0", ntime.Sunday)
	}
	if ntime.Monday != 1 {
		t.Errorf("Monday = %d, want 1", ntime.Monday)
	}
	if ntime.Tuesday != 2 {
		t.Errorf("Tuesday = %d, want 2", ntime.Tuesday)
	}
	if ntime.Wednesday != 3 {
		t.Errorf("Wednesday = %d, want 3", ntime.Wednesday)
	}
	if ntime.Thursday != 4 {
		t.Errorf("Thursday = %d, want 4", ntime.Thursday)
	}
	if ntime.Friday != 5 {
		t.Errorf("Friday = %d, want 5", ntime.Friday)
	}
	if ntime.Saturday != 6 {
		t.Errorf("Saturday = %d, want 6", ntime.Saturday)
	}
}
