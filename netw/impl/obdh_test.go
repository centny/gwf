package impl

import (
	"fmt"
	"github.com/Centny/gwf/netw"
	"github.com/Centny/gwf/pool"
	"github.com/Centny/gwf/util"
	"net"
	"runtime"
	"testing"
	"time"
)

type obdh_c_c struct {
}
type obdh_c struct {
}

func (t *obdh_c_c) OnConn(c netw.Con) bool {
	tcd := c.(*OBDH_Con)
	tcd.Con.Writeb([]byte("jjj"))
	tcd.Con.Writeb([]byte{2})
	c.Writeb([]byte("start"))
	c.Writev(nil)
	c.Writev("sss")
	c.Exec(nil, nil)
	return true
}
func (t *obdh_c_c) OnClose(c netw.Con) {

}
func (t *obdh_c) OnCmd(c netw.Cmd) int {
	fmt.Println("S->" + string(c.Data()))
	c.Writeb([]byte("C-A"))
	time.Sleep(100 * time.Millisecond)
	c.Done()
	return 0
}

type obdh_s_c struct {
	i int
}
type obdh_s struct {
}

func (t *obdh_s_c) OnConn(c netw.Con) bool {
	c.SetWait(true)
	return true
}
func (t *obdh_s_c) OnClose(c netw.Con) {

}
func (t *obdh_s) OnCmd(c netw.Cmd) int {
	c.V(nil)
	c.Writev(nil)
	fmt.Println("C->" + string(c.Data()))
	c.Writeb([]byte("S-A"))
	time.Sleep(100 * time.Millisecond)
	c.Done()
	// c.Err(1, "wwwsss---->")
	return 0
}

func TestOBDM(t *testing.T) {
	runtime.GOMAXPROCS(util.CPU())
	p := pool.NewBytePool(8, 1024)
	obdh := NewOBDH()
	obdh.AddH(1, &obdh_s{})
	l := netw.NewListener2(p, ":7686", netw.NewCCH(&obdh_s_c{}, obdh))
	l.T = 500
	err := l.Run()
	if err != nil {
		t.Error(err.Error())
		return
	}
	tc := &obdh_c{}
	c := netw.NewNConPool2(p, netw.NewCCH(&obdh_c_c{}, tc))
	c.NewCon = func(cp netw.ConPool, p *pool.BytePool, con net.Conn) netw.Con {
		con_ := netw.NewCon_(cp, p, con)
		con_.V2B_ = func(v interface{}) ([]byte, error) {
			if v == nil {
				return nil, util.Err("error")
			} else {
				return []byte{1, 2}, nil
			}
		}
		return NewOBDH_Con(1, con_)
	}
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
