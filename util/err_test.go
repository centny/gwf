package util

import (
	"fmt"
	"strconv"
	"testing"
)

func TestErr(t *testing.T) {
	fmt.Println(NewErr("0", nil))
	fmt.Println(NewErr("0", Err("sdfsd:%v", "dfs")))
	fmt.Println(NewErr2("0", "sdfsd:%v", "dfs"))
	fmt.Println(NewErr2("0", "sdfsd:%v", "dfs").String())
}

func TestUu(t *testing.T) {
	fmt.Println(0755)
	fmt.Println(strconv.ParseUint("0755", 8, 32))
}
