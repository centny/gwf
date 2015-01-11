package handler

import (
	"fmt"
	"github.com/Centny/gwf/netw"
	"github.com/Centny/gwf/pool"
	"runtime"
	"testing"
	"time"
)

type th_c struct {
}

func (t *th_c) OnConn(c *netw.Con) bool {
	c.Write([]byte("start"))
	return true
}
func (t *th_c) OnCmd(c *netw.Cmd) {
	fmt.Println("S->" + string(c.Data))
	c.Con.Write([]byte("C-A"))
	time.Sleep(100 * time.Millisecond)
	c.Done()
}
func (t *th_c) OnClose(c *netw.Con) {

}

type th_s struct {
	i int
}

func (t *th_s) OnConn(c *netw.Con) bool {
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
func (t *th_s) OnCmd(c *netw.Cmd) {
	fmt.Println("C->" + string(c.Data))
	c.Con.Write([]byte("S-A"))
	time.Sleep(100 * time.Millisecond)
	c.Done()
}
func (t *th_s) OnClose(c *netw.Con) {

}
func TestChan(t *testing.T) {
	runtime.GOMAXPROCS(runtime.NumCPU() - 1)
	p := pool.NewBytePool(8, 1024)
	ts := NewChanH(&th_s{})
	ts.Run(5)
	l := netw.NewListener(p, ":7686", ts)
	l.T = 500
	err := l.Run()
	if err != nil {
		t.Error(err.Error())
		return
	}
	tc := &th_c{}
	c := netw.NewNConPool(p, "127.0.0.1:7686", tc)
	err = c.Dail()
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
