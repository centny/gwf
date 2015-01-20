//Package handler provider multi base handler for netw connection.
//
//ChanH: the chan handler provide the feature of distributing command.
//
//RCH_*: the remove command handler provide the feature of remote command call.
package impl

import (
	"github.com/Centny/gwf/log"
	"github.com/Centny/gwf/netw"
	"github.com/Centny/gwf/pool"
	"net"
)

var ShowLog bool = false

func log_d(f string, args ...interface{}) {
	if ShowLog {
		log.D_(1, f, args...)
	}
}

/*


*/
//
func ExecDail(p *pool.BytePool, addr string) (*netw.NConPool, *RC_Con, error) {
	return ExecDail2(p, addr, V2B_Byte, B2V_Copy)
}
func ExecDail2(p *pool.BytePool, addr string, v2b netw.V2Byte, b2v netw.Byte2V) (*netw.NConPool, *RC_Con, error) {
	tc := NewRC_C()
	return ExecDailN(p, addr, tc, tc, v2b, b2v)
}
func ExecDailN(p *pool.BytePool, addr string, h netw.CmdHandler, tc *RC_C, v2b netw.V2Byte, b2v netw.Byte2V) (*netw.NConPool, *RC_Con, error) {
	cch := netw.NewCCH(NewRC_C_H(), h)
	np := netw.NewNConPool(p, cch)
	np.NewCon = func(cp netw.ConPool, p *pool.BytePool, con net.Conn) netw.Con {
		cc := netw.NewCon_(cp, p, con)
		cc.V2B_, cc.B2V_ = v2b, b2v
		rcc := NewRC_Con(cc, tc)
		// cch.Con = rcc
		return rcc
	}
	con, err := np.Dail(addr)
	if err == nil {
		return np, con.(*RC_Con), err
	} else {
		return nil, nil, err
	}
}

func NewExecListener(p *pool.BytePool, port string, h netw.CCHandler) *netw.Listener {
	return NewExecListenerN(p, port, h, V2B_Byte, B2V_Copy)
}
func NewExecListenerN(p *pool.BytePool, port string, h netw.CCHandler, v2b netw.V2Byte, b2v netw.Byte2V) *netw.Listener {
	return netw.NewListenerN(p, port, netw.NewCCH(h, NewRC_S(h)), func(cp netw.ConPool, p *pool.BytePool, con net.Conn) netw.Con {
		cc := netw.NewCon_(cp, p, con)
		cc.V2B_ = v2b
		cc.B2V_ = b2v
		return cc
	})
}

/*


*/

func ExecDail_m(p *pool.BytePool, addr string, v2b netw.V2Byte, b2v netw.Byte2V, na NAV_F) (*netw.NConPool, *RCM_Con, error) {
	tc := NewRC_C()
	return ExecDailN_m(p, addr, tc, tc, v2b, b2v, na)
}
func ExecDailN_m(p *pool.BytePool, addr string, h netw.CmdHandler, tc *RC_C, v2b netw.V2Byte, b2v netw.Byte2V, na NAV_F) (*netw.NConPool, *RCM_Con, error) {
	np, rc, err := ExecDailN(p, addr, h, tc, v2b, b2v)
	return np, NewRCM_Con(rc, na), err
}
func NewExecListener_m(p *pool.BytePool, port string, h netw.ConHandler, nd ND_F, vna VNA_F) (*netw.Listener, *RCM_S) {
	return NewExecListenerN_m(p, port, h, V2B_Byte, B2V_Copy, nd, vna)
}

func NewExecListenerN_m(p *pool.BytePool, port string, h netw.ConHandler, v2b netw.V2Byte, b2v netw.Byte2V, nd ND_F, vna VNA_F) (*netw.Listener, *RCM_S) {
	rc := NewRCM_S(nd, vna)
	return NewExecListenerN_m_r(p, port, h, rc, v2b, b2v), rc
}

func NewExecListenerN_m_r(p *pool.BytePool, port string, h netw.ConHandler, rc *RCM_S, v2b netw.V2Byte, b2v netw.Byte2V) *netw.Listener {
	return netw.NewListenerN(p, port, netw.NewCCH(h, NewRC_S(rc)), func(cp netw.ConPool, p *pool.BytePool, con net.Conn) netw.Con {
		cc := netw.NewCon_(cp, p, con)
		cc.V2B_ = v2b
		cc.B2V_ = b2v
		return cc
	})
}
func NewChanExecListenerN_m_r(p *pool.BytePool, port string, h netw.ConHandler, rc *RCM_S, v2b netw.V2Byte, b2v netw.Byte2V) (*netw.Listener, *ChanH) {
	cc := NewChanH(NewRC_S(rc))
	return netw.NewListenerN(p, port, netw.NewCCH(h, cc), func(cp netw.ConPool, p *pool.BytePool, con net.Conn) netw.Con {
		cc := netw.NewCon_(cp, p, con)
		cc.V2B_ = v2b
		cc.B2V_ = b2v
		return cc
	}), cc
}

/*


*/
func ExecDail_m_j(p *pool.BytePool, addr string) (*netw.NConPool, *RCM_Con, error) {
	tc := NewRC_C()
	return ExecDailN_m(p, addr, tc, tc, Json_V2B, Json_B2V, Json_NAV_)
}

func NewRCM_S_j() *RCM_S {
	return NewRCM_S(Json_ND, Json_VNA)
}
func NewExecListener_m_j(p *pool.BytePool, port string, h netw.ConHandler) (*netw.Listener, *RCM_S) {
	rc := NewRCM_S_j()
	return NewExecListenerN_m_r(p, port, h, rc, Json_V2B, Json_B2V), rc
}
func NewChanExecListener_m_j(p *pool.BytePool, port string, h netw.ConHandler) (*netw.Listener, *ChanH, *RCM_S) {
	rc := NewRCM_S_j()
	l, cc := NewChanExecListenerN_m_r(p, port, h, rc, Json_V2B, Json_B2V)
	return l, cc, rc
}
