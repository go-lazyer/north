package ntime

import "time"

const (
	Sunday    = 0
	Monday    = 1
	Tuesday   = 2
	Wednesday = 3
	Thursday  = 4
	Friday    = 5
	Saturday  = 6
)

// GetDayOfWeekInWeek 返回给定时间所在周的某一天
//
// 参数：
//
//	t: 输入的时间
//	weekday: 要获取的星期几（0=周日, 1=周一, 2=周二, ..., 6=周六）
//	startDayOfWeek: 一周从哪天开始（0=周日, 1=周一, ..., 6=周六）
//
// 返回：所在周的 weekday 对应的日期
func GetDayOfWeek(t time.Time, weekday, startDayOfWeek int) time.Time {
	// 规范化输入（防止传入非法值）
	weekday = (weekday + 7) % 7
	startDayOfWeek = (startDayOfWeek + 7) % 7

	// 当前是星期几 (0=周日, ..., 6=周六)
	current := int(t.Weekday())

	// 计算当前日期距离本周开始（startDayOfWeek）偏移了多少天
	// (current - startDayOfWeek + 7) % 7 就是当前在本周的第几天（从0开始）
	// 所以，从当前时间回退这个天数，就得到本周的 startDayOfWeek
	daysFromStart := (current - startDayOfWeek + 7) % 7
	weekStart := t.AddDate(0, 0, -daysFromStart) // 本周的第一天（startDayOfWeek）

	// 从 weekStart 开始，加上 weekday 是星期几，得到目标日期
	// 例如：weekStart 是周一（1），要找周三（3），则加 3 天？不对！
	//
	// 注意：我们要找的是“统一标准下的 weekday”
	// 所以先算 weekday 距离 startDayOfWeek 有多少天
	daysToAdd := (weekday - startDayOfWeek + 7) % 7

	return weekStart.AddDate(0, 0, daysToAdd)
}
