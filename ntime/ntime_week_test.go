package ntime_test

import (
	"testing"
	"time"

	"github.com/go-lazyer/north/ntime"
)

func TestStartOfWeek(t *testing.T) {
	tests := []struct {
		name string
		date time.Time
		want time.Time
	}{
		{"Monday", time.Date(2026, 5, 18, 0, 0, 0, 0, time.Local), time.Date(2026, 5, 18, 0, 0, 0, 0, time.Local)},
		{"Tuesday", time.Date(2026, 5, 19, 0, 0, 0, 0, time.Local), time.Date(2026, 5, 18, 0, 0, 0, 0, time.Local)},
		{"Wednesday", time.Date(2026, 5, 20, 0, 0, 0, 0, time.Local), time.Date(2026, 5, 18, 0, 0, 0, 0, time.Local)},
		{"Thursday", time.Date(2026, 5, 21, 0, 0, 0, 0, time.Local), time.Date(2026, 5, 18, 0, 0, 0, 0, time.Local)},
		{"Friday", time.Date(2026, 5, 22, 0, 0, 0, 0, time.Local), time.Date(2026, 5, 18, 0, 0, 0, 0, time.Local)},
		{"Saturday", time.Date(2026, 5, 23, 0, 0, 0, 0, time.Local), time.Date(2026, 5, 18, 0, 0, 0, 0, time.Local)},
		{"Sunday", time.Date(2026, 5, 24, 0, 0, 0, 0, time.Local), time.Date(2026, 5, 18, 0, 0, 0, 0, time.Local)},
		{"Monday_with_time", time.Date(2026, 5, 18, 10, 30, 45, 0, time.Local), time.Date(2026, 5, 18, 0, 0, 0, 0, time.Local)},
		{"Sunday_with_time", time.Date(2026, 5, 24, 23, 59, 59, 0, time.Local), time.Date(2026, 5, 18, 0, 0, 0, 0, time.Local)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ntime.StartOfWeek(tt.date); got != tt.want {
				t.Errorf("StartOfWeek(%v) = %v, want %v", tt.date, got, tt.want)
			}
		})
	}
}

func TestEndOfWeek(t *testing.T) {
	tests := []struct {
		name string
		date time.Time
		want time.Time
	}{
		{"Monday", time.Date(2026, 5, 18, 0, 0, 0, 0, time.Local), time.Date(2026, 5, 24, 0, 0, 0, 0, time.Local)},
		{"Tuesday", time.Date(2026, 5, 19, 0, 0, 0, 0, time.Local), time.Date(2026, 5, 24, 0, 0, 0, 0, time.Local)},
		{"Wednesday", time.Date(2026, 5, 20, 0, 0, 0, 0, time.Local), time.Date(2026, 5, 24, 0, 0, 0, 0, time.Local)},
		{"Thursday", time.Date(2026, 5, 21, 0, 0, 0, 0, time.Local), time.Date(2026, 5, 24, 0, 0, 0, 0, time.Local)},
		{"Friday", time.Date(2026, 5, 22, 0, 0, 0, 0, time.Local), time.Date(2026, 5, 24, 0, 0, 0, 0, time.Local)},
		{"Saturday", time.Date(2026, 5, 23, 0, 0, 0, 0, time.Local), time.Date(2026, 5, 24, 0, 0, 0, 0, time.Local)},
		{"Sunday", time.Date(2026, 5, 24, 0, 0, 0, 0, time.Local), time.Date(2026, 5, 24, 0, 0, 0, 0, time.Local)},
		{"Monday_with_time", time.Date(2026, 5, 18, 10, 30, 45, 0, time.Local), time.Date(2026, 5, 24, 0, 0, 0, 0, time.Local)},
		{"Sunday_with_time", time.Date(2026, 5, 24, 23, 59, 59, 0, time.Local), time.Date(2026, 5, 24, 0, 0, 0, 0, time.Local)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ntime.EndOfWeek(tt.date); got != tt.want {
				t.Errorf("EndOfWeek(%v) = %v, want %v", tt.date, got, tt.want)
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
