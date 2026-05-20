package ntime_test

import (
	"testing"
	"time"

	"github.com/go-lazyer/north/ntime"
)

// 基础时间：2024-02-29 15:30:45.123456789 (闰年2月29日，方便测试溢出)
var baseTime = time.Date(2024, 2, 29, 15, 30, 45, 123456789, time.Local)

func TestSetDateTime(t *testing.T) {
	result := ntime.SetDateTime(baseTime, 2023, 12, 25, 10, 20, 30)
	expected := time.Date(2023, 12, 25, 10, 20, 30, 123456789, time.Local)
	if !result.Equal(expected) {
		t.Errorf("SetDateTime expected %v, got %v", expected, result)
	}
}

func TestSetDateTimeMilli(t *testing.T) {
	result := ntime.SetDateTimeMilli(baseTime, 2023, 12, 25, 10, 20, 30, 500)
	expected := time.Date(2023, 12, 25, 10, 20, 30, 500*1e6, time.Local)
	if !result.Equal(expected) {
		t.Errorf("SetDateTimeMilli expected %v, got %v", expected, result)
	}
}

func TestSetDateTimeMicro(t *testing.T) {
	result := ntime.SetDateTimeMicro(baseTime, 2023, 12, 25, 10, 20, 30, 500)
	expected := time.Date(2023, 12, 25, 10, 20, 30, 500*1e3, time.Local)
	if !result.Equal(expected) {
		t.Errorf("SetDateTimeMicro expected %v, got %v", expected, result)
	}
}

func TestSetDateTimeNano(t *testing.T) {
	result := ntime.SetDateTimeNano(baseTime, 2023, 12, 25, 10, 20, 30, 999999999)
	expected := time.Date(2023, 12, 25, 10, 20, 30, 999999999, time.Local)
	if !result.Equal(expected) {
		t.Errorf("SetDateTimeNano expected %v, got %v", expected, result)
	}
}

func TestSetDate(t *testing.T) {
	result := ntime.SetDate(baseTime, 2023, 6, 15)
	expected := time.Date(2023, 6, 15, 15, 30, 45, 123456789, time.Local)
	if !result.Equal(expected) {
		t.Errorf("SetDate expected %v, got %v", expected, result)
	}
}

func TestSetDateMilli(t *testing.T) {
	result := ntime.SetDateMilli(baseTime, 2023, 6, 15, 500)
	expected := time.Date(2023, 6, 15, 15, 30, 45, 500*1e6, time.Local)
	if !result.Equal(expected) {
		t.Errorf("SetDateMilli expected %v, got %v", expected, result)
	}
}

func TestSetDateMicro(t *testing.T) {
	result := ntime.SetDateMicro(baseTime, 2023, 6, 15, 500)
	expected := time.Date(2023, 6, 15, 15, 30, 45, 500*1e3, time.Local)
	if !result.Equal(expected) {
		t.Errorf("SetDateMicro expected %v, got %v", expected, result)
	}
}

func TestSetDateNano(t *testing.T) {
	result := ntime.SetDateNano(baseTime, 2023, 6, 15, 999999999)
	expected := time.Date(2023, 6, 15, 15, 30, 45, 999999999, time.Local)
	if !result.Equal(expected) {
		t.Errorf("SetDateNano expected %v, got %v", expected, result)
	}
}

func TestSetTime(t *testing.T) {
	result := ntime.SetTime(baseTime, 8, 0, 0)
	expected := time.Date(2024, 2, 29, 8, 0, 0, 123456789, time.Local)
	if !result.Equal(expected) {
		t.Errorf("SetTime expected %v, got %v", expected, result)
	}
}

func TestSetTimeMilli(t *testing.T) {
	result := ntime.SetTimeMilli(baseTime, 8, 0, 0, 500)
	expected := time.Date(2024, 2, 29, 8, 0, 0, 500*1e6, time.Local)
	if !result.Equal(expected) {
		t.Errorf("SetTimeMilli expected %v, got %v", expected, result)
	}
}

func TestSetTimeMicro(t *testing.T) {
	result := ntime.SetTimeMicro(baseTime, 8, 0, 0, 500)
	expected := time.Date(2024, 2, 29, 8, 0, 0, 500*1e3, time.Local)
	if !result.Equal(expected) {
		t.Errorf("SetTimeMicro expected %v, got %v", expected, result)
	}
}

func TestSetTimeNano(t *testing.T) {
	result := ntime.SetTimeNano(baseTime, 8, 0, 0, 0)
	expected := time.Date(2024, 2, 29, 8, 0, 0, 0, time.Local)
	if !result.Equal(expected) {
		t.Errorf("SetTimeNano expected %v, got %v", expected, result)
	}
}

func TestSetYear(t *testing.T) {
	// 正常设置年份
	result := ntime.SetYear(baseTime, 2023)
	expected := time.Date(2023, 2, 28, 15, 30, 45, 123456789, time.Local) // 2023非闰年，2月29日溢出到2月28日
	if !result.Equal(expected) {
		t.Errorf("SetYear(闰年2月29日->非闰年) expected %v, got %v", expected, result)
	}

	// 闰年到闰年，不溢出
	leapTime := time.Date(2024, 2, 29, 0, 0, 0, 0, time.Local)
	result = ntime.SetYear(leapTime, 2028)
	expected = time.Date(2028, 2, 29, 0, 0, 0, 0, time.Local)
	if !result.Equal(expected) {
		t.Errorf("SetYear(闰年->闰年) expected %v, got %v", expected, result)
	}

	// 非闰年正常日期，不溢出
	normalTime := time.Date(2023, 3, 15, 0, 0, 0, 0, time.Local)
	result = ntime.SetYear(normalTime, 2025)
	expected = time.Date(2025, 3, 15, 0, 0, 0, 0, time.Local)
	if !result.Equal(expected) {
		t.Errorf("SetYear(正常日期) expected %v, got %v", expected, result)
	}
}

