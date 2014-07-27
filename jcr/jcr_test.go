package jcr

import (
	"fmt"
	"github.com/Centny/Cny4go/util"
	"os"
	"testing"
)

func TestJcr(t *testing.T) {
	StartSrv("abcc")
	StartSrv("jcr_t.jsone")
	go RunSrv("jcr.json")
	url := "http://localhost:8799"
	fmt.Println(util.HGet("%s/jcr/conf", url))
	fmt.Println(util.HGet("%s/jcr/jcr.js", url))
	fmt.Println(util.HGet("%s/jcr/store", url))
	fmt.Println(util.HGet("%s/jcr/store?cover=%s", url, "sssssss"))
	StopSrv()
	os.Remove("JCR_000.json")
}
