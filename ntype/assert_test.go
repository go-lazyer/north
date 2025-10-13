package ntype_test

import (
	"fmt"
	"testing"

	"github.com/go-lazyer/north/ntype"
)

func TestIsNumeric(t *testing.T) {
	fmt.Println(ntype.IsNumeric("1231w.2"))
}