func TestSetMonth(t *testing.T) {
	// 正常设置月份
	result := ntime.SetMonth(baseTime, 6)
	expected := time.Date(2024, 6, 29, 15, 30, 45, 123456789, time.Local)
	if !result.Equal(expected) {
		t.Errorf("SetMonth expected %v, got %v", expected, result)
	}

	// 1月31日 -> 2月，日期溢出，应回退到2月28日（2024闰年29日）
	jan31 := time.Date(2024, 1, 31, 0, 0, 0, 0, time.Local)
	result = ntime.SetMonth(jan31, 2)
	expected = time.Date(2024, 2, 29, 0, 0, 0, 0, time.Local) // 闰年2月有29天
	if !result.Equal(expected) {
		t.Errorf("SetMonth(1月31日->2月, 闰年) expected %v, got %v", expected, result)
	}

	// 1月31日 -> 2月，非闰年溢出
	jan31_2023 := time.Date(2023, 1, 31, 0, 0, 0, 0, time.Local)
	result = ntime.SetMonth(jan31_2023, 2)
	expected = time.Date(2023, 2, 28, 0, 0, 0, 0, time.Local) // 非闰年2月有28天
	if !result.Equal(expected) {
		t.Errorf("SetMonth(1月31日->2月, 非闰年) expected %v, got %v", expected, result)
	}

	// 3月31日 -> 4月，4月只有30天，溢出
	mar31 := time.Date(2024, 3, 31, 0, 0, 0, 0, time.Local)
	result = ntime.SetMonth(mar31, 4)
	expected = time.Date(2024, 4, 30, 0, 0, 0, 0, time.Local)
	if !result.Equal(expected) {
		t.Errorf("SetMonth(3月31日->4月) expected %v, got %v", expected, result)
	}
}

func TestSetDay(t *testing.T) {
	result := ntime.SetDay(baseTime, 15)
	expected := time.Date(2024, 2, 15, 15, 30, 45, 123456789, time.Local)
	if !result.Equal(expected) {
		t.Errorf("SetDay expected %v, got %v", expected, result)
	}
}

func TestSetHour(t *testing.T) {
	result := ntime.SetHour(baseTime, 8)
	expected := time.Date(2024, 2, 29, 8, 30, 45, 123456789, time.Local)
	if !result.Equal(expected) {
		t.Errorf("SetHour expected %v, got %v", expected, result)
	}
}

func TestSetMinute(t *testing.T) {
	result := ntime.SetMinute(baseTime, 0)
	expected := time.Date(2024, 2, 29, 15, 0, 45, 123456789, time.Local)
	if !result.Equal(expected) {
		t.Errorf("SetMinute expected %v, got %v", expected, result)
	}
}

func TestSetSecond(t *testing.T) {
	result := ntime.SetSecond(baseTime, 0)
	expected := time.Date(2024, 2, 29, 15, 30, 0, 123456789, time.Local)
	if !result.Equal(expected) {
		t.Errorf("SetSecond expected %v, got %v", expected, result)
	}
}

func TestSetMillisecond(t *testing.T) {
	result := ntime.SetMillisecond(baseTime, 500)
	expected := time.Date(2024, 2, 29, 15, 30, 45, 500*1e6, time.Local)
	if !result.Equal(expected) {
		t.Errorf("SetMillisecond expected %v, got %v", expected, result)
	}
}

func TestSetMicrosecond(t *testing.T) {
	result := ntime.SetMicrosecond(baseTime, 500)
	expected := time.Date(2024, 2, 29, 15, 30, 45, 500*1e3, time.Local)
	if !result.Equal(expected) {
		t.Errorf("SetMicrosecond expected %v, got %v", expected, result)
	}
}

func TestSetNanosecond(t *testing.T) {
	result := ntime.SetNanosecond(baseTime, 0)
	expected := time.Date(2024, 2, 29, 15, 30, 45, 0, time.Local)
	if !result.Equal(expected) {
		t.Errorf("SetNanosecond expected %v, got %v", expected, result)
	}
}

func TestSetWeekStartsAt(t *testing.T) {
	// 2024-03-07 是周四，以周一为起始，周一应为 2024-03-04
	thu := time.Date(2024, 3, 7, 10, 0, 0, 0, time.Local)
	result := ntime.SetWeekStartsAt(thu, ntime.Monday)
	expected := time.Date(2024, 3, 4, 10, 0, 0, 0, time.Local)
	if !result.Equal(expected) {
		t.Errorf("SetWeekStartsAt(周四, 周一起始) expected %v, got %v", expected, result)
	}

	// 以周日为起始，周日应为 2024-03-03
	result = ntime.SetWeekStartsAt(thu, ntime.Sunday)
	expected = time.Date(2024, 3, 3, 10, 0, 0, 0, time.Local)
	if !result.Equal(expected) {
		t.Errorf("SetWeekStartsAt(周四, 周日起始) expected %v, got %v", expected, result)
	}
}
