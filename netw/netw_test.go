package netw

import (
	"fmt"
	"github.com/Centny/gwf/pool"
	"net"
	"runtime"
	"testing"
	"time"
)

type th_s struct {
	i int
}

func (t *th_s) OnConn(c Con) bool {
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
func (t *th_s) OnCmd(c Cmd) {
	fmt.Println("C->" + string(c.Data()))
	c.Writeb([]byte("S-A"))
	time.Sleep(100 * time.Millisecond)
	c.Done()
}
func (t *th_s) OnClose(c Con) {

}

type th_c struct {
}

func (t *th_c) OnConn(c Con) bool {
	c.Writeb([]byte("start"))
	c.Exec(nil, nil)
	return true
}
func (t *th_c) OnCmd(c Cmd) {
	c.Writev(nil)
	c.V(nil)
	fmt.Println("S->" + string(c.Data()))
	c.Writeb([]byte("C-A"))
	time.Sleep(100 * time.Millisecond)
	c.Done()
}
func (t *th_c) OnClose(c Con) {

}

type th_c2 struct {
	tt bool
}

func (t *th_c2) OnConn(c Con) bool {
	c.R()
	c.Kvs()
	c.W()
	c.Writev(nil)
	return t.tt
}
func (t *th_c2) OnCmd(c Cmd) {
}
func (t *th_c2) OnClose(c Con) {

}

//
func TestNetw(t *testing.T) {
	runtime.GOMAXPROCS(runtime.NumCPU() - 1)
	ShowLog = true
	p := pool.NewBytePool(8, 1024)
	l := NewListener(p, ":7686", &th_s{})
	l.T = 500
	err := l.Run()
	if err != nil {
		t.Error(err.Error())
		return
	}
	tc := &th_c{}
	c := NewNConPool(p, "127.0.0.1:7686", tc)
	_, err = c.Dail()
	if err != nil {
		t.Error(err.Error())
		return
	}
	time.Sleep(200 * time.Millisecond)

	tc2 := &th_c2{
		tt: true,
	}
	tc3 := &th_c2{
		tt: false,
	}
	c2 := NewNConPool(p, "127.0.0.1:7686", tc2)
	c2.NewCon = func(cp ConPool, l *pool.BytePool, con net.Conn) Con {
		cc := NewCon_(cp, p, con)
		cc.V2B_ = func(v interface{}) ([]byte, error) {
			return []byte{1}, nil
		}
		return cc
	}
	c2.Dail()
	c3 := NewNConPool(p, "127.0.0.1:7686", tc3)
	c3.Dail()
	time.Sleep(2 * time.Second)
	c.Close()
	c2.Close()
	c3.Close()
	//
	cc, err := net.Dial("tcp", "127.0.0.1:7686")
	if err != nil {
		t.Error(err.Error())
		return
	}
	cc.Write([]byte("jkk"))
	cc.Write([]byte{0, 0})
	time.Sleep(100 * time.Millisecond)
	cc.Close()
	cc2, err := net.Dial("tcp", "127.0.0.1:7686")
	if err != nil {
		t.Error(err.Error())
		return
	}
	cc2.Write([]byte(H_MOD))
	cc2.Write([]byte{0, 0})
	time.Sleep(100 * time.Millisecond)
	cc2.Close()
	cc3, err := net.Dial("tcp", "127.0.0.1:7686")
	if err != nil {
		t.Error(err.Error())
		return
	}
	cc3.Write([]byte(H_MOD))
	cc3.Write([]byte{1, 0})
	time.Sleep(100 * time.Millisecond)
	cc3.Close()
	//
	l.Close()
	l.Wait()
	//
	NewNConPool(p, "skfjs:dsfs", tc).Dail()
	NewListener(p, "", &th_s{}).Run()
	NewListener(p, "ssfs:dsf", &th_s{}).Run()
	Dail(p, "addr", nil)
}
