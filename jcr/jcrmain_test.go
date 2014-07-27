package jcr

import (
	"fmt"
	"github.com/Centny/Cny4go/util"
	"os"
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
}
func TestMain2(t *testing.T) {
	os.Mkdir("out", os.ModePerm)
	os.Args = []string{"jcr", "start"}
	Run()
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
