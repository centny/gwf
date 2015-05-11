package util

import (
	"fmt"
	"testing"
	"time"
)

func TestArray(t *testing.T) {
	ary := &Array{}
	for i := 0; i < 10; i++ {
		ary.Add(i)
	}
	for i := 0; i < 10; i++ {
		fmt.Println(i, ":", ary.At(i))
	}
	for i := 0; i < 10; i++ {
		ary.Del(0)
		fmt.Println("len:", ary.Ary())
	}
	fmt.Println("len:", CreateArray(10).Len())
	time.Sleep(2 * time.Second)
	fmt.Println(Err("aaa:%v", "kkk"))
}

func TestParseInt(t *testing.T) {
	val, err := ParseInt("10")
	if err != nil {
		t.Error(err.Error())
		return
	}
	fmt.Println(val)
	fmt.Println(ParseInt("sfs"))
}
