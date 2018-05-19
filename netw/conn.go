package netw

import (
	"net"
	"sync"
	"time"

	"github.com/Centny/gwf/log"
	"github.com/Centny/gwf/pool"
	"github.com/Centny/gwf/util"
)

//the client connection pool.
type NConPool struct {
	*LConPool //base connection pool.
	Dailer    *AutoDailer
	DailAddr  func(addr string) (net.Conn, error)
}

//new client connection pool.
func NewNConPool(p *pool.BytePool, h CCHandler, n string) *NConPool {
	var dailer = NewAutoDailer()
	var ch = NewCCH(NewQueueConH(dailer, h), h)
	var cp = &NConPool{
		LConPool: NewLConPoolV(p, ch, n, NewConH),
		Dailer:   dailer,
		DailAddr: func(addr string) (net.Conn, error) {
			return net.Dial("tcp", addr)
		},
	}
	cp.Dailer.Dail = cp.Dail
	return cp
}
func NewNConPool2(p *pool.BytePool, h CCHandler) *NConPool {
	return NewNConPool(p, h, "C-")
}

//dail one connection.
func (n *NConPool) Dail(addr string) (Con, error) {
	con, err := n.DailAddr(addr)
	if err != nil {
		return nil, err
	}
	cc, err := n.NewCon(n, n.P, con)
	if err != nil {
		con.Close()
		return nil, util.Err("NewCon return error(%v)", err)
	}
	if !n.H.OnConn(cc) {
		return nil, util.Err("OnConn return false for %v", addr)
	}
	n.RunC(cc)
	return cc, nil
}

