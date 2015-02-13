//Package handler provider multi base handler for netw connection.
//
//ChanH: the chan handler provide the feature of distributing command.
//
//RCH_*: the remove command handler provide the feature of remote command call.
package impl

import (
	"github.com/Centny/gwf/log"
	"github.com/Centny/gwf/netw"
	"github.com/Centny/gwf/pool"
	"github.com/Centny/gwf/util"
	"net"
	"sync/atomic"
	"time"
)

var ShowLog bool = false

func log_d(f string, args ...interface{}) {
	if ShowLog {
		log.D_(1, f, args...)
	}
}

/*


*/
func NewNConRunner_j(bp *pool.BytePool, addr string, h netw.CmdHandler) *netw.NConRunner {
	return netw.NewNConRunnerN(bp, addr, h, Json_NewCon)
}

//
func ExecDail(p *pool.BytePool, addr string, h netw.ConHandler) (*netw.NConPool, *RC_Con, error) {
	return ExecDail2(p, addr, h, V2B_Byte, B2V_Copy)
}
func ExecDail2(p *pool.BytePool, addr string, h netw.ConHandler, v2b netw.V2Byte, b2v netw.Byte2V) (*netw.NConPool, *RC_Con, error) {
	tc := NewRC_C()
	return ExecDailN(p, addr, netw.NewCCH(h, tc), tc, v2b, b2v)
}
func ExecDailN(p *pool.BytePool, addr string, h netw.CCHandler, tc *RC_C, v2b netw.V2Byte, b2v netw.Byte2V) (*netw.NConPool, *RC_Con, error) {
	// cch := netw.NewCCH(NewRC_C_H(), h)
	np := netw.NewNConPool2(p, h)
	np.NewCon = func(cp netw.ConPool, p *pool.BytePool, con net.Conn) netw.Con {
		cc := netw.NewCon_(cp, p, con)
		cc.V2B_, cc.B2V_ = v2b, b2v
		rcc := NewRC_Con(cc, tc)
		// cch.Con = rcc
		return rcc
	}
	con, err := np.Dail(addr)
	if err == nil {
		return np, con.(*RC_Con), err
	} else {
		return nil, nil, err
	}
}

func NewExecListener(p *pool.BytePool, port string, h netw.CCHandler) *netw.Listener {
	return NewExecListenerN(p, port, h, V2B_Byte, B2V_Copy)
}
func NewExecListenerN(p *pool.BytePool, port string, h netw.CCHandler, v2b netw.V2Byte, b2v netw.Byte2V) *netw.Listener {
	return netw.NewListenerN2(p, port, netw.NewCCH(h, NewRC_S(h)), func(cp netw.ConPool, p *pool.BytePool, con net.Conn) netw.Con {
		cc := netw.NewCon_(cp, p, con)
		cc.V2B_ = v2b
		cc.B2V_ = b2v
		return cc
	})
}

/*


*/

func ExecDail_m(p *pool.BytePool, addr string, h netw.ConHandler, v2b netw.V2Byte, b2v netw.Byte2V, na NAV_F) (*netw.NConPool, *RCM_Con, error) {
	tc := NewRC_C()
	return ExecDailN_m(p, addr, netw.NewCCH(h, tc), tc, v2b, b2v, na)
}
func ExecDailN_m(p *pool.BytePool, addr string, h netw.CCHandler, tc *RC_C, v2b netw.V2Byte, b2v netw.Byte2V, na NAV_F) (*netw.NConPool, *RCM_Con, error) {
	np, rc, err := ExecDailN(p, addr, h, tc, v2b, b2v)
	return np, NewRCM_Con(rc, na), err
}
func NewExecListener_m(p *pool.BytePool, port string, h netw.ConHandler, nd ND_F, vna VNA_F) (*netw.Listener, *RCM_S) {
	return NewExecListenerN_m(p, port, h, V2B_Byte, B2V_Copy, nd, vna)
}

func NewExecListenerN_m(p *pool.BytePool, port string, h netw.ConHandler, v2b netw.V2Byte, b2v netw.Byte2V, nd ND_F, vna VNA_F) (*netw.Listener, *RCM_S) {
	rc := NewRCM_S(nd, vna)
	return NewExecListenerN_m_r(p, port, h, rc, v2b, b2v), rc
}

