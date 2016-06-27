package netw

import (
	"github.com/Centny/gwf/log"
	"github.com/Centny/gwf/pool"
	"github.com/Centny/gwf/util"
	"net"
	"sync"
	"time"
)

//the client connection pool.
type NConPool struct {
	*LConPool //base connection pool.
	DailAddr  func(addr string) (net.Conn, error)
}

//new client connection pool.
func NewNConPool(p *pool.BytePool, h CCHandler, n string) *NConPool {
	return &NConPool{
		LConPool: NewLConPoolV(p, h, n, NewConH),
		DailAddr: func(addr string) (net.Conn, error) {
			return net.Dial("tcp", addr)
		},
	}
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
	Addrs   map[Con]string
	a_lck   sync.RWMutex
	Running bool
	Dail    func(addr string) (Con, error)
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
