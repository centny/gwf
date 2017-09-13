package tutil

import (
	"fmt"
	"runtime"
	"testing"
	"time"

	"github.com/Centny/gwf/log"
	"github.com/Centny/gwf/util"
)

func TestPerf(t *testing.T) {
	runtime.GOMAXPROCS(util.CPU())
	DoPerfV(20, 10, "", func(v int) {
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

func TestAutoPerf(t *testing.T) {
	runtime.GOMAXPROCS(util.CPU())
	used, max, err := DoAutoPerfV(1000, 10, 10, "t2.log", 100,
		func(idx, running int) error {
			log.D("running->%d,%d", running, idx)
			if running < 100 {
				return nil
			}
			return FullError
		}, func(v int) error {
			time.Sleep(100 * time.Millisecond)
			log.D("doing->%d", v)
			return nil
		})
	fmt.Println("used->", used, max, err)
}
