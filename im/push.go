package im

import (
	"github.com/Centny/gwf/netw"
	"github.com/Centny/gwf/netw/impl"
	"github.com/Centny/gwf/pool"
	"github.com/Centny/gwf/util"
)

type PushSrv struct {
	*netw.Listener
}

func (p *PushSrv) Notify(mid string) int {
	return p.Writev(&util.Map{
		"MID": mid,
	})
}

func NewPushSrv(p *pool.BytePool, port string, n string, h netw.CCHandler) *PushSrv {
	return NewPushSrvN(p, port, n, h, impl.Json_NewCon)
}

func NewPushSrvN(p *pool.BytePool, port string, n string, h netw.CCHandler, ncf netw.NewConF) *PushSrv {
	return &PushSrv{
		Listener: netw.NewListenerN(p, port, n, h, ncf),
	}
}
