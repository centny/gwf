package util

import (
	"fmt"
	"testing"
)

func TestStr2Ints(t *testing.T) {
	is, err := Str2Ints("abc")
	if err == nil {
		t.Error("error")
		return
	}
	is, err = Str2Ints("11")
	if err != nil || is[0] != 11 {
		fmt.Println(err, is)
		t.Error("error")
		return
	}
	is, err = Str2Ints("11,22")
	if err != nil || is[0] != 11 || is[1] != 22 {
		fmt.Println(err, is)
		t.Error("error")
		return
	}
	is, err = Str2Ints("11,")
	if err != nil || is[0] != 11 {
		fmt.Println(err, is)
		t.Error("error")
		return
	}
	is, err = Str2Ints("11,ssd")
	if err == nil {
		t.Error("error")
		return
	}
	vals, err := Str2IntsSeq("0/1/2/3", "/")
	if err != nil {
		t.Error(err.Error())
		return
	}
	for idx, val := range vals {
		if idx != val {
			t.Error("error")
			return
		}
	}
	_, err = Str2IntsSeq("0/1x/2/3", "/")
	if err == nil {
		t.Error("not error")
		return
	}
	//
	if Vals2Str(1, 2, 3, 4) != "1,2,3,4" {
		t.Error("error")
		return
	}
	vals2 := Vals2Strs(0, 1, 2, 3, 4)
	for idx, val := range vals2 {
		if fmt.Sprintf("%v", idx) != val {
			fmt.Println(idx, val)
			t.Error("error")
			return
		}
	}
	fmt.Println("test str2int done...")
}
