package util

import (
	"fmt"
	"testing"
	"time"
)

func TestTimestamp(t *testing.T) {
	tt := Timestamp(time.Now())
	bt := Time(tt)
	t2 := Timestamp(bt)
	fmt.Println(1392636938688)
	fmt.Println(tt)
	fmt.Println(t2)
	if tt != t2 {
		t.Error("convert invalid")
		return
	}
	fmt.Println(Now())
	fmt.Println(NowSec())
}

func TestUtil(t *testing.T) {
	fmt.Println(CPU())
}
