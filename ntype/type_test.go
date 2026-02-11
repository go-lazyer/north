package ntype_test

import (
	"fmt"
	"testing"

	"github.com/go-lazyer/north/ntype"
)

func TestJoinAny(t *testing.T) {

	fmt.Println(ntype.JoinAny([]any{"1", "2", "3"}, ","))
	fmt.Println(ntype.JoinAny([]any{1, 2, "3"}, ","))
}
