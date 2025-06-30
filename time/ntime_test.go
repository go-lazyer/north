package ntime

import (
	"fmt"
	"testing"
)

func TestParse(t *testing.T) {
	fmt.Println(Parse("12:00:00").Format("15:04:05"))
}
