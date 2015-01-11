package handler

import (
	"fmt"
	"github.com/Centny/gwf/netw"
	"github.com/Centny/gwf/pool"
	"math"
	"runtime"
	"sync"
	"testing"
)

type trcmd_s struct {
}

func (t *trcmd_s) OnConn(c *netw.Con) bool {
	return true
}
func (t *trcmd_s) OnCmd(c *RC_Cmd) {
	if string(c.Data) == "A-0" {
		c.Cmd.Write([]byte{1})
	} else if string(c.Data) == "A-1" {
		c.Cmd.Write([]byte{255, 244})
	}
	c.Write(c.Data)
}
func (t *trcmd_s) OnClose(c *netw.Con) {
}
func TestRcmd(t *testing.T) {
	runtime.GOMAXPROCS(runtime.NumCPU() - 1)
	p := pool.NewBytePool(8, 1024)
	ts := NewChanH(NewRC_S(&trcmd_s{}))
	ts.Run(5)
	l := netw.NewListener(p, ":7686", ts)
	l.T = 500
	err := l.Run()
	if err != nil {
		t.Error(err.Error())
		return
	}
	tc := NewRC_C()
	fmt.Println(tc.Exec([]byte{}))
	fmt.Println(tc.Exec([]byte{1}))
	c := netw.NewNConPool(p, "127.0.0.1:7686", tc)
	err = c.Dail()
	if err != nil {
		t.Error(err.Error())
		return
	}
	tc.Start()
	wg := sync.WaitGroup{}
	for i := 0; i < 100; i++ {
		go func(tv int) {
			wg.Add(1)
			defer wg.Done()
			sv := fmt.Sprintf("A-%v", tv)
			bv, err := tc.Exec([]byte(sv))
			if err != nil {
				t.Error(err.Error())
				return
			}
			if string(bv.Data) != sv {
				t.Error(fmt.Sprintf("not equal %v:%v", string(bv.Data), tv))
			}
			fmt.Println(sv)
		}(i)
	}
	go func() {
		wg.Wait()
		c.Close()
		fmt.Println(tc.Exec([]byte("sssss")))
		tc.Stop()
		l.Close()
		ts.Stop()
	}()
	fmt.Println(ts.Running())
	fmt.Println(ts.Count())

	//
	l.Wait()
	ts.Wait()

	//
	tc.exec_c = math.MaxUint16
	tc.Exec([]byte{1})
}
