package impl

import (
	"encoding/json"
	"net"
	"reflect"

	"github.com/Centny/gwf/log"
	"github.com/Centny/gwf/netw"
	"github.com/Centny/gwf/pool"
	"github.com/Centny/gwf/util"
)

var ErrTooLarge = util.Err("data too large")

func Json_NAV(rc *RCM_Con, name string, args interface{}) (interface{}, error) {
	return util.Map{
		"name": name,
		"args": args,
	}, nil
}

func Json_V2B(v interface{}) (data []byte, err error) {
	switch v.(type) {
	case string:
		data = []byte(v.(string))
	case []byte:
		data = v.([]byte)
	default:
		data, err = json.Marshal(v)
	}
	if err == nil && len(data) > 60000 {
		log.E("Json_V2B write data(%v) too large by limit(60000)", len(data))
		err = ErrTooLarge
	}
	return
}

func Json_B2V(bys []byte, v interface{}) (data interface{}, err error) {
	switch v.(type) {
	case *string:
		val := v.(*string)
		*val = string(bys)
	case *[]byte:
		val := v.(*[]byte)
		*val = make([]byte, len(bys))
		copy(*val, bys)
	default:
		err = json.Unmarshal(bys, v)
	}
	data = v
	if err != nil {
		err = util.Err("Json_B2V to %v err(%v) by data:%v", reflect.TypeOf(v), err.Error(), string(bys))
	}
	return
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

func Json_NewCon(cp netw.ConPool, p *pool.BytePool, con net.Conn) (netw.Con, error) {
	cc := netw.NewCon_(cp, p, con)
	cc.V2B_ = Json_V2B
	cc.B2V_ = Json_B2V
	return cc, nil
}
