package util

import (
	"fmt"
	"os"
	"os/exec"
	"testing"
)

func TestEnvReplace(t *testing.T) {
	f := &Fcfg{}
	f.SetVal("a", "b111111")
	fmt.Println(f.EnvReplace("sss${a} ${abc} ${da} ${HOME} ${}"))
}

func TestInit(t *testing.T) {
	f := &Fcfg{}
	err := f.InitWithFilePath("not_found.properties")
	if err == nil {
		panic("init error")
	}
	err = f.InitWithFilePath("fcfg_data.properties")
	if err != nil {
		t.Error(err.Error())
		return
	}
	for key, val := range *f {
		fmt.Println(key, ":", val)
	}
	fmt.Println(f.Val("inta"))
	fmt.Println(f.Val("nfound"))
	fmt.Println(f.IntVal("inta"))
	fmt.Println(f.IntVal("nfound"))
	fmt.Println(f.IntVal("a"))
	fmt.Println(f.FloatVal("floata"))
	fmt.Println(f.FloatVal("nfound"))
	fmt.Println(f.FloatVal("a"))
	f.Del("nfound")
	f.Del("a")
	fmt.Println(f.Show())
}
func TestOpenError(t *testing.T) {
	f := &Fcfg{}
	fmt.Println(exec.Command("touch", "/tmp/fcg").Run())
	fmt.Println(exec.Command("chmod", "000", "/tmp/fcg").Run())
	fi, e := os.Open("/tmp/fcg")
	fmt.Println(fi, e)
	err := f.InitWithFilePath("/tmp/fcg")
	if err == nil {
		panic("init error")
	}
	fmt.Println(exec.Command("rm", "-f", "/tmp/fcg").Run())
}
func TestValType(t *testing.T) {
	f := &Fcfg{}
	err := f.InitWithFilePath("fcfg_data.properties")
	if err != nil {
		fmt.Println("error:", err)
	}
	fmt.Println(f.FloatVal("floata"))
	fmt.Println(f.FloatVal("floatb"))
	fmt.Println(f.FloatVal("inta"))
}

func TestLoad(t *testing.T) {
	cfg, err := NewFcfg2("@l:http://127.0.0.1:65432/fcfg_data.properties")
	if err != nil {
		t.Error(err.Error())
		return
	}
	cfg.Print()
	cfg, err = NewFcfg2("@l:ssd.sss")
	if err == nil {
		t.Error("not error")
		return
	}
	NewFcfg2("@l:http://127.0.0.1:6x")
	cfg.Merge(nil)
}
