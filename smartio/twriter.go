package smartio

import (
	"bufio"
	"github.com/Centny/Cny4go/log"
	"io"
	"sync"
	"time"
)

//bytes
const FLOG_BUF_SIZE int = 1024000

//millisecond
const FLOG_CLOCK_DELAY time.Duration = 10000

var wg *sync.WaitGroup = &sync.WaitGroup{}

type TimeFlushWriter struct {
	sw      io.Writer
	bsize   int
	cdelay  time.Duration
	rdelay  time.Duration
	running bool
	*bufio.Writer
}

func NewTWriter(sw io.Writer) *TimeFlushWriter {
	return NewTimeWriter(sw, FLOG_BUF_SIZE, FLOG_CLOCK_DELAY)
}
func NewTimeWriter(sw io.Writer, bsize int, cdelay time.Duration) *TimeFlushWriter {
	fl := &TimeFlushWriter{}
	//
	fl.bsize = bsize
	fl.sw = sw
	fl.Writer = bufio.NewWriterSize(sw, fl.bsize)
	//
	fl.cdelay = cdelay
	fl.rdelay = 1000
	fl.running = true
	go fl.runClock()
	//
	return fl
}
func (t *TimeFlushWriter) runClock() {
	// fmt.Println("TimeWriter clock start...")
	wg.Add(1)
	var ttime time.Duration = 0
	for t.running {
		if ttime >= t.cdelay && t.Buffered() > 0 {
			err := t.Flush()
			if err != nil {
				log.E("flush error for wirter(%v) info(%v,%v):%v", t.sw, t.Available(), t.Buffered(), err.Error())
			}
			ttime = 0
		}
		ttime += t.rdelay
		time.Sleep(t.rdelay * time.Millisecond)
	}
	wg.Done()
	// fmt.Println("TimeWriter clock end...")
}
func (t *TimeFlushWriter) Stop() {
	// fmt.Println("Stop TimeWriter")
	t.Flush()
	t.running = false
}
func TimeWriterWait() {
	// fmt.Println("Waiting all TimeWriter stop")
	wg.Wait()
}
