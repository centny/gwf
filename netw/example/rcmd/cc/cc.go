package cc

import (
	"github.com/Centny/gwf/netw"
	"github.com/Centny/gwf/netw/example/rcmd/common"
	"github.com/Centny/gwf/netw/impl"
	"github.com/Centny/gwf/pool"
	"github.com/Centny/gwf/util"
)

var lc *netw.NConPool = nil
var rcm *impl.RCM_Con = nil

func RunC() {
	p := pool.NewBytePool(8, 1024) //memory pool.
	l, con, err := impl.ExecDail_m_j(p, "127.0.0.1:8797")
	if err != nil {
		panic(err.Error())
	}
	con.Start() //start handing the server reponse.if not start, exec will hang up.
	lc, rcm = l, con
}
func Stop() {
	if rcm != nil {
		rcm.Stop()
	}
	if lc != nil {
		lc.Close()
	}
}

func List() ([]common.Val, error) {
	if rcm == nil {
		return nil, util.Err("RCM not inital,call RunC first")
	}
	var vs []common.Val
	_, err := rcm.Exec("list", map[string]interface{}{
		"tv": "abc",
	}, &vs)
	return vs, err
}
