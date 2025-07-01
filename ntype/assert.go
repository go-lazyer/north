package ntype

func IsBool(v any) bool {
	switch v.(type) {
	case bool:
		return true
	default:
		return false
	}
}

func IsInt(v any) bool {
	switch v.(type) {
	case int, int8, int16, int32, int64:
		return true
	default:
		return false
	}
}

func IsMap(v any) bool {
	switch v.(type) {
	case map[any]any, map[string]int, map[int]string: // 常见 map 类型
		return true
	default:
		return false
	}
}

func IsSlice(v any) bool {
	switch v.(type) {
	case []any, []int, []string: // 常见 slice 类型
		return true
	default:
		return false
	}
}

func IsPointer(v any) bool {
	switch v.(type) {
	case *int, *string, *struct{}: // 常见指针类型
		return true
	default:
		return false
	}
}

// IsNumeric 检查一个值是否为数字类型
func IsNumeric(inter any) bool {
	if inter == nil {
		return false
	}
	switch inter.(type) {
	case int:
		return true
	case int8:
		return true
	case int16:
		return true
	case int32:
		return true
	case int64:
		return true
	case float32:
		return true
	case float64:
		return true
	default:
		return false
	}
}
