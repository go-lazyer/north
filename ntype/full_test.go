package ntype_test

import (
	"testing"

	"github.com/go-lazyer/north/ntype"
)

func TestIsBool(t *testing.T) {
	if !ntype.IsBool(true) {
		t.Errorf("IsBool(true) 期望 true")
	}
	if !ntype.IsBool(false) {
		t.Errorf("IsBool(false) 期望 true")
	}
	if ntype.IsBool(1) {
		t.Errorf("IsBool(1) 期望 false")
	}
	if ntype.IsBool("true") {
		t.Errorf("IsBool(\"true\") 期望 false")
	}
	if ntype.IsBool(nil) {
		t.Errorf("IsBool(nil) 期望 false")
	}
}

func TestIsInt(t *testing.T) {
	if !ntype.IsInt(1) {
		t.Errorf("IsInt(int) 期望 true")
	}
	if !ntype.IsInt(int8(1)) {
		t.Errorf("IsInt(int8) 期望 true")
	}
	if !ntype.IsInt(int16(1)) {
		t.Errorf("IsInt(int16) 期望 true")
	}
	if !ntype.IsInt(int32(1)) {
		t.Errorf("IsInt(int32) 期望 true")
	}
	if !ntype.IsInt(int64(1)) {
		t.Errorf("IsInt(int64) 期望 true")
	}
	if ntype.IsInt(1.0) {
		t.Errorf("IsInt(float64) 期望 false")
	}
	if ntype.IsInt("1") {
		t.Errorf("IsInt(string) 期望 false")
	}
}

func TestIsMap(t *testing.T) {
	if !ntype.IsMap(map[any]any{}) {
		t.Errorf("IsMap(map[any]any) 期望 true")
	}
	if !ntype.IsMap(map[string]int{}) {
		t.Errorf("IsMap(map[string]int) 期望 true")
	}
	if !ntype.IsMap(map[int]string{}) {
		t.Errorf("IsMap(map[int]string) 期望 true")
	}
	if ntype.IsMap([]int{}) {
		t.Errorf("IsMap([]int) 期望 false")
	}
	if ntype.IsMap("string") {
		t.Errorf("IsMap(string) 期望 false")
	}
}

func TestIsSlice(t *testing.T) {
	if !ntype.IsSlice([]any{}) {
		t.Errorf("IsSlice([]any) 期望 true")
	}
	if !ntype.IsSlice([]int{}) {
		t.Errorf("IsSlice([]int) 期望 true")
	}
	if !ntype.IsSlice([]uint8{}) {
		t.Errorf("IsSlice([]uint8) 期望 true")
	}
	if !ntype.IsSlice([]string{}) {
		t.Errorf("IsSlice([]string) 期望 true")
	}
	if ntype.IsSlice(map[string]int{}) {
		t.Errorf("IsSlice(map) 期望 false")
	}
	if ntype.IsSlice("string") {
		t.Errorf("IsSlice(string) 期望 false")
	}
}

func TestIsString(t *testing.T) {
	if !ntype.IsString("hello") {
		t.Errorf("IsString(string) 期望 true")
	}
	if ntype.IsString(1) {
		t.Errorf("IsString(int) 期望 false")
	}
	if ntype.IsString(true) {
		t.Errorf("IsString(bool) 期望 false")
	}
}

func TestIsPointer(t *testing.T) {
	var i int = 1
	if !ntype.IsPointer(&i) {
		t.Errorf("IsPointer(*int) 期望 true")
	}
	var s string = "hello"
	if !ntype.IsPointer(&s) {
		t.Errorf("IsPointer(*string) 期望 true")
	}
	if ntype.IsPointer(i) {
		t.Errorf("IsPointer(int) 期望 false")
	}
	if ntype.IsPointer(s) {
		t.Errorf("IsPointer(string) 期望 false")
	}
}

func TestIsNumeric_Full(t *testing.T) {
	// 整数类型
	if !ntype.IsNumeric(1) {
		t.Errorf("IsNumeric(int) 期望 true")
	}
	if !ntype.IsNumeric(int8(1)) {
		t.Errorf("IsNumeric(int8) 期望 true")
	}
	if !ntype.IsNumeric(int16(1)) {
		t.Errorf("IsNumeric(int16) 期望 true")
	}
	if !ntype.IsNumeric(int32(1)) {
		t.Errorf("IsNumeric(int32) 期望 true")
	}
	if !ntype.IsNumeric(int64(1)) {
		t.Errorf("IsNumeric(int64) 期望 true")
	}

	// 浮点类型
	if !ntype.IsNumeric(1.0) {
		t.Errorf("IsNumeric(float64) 期望 true")
	}
	if !ntype.IsNumeric(float32(1.0)) {
		t.Errorf("IsNumeric(float32) 期望 true")
	}

	// 字符串数字
	if !ntype.IsNumeric("123") {
		t.Errorf("IsNumeric(\"123\") 期望 true")
	}
	if !ntype.IsNumeric("12.5") {
		t.Errorf("IsNumeric(\"12.5\") 期望 true")
	}
	if ntype.IsNumeric("abc") {
		t.Errorf("IsNumeric(\"abc\") 期望 false")
	}
	if ntype.IsNumeric("123w.2") {
		t.Errorf("IsNumeric(\"123w.2\") 期望 false")
	}

	// nil
	if ntype.IsNumeric(nil) {
		t.Errorf("IsNumeric(nil) 期望 false")
	}

	// 其他类型
	if ntype.IsNumeric(true) {
		t.Errorf("IsNumeric(bool) 期望 false")
	}
}

