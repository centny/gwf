//Package rc provide the remote command on server and client.
//
//it base netw and netw/impl package.
//
//for example see rc_test.go under netw/rc package.
//
package rc

import (
	"fmt"
	"github.com/Centny/gwf/log"
	"github.com/Centny/gwf/netw"
	"github.com/Centny/gwf/netw/impl"
	"github.com/Centny/gwf/pool"
	"github.com/Centny/gwf/tutil"
	"github.com/Centny/gwf/util"
	"net"
	"sync"
	"sync/atomic"
)

const (
	CMD_S = 10 //server command channel
	MSG_S = 11 //server message channel
	CMD_C = 20 //client command channel
	MSG_C = 21 //client message channel
)

type RC_Login_h interface {
	OnLogin(rc *impl.RCM_Cmd, token string) (string, error)
}

//remote command callback distributed handler
type RC_Cmd_h struct {
	H    netw.ConHandler
	L    *RC_Listener_m
	CRC  map[string]*impl.RC_C    //remote client callback handler by connection id.
	MID  map[string]string        //mapping connection to custom id.
	MCS  map[string]netw.Con      //message connection for client by custom id.
	CCS  map[string]*impl.RCM_Con //command connection for client by custom id.
	rc_l sync.RWMutex
	//
	//
	Token map[string]int
	//
	cid int64
}

//new remote command callback distributed handler.
func NewRC_Cmd_h() *RC_Cmd_h {
	return &RC_Cmd_h{
		CRC:   map[string]*impl.RC_C{},
		MID:   map[string]string{},
		MCS:   map[string]netw.Con{},
		CCS:   map[string]*impl.RCM_Con{},
		Token: map[string]int{},
	}
}

//see netw.CmdHandler.
func (r *RC_Cmd_h) OnCmd(c netw.Cmd) int {
	if rcs, ok := r.CRC[r.MID[c.Id()]]; ok {
		return rcs.OnCmd(c)
	} else {
		log.W("remote client command call back handler not found by %v", c.Id())
		return -1
	}
}

//see netw.ConHandler
func (r *RC_Cmd_h) OnConn(c netw.Con) bool {
	return r.H.OnConn(c)
}

//see netw.ConHandler
func (r *RC_Cmd_h) OnClose(c netw.Con) {
	r.delc(c)
	r.H.OnClose(c)
}

//add connection to connection list by underlying connection.
func (r *RC_Cmd_h) AddC(cid string, con netw.Con) {
	r.rc_l.Lock()
	defer r.rc_l.Unlock()
	//create message connection.
	r.MCS[cid] = impl.NewOBDH_Con(MSG_C, con)
	//create command connection.
	tc := impl.NewRC_C()
	oc := impl.NewOBDH_Con(CMD_C, con)
	rcc := impl.NewRC_Con(oc, tc)
	rcm := impl.NewRCM_Con(rcc, r.L.Na)
	rcm.Start()
	r.CRC[cid] = tc
	r.CCS[cid] = rcm
	r.MID[con.Id()] = cid
}
func (r *RC_Cmd_h) Exist(cid string) bool {
	_, ok := r.CRC[cid]
	return ok
}
func (r *RC_Cmd_h) delc(c netw.Con) {
	r.rc_l.Lock()
	defer r.rc_l.Unlock()
	cid, ok := r.MID[c.Id()]
	if !ok {
		return
	}
	if rcm, ok := r.CCS[cid]; ok {
		rcm.Stop()
	}
	delete(r.MCS, cid)
	delete(r.CCS, cid)
	delete(r.CRC, cid)
	delete(r.MID, c.Id())
}

//find message connection by id.
func (r *RC_Cmd_h) MsgC(cid string) netw.Con {
	if con, ok := r.MCS[cid]; ok {
		return con
	} else {
		return nil
	}
}

//find command connection by id.
func (r *RC_Cmd_h) CmdC(cid string) *impl.RCM_Con {
	if con, ok := r.CCS[cid]; ok {
		return con
	} else {
		return nil
	}
}

