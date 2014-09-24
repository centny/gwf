package jcr

import (
	"fmt"
	"github.com/Centny/gwf/util"
	"os"
	"testing"
)

func TestJcr(t *testing.T) {
	// StartSrv("abcc")
	StartSrv("coverage_jcr", "out", "jsss")
	go RunSrv("coverage_jcr", "out", ":8799")
	url := "http://localhost:8799"
	fmt.Println(util.HGet("%s/jcr/conf", url))
	fmt.Println(util.HGet("%s/jcr/jcr.js", url))
	fmt.Println(util.HGet("%s/jcr/store", url))
	fmt.Println(util.HGet("%s/jcr/store?cover=%s", url, "sssssss"))
	_conf_.Dir = "/tjjks"
	fmt.Println(util.HGet("%s/jcr/store?cover=%s", url, "sssssss"))
	_conf_.Dir = "out"
	fmt.Println(util.HGet("%s/jcr/exit", url))
	StopSrv()
	os.Remove("JCR_000.json")
}
