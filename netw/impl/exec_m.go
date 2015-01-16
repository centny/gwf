package impl

import (
	"github.com/Centny/gwf/netw"
	"github.com/Centny/gwf/util"
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

/*


*/

type ND_F func() interface{}
type VNA_F func(rc *RCM_S, c netw.Cmd, v interface{}) (string, *util.Map, error)
type RC_M_FFUNC func(r *RCM_Cmd) (bool, interface{}, error)
type RC_M_HFUNC func(r *RCM_Cmd) (interface{}, error)

type RCM_Cmd struct {
	netw.Cmd
	*util.Map
	Vv interface{}
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

func (r *RCM_S) OnCmd(c netw.Cmd) int {
	defer c.Done()
	tv, err := c.V(r.ND())
	if err != nil {
		c.Err(1, "cmd convert to V err:%v", err.Error())
		return -1
	}
	fname, args, err := r.VNA(r, c, tv)
	if err != nil {
		c.Err(1, "find func name/args for V error:%v", err.Error())
		return -1
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
		con, vv, err = r.filter_m[reg](rcm)
		if err != nil {
			rcm.Err(1, "exec filter(%v) val(%v) errr:%v", reg.String(), vv, err.Error())
			return -1
		}
		rcm.Vv = vv
		if !con {
			r.writev(rcm, vv)
			return 0
		}
	}
	if h, ok := r.routes_[fname]; ok {
		val, err := h(rcm)
		if err == nil {
			r.writev(rcm, val)
			return 0
		} else {
			rcm.Err(1, "exec handler func(%v) error:%v", fname, err.Error())
			return -1
		}
	} else {
		rcm.Err(1, "function not found by name(%v)", fname)
		return -1
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
