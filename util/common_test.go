package util

import (
	"bufio"
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

type WW struct {
}

func (w *WW) Write(p []byte) (int, error) {
	return 0, Err("dsfdf")
}
func TestWriter(t *testing.T) {
	w := bufio.NewWriter(&WW{})
	for i := 0; i < 1000; i++ {
		w.Write([]byte("sfsdfsddfsfsfssfs"))
	}
	fmt.Println(w.Write([]byte("sfs")))
	fmt.Println(w.Flush())
}
