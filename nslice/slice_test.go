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
