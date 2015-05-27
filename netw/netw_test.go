package netw

import (
	"code.google.com/p/go.net/websocket"
	"encoding/binary"
	"fmt"
	"github.com/Centny/gwf/pool"
	"github.com/Centny/gwf/routing/httptest"
	"github.com/Centny/gwf/util"
	"math"
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
	dn := NewDoNotH()
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
	if t.i < 5 {
		c.Err(1, "sss")
	}
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
	runtime.GOMAXPROCS(util.CPU())
	ShowLog = true
	p := pool.NewBytePool(8, 1024)
	ts := &th_s{C: NewDoNotH()}
	l := NewListener2(p, ":7686", NewCCH(NewQueueConH(ts, NewDoNotH()), ts))
	l.T = 500
	go http.ListenAndServe(":7687", l.WsH())
	err := l.Run()
	if err != nil {
		t.Error(err.Error())
		return
	}
	tc := &th_c{}
	c := NewNConPool2(p, tc)
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
	c2 := NewNConPool2(p, tc2)
	c2.NewCon = func(cp ConPool, l *pool.BytePool, con net.Conn) Con {
		cc := NewCon_(cp, p, con)
		cc.V2B_ = func(v interface{}) ([]byte, error) {
			return []byte{1}, nil
		}
		return cc
	}
	c2.Dail("127.0.0.1:7686")
	c3 := NewNConPool2(p, tc3)
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
	//testing runner.
	fmt.Println("------>")
	ncr := NewNConRunner(p, "127.0.0.1:7686", tc)
	ncr.Retry = 500
	ncr.StopRunner() //only test,will do nothing.
	ncr.ConH = tc
	ncr.StartRunner()
	// ncr.Try()
	fmt.Println("---->")
	time.Sleep(200 * time.Millisecond)
	ncr.C.Close()
	time.Sleep(500 * time.Millisecond)
	ncr.StopRunner()
	time.Sleep(500 * time.Millisecond)

	//
	l.Close()
	l.Wait()
	//
	NewNConPool2(p, tc).Dail("skfjs:dsfs")
	NewListener2(p, "", &th_s{}).Run()
	NewListener2(p, "ssfs:dsf", &th_s{}).Run()
	Dail(p, "addr", nil)
	ncr.StartRunner()
	time.Sleep(500 * time.Millisecond)
	ncr.StopRunner()
	time.Sleep(500 * time.Millisecond)
	ncr.ConH = nil
	ncr.OnConn(nil)
}

func TestPp(t *testing.T) {
	fmt.Println(fmt.Sprintf("%p", &th_s{}))
}

func TestWs(t *testing.T) {
	runtime.GOMAXPROCS(util.CPU())
	ShowLog = true
	p := pool.NewBytePool(8, 1024)
	ts_h := &th_s{C: NewDoNotH()}
	l := NewListener2(p, ":7688", NewCCH(NewQueueConH(ts_h, NewDoNotH()), ts_h))
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
	tv, err := Writen(con, []byte("AAXX"))
	if err != nil {
		t.Error(err.Error())
		return
	}
	fmt.Println(tv)
	time.Sleep(1 * time.Second)
	con.Close()
}

func TestTick(t *testing.T) {
	runtime.GOMAXPROCS(util.CPU())
	ShowLog = true
	p := pool.NewBytePool(8, 1024)
	l := NewListener(p, ":6679", "N", NewDoNotH())
	l.Run()
	defer l.Close()
	con := NewNConRunner(p, "127.0.0.1:6679", NewDoNotH())
	con.ShowLog = true
	con.Tick = 100
	con.StartRunner()
	go func() {
		time.Sleep(100 * time.Millisecond)
		l.Writev("abc")
		l.Writeb([]byte("xxxx"))
	}()
	time.Sleep(time.Second)
	con.StopRunner()
	time.Sleep(100 * time.Millisecond)
}
func TestLoopTimeout(t *testing.T) {
	runtime.GOMAXPROCS(util.CPU())
	ShowLog = true
	p := pool.NewBytePool(8, 1024)
	dn := NewDoNotH()
	dn.W = false
	l := NewListener(p, ":6679", "N", dn)
	l.T = 100
	l.Run()
	defer l.Close()
	nc, _, _ := Dail(p, "127.0.0.1:6679", NewDoNotH())
	time.Sleep(time.Second)
	if len(nc.Cons()) > 0 {
		t.Error("not right")
	}
}

type PanicCmd struct {
}

func (p *PanicCmd) OnCmd(c Cmd) int {
	c.SetErrd(2)
	fmt.Println("PanicCmd------->OnCmd")
	panic("OnCmd")
	return 0
}
func TestConRecover(t *testing.T) {
	runtime.GOMAXPROCS(util.CPU())
	ShowLog = true
	p := pool.NewBytePool(8, 1024)
	l := NewListener(p, ":6679", "N", NewCCH(NewDoNotH(), &PanicCmd{}))
	l.T = 100
	l.Run()
	defer l.Close()
	nc, _, _ := Dail(p, "127.0.0.1:6679", NewDoNotH())
	nc.Writeb([]byte("sfdsfsd"))
	time.Sleep(time.Second)
	if len(nc.Cons()) > 0 {
		t.Error("not right")
	}
}

func TestSome(t *testing.T) {
	dn := NewDoNotH()
	dn.ShowLog = true
	dn.log_d("abccc->%v", 1)
	NewListenerN2(nil, "", nil, nil)
}

func TestBesss(t *testing.T) {
	fmt.Println(math.Pow(2, 16))
	buf := make([]byte, 2)
	binary.BigEndian.PutUint16(buf, 41828)
	fmt.Println(buf)
	util.FWrite2("/tmp/tt.data", buf)
}

// func TestOnlyWrite(t *testing.T) {
// 	runtime.GOMAXPROCS(util.CPU())
// 	l, err := net.Listen("tcp", ":8435")
// 	if err != nil {
// 		t.Error(err.Error())
// 		return
// 	}
// 	defer l.Close()s
// 	go func() {
// 		for {
// 			l.Accept()
// 		}
// 	}()
// 	con, err := net.Dial("tcp", ":8435")
// 	if err != nil {
// 		t.Error(err.Error())
// 		return
// 	}
// 	for {
// 		_, err = con.Write([]byte("123456789012345678901234567890123456789012345678901234567890"))
// 		fmt.Println(err)
// 		if err != nil {
// 			break
// 		}
// 	}
// }
