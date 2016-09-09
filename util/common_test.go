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

type sortStruct struct {
	S string
	I int
	F float64
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
	//

	//test int sort
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
	//
	//test struct sort
	var sVals = []sortStruct{
		sortStruct{
			S: "b",
			I: 2,
			F: 2,
		},
		sortStruct{
			S: "c",
			I: 3,
			F: 3,
		},
		sortStruct{
			S: "a",
			I: 1,
			F: 1,
		},
	}
	NewFieldIntSorter("I", sVals).Sort(false)
	if sVals[2].I != 3 {
		t.Error("error")
		fmt.Println(sVals)
		return
	}
	NewFieldIntSorter("I", sVals).Sort(true)
	if sVals[2].I != 1 {
		t.Error("error")
		fmt.Println(sVals)
		return
	}
	NewFieldFloatSorter("F", sVals).Sort(false)
	if sVals[2].F != 3 {
		t.Error("error")
		fmt.Println(sVals)
		return
	}
	NewFieldFloatSorter("F", sVals).Sort(true)
	if sVals[2].F != 1 {
		t.Error("error")
		fmt.Println(sVals)
		return
	}
	NewFieldStringSorter("S", sVals).Sort(false)
	if sVals[2].S != "c" {
		t.Error("error")
		fmt.Println(sVals)
		return
	}
	NewFieldStringSorter("S", sVals).Sort(true)
	if sVals[2].S != "a" {
		t.Error("error")
		fmt.Println(sVals)
		return
	}
	//
	//test *struct sort
	var sVals2 = []*sortStruct{
		&sortStruct{
			S: "b",
			I: 2,
			F: 2,
		},
		&sortStruct{
			S: "c",
			I: 3,
			F: 3,
		},
		&sortStruct{
			S: "a",
			I: 1,
			F: 1,
		},
	}
	fmt.Println(sVals2)
	NewFieldIntSorter("I", sVals2).Sort(false)
	fmt.Println(sVals2)
	if sVals2[2].I != 3 {
		t.Error("error")
		fmt.Println(sVals2)
		return
	}
	NewFieldIntSorter("I", sVals2).Sort(true)
	if sVals2[2].I != 1 {
		t.Error("error")
		fmt.Println(sVals2)
		return
	}
	NewFieldFloatSorter("F", sVals2).Sort(false)
	if sVals2[2].F != 3 {
		t.Error("error")
		fmt.Println(sVals2)
		return
	}
	NewFieldFloatSorter("F", sVals2).Sort(true)
	if sVals2[2].F != 1 {
		t.Error("error")
		fmt.Println(sVals2)
		return
	}
	NewFieldStringSorter("S", sVals2).Sort(false)
	if sVals2[2].S != "c" {
		t.Error("error")
		fmt.Println(sVals)
		return
	}
	NewFieldStringSorter("S", sVals2).Sort(true)
	if sVals2[2].S != "a" {
		t.Error("error")
		fmt.Println(sVals2)
		return
	}
}