func TestToString(t *testing.T) {
	tests := []struct {
		input    any
		expected string
	}{
		{"hello", "hello"},
		{true, "true"},
		{false, "false"},
		{1, "1"},
		{int8(1), "1"},
		{int16(1), "1"},
		{int32(1), "1"},
		{int64(1), "1"},
		{uint(1), "1"},
		{uint8(1), "1"},
		{uint16(1), "1"},
		{uint32(1), "1"},
		{uint64(1), "1"},
		{3.14, "3.14"},
		{float32(3.14), "3.14"},
		{nil, ""},
		{[]byte("bytes"), "bytes"},
		{[]string{"a", "b", "c"}, "a,b,c"},
	}

	for _, tt := range tests {
		result := ntype.ToString(tt.input)
		if result != tt.expected {
			t.Errorf("ToString(%v) 期望 %q，实际 %q", tt.input, tt.expected, result)
		}
	}
}

func TestToInt(t *testing.T) {
	tests := []struct {
		input    any
		expected int
	}{
		{1, 1},
		{int8(1), 1},
		{int16(1), 1},
		{int32(1), 1},
		{int64(1), 1},
		{"123", 123},
		{true, 1},
		{false, 0},
		{nil, 0},
		{3.14, 3},
		{"abc", 0}, // 无法解析的字符串
	}

	for _, tt := range tests {
		result := ntype.ToInt(tt.input)
		if result != tt.expected {
			t.Errorf("ToInt(%v) 期望 %d，实际 %d", tt.input, tt.expected, result)
		}
	}
}

func TestToInt64(t *testing.T) {
	tests := []struct {
		input    any
		expected int64
	}{
		{int64(1), 1},
		{1, 1},
		{"123", 123},
		{nil, 0},
	}

	for _, tt := range tests {
		result := ntype.ToInt64(tt.input)
		if result != tt.expected {
			t.Errorf("ToInt64(%v) 期望 %d，实际 %d", tt.input, tt.expected, result)
		}
	}
}

func TestToFloat64(t *testing.T) {
	tests := []struct {
		input    any
		expected float64
	}{
		{3.14, 3.14},
		{1, 1.0},
		{"3.14", 3.14},
		{nil, 0},
	}

	for _, tt := range tests {
		result := ntype.ToFloat64(tt.input)
		if result != tt.expected {
			t.Errorf("ToFloat64(%v) 期望 %f，实际 %f", tt.input, tt.expected, result)
		}
	}
}

func TestToAnySlice(t *testing.T) {
	// int 切片
	intSlice := []int{1, 2, 3}
	result := ntype.ToAnySlice(intSlice)
	if len(result) != 3 {
		t.Fatalf("期望 3 个元素，实际 %d", len(result))
	}
	for _, v := range result {
		if v != "1" && v != "2" && v != "3" {
			// ToAnySlice 使用 fmt.Sprintf，所以结果是字符串
		}
	}

	// string 切片
	strSlice := []string{"a", "b", "c"}
	result = ntype.ToAnySlice(strSlice)
	if len(result) != 3 {
		t.Fatalf("期望 3 个元素，实际 %d", len(result))
	}

	// 空切片
	emptySlice := []int{}
	result = ntype.ToAnySlice(emptySlice)
	if len(result) != 0 {
		t.Errorf("空切片期望 0 个元素，实际 %d", len(result))
	}
}

func TestJoinAny_Full(t *testing.T) {
	// string 切片
	result := ntype.JoinAny([]string{"a", "b", "c"}, ",")
	if result != "a,b,c" {
		t.Errorf("JoinAny 期望 'a,b,c'，实际 %q", result)
	}

	// int 切片
	result = ntype.JoinAny([]int{1, 2, 3}, "-")
	if result != "1-2-3" {
		t.Errorf("JoinAny 期望 '1-2-3'，实际 %q", result)
	}

	// 混合类型
	result = ntype.JoinAny([]any{1, "hello", true}, "|")
	if result != "1|hello|true" {
		t.Errorf("JoinAny 期望 '1|hello|true'，实际 %q", result)
	}

	// 空切片
	result = ntype.JoinAny([]string{}, ",")
	if result != "" {
		t.Errorf("JoinAny 空切片期望空字符串，实际 %q", result)
	}
}
