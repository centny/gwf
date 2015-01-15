package srv

import (
	"github.com/Centny/gwf/netw"
	"github.com/Centny/gwf/netw/example/rcmd/common"
	"github.com/Centny/gwf/netw/impl"
	"github.com/Centny/gwf/pool"
	"runtime"
)

func List(rc *impl.RCM_Cmd) (interface{}, error) {
	var tv string
	err := rc.ValidF(`
		tv,R|S,L:0,
		`, &tv)
	if err != nil {
		return nil, err
	}
	return []common.Val{
		common.Val{
			V: tv + "0",
		},
		common.Val{
			V: tv + "1",
		},
		common.Val{
			V: tv + "2",
		},
	}, nil
}

func RunSrv() {
	// netw.ShowLog = true
	// impl.ShowLog = true
	p := pool.NewBytePool(8, 1024) //memory pool.
	l, cc, cms := impl.NewChanExecListener_m_j(p, ":8797", netw.NewCWH(true))
	cms.AddHFunc("list", List)
	cc.Run(runtime.NumCPU() - 1) //start the chan distribution, if not start, sub handler will not receive message
	err := l.Run()               //run the listen server
	if err != nil {
		panic(err.Error())
	}
	l.Wait()
}
