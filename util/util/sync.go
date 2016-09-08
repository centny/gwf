package util

import (
	"sync"
	"sync/atomic"
)

//WaitGroup the same of sync.WaitGroup and adding wait count record
type WaitGroup struct {
	sync.WaitGroup
	c int32
}

//Add adding count
func (w *WaitGroup) Add(i int) {
	w.WaitGroup.Add(i)
	atomic.AddInt32(&w.c, int32(i))
}

//Done done one count
func (w *WaitGroup) Done() {
	w.WaitGroup.Done()
	atomic.AddInt32(&w.c, int32(-1))
}

//Size return current wait size
func (w *WaitGroup) Size() int {
	return int(w.c)
}
