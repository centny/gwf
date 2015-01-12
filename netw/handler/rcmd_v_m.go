package handler

import (
	"github.com/Centny/gwf/netw"
	"github.com/Centny/gwf/util"
	"regexp"
)

type RC_V_M_FNA func(rc *RC_V_M_C, fname string, args interface{}) (interface{}, error)

type RC_V_M_C struct {
	*RC_V_C
	FNA RC_V_M_FNA
}

func NewRC_V_M_C(fna RC_V_M_FNA, v2b V2Byte, b2v Byte2V) *RC_V_M_C {
	return &RC_V_M_C{
		FNA:    fna,
		RC_V_C: NewRC_V_C(v2b, b2v),
	}
}
func NewRC_Json_M_C(fna RC_V_M_FNA) *RC_V_M_C {
	return &RC_V_M_C{
		FNA:    fna,
		RC_V_C: NewRC_Json_C(),
	}
}
func (r *RC_V_M_C) Exec(fname string, args interface{}, dest interface{}) error {
	v, err := r.FNA(r, fname, args)
	if err == nil {
		return r.RC_V_C.Exec(v, dest)
	} else {
		return err
	}
}

type RC_V_M_FFUNC func(r *RC_V_M_S, rc *RC_V_Cmd, args *util.Map, vv interface{}) (bool, interface{}, error)
type RC_V_M_HFUNC func(r *RC_V_M_S, rc *RC_V_Cmd, args *util.Map) (interface{}, error)

//the extended command handler.
type RC_V_M_H interface {
	netw.ConHandler
	FNAME(rc *RC_V_Cmd) (string, error)
	FARGS(rc *RC_V_Cmd) (*util.Map, error)
}

//the remote command server handler.
type RC_V_M_S struct {
	H        RC_V_M_H
	filter_a []*regexp.Regexp
	filter_m map[*regexp.Regexp]RC_V_M_FFUNC
	routes_  map[string]RC_V_M_HFUNC
}

//new remote command server handler.
func NewRC_V_M_S(h RC_V_M_H) *RC_V_M_S {
	return &RC_V_M_S{
		H:        h,
		filter_a: []*regexp.Regexp{},
		filter_m: map[*regexp.Regexp]RC_V_M_FFUNC{},
		routes_:  map[string]RC_V_M_HFUNC{},
	}
}

func (r *RC_V_M_S) OnConn(c *netw.Con) bool {
	return r.H.OnConn(c)
}

func (r *RC_V_M_S) OnClose(c *netw.Con) {
	r.H.OnClose(c)
}

func (r *RC_V_M_S) OnCmd(rc *RC_V_Cmd) (interface{}, error) {
	fname, err := r.H.FNAME(rc)
	if err != nil {
		return nil, err
	}
	args, err := r.H.FARGS(rc)
	if err != nil {
		return nil, err
	}
	var con bool = false
	var vv interface{} = nil
	for _, reg := range r.filter_a {
		if !reg.MatchString(fname) {
			continue
		}
		con, vv, err = r.filter_m[reg](r, rc, args, vv)
		if err != nil || !con {
			return vv, err
		}
	}
	if h, ok := r.routes_[fname]; ok {
		return h(r, rc, args)
	} else {
		return nil, util.Err("function not found by %v", fname)
	}
}
func (r *RC_V_M_S) AddFFunc(reg string, ff RC_V_M_FFUNC) {
	reg_ := regexp.MustCompile(reg)
	r.filter_a = append(r.filter_a, reg_)
	r.filter_m[reg_] = ff
}
func (r *RC_V_M_S) AddHFunc(name string, hf RC_V_M_HFUNC) {
	r.routes_[name] = hf
}
