package netw

import (
	"code.google.com/p/go.net/websocket"
	"fmt"
	"github.com/Centny/gwf/pool"
	"github.com/Centny/gwf/routing/httptest"
	"net"
	"net/http"
	"runtime"
	"strings"
	"testing"
	"time"
)

type th_s struct {
	i int
	C CmdHandler
}

func (t *th_s) OnConn(c Con) bool {
	//testing queue falise reutrn
	dn := NewDoNoH()
	dn.C = false
	NewQueueConH(dn).OnConn(c)

	//
	if t.i == 0 {
		cwh := NewCWH(true)
		cwh.OnConn(c)
		cwh.OnClose(c)
		c.SetWait(true)
		c.SetId("---->1")
		t.i = 1
	} else if t.i == 1 {
		c.SetWait(false)
	} else {
		c.SetWait(true)
	}
	return true
}
func (t *th_s) OnCmd(c Cmd) int {
	fmt.Println("C->" + string(c.Data()))
	c.Writeb([]byte("S-A"))
	time.Sleep(100 * time.Millisecond)
	c.Done()
	t.C.OnCmd(c)
	return 0
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
func (t *th_c) OnCmd(c Cmd) int {
	c.Writev(nil)
	c.V(nil)
	fmt.Println("S->" + string(c.Data()))
	c.Writeb([]byte("C-A"))
	time.Sleep(100 * time.Millisecond)
	c.Done()
	c.CP()
	c.Err(1, "f")
	c.BaseCon()

	return 0
}
func (t *th_c) OnClose(c Con) {

}

type th_c2 struct {
	tt bool
}

func (t *th_c2) OnConn(c Con) bool {
	// c.R()
	c.Kvs()
	// c.W()
	c.Writev(nil)
	return t.tt
}
func (t *th_c2) OnCmd(c Cmd) int {
	return 0
}
func (t *th_c2) OnClose(c Con) {

}

//
func TestNetw(t *testing.T) {
	runtime.GOMAXPROCS(runtime.NumCPU() - 1)
	ShowLog = true
	p := pool.NewBytePool(8, 1024)
	ts := &th_s{C: NewDoNoH()}
	l := NewListener(p, ":7686", NewCCH(NewQueueConH(ts, NewDoNoH()), ts))
	l.T = 500
	go http.ListenAndServe(":7687", l.WsH())
	err := l.Run()
	if err != nil {
		t.Error(err.Error())
		return
	}
	tc := &th_c{}
	c := NewNConPool(p, tc)
	cc_con, err := c.Dail("127.0.0.1:7686")
	if err != nil {
		t.Error(err.Error())
		return
	}
	time.Sleep(200 * time.Millisecond)
	c.SetId("LLL")
	c.Find(cc_con.Id())
	c.Find("----->")
	//
	//
	tc2 := &th_c2{
		tt: true,
	}
	tc3 := &th_c2{
		tt: false,
	}
	c2 := NewNConPool(p, tc2)
	c2.NewCon = func(cp ConPool, l *pool.BytePool, con net.Conn) Con {
		cc := NewCon_(cp, p, con)
		cc.V2B_ = func(v interface{}) ([]byte, error) {
			return []byte{1}, nil
		}
		return cc
	}
	c2.Dail("127.0.0.1:7686")
	c3 := NewNConPool(p, tc3)
	c3.Dail("127.0.0.1:7686")
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
	go func() { //testing add twice for same connection.
		defer func() {
			fmt.Println(recover())
		}()
		c.add_c(cc_con)
		c.add_c(cc_con)
	}()
	go func() { //testing the write forbiden.
		defer func() {
			fmt.Println(recover())
		}()
		cc_con.Write([]byte{1})
	}()
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
	NewNConPool(p, tc).Dail("skfjs:dsfs")
	NewListener(p, "", &th_s{}).Run()
	NewListener(p, "ssfs:dsf", &th_s{}).Run()
	Dail(p, "addr", nil)
}

func TestPp(t *testing.T) {
	fmt.Println(fmt.Sprintf("%p", &th_s{}))
}

func TestWs(t *testing.T) {
	runtime.GOMAXPROCS(runtime.NumCPU() - 1)
	ShowLog = true
	p := pool.NewBytePool(8, 1024)
	ts_h := &th_s{C: NewDoNoH()}
	l := NewListener(p, ":7688", NewCCH(NewQueueConH(ts_h, NewDoNoH()), ts_h))
	err := l.Run()
	if err != nil {
		t.Error(err.Error())
		return
	}
	defer l.Close()
	time.Sleep(100 * time.Millisecond)
	ts := httptest.NewMuxServer()
	ts.Mux.Handler("^/.*$", l.WsH())
	fmt.Println(ts.S.URL)
	origin := ts.S.URL
	url := strings.Replace(origin, "http://", "ws://", -1)
	con, err := websocket.Dial(url, "", origin)
	if err != nil {
		t.Error(err.Error())
		return
	}
	tv, err := Writeb(con, []byte("AAXX"))
	if err != nil {
		t.Error(err.Error())
		return
	}
	fmt.Println(tv)
	time.Sleep(1 * time.Second)
	con.Close()
}
