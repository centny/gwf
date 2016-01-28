package wdoc

import (
	"fmt"
	"github.com/Centny/gwf/routing/httptest"
	"github.com/Centny/gwf/util"
	"os"
	"testing"
	"time"
)

func TestParser(t *testing.T) {
	pp := NewParser()
	var wait = make(chan int)
	go func() {
		pp.LoopParse(os.Getenv("GOPATH")+"/src/github.com/Centny/gwf/wdoc", nil, nil, 1000)
		wait <- 0
	}()
	time.Sleep(2 * time.Second)
	pp.Running = false
	<-wait
	var res = pp.ToMv("", "x1", "test,x")
	if len(res.Pkgs) < 1 {
		t.Error("error")
		return
	}
	res = pp.ToMv("", "xxxx1", "test,x")
	if len(res.Pkgs) > 0 {
		t.Error("error")
		return
	}
	res = pp.ToM("")
	fmt.Println(util.S2Json(res))
	//
	ts := httptest.NewServer2(pp)
	ts.G("")
	//
	//
	//test error
	NewParser().Parse("/sdfk/sds")
	NewParser().ParseDir("/sdfk/sds", nil, nil)
	go pp.LoopParse("/dsfsfd", nil, nil, 1000)
	time.Sleep(2 * time.Second)
	pp.Running = false
	pkgs_l([]Pkg{Pkg{}, Pkg{}}).Swap(0, 1)
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
