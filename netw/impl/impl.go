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
	"github.com/Centny/gwf/tutil"
	"github.com/Centny/gwf/util"
	"net"
	"strings"
	"sync/atomic"
)

var ShowLog bool = false

func log_d(f string, args ...interface{}) {
	if ShowLog {
		log.D_(1, f, args...)
	}
}

/*


 */
// func NewNConRunner_j(bp *pool.BytePool, addr string, h netw.CmdHandler) *netw.NConRunner {
// 	return netw.NewNConRunnerN(bp, addr, h, Json_NewCon)
// }

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

type F_DAIL func(p *pool.BytePool, addr string, h netw.ConHandler) (netw.ConPool, *RCM_Con, error)

//
type RC_Runner_m struct {
	*RCM_Con
	Name      string
	Addr      string
	CC        netw.ConHandler
	BP        *pool.BytePool
	L         *netw.NConPool
	Connected int32
	wc        int32
	wait_     chan byte
	//
	ShowSlow int64
	M        *tutil.Monitor
	//
	Multi bool
	//
	// w_lck chan int
	//
	Dailer *netw.AutoDailer
	TC     *RC_C
	RC     *RC_Con
	Uuid   string
	//
	NAV NAV_F
	V2B netw.V2Byte
	B2V netw.Byte2V
}

func NewRC_Runner_m_base() *RC_Runner_m {
	return &RC_Runner_m{
		wait_: make(chan byte, 1000),
		// w_lck: make(chan int),
	}
}

func NewRC_Runner_m(addr string, bp *pool.BytePool, na NAV_F, v2b netw.V2Byte, b2v netw.Byte2V) *RC_Runner_m {
	var runner = &RC_Runner_m{
		Addr:      addr,
		BP:        bp,
		Connected: 0,
		wait_:     make(chan byte, 1000),
		// w_lck:     make(chan int),
		Uuid: strings.ToUpper(util.UUID()),
		NAV:  na,
		V2B:  v2b,
		B2V:  b2v,
	}
	runner.TC = NewRC_C()
	runner.Dailer = netw.NewAutoDailer()
	runner.RC = NewRC_Con(nil, runner.TC)
	runner.RCM_Con = NewRCM_Con(runner.RC, na)
	runner.L = netw.NewNConPool(bp, netw.NewCCH(netw.NewQueueConH(runner, runner.Dailer), runner.TC), "RC-")
	runner.L.NewCon = runner.NewCon
	runner.Dailer.Dail = runner.L.Dail
	runner.L.DailAddr = runner.DailAddr
	return runner
}

func (r *RC_Runner_m) NewCon(cp netw.ConPool, p *pool.BytePool, con net.Conn) netw.Con {
	cc := netw.NewCon_(cp, p, con)
	cc.V2B_, cc.B2V_ = r.V2B, r.B2V
	r.RC.Con = cc
	return r.RC
}
func (r *RC_Runner_m) Start() {
	if r.Multi {
		r.Dailer.DailAll(strings.Split(r.Addr, ","))
	} else {
		r.Dailer.DailAll([]string{r.Addr})
	}
}
func (r *RC_Runner_m) DailAddr(addr string) (net.Conn, error) {
	con, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}
	return con, nil
}

func (r *RC_Runner_m) Stop() {
	if r.Dailer != nil {
		r.Dailer.Stop()
	}
	if r.L != nil {
		r.L.Close()
	}
	r.Timeout()
	log.D("RC Runner is stopping")
}
func (r *RC_Runner_m) Wait() {
	log.D("RC Runner wait stop")
	// <-r.w_lck
}
func (r *RC_Runner_m) OnConn(c netw.Con) bool {
	r.RC.Start()
	atomic.AddInt32(&r.Connected, 1)
	tlen := r.wc
	var i int32
	for i = 0; i < tlen; i++ {
		r.wait_ <- byte(0)
	}
	atomic.AddInt32(&r.wc, -tlen)
	if r.CC != nil {
		return r.CC.OnConn(c)
	}
	return true
}
func (r *RC_Runner_m) OnClose(c netw.Con) {
	atomic.AddInt32(&r.Connected, -1)
	if r.CC != nil {
		r.CC.OnClose(c)
	}
}

func (r *RC_Runner_m) Valid() error {
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
func (r *RC_Runner_m) VExec(name string, args interface{}, dest interface{}) (interface{}, error) {
	var mid = ""
	var beg = util.Now()
	if r.M != nil {
		mid = r.M.Start(name)
	}
	defer func() {
		if r.M != nil {
			r.M.Done(mid)
		}
		used := util.Now() - beg
		if r.ShowSlow > 0 && used > r.ShowSlow {
			log.W("RC_Runner_m(%v) slow exec(%v) found by args->%v", r.Name, name, util.S2Json(args))
		}
	}()
	err := r.Valid()
	if err == nil {
		return r.Exec(name, args, dest)
	} else {
		return nil, err
	}
}
func (r *RC_Runner_m) VExecRes(name string, args interface{}) (*RCM_CRes, error) {
	var mid = ""
	var beg = util.Now()
	if r.M != nil {
		mid = r.M.Start(name)
	}
	defer func() {
		if r.M != nil {
			r.M.Done(mid)
		}
		used := util.Now() - beg
		if r.ShowSlow > 0 && used > r.ShowSlow {
			log.W("RC_Runner_m(%v) slow exec(%v) found by args->%v", r.Name, name, util.S2Json(args))
		}
	}()
	err := r.Valid()
	if err == nil {
		return r.ExecRes(name, args)
	} else {
		return nil, err
	}
}
func (r *RC_Runner_m) VExec_m(name string, args interface{}) (util.Map, error) {
	var res util.Map
	_, err := r.VExec(name, args, &res)
	return res, err
}
func (r *RC_Runner_m) VExec_s(name string, args interface{}) (string, error) {
	var mid = ""
	var beg = util.Now()
	if r.M != nil {
		mid = r.M.Start(name)
	}
	defer func() {
		if r.M != nil {
			r.M.Done(mid)
		}
		used := util.Now() - beg
		if r.ShowSlow > 0 && used > r.ShowSlow {
			log.W("RC_Runner_m(%v) slow exec(%v) found by args->%v", r.Name, name, util.S2Json(args))
		}
	}()
	err := r.Valid()
	if err == nil {
		return r.Exec_s(name, args)
	} else {
		return "", err
	}
}
func (r *RC_Runner_m) Timeout() {
	var i int32
	tlen := r.wc
	log.D("sending timeout to %v waiting", tlen)
	for i = 0; i < tlen; i++ {
		r.wait_ <- byte(1)
	}
	atomic.AddInt32(&r.wc, -tlen)
}
func (r *RC_Runner_m) Waitingc() int {
	return int(r.wc)
}

type RC_Runner_m_j struct {
	*RC_Runner_m
}

func NewRC_Runner_m_j(addr string, bp *pool.BytePool) *RC_Runner_m_j {
	return &RC_Runner_m_j{
		RC_Runner_m: NewRC_Runner_m(addr, bp, Json_NAV, Json_V2B, Json_B2V),
	}
}
