package util

import (
	"fmt"
	"testing"
)

func TestErr(t *testing.T) {
	fmt.Println(NewErr("0", nil))
	fmt.Println(NewErr("0", Err("sdfsd:%v", "dfs")))
	fmt.Println(NewErr2("0", "sdfsd:%v", "dfs"))
	fmt.Println(NewErr2("0", "sdfsd:%v", "dfs").String())
}
