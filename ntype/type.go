package ntype

import (
	"encoding/json"
	"fmt"
	"html/template"
	"strconv"
	"strings"
	"time"
)

type Number interface {
	int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64 | float32 | float64
}

type integer interface {
	int | int8 | int16 | int32 | int64
}

type unsigned interface {
	uint | uint8 | uint16 | uint32 | uint64
}

type float interface {
	float32 | float64
}

// ToAnySlice 将一个切片转换为 []any 类型
func ToAnySlice[T comparable](slice []T) []any {
	result := make([]any, len(slice))
	for i, v := range slice {
		switch val := any(v).(type) {
		case fmt.Stringer:
			result[i] = val.String() // 转换实现了 Stringer 接口的类型
		default:
			result[i] = fmt.Sprintf("%v", val) // 通用格式化
		}
	}
	return result
}

func ToString(i any) string {
	switch s := i.(type) {
	case string:
		return s
	case []string:
		return strings.Join(s, ",")
	case bool:
		return strconv.FormatBool(s)
	case float64:
		return strconv.FormatFloat(s, 'f', -1, 64)
	case float32:
		return strconv.FormatFloat(float64(s), 'f', -1, 32)
	case int:
		return strconv.Itoa(s)
	case int8:
		return strconv.FormatInt(int64(s), 10)
	case int16:
		return strconv.FormatInt(int64(s), 10)
	case int32:
		return strconv.FormatInt(int64(s), 10)
	case int64:
		return strconv.FormatInt(s, 10)
	case uint:
		return strconv.FormatUint(uint64(s), 10)
	case uint8:
		return strconv.FormatUint(uint64(s), 10)
	case uint16:
		return strconv.FormatUint(uint64(s), 10)
	case uint32:
		return strconv.FormatUint(uint64(s), 10)
	case uint64:
		return strconv.FormatUint(s, 10)
	case json.Number:
		return s.String()
	case []byte:
		return string(s)
	case template.HTML:
		return string(s)
	case template.URL:
		return string(s)
	case template.JS:
		return string(s)
	case template.CSS:
		return string(s)
	case template.HTMLAttr:
		return string(s)
	case nil:
		return ""
	case fmt.Stringer:
		return s.String()
	case error:
		return s.Error()
	default:
		return ""
	}
}

func ToInt(i any) int {
	toNumber, ok := toNumber[int](i)
	if ok {
		return toNumber
	}
	return 0
}

func ToInt32(i any) int32 {
	toNumber, ok := toNumber[int32](i)
	if ok {
		return toNumber
	}
	return 0
}

func ToInt64(i any) int64 {
	toNumber, ok := toNumber[int64](i)
	if ok {
		return toNumber
	}
	return 0
}

func ToFloat32(i any) float32 {
	toNumber, ok := toNumber[float32](i)
	if ok {
		return toNumber
	}
	return 0
}

func ToFloat64(i any) float64 {
	toNumber, ok := toNumber[float64](i)
	if ok {
		return toNumber
	}
	return 0
}

func toNumber[T Number](i any) (T, bool) {
	i, _ = indirect(i)
	switch s := i.(type) {
	case T:
		return s, true
	case int:
		return T(s), true
	case int8:
		return T(s), true
	case int16:
		return T(s), true
	case int32:
		return T(s), true
	case int64:
		return T(s), true
	case uint:
		return T(s), true
	case uint8:
		return T(s), true
	case uint16:
		return T(s), true
	case uint32:
		return T(s), true
	case uint64:
		return T(s), true
	case float32:
		return T(s), true
	case float64:
		return T(s), true
	case string:
		if strings.Contains(s, ".") {
			if num, err := strconv.ParseFloat(s, 64); err == nil {
				return T(num), true
			}
		} else {
			if num, err := strconv.ParseInt(s, 10, 64); err == nil {
				return T(num), true
			}
		}
		return 0, false
	case bool:
		if s {
			return 1, true
		}
		return 0, true
	case nil:
		return 0, true
	case time.Weekday:
		return T(s), true
	case time.Month:
		return T(s), true
	}
	return 0, false
}
func JoinAny[T any](slice []T, sep string) string {
	parts := make([]string, len(slice))
	for i, v := range slice {
		parts[i] = fmt.Sprint(v)
	}
	return strings.Join(parts, sep)
}

func ToIntSlice(slice []string) []int32 {
	parts := make([]int32, len(slice))
	for i, v := range slice {
		parts[i] = ToInt32(v)
	}
	return parts
}
