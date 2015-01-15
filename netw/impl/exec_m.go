package impl

import (
	"github.com/Centny/gwf/netw"
	"github.com/Centny/gwf/pool"
	"github.com/Centny/gwf/util"
	"net"
	"regexp"
)

type NAV_F func(rc *RCM_Con, name string, args interface{}) (interface{}, error)

type RCM_Con struct {
	*RC_Con
	NAV NAV_F
}

func NewRCM_Con(con *RC_Con, na NAV_F) *RCM_Con {
	return &RCM_Con{
		RC_Con: con,
		NAV:    na,
	}
}
func (r *RCM_Con) Exec(name string, args interface{}, dest interface{}) (interface{}, error) {
	vv, err := r.NAV(r, name, args)
	if err != nil {
		return nil, err
	}
	return r.RC_Con.Exec(vv, dest)
}

func ExecDail_m(p *pool.BytePool, addr string, v2b netw.V2Byte, b2v netw.Byte2V, na NAV_F) (*netw.NConPool, *RCM_Con, error) {
	tc := NewRC_C()
	return ExecDailN_m(p, addr, tc, tc, v2b, b2v, na)
}
func ExecDailN_m(p *pool.BytePool, addr string, h netw.CmdHandler, tc *RC_C, v2b netw.V2Byte, b2v netw.Byte2V, na NAV_F) (*netw.NConPool, *RCM_Con, error) {
	np, rc, err := ExecDailN(p, addr, h, tc, v2b, b2v)
	return np, NewRCM_Con(rc, na), err
}

type ND_F func() interface{}
type VNA_F func(rc *RCM_S, c netw.Cmd, v interface{}) (string, *util.Map, error)
type RC_M_FFUNC func(r *RCM_Cmd, vv interface{}) (bool, interface{}, error)
type RC_M_HFUNC func(r *RCM_Cmd, vv interface{}) (interface{}, error)

type RCM_Cmd struct {
	netw.Cmd
	*util.Map
}

// //the remote command server handler.
type RCM_S struct {
	ND       ND_F
	VNA      VNA_F
	filter_a []*regexp.Regexp
	filter_m map[*regexp.Regexp]RC_M_FFUNC
	routes_  map[string]RC_M_HFUNC
}

//new remote command server handler.
func NewRCM_S(nd ND_F, vna VNA_F) *RCM_S {
	return &RCM_S{
		ND:       nd,
		VNA:      vna,
		filter_a: []*regexp.Regexp{},
		filter_m: map[*regexp.Regexp]RC_M_FFUNC{},
		routes_:  map[string]RC_M_HFUNC{},
	}
}

func (r *RCM_S) OnCmd(c netw.Cmd) {
	defer c.Done()
	tv, err := c.V(r.ND())
	if err != nil {
		c.Err(1, "cmd convert to V err:%v", err.Error())
		return
	}
	fname, args, err := r.VNA(r, c, tv)
	if err != nil {
		c.Err(1, "find func name/args for V error:%v", err.Error())
		return
	}
	log_d("ExecM name(%v) args(%v)", fname, args)
	var con bool = false
	var vv interface{} = nil
	rcm := &RCM_Cmd{
		Cmd: c,
		Map: args,
	}
	for _, reg := range r.filter_a {
		if !reg.MatchString(fname) {
			continue
		}
		con, vv, err = r.filter_m[reg](rcm, vv)
		if err != nil {
			rcm.Err(1, "exec filter(%v) val(%v) errr:%v", reg.String(), vv, err.Error())
			return
		}
		if !con {
			r.writev(rcm, vv)
			return
		}
	}
	if h, ok := r.routes_[fname]; ok {
		val, err := h(rcm, vv)
		if err == nil {
			r.writev(rcm, val)
		} else {
			rcm.Err(1, "exec handler func(%v) error:%v", fname, err.Error())
		}
	} else {
		rcm.Err(1, "function not found by name(%v)", fname)
	}
}
func (r *RCM_S) writev(c *RCM_Cmd, val interface{}) {
	if _, err := c.Writev(val); err != nil {
		c.Err(1, "server sending return value err(%v)", err.Error())
	}
}
func (r *RCM_S) AddFFunc(reg string, ff RC_M_FFUNC) {
	reg_ := regexp.MustCompile(reg)
	r.filter_a = append(r.filter_a, reg_)
	r.filter_m[reg_] = ff
}
func (r *RCM_S) AddHFunc(name string, hf RC_M_HFUNC) {
	r.routes_[name] = hf
}

func NewExecListener_m(p *pool.BytePool, port string, h netw.ConHandler, nd ND_F, vna VNA_F) (*netw.Listener, *RCM_S) {
	return NewExecListenerN_m(p, port, h, V2B_Byte, B2V_Copy, nd, vna)
}

func NewExecListenerN_m(p *pool.BytePool, port string, h netw.ConHandler, v2b netw.V2Byte, b2v netw.Byte2V, nd ND_F, vna VNA_F) (*netw.Listener, *RCM_S) {
	rc := NewRCM_S(nd, vna)
	return NewExecListenerN_m_r(p, port, h, rc, v2b, b2v), rc
}

func NewExecListenerN_m_r(p *pool.BytePool, port string, h netw.ConHandler, rc *RCM_S, v2b netw.V2Byte, b2v netw.Byte2V) *netw.Listener {
	return netw.NewListenerN(p, port, netw.NewCCH(h, NewRC_S(rc)), func(cp netw.ConPool, p *pool.BytePool, con net.Conn) netw.Con {
		cc := netw.NewCon_(cp, p, con)
		cc.V2B_ = v2b
		cc.B2V_ = b2v
		return cc
	})
}
