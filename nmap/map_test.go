package nmap

import (
	"encoding/json"
	"testing"
)

func TestToString_Nil(t *testing.T) {
	result := ToString(nil)
	if result != "" {
		t.Fatalf("期望空字符串，实际为 %s", result)
	}
}

func TestToString_Empty(t *testing.T) {
	result := ToString(map[string]any{})
	if result != "" {
		t.Fatalf("期望空字符串，实际为 %s", result)
	}
}

func TestToString_Normal(t *testing.T) {
	m := map[string]any{
		"city":       "北京",
		"manualCity": "北京",
		"token":      "abc123",
	}
	result := ToString(m)
	if !json.Valid([]byte(result)) && result != "" {
		// 非JSON格式，检查是否包含预期键值
	}
	if result == "" {
		t.Fatal("期望非空字符串")
	}
	// 检查包含所有键
	if !containsKey(result, "city") {
		t.Fatalf("期望包含 city，实际为 %s", result)
	}
	if !containsKey(result, "manualCity") {
		t.Fatalf("期望包含 manualCity，实际为 %s", result)
	}
	if !containsKey(result, "token") {
		t.Fatalf("期望包含 token，实际为 %s", result)
	}
}

func TestToString_NilValue(t *testing.T) {
	m := map[string]any{
		"empty": nil,
	}
	result := ToString(m)
	if result != "empty=" {
		t.Fatalf("期望 empty=，实际为 %s", result)
	}
}

func containsKey(s, key string) bool {
	return len(s) >= len(key) && s[:len(key)] == key || len(s) > len(key) && containsKey(s[1:], key)
}

func TestKeys(t *testing.T) {
	m := map[string]int{
		"a": 1,
		"b": 2,
		"c": 3,
	}
	keys := Keys(m)
	if len(keys) != 3 {
		t.Fatalf("期望 3 个 key，实际为 %d", len(keys))
	}
	// 检查所有 key 都存在
	keySet := map[string]bool{}
	for _, k := range keys {
		keySet[k] = true
	}
	for k := range m {
		if !keySet[k] {
			t.Fatalf("期望包含 key %s", k)
		}
	}
}

func TestKeys_Empty(t *testing.T) {
	keys := Keys(map[string]int{})
	if len(keys) != 0 {
		t.Fatalf("期望 0 个 key，实际为 %d", len(keys))
	}
}

func TestKeys_Int(t *testing.T) {
	m := map[int]string{
		1: "one",
		2: "two",
	}
	keys := Keys(m)
	if len(keys) != 2 {
		t.Fatalf("期望 2 个 key，实际为 %d", len(keys))
	}
	keySet := map[int]bool{}
	for _, k := range keys {
		keySet[k] = true
	}
	for k := range m {
		if !keySet[k] {
			t.Fatalf("期望包含 key %d", k)
		}
	}
}

type testStruct struct {
	Name  string `json:"name"`
	Age   int    `json:"age"`
	Email string `json:"email"`
}

func TestToMap(t *testing.T) {
	obj := testStruct{Name: "张三", Age: 25, Email: "test@example.com"}
	m, err := ToMap(obj)
	if err != nil {
		t.Fatalf("期望无错误，实际为 %v", err)
	}
	if m["name"] != "张三" {
		t.Fatalf("期望 name=张三，实际为 %v", m["name"])
	}
	if m["email"] != "test@example.com" {
		t.Fatalf("期望 email=test@example.com，实际为 %v", m["email"])
	}
	// Age 应为 json.Number 而非 float64
	ageNum, ok := m["age"].(json.Number)
	if !ok {
		t.Fatalf("期望 age 为 json.Number，实际为 %T", m["age"])
	}
	if ageNum.String() != "25" {
		t.Fatalf("期望 age=25，实际为 %s", ageNum.String())
	}
}

func TestToMap_NilInput(t *testing.T) {
	m, err := ToMap(nil)
	if err != nil {
		t.Fatalf("期望无错误，实际为 %v", err)
	}
	if len(m) != 0 {
		t.Fatalf("期望空 map，实际长度为 %d", len(m))
	}
}
