package ntime

import (
	"fmt"
	"time"
)

const (
	DEFAULT_TIMEZONE = "Asia/Shanghai"
)

func Parse(value string, timezone ...string) time.Time {
	if value == "" {
		return time.Time{}
	}

	var (
		tz  string
		loc *time.Location
		err error
	)

	if len(timezone) > 0 {
		tz = timezone[0]
	} else {
		tz = DEFAULT_TIMEZONE
	}
	if loc, err = parseTimezone(tz); err != nil {
		return time.Time{}
	}
	switch value {
	case "now":
		return time.Now().In(loc)
	case "yesterday":
		return time.Now().AddDate(0, 0, -1).In(loc)
	case "tomorrow":
		return time.Now().AddDate(0, 0, 1).In(loc)
	}
	for i := range defaultLayouts {
		if tt, err := time.ParseInLocation(defaultLayouts[i], value, loc); err == nil {
			return tt
		}
	}
	return time.Time{}
}
func parseTimezone(timezone string) (loc *time.Location, err error) {
	if timezone == "" {
		return nil, fmt.Errorf("timezone cannot be empty")
	}
	if loc, err = time.LoadLocation(timezone); err != nil {
		err = fmt.Errorf("%w: %w", fmt.Errorf("invalid timezone %q, please see the file %q for all valid timezones", timezone, "$GOROOT/lib/time/zoneinfo.zip"), err)
	}
	return
}
