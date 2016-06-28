package rctest

import (
	"fmt"
	"github.com/Centny/gwf/netw"
	"github.com/Centny/gwf/netw/impl"
	"github.com/Centny/gwf/netw/rc"
	"github.com/Centny/gwf/pool"
)

type RCTest struct {
	L    *rc.RC_Listener_m
	R    *rc.RC_Runner_m
	Addr string
}

func NewRCTest(p *pool.BytePool, port string, sh netw.CCHandler, ch netw.CCHandler,
	nd impl.ND_F, vna impl.VNA_F, v2b netw.V2Byte, b2v netw.Byte2V, na impl.NAV_F) *RCTest {
	// rcm_s := impl.NewRCM_S(nd, vna)
	rcm_c := impl.NewRCM_S(nd, vna)
	addr := fmt.Sprintf("127.0.0.1%v", port)
	rct := &RCTest{
		L:    rc.NewRC_Listener_m(p, port, sh, nd, vna, v2b, b2v, na),
		R:    rc.NewRC_Runner_m(p, addr, ch, rcm_c, v2b, b2v, na),
		Addr: addr,
	}
	err := rct.Run()
	if err != nil {
		panic(err.Error())
	}
	return rct
}

func NewRCTest_j(p *pool.BytePool, port string, sh netw.CCHandler, ch netw.CCHandler) *RCTest {
	return NewRCTest(p, port, sh, ch, impl.Json_ND, impl.Json_VNA, impl.Json_V2B, impl.Json_B2V, impl.Json_NAV)
}

func NewRCTest_j2(addr string) *RCTest {
	return NewRCTest_j(pool.BP, addr, netw.NewDoNotH(), netw.NewDoNotH())
}

func (r *RCTest) Run() error {
	err := r.L.Run()
	if err == nil {
		r.R.Start()
	}
	return err
}

func (r *RCTest) Close() {
	r.R.Stop()
	r.L.Close()
}

func (r *RCTest) ShowLog(v bool) {
	netw.ShowLog = v
	impl.ShowLog = v
}

func (r *RCTest) Runner() *rc.RC_Runner_m {
	return r.R
}

func (r *RCTest) Listener() *rc.RC_Listener_m {
	return r.L
}
