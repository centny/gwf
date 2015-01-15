package impl

import (
	"fmt"
	"github.com/Centny/gwf/netw"
	"github.com/Centny/gwf/pool"
	"github.com/Centny/gwf/util"
	"math"
	"runtime"
	"testing"
	"time"
)

type exec_s_c struct {
	i int
}
type exec_s struct {
}

func (t *exec_s_c) OnConn(c netw.Con) bool {
	c.SetWait(true)
	return true
}
func (t *exec_s_c) OnClose(c netw.Con) {

}
func (t *exec_s) OnCmd(c netw.Cmd) {
	defer c.Done()
	if c.Data()[0] == 1 {
		(c.(*rc_h_cmd)).Cmd.Writeb([]byte{1})
		return
	} else if c.Data()[0] == 2 {
		c.Writeb([]byte{11, 22, 3})
		c.Writeb([]byte{11, 22, 3, 33, 33})
		return
	} else if c.Data()[0] == 3 {
		c.Err(1, "abcc")
		return
	} else if c.Data()[0] == 4 {
		return
	}

	c.V(nil)
	c.Writeb([]byte("S-A"))
	time.Sleep(100 * time.Millisecond)
}

func TestExec(t *testing.T) {
	runtime.GOMAXPROCS(runtime.NumCPU() - 1)
	p := pool.NewBytePool(8, 1024)
	l := NewExecListener(p, ":7686", netw.NewCCH(&exec_s_c{}, &exec_s{}))
	l.T = 500
	err := l.Run()
	if err != nil {
		t.Error(err.Error())
		return
	}
	np, con, err := ExecDail(p, "127.0.0.1:7686")
	if err != nil {
		t.Error(err.Error())
		return
	}
	con.Start()
	time.Sleep(100 * time.Millisecond)
	go func() {
		go con.Con.Writeb([]byte{1})
		go con.Exec([]byte{1}, nil)
		go con.Exec([]byte{2}, nil)
		go con.Exec([]byte{3}, nil)
		go con.Exec([]byte{4}, nil)
		go con.Exec(nil, nil)
		for i := 0; i < 10; i++ {
			fmt.Println(con.Exec([]byte{22}, nil))
		}
	}()
	time.Sleep(time.Second)
	con.Close()
	go con.Exec([]byte{22}, nil)
	go func() {
		time.Sleep(2 * time.Second)
		np.Close()
		con.Stop()
		l.Close()
	}()
	l.Wait()
	time.Sleep(400 * time.Millisecond)
	(con.Con.(*netw.Con_)).V2B_ = func(v interface{}) ([]byte, error) {
		return nil, util.Err("error")
	}
	con.Exec([]byte{22}, nil)
	con.exec_c = math.MaxUint16
	con.Exec([]byte{22}, nil)
	con.Con = nil
	con.Exec([]byte{22}, nil)
	ExecDail(p, "addr")
	V2B_Byte("sfdsdfs")
}
