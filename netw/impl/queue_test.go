package impl

import (
	"github.com/Centny/gwf/netw"
	"github.com/Centny/gwf/pool"
	"runtime"
	"testing"
	"time"
)

type th_s2 struct {
	i int
}

func (t *th_s2) OnCmd(c netw.Cmd) int {
	t.i++
	if t.i%3 == 0 {
		return 1
	}
	return 0
}

func TestQueue(t *testing.T) {
	runtime.GOMAXPROCS(util.CPU())
	p := pool.NewBytePool(8, 1024)
	ts := NewQueueH()
	ts.AddH(&th_s2{})
	ts.AddH(&th_s{})
	l := netw.NewListener2(p, ":7686", netw.NewCCH(&th_s_c{}, ts))
	l.T = 500
	err := l.Run()
	if err != nil {
		t.Error(err.Error())
		return
	}
	tc := &th_c{}
	c := netw.NewNConPool2(p, netw.NewCCH(&th_c_c{}, tc))
	_, err = c.Dail("127.0.0.1:7686")
	if err != nil {
		t.Error(err.Error())
		return
	}
	go func() {
		time.Sleep(2 * time.Second)
		c.Close()
		l.Close()
	}()
	l.Wait()
}
