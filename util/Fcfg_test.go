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
	if f.Val("abc_conf") != `1
2
3` {
		fmt.Println("->\n", f.Val("abc_conf"))
		t.Error("error")
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
	fmt.Println(f.Val("abc_conf"))
	fmt.Println("\n\n\n\n")
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
	// NewFcfg2("@l:https://cfg:!DyCfg_321@192.168.1.14/WebDAV/cfg/www.properties")
}

func TestAccessLoad(t *testing.T) {
	go func() {
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			username, password, exist := r.BasicAuth()
			if exist && username == "test" && password == "123" {
				w.Write([]byte("a=1"))
			} else {
				w.WriteHeader(500)
			}
		})
		http.ListenAndServe(":8335", nil)
	}()
	cfg := NewFcfg3()
	err := cfg.InitWithURL2("http://txest:123@127.0.0.1:8335", false)
	if err == nil {
		t.Error("not error")
		return
	}
	err = cfg.InitWithURL2("http://test:123@127.0.0.1:8335", false)
	if err != nil {
		t.Error(err.Error())
		return
	}
	if cfg.IntValV("a", 0) != 1 {
		t.Error("fail")
		return
	}
}

func TestAccessLoad2(t *testing.T) {
	cfg := NewFcfg3()
	err := cfg.InitWithURL2("http://txest:123@loc.m/test.properties", false)
	if err == nil {
		t.Error("not error")
		return
	}
	err = cfg.InitWithURL2("http://test:123@127.0.0.1:8335", false)
	if err != nil {
		t.Error(err.Error())
		return
	}
	if cfg.IntValV("a", 0) != 1 {
		t.Error("fail")
		return
	}
}

func TestSection(t *testing.T) {
	f := NewFcfg3()
	err := f.InitWithFilePath("fcfg_data.properties?ukk=123")
	if err != nil {
		fmt.Println("error:", err)
	}
	if f.Val("ukk") != "123" {
		fmt.Println(f.Val("ukk"))
		t.Error("not right")
		return
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

func TestCfgEscape(t *testing.T) {
	var cfg, _ = NewFcfg2(`
		kk=${ss}
		${xx}=1
		`)
	cfg.Print()
	fmt.Println(cfg.Val("kk"))
}

func TestSectionMerge(t *testing.T) {
	var cfga, cfgb = NewFcfg3(), NewFcfg3()
	cfga.InitWithFilePath2("fcfg_a.properties", true)
	cfgb.InitWithFilePath2("fcfg_b.properties", true)
	if len(cfga.Seces) != 2 {
		t.Error("error")
		return
	}
	cfga.Merge(cfgb)
	if len(cfga.Seces) != 4 {
		t.Error("error")
		return
	}
	fmt.Println(cfga.EnvReplaceV("${C_PWD}", false))
}

func TestFileMode(t *testing.T) {
	var cfg = NewFcfg3()
	cfg.SetVal("abc", "0644")
	cfg.SetVal("abc2", "06l44")
	fmt.Println(cfg.FileModeV("abc", os.ModePerm))
	fmt.Println(cfg.FileModeV("xx", os.ModePerm))
	fmt.Println(cfg.FileModeV("abc2", os.ModePerm))
}
