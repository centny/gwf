package timer

import (
	"fmt"
	"github.com/Centny/gwf/util"
	"testing"
	"time"
)

func on_time(i uint64) error {
	fmt.Println("on time", i)
	return nil
}
func on_time_e(i uint64) error {
	fmt.Println("on time error", i)
	return util.Err("error")
}

type time_p struct {
}

func (t *time_p) OnTime(i uint64) error {
	panic("testing")
}

func (t *time_p) Name() string {
	return "xxxx"
}

func TestTimer(t *testing.T) {
	ShowLog = true
	tp := &time_p{}
	Register2(100, on_time)
	Register2(100, on_time_e)
	Register(100, tp)
	time.Sleep(time.Second)
	Remove2(on_time)
	Remove2(on_time_e)
	Remove(tp)
	time.Sleep(time.Second)
	Stop()
}
