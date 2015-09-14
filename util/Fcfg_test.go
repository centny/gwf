package util

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"testing"
	"time"
)

func TestEnvReplace(t *testing.T) {
	f := NewFcfg3()
	f.SetVal("a", "b111111")
	fmt.Println(f.EnvReplace("sss${a} ${abc} ${da} ${HOME} ${}"))
}

func TestInit(t *testing.T) {
	f := NewFcfg3()
	err := f.InitWithFilePath2("not_found.properties", false)
	if err == nil {
		panic("init error")
	}
	err = f.InitWithFilePath("fcfg_data.properties")
	if err != nil {
		t.Error(err.Error())
		return
	}
	for key, val := range f.Map {
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
	f := NewFcfg3()
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
	f := NewFcfg3()
	err := f.InitWithFilePath("fcfg_data.properties?ukk=123")
	if err != nil {
		fmt.Println("error:", err)
	}
	fmt.Println(f.FloatVal("floata"))
	fmt.Println(f.FloatVal("floatb"))
	fmt.Println(f.FloatVal("inta"))
	fmt.Println(f.FloatVal("ukk"))
}

func TestLoad(t *testing.T) {
	os.Remove("ssd.sss")
	cfg, err := NewFcfg2("@l:http://127.0.0.1:65432/fcfg_data.properties")
	if err != nil {
		t.Error(err.Error())
		return
	}
	cfg.Print()
	go func() {
		time.Sleep(time.Second)
		FWrite("ssd.sss", "a=1")
	}()
	cfg, err = NewFcfg2("@l:ssd.sss")
	if err != nil {
		t.Error(err.Error())
		return
	}
	go func() {
		time.Sleep(time.Second)
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("a=1"))
		})
		http.ListenAndServe(":8334", nil)
	}()
	cfg, err = NewFcfg2("@l:http://127.0.0.1:8334")
	if err != nil {
		t.Error(err.Error())
		return
	}
	cfg.Merge(nil)
	cfg = NewFcfg3()
	cfg.InitWithFilePath2("sfdsfsd", false)
	cfg.InitWithURL2("sdfsfs", false)
	NewFcfg2("@l:https://cfg:!DyCfg_321@192.168.1.14/WebDAV/cfg/www.properties")
}

func TestSection(t *testing.T) {
	f := NewFcfg3()
	err := f.InitWithFilePath("fcfg_data.properties?ukk=123")
	if err != nil {
		fmt.Println("error:", err)
	}
	if f.Val("abc/txabc") != "1" {
		t.Error("not right")
		return
	}
	if f.Val("abd/dxabc") != "1" {
		t.Error("not right")
		return
	}
	fmt.Println("%v", f)
	f.Exist("kjuu")
	f.Val("kjuu")
	fmt.Println(f.Seces)
	os.Remove("tt.properties")
	err = f.Store("abc", "tt.properties", "xx")
	if err != nil {
		t.Error(err.Error())
		return
	}
	f.Store("adkkdbc", "tt.properties", "xx")
	f.Store("abc", "/tt.properties", "xx")

	fmt.Println(f.Val("wwwk"))
	fmt.Println(f.Val("wxk"))
}
