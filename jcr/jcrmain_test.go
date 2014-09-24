package jcr

import (
	"fmt"
	"github.com/Centny/gwf/util"
	"os"
	"time"
	// "os"
	"testing"
)

func TestMain(t *testing.T) {
	os.Args = []string{"jcr"}
	Run()
	os.Mkdir("out", os.ModePerm)
	os.Args = []string{"jcr", "app"}
	Run()
	os.Args = []string{
		"jcr", "app",
		"-d", "t",
		"-ex", ".*\\.json",
		"-in", ".*\\.html",
		"-o", "out",
		"-js", "http://localhost/jcr.js",
		"-inval", "jj",
	}
	Run()
	os.Args = []string{
		"jcr", "app",
		"-d", "t",
		"-ex", "adge.*\\.json",
		"-in", ".*\\.html",
		"-o", "out",
		"-js", "http://localhost/jcr.js",
		"-inval", "jj",
	}
	Run()
	os.Args = []string{
		"jcr", "app",
		"-d", "t",
		"-ex", ".*\\.html",
		"-in", "a.*\\.html",
		"-o", "out",
		"-js", "http://localhost/jcr.js",
		"-inval", "jj",
	}
	Run()
	os.Args = []string{
		"jcr", "app",
		"-d", "t",
		"-ex", ".*\\.json",
		"-in", "a.*\\.html",
		"-o", "out",
		"-js", "http://localhost/jcr.js",
		"-inval", "jj",
	}
	Run()
	os.Args = []string{
		"jcr", "app",
		"-d", "t",
		"-ex", ".*\\.json",
		"-in", ".*\\.html",
		"-o", "/srv",
		"-js", "http://localhost/jcr.js",
		"-inval", "jj",
	}
	Run()
	os.Args = []string{
		"jcr", "app",
		"-d", "t",
		"-ex", ".*\\.json",
		"-in", ".*\\.html",
		"-o", "/kkdfs/dsd",
		"-js", "http://localhost/jcr.js",
		"-inval", "jj",
	}
	Run()
}
func TestMain2(t *testing.T) {
	os.Mkdir("out", os.ModePerm)
	os.Args = []string{"jcr", "start", "-p", ":8798", "-n", "jjj"}
	go Run()
	time.Sleep(200)
	StopSrv()
	time.Sleep(100)
	//
	os.Args = []string{
		"jcr", "start",
		"-f", "jcr.json",
	}
	go Run()
	url := "http://localhost:8799"
	fmt.Println(util.HGet("%s/jcr/conf", url))
	fmt.Println(util.HGet("%s/jcr/jcr.js", url))
	fmt.Println(util.HGet("%s/jcr/store", url))
	fmt.Println(util.HGet("%s/jcr/store?cover=%s", url, "sssssss"))
	StopSrv()
}

func TestMatch(t *testing.T) {
	match([]string{""}, "//jjj")
}
