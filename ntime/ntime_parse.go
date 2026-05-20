package ntime

import (
	"fmt"
	"time"
)

const (
	DEFAULT_TIMEZONE = "Asia/Shanghai"
)

func Parse(value any, timezone ...string) time.Time {
	if value == "" {
		return time.Time{}
	}
	// 如果value是time.Time类型，直接返回
	if t, ok := value.(time.Time); ok {
		return t
	}
	// 如果value是指针类型，且指向time.Time类型，直接返回
	if t, ok := value.(*time.Time); ok {
		if t != nil {
			return *t
		}
		return time.Time{}
	}
	// 如果是int64类型，认为是时间戳，直接转换成time.Time
	if ts, ok := value.(int64); ok {
		return time.Unix(ts, 0)
	}
	if ts, ok := value.(int); ok {
		return time.Unix(int64(ts), 0)
	}
	//如果非string类型，返回空time.Time
	strValue, ok := value.(string)
	if !ok {
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
	switch strValue {
	case "now":
		return time.Now().In(loc)
	case "yesterday":
		return time.Now().AddDate(0, 0, -1).In(loc)
	case "tomorrow":
		return time.Now().AddDate(0, 0, 1).In(loc)
	}
	for i := range defaultLayouts {
		if tt, err := time.ParseInLocation(defaultLayouts[i], strValue, loc); err == nil {
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
