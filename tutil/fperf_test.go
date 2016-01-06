package tutil

import (
	"fmt"
	"os"
	"testing"
)

func TestFPerf(t *testing.T) {
	os.RemoveAll("abc")
	os.Mkdir("abc", os.ModePerm)
	fp := NewFPerf("./abc")
	used, err := fp.Perf4MultiRw("a", "", 100, 5, 10240, 1)
	if err != nil {
		t.Error(err.Error())
		return
	}
	fmt.Println(used)
	used, _, err = fp.Perf4MultiR("a", "", 0, 100, 5)
	if err != nil {
		t.Error(err.Error())
		return
	}
	fmt.Println(used)
	used, err = fp.Perf4MultiW("ab", "", 100, 5, 10240, 1)
	if err != nil {
		t.Error(err.Error())
		return
	}
	fmt.Println(used)
	//
	fp = NewFPerf("/fsfsdfs")
	_, err = fp.Read("ssdd")
	if err == nil {
		t.Error("error")
		return
	}
	err = fp.Write("sss", 10240, 1)
	if err == nil {
		t.Error("error")
		return
	}
	err = fp.Rw("sff", 10240, 1)
	if err == nil {
		t.Error("error")
		return
	}
	_, err = fp.Perf4MultiW("sfsf", "", 10, 10, 10, 1)
	if err == nil {
		t.Error("error")
		return
	}
	_, err = fp.Perf4MultiRw("sfsf", "", 10, 10, 10, 1)
	if err == nil {
		t.Error("error")
		return
	}
	_, _, err = fp.Perf4MultiR("xdsd", "", 0, 100, 5)
	if err == nil {
		t.Error("error")
		return
	}
}
