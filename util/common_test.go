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

func TestSorter(t *testing.T) {
	// A := []int{1, 2}
	// V := reflect.ValueOf(A)
	// x, y := V.Index(0).Interface(), V.Index(1).Interface()
	// V.Index(0).Set(reflect.ValueOf(y))
	// V.Index(1).Set(reflect.ValueOf(x))
	// fmt.Println(A)
	//
	// A := []string{"world", "hello"}
	// V := reflect.ValueOf(A)
	// x, y := V.Index(0).Interface(), V.Index(1).Interface()
	// V.Index(0).Set(reflect.ValueOf(y))
	// V.Index(1).Set(reflect.ValueOf(x))
	// fmt.Println(A)
	//
	var intVals = []int{1, 5, 0, 2, 4, 3}
	NewIntSorter(intVals).Sort(false)
	for idx, val := range intVals {
		if idx != val {
			fmt.Println(idx, val, intVals)
			t.Error("error")
			return
		}
	}
	NewIntSorter(intVals).Sort(true)
	for idx, val := range intVals {
		if val != 5-idx {
			fmt.Println(5-idx, val, intVals)
			t.Error("error")
			return
		}
	}
}
