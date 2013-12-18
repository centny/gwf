package twriter

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

type TimeWriter struct {
	sw      io.Writer
	w       *bufio.Writer
	bsize   int
	cdelay  time.Duration
	running bool
}

func NewTWriter(sw io.Writer) *TimeWriter {
	return NewTimeWriter(sw, FLOG_BUF_SIZE, FLOG_CLOCK_DELAY)
}
func NewTimeWriter(sw io.Writer, bsize int, cdelay time.Duration) *TimeWriter {
	fl := &TimeWriter{}
	//
	fl.bsize = bsize
	fl.sw = sw
	fl.w = bufio.NewWriterSize(sw, fl.bsize)
	//
	fl.cdelay = cdelay
	fl.running = true
	go fl.runClock()
	wg.Add(1)
	//
	return fl
}
func (t *TimeWriter) runClock() {
	log.D("TimeWriter clock start...")
	for t.running {
		t.w.Flush()
		time.Sleep(t.cdelay * time.Millisecond)
	}
	wg.Done()
	log.D("TimeWriter clock end...")
}
func (t *TimeWriter) Stop() {
	t.w.Flush()
	t.running = false
}
func (t *TimeWriter) Writer() *bufio.Writer {
	return t.w
}
func Wait() {
	log.D("Waiting all TimeWriter stop")
	wg.Wait()
}