//adding token
func (r *RC_Cmd_h) AddToken(ts map[string]int) {
	for k, v := range ts {
		r.Token[k] = v
	}
}
func (r *RC_Cmd_h) AddToken2(ts []string) {
	for idx, v := range ts {
		r.Token[v] = idx + 1
	}
}
func (r *RC_Cmd_h) AddToken3(token string, v int) {
	r.Token[token] = v
}
func (r *RC_Cmd_h) TokenVal(token string) int {
	return r.Token[token]
}
func (r *RC_Cmd_h) OnLogin(rc *impl.RCM_Cmd, token string) (string, error) {
	cid := atomic.AddInt64(&r.cid, 1)
	return fmt.Sprintf("N-%v", cid), nil
}

func AuthF(r *impl.RCM_Cmd) (bool, interface{}, error) {
	if r.Having("cid") {
		return true, nil, nil
	} else {
		return false, nil, util.Err("please login first")
	}
}

//remote command listener.
type RC_Listener_m struct {
	*netw.Listener //listener
	*impl.RCM_S    //remote command handler.
	//
	OH  *impl.OBDH  //OBDH by CMD_S/CMD_C/MSG_S/MSG_C
	Na  impl.NAV_F  //remote command function name.
	CH  *impl.ChanH //process chan.
	RCH *RC_Cmd_h   //remote client command call back handler.
	LCH RC_Login_h
}

//new remote command listener by common convert function
func NewRC_Listener_m(p *pool.BytePool, port string, h netw.CCHandler, rc *impl.RCM_S, v2b netw.V2Byte, b2v netw.Byte2V, na impl.NAV_F) *RC_Listener_m {
	obdh := impl.NewOBDH()
	obdh.AddH(CMD_S, impl.NewRC_S(rc))
	obdh.AddH(MSG_S, h)
	//
	rch := NewRC_Cmd_h()
	rch.H = h
	obdh.AddH(CMD_C, rch)
	//
	ch := impl.NewChanH(obdh)
	l := netw.NewListenerN2(p, port, netw.NewCCH(rch, ch), func(cp netw.ConPool, p *pool.BytePool, con net.Conn) netw.Con {
		cc := netw.NewCon_(cp, p, con)
		cc.V2B_ = v2b
		cc.B2V_ = b2v
		return cc
	})
	rcl := &RC_Listener_m{
		Listener: l,
		OH:       obdh,
		CH:       ch,
		RCH:      rch,
		Na:       na,
		RCM_S:    rc,
		LCH:      rch,
	}
	rch.L = rcl
	rcl.AddHFunc("login_", rcl.Login_)
	return rcl
}

//new remote command listener by json convert function
func NewRC_Listener_m_j(p *pool.BytePool, port string, h netw.CCHandler) *RC_Listener_m {
	rcm := impl.NewRCM_S_j()
	return NewRC_Listener_m(p, port, h, rcm, impl.Json_V2B, impl.Json_B2V, impl.Json_NAV)
}

//start listener
func (r *RC_Listener_m) Run() error {
	r.CH.Run(util.CPU())
	return r.Listener.Run()
}

//check if exit by client id
func (r *RC_Listener_m) Exist(cid string) bool {
	return r.RCH.Exist(cid)
}

//add connection to connection list by message command.
func (r *RC_Listener_m) AddC_c(cid string, c netw.Cmd) {
	r.RCH.AddC(cid, c.BaseCon().(*netw.Con_))
}

//add connection to connection list by remote command.
func (r *RC_Listener_m) AddC_rc(cid string, rc *impl.RCM_Cmd) {
	r.RCH.AddC(cid, rc.BaseCon().(*netw.Con_))
}

func (r *RC_Listener_m) AddToken(ts map[string]int) {
	r.RCH.AddToken(ts)
}
func (r *RC_Listener_m) AddToken2(ts []string) {
	r.RCH.AddToken2(ts)
}
func (r *RC_Listener_m) AddToken3(token string, v int) {
	r.RCH.AddToken3(token, v)
}

