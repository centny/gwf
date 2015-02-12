package im

import (
	"fmt"
	"github.com/Centny/gwf/im/pb"
	"github.com/Centny/gwf/log"
	"github.com/Centny/gwf/netw"
	"github.com/Centny/gwf/netw/impl"
	"github.com/Centny/gwf/util"
	"strings"
	"sync/atomic"
)

const (
	MK_NRC_LI = 10
	MK_NRC_LO = 20
)

//
type NIM_Rh struct {
	Db  DbH
	SS  Sender
	DS  Sender
	DC  uint64
	idc int64
}

func (n *NIM_Rh) OnConn(c netw.Con) bool {
	return n.Db.OnConn(c)
}
func (n *NIM_Rh) OnClose(c netw.Con) {
	n.Db.OnCloseCon(c, n.SS.Id(), c.Id(), CT_TCP)
	n.Db.OnClose(c)
}

func (n *NIM_Rh) OnCmd(c netw.Cmd) int {
	defer c.Done()
	// log_d("NIM_Rh receive data:%v", string(c.Data()))
	sid, tn := n.Db.FUsrR(c), util.Now()
	if len(sid) < 1 {
		log.W("receive message for not login connect->%v", c.Data())
		return -1
	}
	var mc Msg
	_, err := c.V(&mc.ImMsg)
	if err != nil {
		log.E("convert values(%v) to IM msg error:%v", c.Data(), err.Error())
		return -1
	}
	mc.S = &sid
	mc.Cmd = c
	mc.Ms = map[string]string{}
	mc.Time = &tn
	return n.OnMsg(&mc)
}
func (n *NIM_Rh) DoRobot(mc *Msg) int {
	if len(mc.R) < 1 {
		log.E("empty R(%v) from:%v", mc.R, mc.RemoteAddr().String())
		return -1
	}
	ss := mc.R[0]
	if !strings.HasPrefix(ss, "S-Robot") {
		return 0
	}
	mi := fmt.Sprintf("RMI-%v", atomic.AddInt64(&n.idc, 1))
	mc.I = &mi
	mc.D = mc.S
	mc.S = &ss
	mc.R = []string{*mc.D}
	err := n.SS.Send(mc.Id(), &mc.ImMsg)
	if err == nil {
		return 1
	}
	return 0
}
func (n *NIM_Rh) OnMsg(mc *Msg) int {
	if dr := n.DoRobot(mc); dr != 0 {
		return dr
	}
	gr, ur, err := n.Db.Sift(mc.R)
	if err != nil {
		log.E("sift R(%v) err:%v", mc.R, err.Error())
		return -1
	}
	if len(gr) > 0 {
		gur, err := n.Db.ListUsrR(gr)
		if err != nil {
			log.E("list user R for group(%v) err:%v", gr, err.Error())
			return -1
		}
		ur = append(ur, gur...)
	}
	if len(ur) < 1 {
		log.E("receive empty R message(%v)", mc)
		return -1
	}
	log_d("receive message(%v) to RS(%v) in S(%v)", mc, ur, n.SS.Id())
	//
	cons, err := n.Db.ListCon(ur)
	if err != nil {
		log.E("list Con by R(%v) err:%v", ur, err.Error())
		return -1
	}
	log_d("found %v online user for RS(%v) in S(%v)", len(cons), ur, n.SS.Id())
	mid := n.Db.NewMid()
	mc.I = &mid
	c_sid := n.SS.Id()             //current server id.
	sr_ed := map[string]byte{}     //already exec
	dr_rc := map[string][]*pb.RC{} //
	for _, con := range cons {     //do online user
		sr_ed[con.R] = 1
		if con.Sid == c_sid { //in current server
			mc.D = &con.R                       //setting current receive user R.
			err = n.SS.Send(con.Cid, &mc.ImMsg) //send message to client.
			if err == nil {
				atomic.AddUint64(&n.DC, 1)
				mc.Ms[con.R] = MS_DONE //mark done
			} else {
				log.E("sending message(%v) to R(%v) err:%v", mc.ImMsg, con.R, err.Error())
				mc.Ms[con.R] = MS_ERR + err.Error() //mark send error.
			}
		} else { //in other distribution server
			mc.Ms[con.R] = MS_PENDING //mark to pending.
			tr, tc := con.R, con.Cid
			if _, ok := dr_rc[con.Sid]; ok {
				dr_rc[con.Sid] = append(dr_rc[con.Sid],
					&pb.RC{
						R: &tr,
						C: &tc,
					})
			} else {
				dr_rc[con.Sid] = []*pb.RC{
					&pb.RC{
						R: &tr,
						C: &tc,
					},
				}
			}
		}
	}
	// log_d("sr_ed---->%v in S(%v)", sr_ed, c_sid)
	for _, r := range ur { //do offline user
		if _, ok := sr_ed[r]; ok {
			continue
		}
		mc.Ms[r] = MS_PENDING
	}
	if len(ur) > len(mc.Ms) {
		log.W("duplicate R(%v) found for message(%v)", ur, mc)
	}
	err = n.Db.Store(mc) //store mesage.
	if err != nil {
		log.E("store message(%v) err:%v", mc, err.Error())
		return -1
	}
	if n.DS == nil {
		return 0
	}
	// log_d("sending in S(%v) DR_RC->%v", c_sid, dr_rc)
	for dr, rc := range dr_rc { //if having distribution message.
		dmc := &pb.DsMsg{
			M:  &mc.ImMsg,
			Rc: rc,
		}
		err = n.DS.Send(dr, dmc)
		if err == nil { //if not err,the other distribution server will makr result.
			continue
		} else {
			log.E("sending message(%v) to distribution server(%v) err:%v", mc.ImMsg, dr, err.Error())
		}
	}
	return 0
}
func (n *NIM_Rh) H(obdh *impl.OBDH) {
	obdh.AddF(MK_NRC_LI, n.LI)
	obdh.AddF(MK_NRC_LO, n.LO)
}

