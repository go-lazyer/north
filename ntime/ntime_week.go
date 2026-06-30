package ntime

import "time"

const (
	Monday    = 1 //周一
	Tuesday   = 2 //周二
	Wednesday = 3 //周三
	Thursday  = 4 //周四
	Friday    = 5 //周五
	Saturday  = 6 //周六
	Sunday    = 0 //周日
)

// StartOfWeek 返回t所在周的开始时间，默认周一为一周的第一天
func StartOfWeek(t time.Time) time.Time {
	weekday := int(t.Weekday())
	d := t.AddDate(0, 0, -(weekday-1+7)%7)
	return time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, d.Location())
}

// EndOfWeek 返回t所在周的结束时间，默认周一为一周的第一天
func EndOfWeek(t time.Time) time.Time {
	weekday := int(t.Weekday())
	d := t.AddDate(0, 0, (7-weekday)%7)
	return time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, d.Location())
}
