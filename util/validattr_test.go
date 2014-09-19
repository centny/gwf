package util

import (
	"fmt"
	"testing"
)

func TestValidAttr(t *testing.T) {
	v, err := ValidAttrT("测试", "R|S", "L:~10", true)
	if err != nil {
		t.Error(err.Error())
		return
	}
	fmt.Println(v)
	v, err = ValidAttrT("测试测试测试测试", "R|S", "L:~10", true)
	if err == nil {
		t.Error("not error")
		return
	}
	//
	v, err = ValidAttrT("男", "R|S", "O:男~女", true)
	if err != nil {
		t.Error(err.Error())
		return
	}
	fmt.Println(v)
	v, err = ValidAttrT("男ks", "R|S", "O:男~女", true)
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
	v, err = ValidAttrT("8", "O|I", "R:5~10", true)
	if err != nil {
		t.Error(err.Error())
		return
	}
	fmt.Println(v)
	v, err = ValidAttrT("8", "O|I", "R:5~", true)
	if err != nil {
		t.Error(err.Error())
		return
	}
	fmt.Println(v)
	v, err = ValidAttrT("12", "O|I", "R:5~10", true)
	if err == nil {
		t.Error("not error")
		return
	}
	//
	v, err = ValidAttrT("8", "O|F", "R:5~10", true)
	if err != nil {
		t.Error(err.Error())
		return
	}
	fmt.Println(v)
	v, err = ValidAttrT("8", "O|F", "R:5~", true)
	if err != nil {
		t.Error(err.Error())
		return
	}
	fmt.Println(v)
	v, err = ValidAttrT("12", "O|F", "R:5~10", true)
	if err == nil {
		t.Error("not error")
		return
	}
	//
	v, err = ValidAttrT("测", "O|S", "L:~8", true)
	if err != nil {
		t.Error(err.Error())
		return
	}
	fmt.Println(v)
	v, err = ValidAttrT("测", "O|S", "L:2~", true)
	if err != nil {
		t.Error(err.Error())
		return
	}
	fmt.Println(v)
	v, err = ValidAttrT("测度测度测度测度测度", "O|S", "L:~8", true)
	if err == nil {
		t.Error("not error")
		return
	}
	//
	v, err = ValidAttrT("测", "O|S", "L:2~8", true)
	if err != nil {
		t.Error(err.Error())
		return
	}
	fmt.Println(v)
	v, err = ValidAttrT("a", "O|S", "L:2~8", true)
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
	v, err = ValidAttrT("1", "O|I", "O:1~2~3~4~5", true)
	if err != nil {
		t.Error(err.Error())
		return
	}
	fmt.Println(v)
	v, err = ValidAttrT("11", "O|I", "O:1~2~3~4~5", true)
	if err == nil {
		t.Error("not error")
		return
	}
	//
	v, err = ValidAttrT("1.1", "O|F", "O:1.1~2.2~3.3~4~5", true)
	if err != nil {
		t.Error(err.Error())
		return
	}
	fmt.Println(v)
	v, err = ValidAttrT("11", "O|F", "O:1~2~3~4~5", true)
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
	v, err = ValidAttrT("测", "O|I", "R:8~9", true)
	if err == nil {
		t.Error("not error")
		return
	}
	v, err = ValidAttrT("测", "O|F", "R:8~9", true)
	if err == nil {
		t.Error("not error")
		return
	}
	v, err = ValidAttrT("测", "O|N", "R:8~9", true)
	if err == nil {
		t.Error("not error")
		return
	}
	v, err = ValidAttrT("5", "R|I", "R:~1", true)
	if err == nil {
		t.Error("not error")
		return
	}
	v, err = ValidAttrT("5", "R|I", "R:a~10", true)
	if err == nil {
		t.Error("not error")
		return
	}
	v, err = ValidAttrT("5", "R|I", "R:1~a", true)
	if err == nil {
		t.Error("not error")
		return
	}
	v, err = ValidAttrT("5", "R|F", "R:~1", true)
	if err == nil {
		t.Error("not error")
		return
	}
	v, err = ValidAttrT("5", "R|F", "R:a~10", true)
	if err == nil {
		t.Error("not error")
		return
	}
	v, err = ValidAttrT("5", "R|F", "R:1~a", true)
	if err == nil {
		t.Error("not error")
		return
	}
	v, err = ValidAttrT("5", "R|I", "M:1~a", true)
	if err == nil {
		t.Error("not error")
		return
	}
	v, err = ValidAttrT("5", "R|I", "O:1~a", true)
	if err == nil {
		t.Error("not error")
		return
	}
	v, err = ValidAttrT("5", "R|F", "O:1~a", true)
	if err == nil {
		t.Error("not error")
		return
	}
	v, err = ValidAttrT("5", "R|F", "M:1~k", true)
	if err == nil {
		t.Error("not error")
		return
	}
	v, err = ValidAttrT("5", "R|I", "O", true)
	if err == nil {
		t.Error("not error")
		return
	}
	v, err = ValidAttrT("5", "R", "O:1~10", true)
	if err == nil {
		t.Error("not error")
		return
	}
	v, err = ValidAttrT("", "R|I", "O:1~10", true)
	if err == nil {
		t.Error("not error")
		return
	}
	v, err = ValidAttrT("", "O|I", "O:1~10", true)
	if err != nil {
		t.Error(err.Error())
		return
	}
	v, err = ValidAttrT("a", "O|S", "L:a~8", true)
	if err == nil {
		t.Error("not error")
		return
	}
	v, err = ValidAttrT("a", "O|S", "L:2~a", true)
	if err == nil {
		t.Error("not error")
		return
	}
}