// func (n *NIM_Rh) Exec(r *impl.RCM_Cmd) (interface{}, error) {
// 	log_d("call action(%v)", r.Name)
// 	switch r.Name {
// 	case "LI":
// 		return n.LI(r)
// 	case "LO":
// 		return n.LO(r)
// 	}
// 	return nil, util.Err("action not found by name(%v)", r.Name)
// }
func (n *NIM_Rh) writev_c(c netw.Cmd, res interface{}) int {
	c.Writev(util.Map{
		"res":  res,
		"code": 0,
	})
	return 0
}
func (n *NIM_Rh) writev_ce(c netw.Cmd, err string) int {
	c.Writev(util.Map{
		"err":  err,
		"code": 1,
	})
	return 0
}
func (n *NIM_Rh) LI(r netw.Cmd) int {
	var args util.Map
	_, err := r.V(&args)
	if err != nil {
		log.W("login V fail:%v", err.Error())
		return n.writev_ce(r, err.Error())
	}
	rv, ct, err := n.Db.OnLogin(r, &args)
	if err != nil {
		log.W("login OnLogin fail:%v", err.Error())
		return n.writev_ce(r, err.Error())
	}
	con := &Con{
		Sid: n.SS.Id(),
		Cid: r.Id(),
		R:   rv,
		S:   "N",
		T:   CT_TCP,
		C:   ct,
	}
	err = n.Db.AddCon(con)
	if err != nil {
		log.W("login AddCon fail:%v", err.Error())
		return n.writev_ce(r, err.Error())
	}
	r.SetWait(true)
	// r.Kvs().SetVal("R", rv)
	// con.Sid = ""
	res := n.writev_c(r, con)
	go SendUnread(n.SS, n.Db, r, rv, ct)
	return res
}
func (n *NIM_Rh) LO(r netw.Cmd) int {
	var args util.Map
	_, err := r.V(&args)
	if err != nil {
		return n.writev_ce(r, err.Error())
	}
	rv, ct, w, err := n.Db.OnLogout(r, &args)
	if err != nil {
		log.W("LO OnLogout fail:%v", err.Error())
		return n.writev_ce(r, err.Error())
	}
	if !w {
		r.SetWait(false)
	}
	con, err := n.Db.DelCon(n.SS.Id(), r.Id(), rv, CT_TCP, ct)
	if err != nil {
		log.W("LO DelCon fail:%v", err.Error())
		return n.writev_ce(r, err.Error())
	}
	return n.writev_c(r, con)
}

// func (n *NIM_Rh) onlo(con netw.Con) error {
// 	err := n.Db.DelCon(n.SS.Id(), con.Id(), con.Kvs().StrVal("R"), CT_TCP)
// 	con.SetWait(false)
// 	return err
// }
