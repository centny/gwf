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
	mv["ary"] = "1,2,3,4,5"
	mv["ary2"] = "1,2,3,,4,5"
	var a string
	var i int64
	var k string
	var ks []string
	var f float64
	var iv1 int
	var iv1_a []int
	var iv2 int16
	var iv3 int32
	var iv4 int64
	var iv5 uint
	var iv6 uint16
	var iv7 uint32
	var iv8 uint64
	var iv9 float32
	var iv10 float64
	var iv10_a []float64
	var iv11 string
	var iv12 int64
	var a_i []int
	var a_s []string
	var a_f []float64
	err := ValidAttrF(`//abc
		a,R|S,L:~5;//abc
		i,R|I,R:1~20;
		i,O|I,R:1~20;//sfdsj
		i,O|I,R:1~20;//sfdsj
		f,R|F,R:1.5~20;
		i,R|I,R:0;
		i,R|I,R:0;
		i,R|I,R:0;
		i,R|I,R:0;
		i,R|I,R:0;
		i,R|I,R:0;
		i,R|I,R:0;
		i,R|I,R:0;
		i,R|I,R:0;
		i,R|I,R:0;
		i,R|I,R:0;
		i,R|I,R:0;
		i,R|I,R:0;
		i,R|S,L:0;
		ary,R|S,L:0;
		ary,R|I,R:0;
		ary,R|F,R:0;
		`, func(key string) string {
		return mv[key]
	}, true, &a, &i, &k, &ks, &f,
		&iv1, &iv1_a, &iv2, &iv3, &iv4, &iv5,
		&iv6, &iv7, &iv8, &iv9, &iv10, &iv10_a,
		&iv11, &iv12, &a_s, &a_i, &a_f)
	if err != nil {
		t.Error(err.Error())
		return
	}
	fmt.Println(k, ks, len(iv1_a), len(iv10_a))
	if k != "10" || ks[0] != "10" || iv1 != 10 || iv1_a[0] != 10 || iv10 != 10 || iv10_a[0] != 10 {
		t.Error("error")
		return
	}
	fmt.Println(len(a_s), len(a_i), len(a_f))
	if len(a_s) != 5 || len(a_i) != 5 || len(a_f) != 5 {
		t.Error("error")
		return
	}
	fmt.Println(a, i, k, f)
	//
	//test array
	a_s, a_i, a_f = nil, nil, nil
	err = ValidAttrF(`
		ary2,R|S,L:0;
		ary2,R|I,R:0;
		ary2,R|F,R:0;
		`, func(key string) string {
		return mv[key]
	}, true, &a_s, &a_i, &a_f)
	if err == nil {
		t.Error("error")
		return
	}
	a_s, a_i, a_f = nil, nil, nil
	err = ValidAttrF(`
		ary2,O|S,L:0;
		ary2,O|I,R:0;
		ary2,O|F,R:0;
		`, func(key string) string {
		return mv[key]
	}, true, &a_s, &a_i, &a_f)
	if err != nil {
		t.Error("error")
		return
	}
	if len(a_s) != 5 || len(a_i) != 5 || len(a_f) != 5 {
		t.Error("error")
		return
	}
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

func TestValidAttrFAddr(t *testing.T) {
	mv := map[string]string{}
	mv["a"] = "abc"
	mv["i"] = "10"
	mv["f"] = "10.3"
	mv["ef"] = "20.3"
	mv["len"] = "11111111"
	mv["ary"] = "1,2,3,4,5"
	mv["ary2"] = "1,2,3,,4,5"
	var a *string
	var i *int64
	var k *string
	var ks []*string
	var f *float64
	var iv1 *int
	var iv1_a []*int
	var iv2 *int16
	var iv3 *int32
	var iv4 *int64
	var iv5 *uint
	var iv6 *uint16
	var iv7 *uint32
	var iv8 *uint64
	var iv9 *float32
	var iv10 *float64
	var iv10_a []*float64
	var iv11 *string
	var iv12 *int64
	var a_i []*int
	var a_s []*string
	var a_f []*float64
	err := ValidAttrF(`//abc
		a,R|S,L:~5;//abc
		i,R|I,R:1~20;
		i,O|I,R:1~20;//sfdsj
		i,O|I,R:1~20;//sfdsj
		f,R|F,R:1.5~20;
		i,R|I,R:0;
		i,R|I,R:0;
		i,R|I,R:0;
		i,R|I,R:0;
		i,R|I,R:0;
		i,R|I,R:0;
		i,R|I,R:0;
		i,R|I,R:0;
		i,R|I,R:0;
		i,R|I,R:0;
		i,R|I,R:0;
		i,R|I,R:0;
		i,R|I,R:0;
		i,R|S,L:0;
		ary,R|S,L:0;
		ary,R|I,R:0;
		ary,R|F,R:0;
		`, func(key string) string {
		return mv[key]
	}, true, &a, &i, &k, &ks, &f,
		&iv1, &iv1_a, &iv2, &iv3, &iv4, &iv5,
		&iv6, &iv7, &iv8, &iv9, &iv10, &iv10_a,
		&iv11, &iv12, &a_s, &a_i, &a_f)
	if err != nil {
		t.Error(err.Error())
		return
	}
	fmt.Println(k, ks, len(iv1_a), len(iv10_a))
	if *k != "10" || *ks[0] != "10" || *iv1 != 10 || *iv1_a[0] != 10 || *iv10 != 10 || *iv10_a[0] != 10 {
		t.Error("error")
		return
	}
	fmt.Println(len(a_s), len(a_i), len(a_f))
	if len(a_s) != 5 || len(a_i) != 5 || len(a_f) != 5 {
		t.Error("error")
		return
	}
	fmt.Println(a, i, k, f)
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
