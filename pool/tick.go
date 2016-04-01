package pool

import (
	"container/list"
	"sync"
	"time"
)

var Ticks = map[time.Duration]*list.List{}
var TickL = sync.RWMutex{}

func NewTick(d time.Duration) <-chan time.Time {
	TickL.Lock()
	defer TickL.Unlock()
	var ls = Ticks[d]
	if ls == nil || ls.Len() < 1 {
		return time.Tick(d)
	}
	var val = ls.Front()
	ls.Remove(val)
	return val.Value.(<-chan time.Time)
}

func PutTick(d time.Duration, t <-chan time.Time) {
	TickL.Lock()
	defer TickL.Unlock()
	var ls = Ticks[d]
	if ls == nil {
		ls = list.New()
	}
	ls.PushBack(t)
	Ticks[d] = ls
}
