//Package rc provide the remote command on server and client.
//
//it base netw and netw/impl package.
//
//for example see rc_test.go under netw/rc package.
//
package rc

import (
	"github.com/Centny/gwf/log"
	"github.com/Centny/gwf/netw"
	"github.com/Centny/gwf/netw/impl"
	"github.com/Centny/gwf/pool"
	"github.com/Centny/gwf/util"
	"net"
	"sync"
)

const (
	CMD_S = 10 //server command channel
	MSG_S = 11 //server message channel
	CMD_C = 20 //client command channel
	MSG_C = 21 //client message channel
)

//remote command callback distributed handler
type RC_Cmd_h struct {
	H    netw.ConHandler
	L    *RC_Listener_m
	CRC  map[string]*impl.RC_C    //remote client callback handler by connection id.
	MID  map[string]string        //mapping connection to custom id.
	MCS  map[string]netw.Con      //message connection for client by custom id.
	CCS  map[string]*impl.RCM_Con //command connection for client by custom id.
	rc_l sync.RWMutex
}

//new remote command callback distributed handler.
func NewRC_Cmd_h() *RC_Cmd_h {
	return &RC_Cmd_h{
		CRC: map[string]*impl.RC_C{},
		MID: map[string]string{},
		MCS: map[string]netw.Con{},
		CCS: map[string]*impl.RCM_Con{},
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

//remote command listener.
type RC_Listener_m struct {
	*netw.Listener //listener
	*impl.RCM_S    //remote command handler.
	//
	OH  *impl.OBDH  //OBDH by CMD_S/CMD_C/MSG_S/MSG_C
	Na  impl.NAV_F  //remote command function name.
	CH  *impl.ChanH //process chan.
	RCH *RC_Cmd_h   //remote client command call back handler.
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
	}
	rch.L = rcl
	return rcl
}

//new remote command listener by json convert function
func NewRC_Listener_m_j(p *pool.BytePool, port string, h netw.CCHandler) *RC_Listener_m {
	rcm := impl.NewRCM_S_j()
	return NewRC_Listener_m(p, port, h, rcm, impl.Json_V2B, impl.Json_B2V, impl.Json_NAV)
}
func (r *RC_Listener_m) Run() error {
	r.CH.Run(util.CPU())
	return r.Listener.Run()
}

//add connection to connection list by message command.
func (r *RC_Listener_m) AddC_c(cid string, c netw.Cmd) {
	r.RCH.AddC(cid, c.BaseCon().(*netw.Con_))
}

//add connection to connection list by remote command.
func (r *RC_Listener_m) AddC_rc(cid string, rc *impl.RCM_Cmd) {
	r.RCH.AddC(cid, rc.BaseCon().(*netw.Con_))
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

//remote command client runner.
type RC_Runner_m struct {
	*impl.RC_Runner_m
	*impl.RCM_S
	CC  netw.CCHandler //command message and connection event handler.
	TC  *impl.RC_C     //remote callback handler.
	CH  *impl.ChanH    //process chan
	OH  *impl.OBDH     //OBDH by CMD_S/CMD_C/MSG_S/MSG_C
	V2b netw.V2Byte    //common convert function
	B2v netw.Byte2V    //common convert function
	Na  impl.NAV_F     //remote command function name.
	MC  netw.Con       //message connection.
	RC  *impl.RCM_Con  //remote command connection.
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
	cc := netw.NewCCH(h, ch)
	run := &RC_Runner_m{
		TC:    tc,
		CH:    ch,
		OH:    obdh,
		CC:    cc,
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
	np := netw.NewNConPool2(p, netw.NewCCH(h, r.OH))
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
