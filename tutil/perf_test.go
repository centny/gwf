package tutil

import (
	"fmt"
	"github.com/Centny/gwf/log"
	"github.com/Centny/gwf/util"
	"runtime"
	"testing"
	"time"
)

func TestPerf(t *testing.T) {
	runtime.GOMAXPROCS(util.CPU())
	DoPerf(10, "", func(v int) {
		time.Sleep(100 * time.Millisecond)
		log.D("doing->%d", v)
	})
	used, err := DoPerf(1000, "t.log", func(v int) {
		time.Sleep(100 * time.Millisecond)
		log.D("doing->%d", v)
		fmt.Println(v)
	})
	if err != nil {
		t.Error(err.Error())
		return
	}
	fmt.Println("used->", used)

}
