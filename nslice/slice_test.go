package nslice

import (
	"fmt"
	"testing"
)

func TestCut(t *testing.T) {
	num := make([]int, 0)
	for i := 0; i < 101; i++ {
		num = append(num, i)
	}

	res := Cut(num, 10)
	for _, re := range res {
		fmt.Println(len(re))
	}
}

func TestRemove(t *testing.T) {
	data := []string{"apple", "banana", "cherry", "date", "elderberry"}

	// 删除索引 1 和 3（"banana" 和 "date"）
	result := Remove(data, 1)

	fmt.Println(result) // 输出: [apple cherry elderberry]
}

func TestRemoves(t *testing.T) {
	data := []string{"apple", "banana", "cherry", "date", "elderberry"}

	// 删除索引 1 和 3（"banana" 和 "date"）
	result := Removes(data, []int{1, 3})

	fmt.Println(result) // 输出: [apple cherry elderberry]
}
