package nstring

import (
	"testing"
)

func TestToUpperCamelCase(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"hello_world", "HelloWorld"},
		{"user_name", "UserName"},
		{"a_b_c", "ABC"},
		{"single", "Single"},
		{"", ""},
		{"already", "Already"},
		{"_leading", "Leading"},
		{"trailing_", "Trailing"},
	}

	for _, tt := range tests {
		result := ToUpperCamelCase(tt.input)
		if result != tt.expected {
			t.Errorf("ToUpperCamelCase(%q) 期望 %q，实际 %q", tt.input, tt.expected, result)
		}
	}
}

func TestToLowerCamelCase(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"hello_world", "helloWorld"},
		{"user_name", "userName"},
		{"a_b_c", "aBC"},
		{"single", "single"},
		{"", ""},
		{"already", "already"},
		{"_leading", "Leading"},
		{"trailing_", "trailing"},
	}

	for _, tt := range tests {
		result := ToLowerCamelCase(tt.input)
		if result != tt.expected {
			t.Errorf("ToLowerCamelCase(%q) 期望 %q，实际 %q", tt.input, tt.expected, result)
		}
	}
}

func TestToUpperFirst(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"hello", "Hello"},
		{"world", "World"},
		{"a", "A"},
		{"", ""},
		{"Already", "Already"},
		{"123abc", "123abc"},
	}

	for _, tt := range tests {
		result := ToUpperFirst(tt.input)
		if result != tt.expected {
			t.Errorf("ToUpperFirst(%q) 期望 %q，实际 %q", tt.input, tt.expected, result)
		}
	}
}

func TestCompactJson(t *testing.T) {
	tests := []struct {
		input    string
		expected string
		hasError bool
	}{
		// 正常 JSON 压缩
		{
			`{"name": "张三", "age": 25}`,
			`{"name":"张三","age":25}`,
			false,
		},
		// 带换行的 JSON
		{
			"{\n  \"key\": \"value\"\n}",
			`{"key":"value"}`,
			false,
		},
		// 空字符串
		{
			"",
			"",
			false,
		},
		// 无效 JSON
		{
			"not json",
			"",
			true,
		},
		// 嵌套 JSON
		{
			`{"a": {"b": 1}}`,
			`{"a":{"b":1}}`,
			false,
		},
	}

	for _, tt := range tests {
		result, err := CompactJson(tt.input)
		if tt.hasError {
			if err == nil {
				t.Errorf("CompactJson(%q) 期望返回错误，实际无错误", tt.input)
			}
		} else {
			if err != nil {
				t.Errorf("CompactJson(%q) 未期望错误，实际 %v", tt.input, err)
			}
			if result != tt.expected {
				t.Errorf("CompactJson(%q) 期望 %q，实际 %q", tt.input, tt.expected, result)
			}
		}
	}
}