func NewExecListenerN_m_r(p *pool.BytePool, port string, h netw.ConHandler, rc *RCM_S, v2b netw.V2Byte, b2v netw.Byte2V) *netw.Listener {
	return netw.NewListenerN2(p, port, netw.NewCCH(h, NewRC_S(rc)), func(cp netw.ConPool, p *pool.BytePool, con net.Conn) netw.Con {
		cc := netw.NewCon_(cp, p, con)
		cc.V2B_ = v2b
		cc.B2V_ = b2v
		return cc
	})
}
func NewChanExecListenerN_m_r(p *pool.BytePool, port string, h netw.ConHandler, rc *RCM_S, v2b netw.V2Byte, b2v netw.Byte2V) (*netw.Listener, *ChanH) {
	cc := NewChanH(NewRC_S(rc))
	return netw.NewListenerN2(p, port, netw.NewCCH(h, cc), func(cp netw.ConPool, p *pool.BytePool, con net.Conn) netw.Con {
		cc := netw.NewCon_(cp, p, con)
		cc.V2B_ = v2b
		cc.B2V_ = b2v
		return cc
	}), cc
}

/*


*/
func ExecDail_m_j(p *pool.BytePool, addr string, h netw.ConHandler) (*netw.NConPool, *RCM_Con, error) {
	tc := NewRC_C()
	return ExecDailN_m(p, addr, netw.NewCCH(h, tc), tc, Json_V2B, Json_B2V, Json_NAV)
}

func NewRCM_S_j() *RCM_S {
	return NewRCM_S(Json_ND, Json_VNA)
}
func NewExecListener_m_j(p *pool.BytePool, port string, h netw.ConHandler) (*netw.Listener, *RCM_S) {
	rc := NewRCM_S_j()
	return NewExecListenerN_m_r(p, port, h, rc, Json_V2B, Json_B2V), rc
}
func NewChanExecListener_m_j(p *pool.BytePool, port string, h netw.ConHandler) (*netw.Listener, *ChanH, *RCM_S) {
	rc := NewRCM_S_j()
	l, cc := NewChanExecListenerN_m_r(p, port, h, rc, Json_V2B, Json_B2V)
	return l, cc, rc
}

//
type RC_Runner_m_j struct {
	*RCM_Con
	Addr      string
	BP        *pool.BytePool
	L         *netw.NConPool
	R         bool
	Delay     int64
	Connected int32
	wc        int32
	wait_     chan byte
}

func NewRC_Runner_m_j(addr string, bp *pool.BytePool) *RC_Runner_m_j {
	return &RC_Runner_m_j{
		Addr:      addr,
		BP:        bp,
		Delay:     3000,
		Connected: 0,
		wait_:     make(chan byte, 1000),
	}
}
func (r *RC_Runner_m_j) Start() {
	r.R = true
	go r.Try()
}
func (r *RC_Runner_m_j) Start_() error {
	r.R = true
	return r.Run()
}
func (r *RC_Runner_m_j) Stop() {
	r.R = false
	if r.L != nil {
		r.L.Close()
	}
}
func (r *RC_Runner_m_j) Run() error {
	atomic.StoreInt32(&r.Connected, 0)
	var err error
	r.L, r.RCM_Con, err = ExecDail_m_j(r.BP, r.Addr, r)
	if err != nil {
		return err
	}
	r.RCM_Con.Start()
	atomic.StoreInt32(&r.Connected, 1)
	var i int32
	tlen := r.wc
	for i = 0; i < tlen; i++ {
		r.wait_ <- byte(0)
	}
	atomic.AddInt32(&r.wc, -tlen)
	return nil
}
func (r *RC_Runner_m_j) OnConn(c netw.Con) bool {
	c.SetWait(true)
	log.D("RC Runner connect to %v success", r.Addr)
	return true
}
func (r *RC_Runner_m_j) Try() {
	atomic.StoreInt32(&r.Connected, 0)
	var last, now int64 = util.Now(), 0
	var t int = 0
	for r.R {
		t++
		err := r.Run()
		if err == nil {
			break
		}
		now = util.Now()
		if now-last < r.Delay {
			log.E("RC connect server err:%v, will retry(%v) after %v ms", err.Error(), t, r.Delay)
			time.Sleep(time.Duration(r.Delay) * time.Millisecond)
		}
		last = now
	}
}
func (r *RC_Runner_m_j) OnClose(c netw.Con) {
	atomic.StoreInt32(&r.Connected, 0)
	r.RCM_Con.Stop()
	if r.R {
		log.W("RC connection  is closed, Runner will retry connect to %v", r.Addr)
		go r.Try()
	}
}

func (r *RC_Runner_m_j) Valid() error {
	if atomic.LoadInt32(&r.Connected) > 0 {
		return nil
	}
	atomic.AddInt32(&r.wc, 1)
	if v := <-r.wait_; v > 0 {
		return util.Err("time out")
	} else {
		return nil
	}
}
