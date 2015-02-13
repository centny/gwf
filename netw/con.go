package netw

import (
	"github.com/Centny/gwf/log"
	"github.com/Centny/gwf/pool"
	"github.com/Centny/gwf/util"
	"net"
	"time"
)

// type NCon struct {
// 	Addr string
// 	Con
// }

// func NewNCon(con Con, addr string) Con {
// 	return &NCon{
// 		Con:  con,
// 		Addr: addr,
// 	}
// }

//the client connection pool.
type NConPool struct {
	*LConPool //base connection pool.
}

//new client connection pool.
func NewNConPool(p *pool.BytePool, h CCHandler, n string) *NConPool {
	return &NConPool{
		LConPool: NewLConPool(p, h, n),
	}
}
func NewNConPool2(p *pool.BytePool, h CCHandler) *NConPool {
	return NewNConPool(p, h, "C-")
}

//dail one connection.
func (n *NConPool) Dail(addr string) (Con, error) {
	con, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}
	cc := n.NewCon(n, n.P, con)
	if !n.H.OnConn(cc) {
		return nil, util.Err("OnConn return false for %v", addr)
	}
	n.RunC(cc)
	return cc, nil
}

func Dail(p *pool.BytePool, addr string, h CCHandler) (*NConPool, Con, error) {
	return DailN(p, addr, h, NewCon)
}
func DailN(p *pool.BytePool, addr string, h CCHandler, ncf NewConF) (*NConPool, Con, error) {
	nc := NewNConPool2(p, h)
	nc.NewCon = ncf
	cc, err := nc.Dail(addr)
	return nc, cc, err
}

type NConRunner struct {
	*NConPool
	C         Con
	ConH      ConHandler
	Connected bool
	Running   bool
	Retry     time.Duration
	Tick      time.Duration
	//
	Addr string
	NCF  NewConF
	BP   *pool.BytePool
	CmdH CmdHandler
}

func (n *NConRunner) OnConn(c Con) bool {
	if n.ConH == nil {
		return true
	}
	return n.ConH.OnConn(c)
}
func (n *NConRunner) OnClose(c Con) {
	if n.Running {
		go n.Try()
	}
	if n.ConH != nil {
		n.ConH.OnClose(c)
	}
}
func (n *NConRunner) StartRunner() {
	go n.Try()
}
func (n *NConRunner) StopRunner() {
	n.Running = false
	if n.NConPool == nil {
		return
	}
	n.NConPool.Close()
}
func (n *NConRunner) StartTick() {
	go n.RunTick_()
}
func (n *NConRunner) RunTick_() {
	tk := time.Tick(n.Tick * time.Millisecond)
	c := n.C
	n.Running = true
	for n.Running {
		select {
		case <-tk:
			c = n.C
			if c != nil {
				c.Writeb([]byte("Tick\n"))
				log_d("sending tick message to Push Server")
			}
		}
	}
}
func (n *NConRunner) Try() {
	n.Running = true
	for n.Running {
		err := n.Dail()
		if err == nil {
			break
		}
		log.D("try connect to server(%v) err:%v,will retry after %v ms", n.Addr, err.Error(), 5000)
		time.Sleep(n.Retry * time.Millisecond)
	}
	log.D("connect try stopped")
}
func (n *NConRunner) Dail() error {
	n.Connected = false
	nc, cc, err := DailN(n.BP, n.Addr, NewCCH(n, n.CmdH), n.NCF)
	if err != nil {
		return err
	}
	n.NConPool = nc
	n.C = cc
	n.Connected = true
	return nil
}

func NewNConRunnerN(bp *pool.BytePool, addr string, h CmdHandler, ncf NewConF) *NConRunner {
	return &NConRunner{
		Addr:  addr,
		NCF:   ncf,
		BP:    bp,
		CmdH:  h,
		Retry: 5000,
		Tick:  30000,
	}
}
func NewNConRunner(bp *pool.BytePool, addr string, h CmdHandler) *NConRunner {
	return NewNConRunnerN(bp, addr, h, NewCon)
}