func (n *NConPool) Close() {
	n.Dailer.Stop()
	n.LConPool.Close()
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

type AutoDailer struct {
	Addrs      map[Con]string
	a_lck      sync.RWMutex
	Running    bool
	Dail       func(addr string) (Con, error)
	OnDailFail func(addr string, err error)
}

func NewAutoDailer() *AutoDailer {
	return &AutoDailer{
		Addrs: map[Con]string{},
	}
}

func (a *AutoDailer) OnConn(c Con) bool {
	return true
}

func (a *AutoDailer) OnClose(c Con) {
	if a.Running {
		go a.Try(a.Addrs[c])
	}
	a.a_lck.Lock()
	delete(a.Addrs, c)
	a.a_lck.Unlock()
}
func (a *AutoDailer) DailAll(addrs []string) {
	for _, addr := range addrs {
		go a.Try(addr)
	}
}

func (a *AutoDailer) Stop() {
	a.Running = false
}

func (a *AutoDailer) Try(addr string) error {
	a.Running = true
	var con Con
	var err error
	var tempDelay time.Duration
	for a.Running {
		con, err = a.Dail(addr)
		log.D("NConRunner dail to server(%v) success", addr)
		if err == nil {
			a.a_lck.Lock()
			a.Addrs[con] = addr
			a.a_lck.Unlock()
			break
		}
		if a.OnDailFail != nil {
			a.OnDailFail(addr, err)
		}
		if tempDelay == 0 {
			tempDelay = 5 * time.Millisecond
		} else {
			tempDelay *= 2
		}
		if max := 8 * time.Second; tempDelay > max {
			tempDelay = max
		}
		log.D("NConRunner try dail to server(%v) err:%v,will retry after %v", addr, err.Error(), tempDelay)
		time.Sleep(tempDelay)
	}
	return err
}

type NConRunner struct {
	*NConPool
	C         Con
	ConH      ConHandler
	Connected bool
	Running   bool
	Retry     time.Duration
	Tick      time.Duration
	TickData  []byte
	//
	Addr    string
	NCF     NewConF
	BP      *pool.BytePool
	CmdH    CmdHandler
	ShowLog bool //setting the ShowLog to Con_
	TickLog bool //if show the tick log.
	wg      sync.WaitGroup
	//
	lastConnTime int64
	DailAddr     func(addr string) (net.Conn, error)
}

func (n *NConRunner) OnConn(c Con) bool {
	n.lastConnTime = util.Now()
	if n.ConH == nil {
		return true
	}
	return n.ConH.OnConn(c)
}

func (n *NConRunner) OnClose(c Con) {
	if n.ConH != nil {
		n.ConH.OnClose(c)
	}
	if n.Running {
		n.wg.Add(1)
		go n.Try()
	}
}
func (n *NConRunner) StartRunner() {
	n.wg.Add(1)
	go n.Try()
	//
	go n.StartTick()
	log.D("starting runner...")
}
func (n *NConRunner) StopRunner() {
	//n.wg.Add(1)
	n.Running = false
	if n.NConPool != nil {
		n.NConPool.Close()
		n.NConPool.Wait()
	}
	log.D("stopping runner...")
	n.wg.Wait()
}
func (n *NConRunner) StartTick() {
	if len(n.TickData) < 1 {
		return
	}
	n.wg.Add(1)
	go n.RunTick_()
}
func (n *NConRunner) write_tick() {
	c := n.C
	if c == nil {
		log.D("sending tick message err: the connection is nil")
		return
	}
	_, err := c.Writeb(n.TickData)
	if err != nil {
		log.W("send tck message err:%v", err)
		return
	}
	if n.TickLog {
		log.D("sending tick message to Push Server")
	}
}
func (n *NConRunner) RunTick_() {
	tk := time.Tick(n.Tick * time.Millisecond)
	n.Running = true
	log.I("starting tick(%vms) to server(%v)", int(n.Tick), n.Addr)
	for n.Running {
		select {
		case <-tk:
			n.write_tick()
		}
	}
	log.I("tick to server(%v) will stop", n.Addr)
	n.wg.Done()
}
func (n *NConRunner) Try() {
	n.Running = true
	if util.Now()-n.lastConnTime < 1000 {
		time.Sleep(5 * time.Second)
	}
	for n.Running {
		err := n.Dail()
		log.D("connect to server(%v) success", n.Addr)
		if err == nil {
			break
		}
		log.D("try connect to server(%v) err:%v,will retry after %v ms", n.Addr, err.Error(), int64(n.Retry))
		time.Sleep(n.Retry * time.Millisecond)
	}
	n.wg.Done()
	// log.D("connect try stopped")
}
func (n *NConRunner) Dail() error {
	n.Connected = false
	//
	nc := NewNConPool2(n.BP, NewCCH(n, n.CmdH))
	nc.NewCon = n.NCF
	nc.DailAddr = n.DailAddr
	cc, err := nc.Dail(n.Addr)
	if err != nil {
		return err
	}
	n.NConPool = nc
	n.C = cc
	n.Connected = true
	cc.(*Con_).ShowLog = n.ShowLog
	return nil
}

func NewNConRunnerN(bp *pool.BytePool, addr string, h CmdHandler, ncf NewConF) *NConRunner {
	return &NConRunner{
		Addr:     addr,
		NCF:      ncf,
		BP:       bp,
		CmdH:     h,
		Retry:    5000,
		Tick:     30000,
		TickData: []byte("Tick\n"),
		wg:       sync.WaitGroup{},
		DailAddr: func(addr string) (net.Conn, error) {
			return net.Dial("tcp", addr)
		},
	}
}
func NewNConRunner(bp *pool.BytePool, addr string, h CmdHandler) *NConRunner {
	return NewNConRunnerN(bp, addr, h, NewCon)
}

// func NewNConRunnerN(bp *pool.BytePool, addr string, h CmdHandler, ncf NewConF) *NConRunner {
// 	return &NConRunner{
// 		Addr:     addr,
// 		NCF:      ncf,
// 		BP:       bp,
// 		CmdH:     h,
// 		Retry:    5000,
// 		Tick:     30000,
// 		TickData: []byte("Tick\n"),
// 		wg:       sync.WaitGroup{},
// 	}
// }
// func NewNConRunner(bp *pool.BytePool, addr string, h CmdHandler) *NConRunner {
// 	return NewNConRunnerN(bp, addr, h, NewCon)
// }
