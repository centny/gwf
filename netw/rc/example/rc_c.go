package main

import (
	"github.com/Centny/gwf/log"
	"github.com/Centny/gwf/netw"
	"github.com/Centny/gwf/netw/impl"
	"github.com/Centny/gwf/netw/rc"
	"github.com/Centny/gwf/pool"
	"github.com/Centny/gwf/util"
)

func CallTestSrv(name string) (util.Map, error) {
	return Runner.VExec_m("test_srv", util.Map{
		"name": name,
	})
}

func TestClient(rc *impl.RCM_Cmd) (interface{}, error) {
	var value string
	var err = rc.ValidF(`
		value,R|S,L:0;
		`, &value)
	if err == nil {
		log.D("TestSrv receive value(%v)", value)
		return util.Map{
			"code": 0,
			"data": value,
		}, nil
	} else {
		log.D("TestSrv receive error(%v)", err)
		return util.Map{
			"code": -1,
			"err":  err.Error(),
		}, nil
	}
}

func Hand_C(rc *rc.RC_Runner_m) {
	rc.AddHFunc("test_client", TestClient)
}

var Runner *rc.RC_Runner_m

func StartRCClient(addr string) {
	var auto = &rc.AutoLoginH{}
	auto.Token = "abc"
	Runner = rc.NewRC_Runner_m_j(pool.BP, addr, netw.NewCCH(auto, netw.NewDoNotH()))
	auto.Runner = Runner
	Hand_C(Runner)
	Runner.Start()
}