func TestValidAttrF(t *testing.T) {
	mv := map[string]string{}
	mv["a"] = "abc"
	mv["i"] = "10"
	mv["f"] = "10.3"
	mv["ef"] = "20.3"
	mv["len"] = "11111111"
	var a string
	var i int64
	var k string
	var f float64
	err := ValidAttrF(`//abc
		a,R|S,L:~5;//abc
		i,R|I,R:1~20;
		k,O|I,R:1~20;//sfdsj
		f,R|F,R:1.5~20;
		`, func(key string) string {
		return mv[key]
	}, true, &a, &i, &k, &f)
	if err != nil {
		t.Error(err.Error())
		return
	}
	fmt.Println(a, i, k, f)
	//
	err = ValidAttrF(`
		a,R|S L:~5;
		`, func(key string) string {
		return mv[key]
	}, true, &a)
	if err == nil {
		t.Error("not error")
		return
	}
	//
	err = ValidAttrF(`
		len,R|S,L:~5;
		`, func(key string) string {
		return mv[key]
	}, true, &a)
	if err == nil {
		t.Error("not error")
		return
	}
	//
	var ea float32
	err = ValidAttrF(`
		a,R|S,L:~5;
		`, func(key string) string {
		return mv[key]
	}, true, &ea)
	if err == nil {
		t.Error("not error")
		return
	}
	fmt.Println(err.Error())
	//
	err = ValidAttrF(``, func(key string) string {
		return mv[key]
	}, true, &a)
	if err == nil {
		t.Error("not error")
		return
	}
	fmt.Println(err.Error())
	//
	err = ValidAttrF(`
		len,R|S,L:~5;
		len,R|S,L:~5;
		`, func(key string) string {
		return mv[key]
	}, true, &a)
	if err == nil {
		t.Error("not error")
		return
	}
	fmt.Println(err.Error())
	err = ValidAttrF(`
		len,R|S,L:~5,this is error message;
		`, func(key string) string {
		return mv[key]
	}, true, &a)
	if err == nil {
		t.Error("not error")
		return
	}
	fmt.Println(err.Error())
}
func TestEscape(t *testing.T) {
	//
	var a string
	err := ValidAttrF(`
		len,R|S,P:[^%N]*%N.*$;
		`, func(key string) string {
		return "abc,ddf"
	}, true, &a)
	if err != nil {
		t.Error(err.Error())
		return
	}
}
