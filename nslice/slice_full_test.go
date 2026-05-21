package nslice

import (
	"reflect"
	"testing"
)

func TestCut_Normal(t *testing.T) {
	nums := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	result := Cut(nums, 3)

	expected := [][]int{
		{1, 2, 3},
		{4, 5, 6},
		{7, 8, 9},
		{10},
	}
	if len(result) != len(expected) {
		t.Fatalf("期望 %d 组，实际 %d 组", len(expected), len(result))
	}
	for i, group := range result {
		if !reflect.DeepEqual(group, expected[i]) {
			t.Errorf("第 %d 组期望 %v，实际 %v", i, expected[i], group)
		}
	}
}

func TestCut_ExactDivision(t *testing.T) {
	nums := []int{1, 2, 3, 4, 5, 6}
	result := Cut(nums, 3)

	expected := [][]int{
		{1, 2, 3},
		{4, 5, 6},
	}
	if len(result) != len(expected) {
		t.Fatalf("期望 %d 组，实际 %d 组", len(expected), len(result))
	}
	for i, group := range result {
		if !reflect.DeepEqual(group, expected[i]) {
			t.Errorf("第 %d 组期望 %v，实际 %v", i, expected[i], group)
		}
	}
}

func TestCut_EmptySlice(t *testing.T) {
	result := Cut([]int{}, 3)
	if len(result) != 0 {
		t.Errorf("空切片期望 0 组，实际 %d 组", len(result))
	}
}

func TestCut_InvalidShareNums(t *testing.T) {
	result := Cut([]int{1, 2, 3}, 0)
	if len(result) != 0 {
		t.Errorf("shareNums=0 期望 0 组，实际 %d 组", len(result))
	}

	result = Cut([]int{1, 2, 3}, -1)
	if len(result) != 0 {
		t.Errorf("shareNums=-1 期望 0 组，实际 %d 组", len(result))
	}
}

func TestCut_SingleElement(t *testing.T) {
	result := Cut([]int{1}, 5)
	if len(result) != 1 {
		t.Fatalf("期望 1 组，实际 %d 组", len(result))
	}
	if !reflect.DeepEqual(result[0], []int{1}) {
		t.Errorf("期望 [1]，实际 %v", result[0])
	}
}

func TestRemove_Indices(t *testing.T) {
	// Remove 是原地修改切片，所以每次测试需要用新切片
	data1 := []string{"apple", "banana", "cherry", "date", "elderberry"}
	result := Remove(data1, 1, 3)
	expected := []string{"apple", "cherry", "elderberry"}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Remove(1,3) 期望 %v，实际 %v", expected, result)
	}

	// 删除索引 0
	data2 := []string{"apple", "banana", "cherry", "date", "elderberry"}
	result = Remove(data2, 0)
	expected = []string{"banana", "cherry", "date", "elderberry"}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Remove(0) 期望 %v，实际 %v", expected, result)
	}

	// 无索引删除
	data3 := []string{"apple", "banana", "cherry", "date", "elderberry"}
	result = Remove(data3)
	expected = []string{"apple", "banana", "cherry", "date", "elderberry"}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Remove() 无索引期望原切片不变，实际 %v", result)
	}
}

func TestRemoveByValue_Values(t *testing.T) {
	// RemoveByValue 是原地修改切片，所以每次测试需要用新切片
	data1 := []string{"apple", "banana", "cherry", "date", "elderberry"}
	result := RemoveByValue(data1, "banana", "elderberry")
	expected := []string{"apple", "cherry", "date"}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("RemoveByValue 期望 %v，实际 %v", expected, result)
	}

	// 删除不存在的值
	data2 := []string{"apple", "banana", "cherry", "date", "elderberry"}
	result = RemoveByValue(data2, "grape")
	expected = []string{"apple", "banana", "cherry", "date", "elderberry"}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("RemoveByValue 不存在的值期望原切片不变，实际 %v", result)
	}

	// 无值删除
	data3 := []string{"apple", "banana", "cherry", "date", "elderberry"}
	result = RemoveByValue(data3)
	expected = []string{"apple", "banana", "cherry", "date", "elderberry"}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("RemoveByValue() 无值期望原切片不变，实际 %v", result)
	}

	// int 类型
	intData := []int{1, 2, 3, 2, 4}
	resultInt := RemoveByValue(intData, 2)
	expectedInt := []int{1, 3, 4}
	if !reflect.DeepEqual(resultInt, expectedInt) {
		t.Errorf("RemoveByValue[int] 期望 %v，实际 %v", expectedInt, resultInt)
	}
}

func TestReplace(t *testing.T) {
	// 替换字符串
	s := []string{"a", "b", "c", "b", "d"}
	result := Replace(s, "b", "x")
	expected := []string{"a", "x", "c", "x", "d"}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Replace 期望 %v，实际 %v", expected, result)
	}

	// 替换不存在的值
	result = Replace(s, "z", "y")
	if !reflect.DeepEqual(result, s) {
		t.Errorf("Replace 不存在的值期望原切片不变")
	}

	// int 类型
	intSlice := []int{1, 2, 3, 2, 4}
	resultInt := Replace(intSlice, 2, 5)
	expectedInt := []int{1, 5, 3, 5, 4}
	if !reflect.DeepEqual(resultInt, expectedInt) {
		t.Errorf("Replace[int] 期望 %v，实际 %v", expectedInt, resultInt)
	}
}

func TestToCsv(t *testing.T) {
	// 正常数据
	data := []map[string]any{
		{"name": "张三", "age": 25, "city": "北京"},
		{"name": "李四", "age": 30, "city": "上海"},
	}

	result := ToCsv(data)
	if len(result) != 3 { // 1 header + 2 data rows
		t.Fatalf("期望 3 行，实际 %d 行", len(result))
	}

	// 第一行是标题行
	header := result[0]
	if len(header) != 3 {
		t.Fatalf("期望标题行 3 列，实际 %d 列", len(header))
	}
	// 检查标题包含所有 key
	for _, key := range []string{"name", "age", "city"} {
		found := false
		for _, h := range header {
			if h == key {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("标题行应包含 %q", key)
		}
	}

	// 空数据
	emptyResult := ToCsv([]map[string]any{})
	if len(emptyResult) != 0 {
		t.Errorf("空数据期望 0 行，实际 %d 行", len(emptyResult))
	}

	// nil 值处理
	dataWithNil := []map[string]any{
		{"name": "王五", "age": nil},
	}
	resultNil := ToCsv(dataWithNil)
	if len(resultNil) != 2 {
		t.Fatalf("期望 2 行，实际 %d 行", len(resultNil))
	}
	// 检查 nil 值转为空字符串
	for i, cell := range resultNil[1] {
		if resultNil[0][i] == "age" && cell != "" {
			t.Errorf("nil 值期望空字符串，实际 %q", cell)
		}
	}
}
