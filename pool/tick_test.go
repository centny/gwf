package pool

import (
	"fmt"
	"github.com/Centny/gwf/util"
	"testing"
	"time"
)

func TestTick(t *testing.T) {
	var tick = NewTick(100 * time.Millisecond)
	var beg = util.Now()
	for i := 0; i < 10; i++ {
		<-tick
	}
	var used = util.Now() - beg
	if used > 1100 || used < 900 {
		t.Error("error")
		return
	}
	fmt.Println(used)
	//
	//
	PutTick(100*time.Millisecond, tick)
	if len(Ticks) < 1 || Ticks[100*time.Millisecond].Len() < 1 {
		t.Error("error")
		return
	}
	time.Sleep(time.Second)
	//
	tick = NewTick(100 * time.Millisecond)
	if Ticks[100*time.Millisecond].Len() > 0 {
		t.Error("error")
		return
	}
	beg = util.Now()
	for i := 0; i < 10; i++ {
		<-tick
	}
	used = util.Now() - beg
	if used > 1100 || used < 800 {
		t.Error(used)
		return
	}
	fmt.Println(used)

}
