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

// 返回t所在周的weekday天是第几天
func DayOfWeek(t time.Time) int {
	weekday := int(t.Weekday())
	if weekday == 0 {
		return 7
	}
	return weekday
}
