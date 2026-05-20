package nmap

import (
	"strings"
	"testing"
)

func TestNewLinkMap(t *testing.T) {
	om := NewLinkMap[string, int]()
	if om == nil {
		t.Fatal("NewLinkMap 返回 nil")
	}
	if len(om.keys) != 0 {
		t.Fatalf("期望 keys 长度为 0，实际为 %d", len(om.keys))
	}
	if len(om.m) != 0 {
		t.Fatalf("期望 map 长度为 0，实际为 %d", len(om.m))
	}
}

func TestLinkMap_Set_NewKey(t *testing.T) {
	om := NewLinkMap[string, int]()
	om.Set("a", 1)
	om.Set("b", 2)
	om.Set("c", 3)

	if len(om.keys) != 3 {
		t.Fatalf("期望 keys 长度为 3，实际为 %d", len(om.keys))
	}
	if om.keys[0] != "a" || om.keys[1] != "b" || om.keys[2] != "c" {
		t.Fatalf("期望 keys 顺序为 [a,b,c]，实际为 %v", om.keys)
	}
}

func TestLinkMap_Set_UpdateExisting(t *testing.T) {
	om := NewLinkMap[string, int]()
	om.Set("a", 1)
	om.Set("b", 2)
	om.Set("a", 10) // 更新已有 key，顺序不变

	if len(om.keys) != 2 {
		t.Fatalf("期望 keys 长度为 2，实际为 %d", len(om.keys))
	}
	if om.keys[0] != "a" || om.keys[1] != "b" {
		t.Fatalf("期望 keys 顺序为 [a,b]，实际为 %v", om.keys)
	}
	v, ok := om.Get("a")
	if !ok || v != 10 {
		t.Fatalf("期望 a 的值为 10，实际为 %d", v)
	}
}

func TestLinkMap_Get_Existing(t *testing.T) {
	om := NewLinkMap[string, int]()
	om.Set("a", 1)
	v, ok := om.Get("a")
	if !ok {
		t.Fatal("期望 Get 返回 ok=true")
	}
	if v != 1 {
		t.Fatalf("期望值为 1，实际为 %d", v)
	}
}

func TestLinkMap_Get_NotExisting(t *testing.T) {
	om := NewLinkMap[string, int]()
	v, ok := om.Get("x")
	if ok {
		t.Fatal("期望 Get 返回 ok=false")
	}
	if v != 0 {
		t.Fatalf("期望零值为 0，实际为 %d", v)
	}
}

func TestLinkMap_Delete_Existing(t *testing.T) {
	om := NewLinkMap[string, int]()
	om.Set("a", 1)
	om.Set("b", 2)
	om.Set("c", 3)
	om.Delete("b")

	if len(om.keys) != 2 {
		t.Fatalf("期望 keys 长度为 2，实际为 %d", len(om.keys))
	}
	if om.keys[0] != "a" || om.keys[1] != "c" {
		t.Fatalf("期望 keys 顺序为 [a,c]，实际为 %v", om.keys)
	}
	_, ok := om.Get("b")
	if ok {
		t.Fatal("期望删除后 Get 返回 ok=false")
	}
}

func TestLinkMap_Delete_NotExisting(t *testing.T) {
	om := NewLinkMap[string, int]()
	om.Set("a", 1)
	om.Delete("x") // 删除不存在的 key，无影响

	if len(om.keys) != 1 {
		t.Fatalf("期望 keys 长度为 1，实际为 %d", len(om.keys))
	}
}

func TestLinkMap_Range_Order(t *testing.T) {
	om := NewLinkMap[string, int]()
	om.Set("b", 2)
	om.Set("a", 1)
	om.Set("c", 3)

	var result []string
	om.Range(func(k string, v int) {
		result = append(result, k)
	})

	expected := []string{"b", "a", "c"}
	if len(result) != len(expected) {
		t.Fatalf("期望遍历顺序长度为 %d，实际为 %d", len(expected), len(result))
	}
	for i, k := range result {
		if k != expected[i] {
			t.Fatalf("位置 %d：期望 %s，实际 %s", i, expected[i], k)
		}
	}
}

func TestLinkMap_Range_Empty(t *testing.T) {
	om := NewLinkMap[string, int]()
	called := false
	om.Range(func(k string, v int) {
		called = true
	})
	if called {
		t.Fatal("期望空 LinkMap 不调用 Range 回调")
	}
}

func TestLinkMap_ToString_DefaultSeps(t *testing.T) {
	om := NewLinkMap[string, int]()
	om.Set("a", 1)
	om.Set("b", 2)
	om.Set("c", 3)

	result := om.ToString()
	if !strings.Contains(result, "a:1") {
		t.Fatalf("期望包含 a:1，实际为 %s", result)
	}
	if !strings.Contains(result, "b:2") {
		t.Fatalf("期望包含 b:2，实际为 %s", result)
	}
	if !strings.Contains(result, "c:3") {
		t.Fatalf("期望包含 c:3，实际为 %s", result)
	}
}

func TestLinkMap_ToString_CustomSeps(t *testing.T) {
	om := NewLinkMap[string, int]()
	om.Set("a", 1)
	om.Set("b", 2)

	result := om.ToString("=", "&")
	if !strings.Contains(result, "a=1") {
		t.Fatalf("期望包含 a=1，实际为 %s", result)
	}
	if !strings.Contains(result, "b=2") {
		t.Fatalf("期望包含 b=2，实际为 %s", result)
	}
	if !strings.Contains(result, "&") {
		t.Fatalf("期望包含 & 分隔符，实际为 %s", result)
	}
}

func TestLinkMap_ToString_OneSep(t *testing.T) {
	om := NewLinkMap[string, int]()
	om.Set("a", 1)

	result := om.ToString("=")
	if result != "a=1" {
		t.Fatalf("期望 a=1，实际为 %s", result)
	}
	// 多项时 itemSep 应为默认的 ";"，且不在末尾追加
	om.Set("b", 2)
	result = om.ToString("=")
	if result != "a=1;b=2" {
		t.Fatalf("期望 a=1;b=2，实际为 %s", result)
	}
}

func TestLinkMap_IntKey(t *testing.T) {
	om := NewLinkMap[int, string]()
	om.Set(1, "one")
	om.Set(2, "two")
	om.Set(3, "three")

	v, ok := om.Get(2)
	if !ok || v != "two" {
		t.Fatalf("期望 2 -> two，实际为 %s", v)
	}

	om.Delete(1)
	if len(om.keys) != 2 {
		t.Fatalf("期望 keys 长度为 2，实际为 %d", len(om.keys))
	}
}