//find message connection by id.
func (r *RC_Listener_m) MsgC(cid string) netw.Con {
	return r.RCH.MsgC(cid)
}
func (r *RC_Listener_m) MsgCs() map[string]netw.Con {
	return r.RCH.MCS
}

//find command connection by id.
func (r *RC_Listener_m) CmdC(cid string) *impl.RCM_Con {
	return r.RCH.CmdC(cid)
}
func (r *RC_Listener_m) CmdCs() map[string]*impl.RCM_Con {
	return r.RCH.CCS
}

//find user id by connection id
func (r *RC_Listener_m) ConCid(c netw.Con) string {
	return r.RCH.MID[c.Id()]
}

func (r *RC_Listener_m) TokenVal(token string) int {
	return r.RCH.TokenVal(token)
}

func (r *RC_Listener_m) Login_(rc *impl.RCM_Cmd) (interface{}, error) {
	var token string
	err := rc.ValidF(`
		token,R|S,L:0;
		`, &token)
	if err != nil {
		return nil, err
	}
	otk := rc.Kvs().StrVal("token")
	log.D("RC_Listener_m(%v) login by token(%v),old(%v)", r.Name, token, otk)
	tval := r.TokenVal(token)
	if tval < 1 {
		log.W("RC_Listener_m(%v) login by token(%v) fail->token is not found", r.Name, token)
		return util.Map{"code": -2, "err": fmt.Sprintf("token(%v) is not found", token)}, nil
	}
	if tval == 1 && len(otk) > 0 {
		log.W("RC_Listener_m(%v) login by token(%v) fail->token is logined", r.Name, token)
		return util.Map{"code": -3, "err": fmt.Sprintf("token(%v) is logined", token)}, nil
	}
	tcid, err := r.LCH.OnLogin(rc, token)
	if err != nil {
		log.W("RC_Listener_m(%v) login by token(%v) fail->call OnLogin fail->%v", r.Name, token, err)
		return util.Map{"code": -4, "err": err.Error()}, nil
	}
	r.AddC_rc(tcid, rc)
	rc.Kvs().SetVal("token", token)
	rc.Kvs().SetVal("cid", tcid)
	log.D("RC_Listener_m(%v) login by token(%v),old(%v) success with cid(%v)", r.Name, token, otk, tcid)
	return util.Map{"code": 0}, nil
}

//set the show slow log
func (r *RC_Listener_m) SetShowSlow(v int64) {
	r.RCM_S.ShowSlow = v
}

//start monitor
func (r *RC_Listener_m) StartMonitor() {
	r.RCM_S.M = tutil.NewMonitor()
}

func (r *RC_Listener_m) State() (interface{}, error) {
	if r.M == nil {
		return nil, nil
	} else {
		return r.M.State()
	}
}

//remote command client runner.
type RC_Runner_m struct {
	*impl.RC_Runner_m
	*impl.RCM_S
	CC  netw.ConHandler //command message and connection event handler.
	TC  *impl.RC_C      //remote callback handler.
	CH  *impl.ChanH     //process chan
	OH  *impl.OBDH      //OBDH by CMD_S/CMD_C/MSG_S/MSG_C
	V2b netw.V2Byte     //common convert function
	B2v netw.Byte2V     //common convert function
	Na  impl.NAV_F      //remote command function name.
	MC  netw.Con        //message connection.
	RC  *impl.RCM_Con   //remote command connection.
	BC  *netw.Con_
}

//new remote command client runner by common convert function.
func NewRC_Runner_m(p *pool.BytePool, addr string, h netw.CCHandler, rc *impl.RCM_S, v2b netw.V2Byte, b2v netw.Byte2V, na impl.NAV_F) *RC_Runner_m {
	tc := impl.NewRC_C()
	obdh := impl.NewOBDH()
	obdh.AddH(CMD_S, tc)
	obdh.AddH(MSG_C, h)
	obdh.AddH(CMD_C, impl.NewRC_S(rc))
	ch := impl.NewChanH(obdh)
	ch.Run(util.CPU())
	// cc := netw.NewCCH(h, ch)
	run := &RC_Runner_m{
		TC:    tc,
		CH:    ch,
		OH:    obdh,
		CC:    h,
		V2b:   v2b,
		B2v:   b2v,
		Na:    na,
		RCM_S: rc,
	}
	run.RC_Runner_m = impl.NewRC_Runner_m(addr, p, run.Dail)
	return run
}

