package impl

import (
	"github.com/Centny/gwf/netw"
	"sync"
	"sync/atomic"
	"time"
)

//the chan handler struct.
type ChanH struct {
	H       netw.CmdHandler
	cc      chan netw.Cmd
	Wg      sync.WaitGroup
	Sleep   time.Duration
	count_  int32
	running bool
}

//new one chan handler.
func NewChanH(h netw.CmdHandler) *ChanH {
	return &ChanH{
		H:     h,
		cc:    make(chan netw.Cmd, 100),
		Sleep: 300,
	}
}
func NewChanH2(h netw.CmdHandler, gc int) *ChanH {
	ch := &ChanH{
		H:     h,
		cc:    make(chan netw.Cmd, 100),
		Sleep: 300,
	}
	ch.Run(gc)
	return ch
}

//running gorutine count
func (ch *ChanH) Count() int {
	return int(ch.count_)
}

//whether running or not.
func (ch *ChanH) Running() bool {
	return ch.running
}
func (ch *ChanH) OnCmd(c netw.Cmd) int {
	ch.cc <- c
	return 0
}

//run target number of gorutine to process command.
func (ch *ChanH) Run(gc int) {
	if ch.running {
		return
	}
	for i := 0; i < gc; i++ {
		go ch.run_c()
	}
}
func (ch *ChanH) run_c() {
	ch.Wg.Add(1)
	defer ch.Wg.Done()
	atomic.AddInt32(&ch.count_, 1)
	ch.running = true
	var cmd netw.Cmd = nil
	//tk := pool.NewTick(ch.Sleep * time.Millisecond)
	for ch.running {
		select {
		case cmd = <-ch.cc:
			if cmd == nil {
				break
			}
			ch.H.OnCmd(cmd)
			//case <-tk:
		}
	}
	atomic.AddInt32(&ch.count_, -1)
	if ch.count_ < 1 {
		ch.running = false
	}
	//pool.PutTick(ch.Sleep*time.Millisecond, tk)
}

//stop gorutine
func (ch *ChanH) Stop() {
	ch.running = false
	close(ch.cc)
}

//wait gorutine done.
func (ch *ChanH) Wait() {
	ch.Wg.Wait()
}
