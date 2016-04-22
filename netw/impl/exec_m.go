package impl

import (
	"github.com/Centny/gwf/log"
	"github.com/Centny/gwf/netw"
	"github.com/Centny/gwf/util"
	"regexp"
)

type NAV_F func(rc *RCM_Con, name string, args interface{}) (interface{}, error)

type RCM_CRes struct {
	Code int         `json:"code"`
	Res  interface{} `json:"res"`
}
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
func (r *RCM_Con) Exec2(name string, args interface{}) ([]byte, error) {
	vv, err := r.NAV(r, name, args)
	if err != nil {
		return nil, err
	}
	return r.RC_Con.Exec2(vv)
}
func (r *RCM_Con) ExecRes(name string, args interface{}) (*RCM_CRes, error) {
	var crs RCM_CRes
	_, err := r.Exec(name, args, &crs)
	return &crs, err
}
func (r *RCM_Con) Exec_m(name string, args interface{}) (util.Map, error) {
	var res util.Map
	_, err := r.Exec(name, args, &res)
	return res, err
}
func (r *RCM_Con) Exec_s(name string, args interface{}) (string, error) {
	var bys, err = r.Exec2(name, args)
	return string(bys), err
}

/*


 */

type ND_F func() interface{}
type VNA_F func(rc *RCM_S, c netw.Cmd, v interface{}) (string, *util.Map, error)
type RC_M_FH interface {
	Exec(r *RCM_Cmd) (bool, interface{}, error)
}
type RC_M_FFUNC func(r *RCM_Cmd) (bool, interface{}, error)

func (rf RC_M_FFUNC) Exec(r *RCM_Cmd) (bool, interface{}, error) {
	return rf(r)
}

type RC_M_HH interface {
	Exec(r *RCM_Cmd) (interface{}, error)
}
type RC_M_HFUNC func(r *RCM_Cmd) (interface{}, error)
type RC_M_HFUNCv func(v util.Validable) (interface{}, error)

func (rh RC_M_HFUNC) Exec(r *RCM_Cmd) (interface{}, error) {
	return rh(r)
}
func (rh RC_M_HFUNCv) Exec(r *RCM_Cmd) (interface{}, error) {
	return rh(r)
}

type RCM_Cmd struct {
	Name string
	netw.Cmd
	*util.Map
	Vv interface{}
}

func (r *RCM_Cmd) CRes(code int, v interface{}) (interface{}, error) {
	return &RCM_CRes{
		Code: code,
		Res:  v,
	}, nil
}

// //the remote command server handler.
type RCM_S struct {
	ND       ND_F
	VNA      VNA_F
	filter_a []*regexp.Regexp
	filter_m map[*regexp.Regexp]RC_M_FH
	routes_  map[string]RC_M_HH
}

//new remote command server handler.
func NewRCM_S(nd ND_F, vna VNA_F) *RCM_S {
	return &RCM_S{
		ND:       nd,
		VNA:      vna,
		filter_a: []*regexp.Regexp{},
		filter_m: map[*regexp.Regexp]RC_M_FH{},
		routes_:  map[string]RC_M_HH{},
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
		Name: fname,
		Cmd:  c,
		Map:  args,
	}
	for _, reg := range r.filter_a {
		if !reg.MatchString(fname) {
			continue
		}
		con, vv, err = r.exec_f(r.filter_m[reg], rcm)
		if err != nil {
			log_d("exec filter(%v) val(%v) errr:%v", reg.String(), vv, err.Error())
			rcm.Err(1, "%v", err.Error())
			return -1
		}
		rcm.Vv = vv
		if !con {
			r.writev(rcm, vv)
			return 0
		}
	}
	if h, ok := r.routes_[fname]; ok {
		val, err := r.exec_h(h, rcm)
		if err == nil {
			r.writev(rcm, val)
			return 0
		} else {
			log_d("exec handler func(%v) error:%v", fname, err.Error())
			rcm.Err(1, "%v", err.Error())
			return -1
		}
	} else {
		rcm.Err(1, "function not found by name(%v)", fname)
		return -1
	}
}
func (r *RCM_S) exec_f(f RC_M_FH, rcm *RCM_Cmd) (con bool, val interface{}, err error) {
	defer func() {
		var terr = recover()
		if terr != nil {
			log.E("RCM_S exec filter(%p) is panic(%v) with call stack->\n%v", f, terr, util.CallStatck())
			err = util.Err("%v", terr)
		}
	}()
	con, val, err = f.Exec(rcm)
	return
}
func (r *RCM_S) exec_h(h RC_M_HH, rcm *RCM_Cmd) (val interface{}, err error) {
	defer func() {
		var terr = recover()
		if terr != nil {
			log.E("RCM_S exec handler(%p) is panic(%v) with call stack->\n%v", h, terr, util.CallStatck())
			err = util.Err("%v", terr)
		}
	}()
	val, err = h.Exec(rcm)
	return
}
func (r *RCM_S) writev(c *RCM_Cmd, val interface{}) {
	if _, err := c.Writev(val); err != nil {
		c.Err(1, "server sending return value err(%v)", err.Error())
	}
}
func (r *RCM_S) AddFFunc(reg string, ff RC_M_FFUNC) {
	r.AddFH(reg, RC_M_FH(ff))
}
func (r *RCM_S) AddFH(reg string, fh RC_M_FH) {
	reg_ := regexp.MustCompile(reg)
	r.filter_a = append(r.filter_a, reg_)
	r.filter_m[reg_] = fh
}
func (r *RCM_S) AddHFunc(name string, hf RC_M_HFUNC) {
	r.AddHH(name, RC_M_HH(hf))
}
func (r *RCM_S) AddHFuncv(name string, hf RC_M_HFUNCv) {
	r.AddHH(name, RC_M_HH(hf))
}
func (r *RCM_S) AddHH(name string, hh RC_M_HH) {
	r.routes_[name] = hh
}
