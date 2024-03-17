package tests

import (
	"bool_interp/bxp"
	"fmt"
	"testing"
)

func TestExpression(t *testing.T) {
	str := "((a|b)&(a|d))&b"
	tt := bxp.NewTruthTable(str)
	tt.Compute()
	fmt.Println(tt.ToString())
}
