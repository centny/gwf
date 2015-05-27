package impl

import (
	"fmt"
	"github.com/Centny/gwf/netw"
	"github.com/Centny/gwf/pool"
	"runtime"
	"testing"
	"time"
)

type th_c_c struct {
}
type th_c struct {
}

func (t *th_c_c) OnConn(c netw.Con) bool {
	c.Writeb([]byte("start"))
	return true
}
func (t *th_c_c) OnClose(c netw.Con) {

}
func (t *th_c) OnCmd(c netw.Cmd) int {
	fmt.Println("S->" + string(c.Data()))
	c.Writeb([]byte("C-A"))
	time.Sleep(100 * time.Millisecond)
	c.Done()
	return 0
}

type th_s_c struct {
	i int
}
type th_s struct {
}

func (t *th_s_c) OnConn(c netw.Con) bool {
	if t.i == 0 {
		c.SetWait(true)
		t.i = 1
	} else if t.i == 1 {
		c.SetWait(false)
	} else {
		c.SetWait(true)
	}
	return true
}
func (t *th_s_c) OnClose(c netw.Con) {

}
func (t *th_s) OnCmd(c netw.Cmd) int {
	fmt.Println("C->" + string(c.Data()))
	c.Writeb([]byte("S-A"))
	time.Sleep(100 * time.Millisecond)
	c.Done()
	return 0
}

func TestChan(t *testing.T) {
	runtime.GOMAXPROCS(util.CPU())
	p := pool.NewBytePool(8, 1024)
	ts := NewChanH(&th_s{})
	ts.Run(5)
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
		ts.Stop()
	}()
	fmt.Println(ts.Running())
	fmt.Println(ts.Count())
	l.Wait()
	ts.Wait()
}
