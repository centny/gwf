package im

import (
	"fmt"
	"github.com/Centny/gwf/log"
	"github.com/Centny/gwf/netw"
	"github.com/Centny/gwf/netw/impl"
	"github.com/Centny/gwf/util"
	"sync"
)

const (
	MK_NDC_NLI = 0
	MK_NDC_ULI = 10
	MK_NDC_ULO = 11
)

//
type NodeRh struct {
	NIM *NIM_Rh
}

func (n *NodeRh) OnCmd(c netw.Cmd) int {
	defer c.Done()
	// log.D("DIM_Rh recieve data:%v", string(c.Data()))
	var mc Msg
	_, err := c.V(&mc)
	if err != nil {
		log.E("convert value(%v) to Msg error:%v", string(c.Data()), err.Error())
		return -1
	}
	if len(mc.S) < 1 {
		log.E("receive not sender Msg(%v) from(%v)", string(c.Data()), c.RemoteAddr().String())
		return -1
	}
	mc.Cmd = c
	mc.Ms = map[string]string{}
	return n.NIM.OnMsg(&mc)
}

type NodeV struct {
	V util.Map `json:"v"`
	B string   `json:"b"`
}
type NodeCmds struct {
	Db   DbH
	SS   Sender
	DS   map[string]netw.Con
	ds_l sync.RWMutex
}

func (n *NodeCmds) OnConn(c netw.Con) bool {
	return true
}
func (n *NodeCmds) OnClose(c netw.Con) {
	n.ds_l.Lock()
	defer n.ds_l.Unlock()
	delete(n.DS, c.Id())
}
func (n *NodeCmds) Find(id string) netw.Con {
	return n.DS[id]
}
func (n *NodeCmds) writev_c(c netw.Cmd, na NodeV, code int, res interface{}) int {
	na.V = util.Map{
		"res":  res,
		"code": code,
	}
	c.Writev(na)
	return 0
}
func (n *NodeCmds) H(obdh *impl.OBDH) {
	obdh.AddF(MK_NDC_NLI, n.NLI)
	obdh.AddF(MK_NDC_ULI, n.ULI)
	obdh.AddF(MK_NDC_ULO, n.ULO)
}
func (n *NodeCmds) NLI(c netw.Cmd) int {
	var na NodeV
	_, err := c.B2V()(c.Data(), &na)
	if err != nil {
		fmt.Println(c.Data())
		log.E("Node Cmd data(%v) to value err:%v", string(c.Data()), err.Error())
		return n.writev_c(c, na, 1, err.Error())
	}

	var token string
	err = na.V.ValidF(`
		token,R|S,L:0,token is empty;
		`, &token)
	if err != nil {
		return n.writev_c(c, na, 1, err.Error())
	}
	srv, err := n.Db.FindSrv(token)
	if err != nil {
		return n.writev_c(c, na, 1, err.Error())
	}
	if srv.Sid != n.SS.Id() {
		errs := fmt.Sprintf("login fail,invalid token(%v) for current server(%v,%v)", token, n.SS.Id(), srv.Token)
		log.W("Node LI login(%v)", errs)
		return n.writev_c(c, na, 1, errs)
	}
	n.ds_l.Lock()
	defer n.ds_l.Unlock()
	n.DS[c.Id()] = c.BaseCon()
	// c.SetId(sid)
	c.SetWait(true)
	log_d("Node server login success from(%v)", c.RemoteAddr().String())
	return n.writev_c(c, na, 0, "OK")
}
func (n *NodeCmds) ULI(c netw.Cmd) int {
	var na NodeV
	_, err := c.B2V()(c.Data(), &na)
	if err != nil {
		log.E("Node Cmd data(%v) to value err:%v", string(c.Data()), err.Error())
		return n.writev_c(c, na, 1, err.Error())
	}
	rv, err := n.Db.OnUsrLogin(c, &na.V)
	if err != nil {
		return n.writev_c(c, na, 1, err.Error())
	}
	con := &Con{
		Sid: n.SS.Id(),
		Cid: c.Id(),
		R:   rv,
		S:   "N",
		T:   CT_WS,
	}
	err = n.Db.AddCon(con)
	if err != nil {
		return n.writev_c(c, na, 1, err.Error())
	} else {
		return n.writev_c(c, na, 0, con)
	}
}
func (n *NodeCmds) ULO(c netw.Cmd) int {
	var na NodeV
	_, err := c.B2V()(c.Data(), &na)
	if err != nil {
		log.E("Node Cmd data(%v) to value err:%v", string(c.Data()), err.Error())
		return n.writev_c(c, na, 1, err.Error())
	}
	var r string
	err = na.V.ValidF(`
		r,R|S,L:0,user R is empty;
		`, &r)
	if err != nil {
		return n.writev_c(c, na, 1, err.Error())
	}
	err = n.Db.DelCon(n.SS.Id(), c.Id(), r, CT_WS)
	if err != nil {
		return n.writev_c(c, na, 1, err.Error())
	} else {
		return n.writev_c(c, na, 0, "OK")
	}
}