//new remote command client runner by json convert function.
func NewRC_Runner_m_j(p *pool.BytePool, addr string, h netw.CCHandler) *RC_Runner_m {
	rcm := impl.NewRCM_S_j()
	return NewRC_Runner_m(p, addr, h, rcm, impl.Json_V2B, impl.Json_B2V, impl.Json_NAV)
}

//dail to server and create remote command connection.
func (r *RC_Runner_m) Dail(p *pool.BytePool, addr string, h netw.ConHandler) (*netw.NConPool, *impl.RCM_Con, error) {
	log.I("RC_Runner_m(%v) start connect to addr(%v)", r.Name, addr)
	cch := netw.NewCCH(netw.NewQueueConH(h, r.CC), r.CH)
	np := netw.NewNConPool2(p, cch)
	np.NewCon = func(cp netw.ConPool, p *pool.BytePool, con net.Conn) netw.Con {
		r.BC = netw.NewCon_(cp, p, con)
		r.BC.V2B_, r.BC.B2V_ = r.V2b, r.B2v
		rcc := impl.NewRC_Con(impl.NewOBDH_Con(CMD_S, r.BC), r.TC)
		r.RCM_Con = impl.NewRCM_Con(rcc, r.Na)
		r.MC = impl.NewOBDH_Con(MSG_S, r.BC)
		return r.BC
	}
	_, err := np.Dail(addr)
	if err == nil {
		return np, r.RCM_Con, nil
	} else {
		return nil, nil, err
	}
}
func (r *RC_Runner_m) Writeb(bys ...[]byte) (int, error) {
	err := r.Valid()
	if err == nil {
		return r.MC.Writeb(bys...)
	} else {
		return 0, err
	}
}
func (r *RC_Runner_m) Writev(val interface{}) (int, error) {
	err := r.Valid()
	if err == nil {
		return r.MC.Writev(val)
	} else {
		return 0, err
	}
}
func (r *RC_Runner_m) Writev2(bys []byte, val interface{}) (int, error) {
	err := r.Valid()
	if err == nil {
		return r.MC.Writev2(bys, val)
	} else {
		return 0, err
	}
}

func (r *RC_Runner_m) Login_(token string) error {
	log.I("RC_Runner_m(%v) login by token(%v)", r.Name, token)
	res, err := r.VExec_m("login_", util.Map{
		"token": token,
	})
	if err != nil {
		return err
	}
	if res.IntVal("code") == 0 {
		log.I("RC_Runner_m(%v) login success by token(%v)", r.Name, token)
		return nil
	} else {
		log.I("RC_Runner_m(%v) login fail by token(%v)->%v", r.Name, token, util.S2Json(res))
		return util.Err("login error->%v", res.StrVal("err"))
	}
}

//set the show slow log
func (r *RC_Runner_m) SetShowSlow(v int64) {
	r.RC_Runner_m.ShowSlow = v
	r.RCM_S.ShowSlow = v
}

//start monitor
func (r *RC_Runner_m) StartMonitor() {
	r.RC_Runner_m.M = tutil.NewMonitor()
	r.RCM_S.M = tutil.NewMonitor()
}

//the runner state
func (r *RC_Runner_m) State() (interface{}, error) {
	var res = util.Map{}
	if r.RC_Runner_m.M != nil {
		val, _ := r.RC_Runner_m.M.State()
		res["exec"] = val
	}
	if r.RCM_S.M != nil {
		val, _ := r.RC_Runner_m.M.State()
		res["hand"] = val
	}
	return res, nil
}

type AutoLoginH struct {
	Runner *RC_Runner_m
	Token  string
}

func (a *AutoLoginH) OnConn(c netw.Con) bool {
	go a.Runner.Login_(a.Token)
	return true
}

func (a *AutoLoginH) OnClose(c netw.Con) {
}
