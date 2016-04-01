package main

import (
	"fmt"
	"github.com/Centny/gwf/log"
	"github.com/Centny/gwf/netw"
	"github.com/Centny/gwf/netw/impl"
	"github.com/Centny/gwf/netw/rc"
	"github.com/Centny/gwf/pool"
	"github.com/Centny/gwf/util"
)

func CallTestClient(value string) {
	var clients = L.CmdCs()
	fmt.Println("calling->", len(clients))
	for name, client := range clients {
		res, err := client.Exec_m("test_client", util.Map{
			"value": value,
		})
		fmt.Printf("call client(%v) result is data(%v),err(%v)\n", name, util.S2Json(res), err)
	}
}

func TestSrv(rc *impl.RCM_Cmd) (interface{}, error) {
	var name string
	L.CmdC(cid)
	var err = rc.ValidF(`
		name,R|S,L:0;
		`, &name)
	if err == nil {
		log.D("TestSrv receive name(%v)", name)
		return util.Map{
			"code": 0,
			"data": "OK",
		}, nil
	} else {
		return util.Map{
			"code": -1,
			"err":  err,
		}, nil
	}
}

func Hand_S(rc *rc.RC_Listener_m) {
	rc.AddHFunc("test_srv", TestSrv)
}

var L *rc.RC_Listener_m

func StartRCSrv(addr string) error {
	L = rc.NewRC_Listener_m_j(pool.BP, addr, netw.NewDoNotH())
	L.AddToken2([]string{"abc"})
	Hand_S(L)
	return L.Run()
}
