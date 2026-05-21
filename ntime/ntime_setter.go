package ntime

import "time"

// SetDateTime 设置年、月、日、时、分、秒
func SetDateTime(t time.Time, year, month, day, hour, minute, second int) time.Time {
	return time.Date(year, time.Month(month), day, hour, minute, second, t.Nanosecond(), t.Location())
}

// SetDateTimeMilli 设置年、月、日、时、分、秒、毫秒
func SetDateTimeMilli(t time.Time, year, month, day, hour, minute, second, millisecond int) time.Time {
	return time.Date(year, time.Month(month), day, hour, minute, second, millisecond*1e6, t.Location())
}

// SetDateTimeMicro 设置年、月、日、时、分、秒、微秒
func SetDateTimeMicro(t time.Time, year, month, day, hour, minute, second, microsecond int) time.Time {
	return time.Date(year, time.Month(month), day, hour, minute, second, microsecond*1e3, t.Location())
}

// SetDateTimeNano 设置年、月、日、时、分、秒、纳秒
func SetDateTimeNano(t time.Time, year, month, day, hour, minute, second, nanosecond int) time.Time {
	return time.Date(year, time.Month(month), day, hour, minute, second, nanosecond, t.Location())
}

// SetDate 设置年、月、日
func SetDate(t time.Time, year, month, day int) time.Time {
	return time.Date(year, time.Month(month), day, t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), t.Location())
}

// SetDateMilli 设置年、月、日、毫秒
func SetDateMilli(t time.Time, year, month, day, millisecond int) time.Time {
	return time.Date(year, time.Month(month), day, t.Hour(), t.Minute(), t.Second(), millisecond*1e6, t.Location())
}

// SetDateMicro 设置年、月、日、微秒
func SetDateMicro(t time.Time, year, month, day, microsecond int) time.Time {
	return time.Date(year, time.Month(month), day, t.Hour(), t.Minute(), t.Second(), microsecond*1e3, t.Location())
}

// SetDateNano 设置年、月、日、纳秒
func SetDateNano(t time.Time, year, month, day, nanosecond int) time.Time {
	return time.Date(year, time.Month(month), day, t.Hour(), t.Minute(), t.Second(), nanosecond, t.Location())
}

// SetTime 设置时、分、秒
func SetTime(t time.Time, hour, minute, second int) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), hour, minute, second, t.Nanosecond(), t.Location())
}

// SetTimeMilli 设置时、分、秒、毫秒
func SetTimeMilli(t time.Time, hour, minute, second, millisecond int) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), hour, minute, second, millisecond*1e6, t.Location())
}

// SetTimeMicro 设置时、分、秒、微秒
func SetTimeMicro(t time.Time, hour, minute, second, microsecond int) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), hour, minute, second, microsecond*1e3, t.Location())
}

// SetTimeNano 设置时、分、秒、纳秒
func SetTimeNano(t time.Time, hour, minute, second, nanosecond int) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), hour, minute, second, nanosecond, t.Location())
}

// SetYear 设置年份（月份不溢出）
func SetYear(t time.Time, year int) time.Time {
	newT := time.Date(year, t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), t.Location())
	// 如果日期溢出（如 2月29日 -> 非闰年），Go 会自动进位，需要回退到该月最后一天
	if newT.Month() != t.Month() {
		newT = time.Date(year, t.Month()+1, 0, t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), t.Location())
	}
	return newT
}

// SetMonth 设置月份（月份不溢出）
func SetMonth(t time.Time, month int) time.Time {
	newT := time.Date(t.Year(), time.Month(month), t.Day(), t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), t.Location())
	// 如果日期溢出（如 1月31日 -> 2月），Go 会自动进位，需要回退到该月最后一天
	if newT.Month() != time.Month(month) {
		newT = time.Date(t.Year(), time.Month(month)+1, 0, t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), t.Location())
	}
	return newT
}

// SetDay 设置日期
func SetDay(t time.Time, day int) time.Time {
	return time.Date(t.Year(), t.Month(), day, t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), t.Location())
}

// SetHour 设置小时
func SetHour(t time.Time, hour int) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), hour, t.Minute(), t.Second(), t.Nanosecond(), t.Location())
}

// SetMinute 设置分钟
func SetMinute(t time.Time, minute int) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), minute, t.Second(), t.Nanosecond(), t.Location())
}

// SetSecond 设置秒
func SetSecond(t time.Time, second int) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), second, t.Nanosecond(), t.Location())
}

// SetMillisecond 设置毫秒
func SetMillisecond(t time.Time, millisecond int) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), millisecond*1e6, t.Location())
}

// SetMicrosecond 设置微秒
func SetMicrosecond(t time.Time, microsecond int) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), microsecond*1e3, t.Location())
}

// SetNanosecond 设置纳秒
func SetNanosecond(t time.Time, nanosecond int) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), nanosecond, t.Location())
}
