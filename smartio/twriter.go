package smartio

import (
	"bufio"
	"fmt"
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
	buf     *bufio.Writer
	LCK     sync.RWMutex
}

func NewTWriter(sw io.Writer) *TimeFlushWriter {
	return NewTimeWriter(sw, FLOG_BUF_SIZE, FLOG_CLOCK_DELAY)
}
func NewTimeWriter(sw io.Writer, bsize int, cdelay time.Duration) *TimeFlushWriter {
	fl := &TimeFlushWriter{}
	//
	fl.bsize = bsize
	fl.sw = sw
	fl.buf = bufio.NewWriterSize(sw, fl.bsize)
	//
	fl.cdelay = cdelay
	fl.rdelay = 1000
	fl.running = true
	go fl.runClock()
	//
	return fl
}
func (t *TimeFlushWriter) Write(p []byte) (nn int, err error) {
	t.LCK.Lock()
	defer t.LCK.Unlock()
	return t.buf.Write(p)
}
func (t *TimeFlushWriter) runClock() {
	slog("TimeWriter clock start...")
	wg.Add(1)
	var ttime time.Duration = 0
	for t.running {
		if ttime >= t.cdelay && t.buf.Buffered() > 0 {
			slog("TimeWriter do flush...")
			t.LCK.Lock()
			err := t.buf.Flush()
			if err != nil {
				fmt.Fprintf(LOG, "flush error for wirter(%v) info(%v,%v):%v\n",
					t.sw, t.buf.Available(), t.buf.Buffered(), err.Error())
			}
			ttime = 0
			t.LCK.Unlock()
		}
		ttime += t.rdelay
		time.Sleep(t.rdelay * time.Millisecond)
	}
	wg.Done()
	slog("TimeWriter clock end...")
}
func (t *TimeFlushWriter) Stop() {
	// fmt.Println("Stop TimeWriter")
	t.buf.Flush()
	t.running = false
}
func (t *TimeFlushWriter) WriteString(s string) (int, error) {
	return t.buf.WriteString(s)
}
func TimeWriterWait() {
	// fmt.Println("Waiting all TimeWriter stop")
	wg.Wait()
}
