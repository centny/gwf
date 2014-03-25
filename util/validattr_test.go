package util

import (
	"fmt"
	"testing"
)

func TestValidAttr(t *testing.T) {
	v, err := ValidAttrT("测试", "R|S", "L:10", true)
	if err != nil {
		t.Error(err.Error())
		return
	}
	fmt.Println(v)
	v, err = ValidAttrT("测试测试测试测试", "R|S", "L:10", true)
	if err == nil {
		t.Error("not error")
		return
	}
	//
	v, err = ValidAttrT("男", "R|S", "O:男-女", true)
	if err != nil {
		t.Error(err.Error())
		return
	}
	fmt.Println(v)
	v, err = ValidAttrT("男ks", "R|S", "O:男-女", true)
	if err == nil {
		t.Error("not error")
		return
	}
	//
	v, err = ValidAttrT("centny@gmail.com", "R|S", "P:^.*\\@.*$", true)
	if err != nil {
		t.Error(err.Error())
		return
	}
	fmt.Println(v)
	v, err = ValidAttrT("ks", "R|S", "P:^.*\\@.*$", true)
	if err == nil {
		t.Error("not error")
		return
	}
	//
	v, err = ValidAttrT("8", "O|I", "R:5-10", true)
	if err != nil {
		t.Error(err.Error())
		return
	}
	fmt.Println(v)
	v, err = ValidAttrT("12", "O|I", "R:5-10", true)
	if err == nil {
		t.Error("not error")
		return
	}
	//
	v, err = ValidAttrT("8", "O|F", "R:5-10", true)
	if err != nil {
		t.Error(err.Error())
		return
	}
	fmt.Println(v)
	v, err = ValidAttrT("12", "O|F", "R:5-10", true)
	if err == nil {
		t.Error("not error")
		return
	}
	//
	v, err = ValidAttrT("测", "O|S", "L:8", true)
	if err != nil {
		t.Error(err.Error())
		return
	}
	fmt.Println(v)
	v, err = ValidAttrT("测度测度测度测度测度", "O|S", "L:8", true)
	if err == nil {
		t.Error("not error")
		return
	}
	//
	v, err = ValidAttrT("centny@gmail.com", "O|S", "P:^.*\\@.*$", true)
	if err != nil {
		t.Error(err.Error())
		return
	}
	fmt.Println(v)
	v, err = ValidAttrT("ks", "O|S", "P:^.*\\@.*$", true)
	if err == nil {
		t.Error("not error")
		return
	}
	//
	v, err = ValidAttrT("1", "O|I", "O:1-2-3-4-5", true)
	if err != nil {
		t.Error(err.Error())
		return
	}
	fmt.Println(v)
	v, err = ValidAttrT("11", "O|I", "O:1-2-3-4-5", true)
	if err == nil {
		t.Error("not error")
		return
	}
	//
	v, err = ValidAttrT("测", "O|S", "L:a", true)
	if err == nil {
		t.Error("not error")
		return
	}
	v, err = ValidAttrT("测", "O|S", "KK:a", true)
	if err == nil {
		t.Error("not error")
		return
	}
	v, err = ValidAttrT("centny@gmail.com", "O|S", "P:*,..", true)
	if err == nil {
		t.Error("not error")
		return
	}
	v, err = ValidAttrT("测", "O|I", "R:8-9", true)
	if err == nil {
		t.Error("not error")
		return
	}
	v, err = ValidAttrT("测", "O|F", "R:8-9", true)
	if err == nil {
		t.Error("not error")
		return
	}
	v, err = ValidAttrT("测", "O|N", "R:8-9", true)
	if err == nil {
		t.Error("not error")
		return
	}
	v, err = ValidAttrT("5", "R|I", "R:1", true)
	if err == nil {
		t.Error("not error")
		return
	}
	v, err = ValidAttrT("5", "R|I", "R:a-10", true)
	if err == nil {
		t.Error("not error")
		return
	}
	v, err = ValidAttrT("5", "R|I", "R:1-a", true)
	if err == nil {
		t.Error("not error")
		return
	}
	v, err = ValidAttrT("5", "R|I", "M:1-a", true)
	if err == nil {
		t.Error("not error")
		return
	}
	v, err = ValidAttrT("5", "R|I", "O:1-a", true)
	if err == nil {
		t.Error("not error")
		return
	}
	v, err = ValidAttrT("5", "R|I", "O", true)
	if err == nil {
		t.Error("not error")
		return
	}
	v, err = ValidAttrT("5", "R", "O:1-10", true)
	if err == nil {
		t.Error("not error")
		return
	}
	v, err = ValidAttrT("", "R|I", "O:1-10", true)
	if err == nil {
		t.Error("not error")
		return
	}
	v, err = ValidAttrT("", "O|I", "O:1-10", true)
	if err != nil {
		t.Error(err.Error())
		return
	}
}
