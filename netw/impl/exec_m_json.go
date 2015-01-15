package impl

import (
	"encoding/json"
	"github.com/Centny/gwf/netw"
	"github.com/Centny/gwf/pool"
	"github.com/Centny/gwf/util"
)

func json_NAV_(rc *RCM_Con, name string, args interface{}) (interface{}, error) {
	return util.Map{
		"name": name,
		"args": args,
	}, nil
}

func json_V2B(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func json_B2V(bys []byte, v interface{}) (interface{}, error) {
	return v, json.Unmarshal(bys, v)
}

func ExecDail_m_j(p *pool.BytePool, addr string) (*netw.NConPool, *RCM_Con, error) {
	tc := NewRC_C()
	return ExecDailN_m(p, addr, tc, tc, json_V2B, json_B2V, json_NAV_)
}

/*

*/
func json_ND() interface{} {
	return &util.Map{}
}

func json_VNA(rc *RCM_S, c netw.Cmd, v interface{}) (string, *util.Map, error) {
	vv := v.(*util.Map)
	name := vv.StrVal("name")
	if len(name) < 1 {
		return "", nil, util.Err(`json_VNA, func name not found,using {"name":"func","args":{}}`)
	}
	args := vv.MapVal("args")
	return name, &args, nil
}

func NewRCM_S_j() *RCM_S {
	return NewRCM_S(json_ND, json_VNA)
}
func NewExecListener_m_j(p *pool.BytePool, port string, h netw.ConHandler) (*netw.Listener, *RCM_S) {
	rc := NewRCM_S_j()
	return NewExecListenerN_m_r(p, port, h, rc, json_V2B, json_B2V), rc
}
