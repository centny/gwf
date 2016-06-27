package impl

import (
	"fmt"
	"github.com/Centny/gwf/log"
	"github.com/Centny/gwf/netw"
	"github.com/Centny/gwf/tutil"
	"github.com/Centny/gwf/util"
	"sync"
	"sync/atomic"
	"time"
)

//the chan handler struct.
type ChanH struct {
	H        netw.CmdHandler
	cc       chan netw.Cmd
	Wg       sync.WaitGroup
	Sleep    time.Duration
	M        *tutil.Monitor
	count_   int32
	process_ int32
	running  bool
	Name     string
	Idle     int
	Max      int
}

//new one chan handler.
func NewChanH(h netw.CmdHandler) *ChanH {
	return &ChanH{
		H:     h,
		cc:    make(chan netw.Cmd, 100),
		Sleep: 300,
		Idle:  util.CPU(),
	}
}
func NewChanH2(h netw.CmdHandler, gc int) *ChanH {
	ch := &ChanH{
		H:     h,
		cc:    make(chan netw.Cmd, 100),
		Sleep: 300,
		Idle:  util.CPU(),
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
	if ch.M != nil {
		ch.M.Start_(fmt.Sprintf("chan/%p", c))
	}
	ch.cc <- c
	process := int(atomic.LoadInt32(&ch.process_))
	running := int(atomic.LoadInt32(&ch.count_))
	if running < ch.Max && (running-process) < len(ch.cc) {
		go ch.run_c()
	}
	return 0
}

//run target number of gorutine to process command.
func (ch *ChanH) Run(gc int) {
	if ch.running {
		return
	}
	if gc < 1 {
		panic("ChanH at last one core is reqquired")
	}
	ch.Max = gc
	log.D("ChanH(%v) start run by %v max core", ch.Name, gc)
	// for i := 0; i < gc; i++ {
	// 	go ch.run_c()
	// }
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
			atomic.AddInt32(&ch.process_, 1)
			data := cmd.Data()
			mid := ""
			if ch.M != nil {
				ch.M.Done(fmt.Sprintf("chan/%p", cmd))
				mid = ch.M.Start(fmt.Sprintf("C->%v", data[0]))
			}
			ch.H.OnCmd(cmd)
			if ch.M != nil {
				ch.M.Done(mid)
			}
			atomic.AddInt32(&ch.process_, -1)
			if len(ch.cc) < 1 && int(atomic.LoadInt32(&ch.count_)) > ch.Idle {
				break
			}
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
