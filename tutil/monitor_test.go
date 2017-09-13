package tutil

import (
	"fmt"
	"runtime"
	"testing"
	"time"

	"github.com/Centny/gwf/util"
)

func TestMonitor(t *testing.T) {
	runtime.GOMAXPROCS(util.CPU())
	var m = NewMonitor()
	_, err := DoPerfV_(200, 30, "",
		func(idx, running int) (int, error) {
			return 1, nil
		},
		func(i int) error {
			if i%10 == 0 {
				m.State()
				return nil
			}
			var id = m.Start(fmt.Sprintf("_%v", i%3))
			time.Sleep(time.Duration(i) * time.Millisecond)
			m.Done(id)
			return nil
		})
	if err != nil {
		t.Error(err)
	}
	val, _ := m.State()
	fmt.Println(util.S2Json(val))
}
