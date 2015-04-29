package impl

import (
	"encoding/json"
	"github.com/Centny/gwf/netw"
	"github.com/Centny/gwf/pool"
	"github.com/Centny/gwf/util"
	"net"
)

func Json_NAV(rc *RCM_Con, name string, args interface{}) (interface{}, error) {
	return util.Map{
		"name": name,
		"args": args,
	}, nil
}

func Json_V2B(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func Json_B2V(bys []byte, v interface{}) (interface{}, error) {
	err := json.Unmarshal(bys, v)
	if err == nil {
		return v, nil
	} else {
		return v, util.Err("Json_B2V err(%v) by data:%v", err.Error(), string(bys))
	}
}

/*

*/
func Json_ND() interface{} {
	return &util.Map{}
}

func Json_VNA(rc *RCM_S, c netw.Cmd, v interface{}) (string, *util.Map, error) {
	vv := v.(*util.Map)
	name := vv.StrVal("name")
	if len(name) < 1 {
		return "", nil, util.Err(`json_VNA, func name not found,using {"name":"func","args":{}}`)
	}
	args := vv.MapVal("args")
	return name, &args, nil
}

func Json_NewCon(cp netw.ConPool, p *pool.BytePool, con net.Conn) netw.Con {
	cc := netw.NewCon_(cp, p, con)
	cc.V2B_ = Json_V2B
	cc.B2V_ = Json_B2V
	return cc
}
