package wdoc

import (
	"fmt"
	"github.com/Centny/gwf/util"
	"testing"
)

func TestParser(t *testing.T) {
	pp := NewParser()
	err := pp.Parse("/Users/vty/vgo/src/github.com/Centny/gwf/wdoc/test")
	if err != nil {
		t.Error(err.Error())
		return
	}
	var res = pp.ToM("")
	fmt.Println(util.S2Json(res))
}

func TestReg(t *testing.T) {
	var ta = "xx	R	sss"
	var tb = "x1	O	sss"
	var tc = "x2	optional	sss"
	var td = "x3	required	sss"
	if !ARG_REG.MatchString(ta) {
		t.Error("error->a")
	}
	if !ARG_REG.MatchString(tb) {
		t.Error("error->b")
	}
	if !ARG_REG.MatchString(tc) {
		t.Error("error->c")
	}
	if !ARG_REG.MatchString(td) {
		t.Error("error->d")
	}
	//
	var texts = []string{
		"xx	S	sss",
		"x1	I	sss",
		"x2	F	sss",
		"x3	A	sss",
		"x3	O	sss",
	}
	for _, text := range texts {
		if !RET_REG.MatchString(text) {
			t.Error("error->" + text)
		}
	}
	fmt.Println("done...")
}
