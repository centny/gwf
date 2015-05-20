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
	MK_NRC_HB = 9
	MK_NRC_LI = 10
	MK_NRC_LO = 20
	MK_NRC_UR = 30
)
const (
	MS_SEQ = "^->"
)

type NMR_Rh struct {
	Db DbH
}

func (n *NMR_Rh) OnCmd(r netw.Cmd) int {
	defer r.Done()
	if r.Closed() {
		return -1
	}
	tr := n.Db.FUsrR(r)
	if len(tr) < 1 {
		return -1
	}
	var args util.Map
	_, err := r.V(&args)
	if err != nil {
		log.W("MR V fail:%v", err.Error())
		return -1
	}
	var i, a string
	err = args.ValidF(`
		i,R|S,L:0;
		a,R|S,L:0;
		`, &i, &a)
	if err != nil {
		log.W("MR args(%v) fail:%v", args, err.Error())
		return -1
	}
	err = n.Db.MarkRecv(tr, a, i)
	if err == nil {
		return 0
	} else {
		log.W("MarkRecv by i(%v) fail:%v", args.StrVal("i"), err.Error())
		return -1
	}
}

//
type NIM_Rh struct {
	Db DbH
	SS Sender
	DS Sender
	// DC       uint64
	idc      int64
	Running  bool
	PushChan chan string
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
	if c.Closed() {
		return -1
	}
	// log_d("NIM_Rh receive data:%v", string(c.Data()))
	sid, tn := n.Db.FUsrR(c), util.Now()
	if len(sid) < 1 {
		log.W("receive message for not login connect->%v", c.Data())
		c.Close()
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
	mc.Ms = map[string][]*MSS{}
	mc.Time = &tn
	return n.OnMsg(&mc)
}
func (n *NIM_Rh) OnMsg(mc *Msg) int {
	if len(mc.R) < 1 {
		log.E("receive empty R from %v", mc.RemoteAddr())
		mc.Close()
		return -1
	}
	for _, tr := range mc.R {
		if len(strings.Trim(tr, "\t ")) < 1 {
			log.E("receive empty string in R from %v", mc.RemoteAddr())
			mc.Close()
			return -1
		}
	}
	if dr := n.DoRobot(mc); dr != 0 {
		return dr
	}
	gr, ur, err := n.Db.Sift(mc.R)
	if err != nil {
		log.E("sift R(%v) err:%v", mc.R, err.Error())
		return -1
	}
	var gur map[string][]string = map[string][]string{}
	if len(gr) > 0 {
		gur, err = n.Db.ListUsrR(gr)
		if err != nil {
			log.E("list user R by gr(%v) err:%v", gr, err.Error())
			return -1
		}
	}
	if len(ur) > 0 {
		gur[mc.GetS()] = ur
	}
	if len(gur) < 1 {
		log.E("receive empty R message(%v)", mc)
		return -1
	}
	log_d("receive message(%v) to R(%v) in S(%v)", mc, gur, n.SS.Id())
	mid := n.Db.NewMid()
	mc.I = &mid
	dr_rc := map[string][]*pb.RC{} //
	var iv int
	for r, ur := range gur {
		iv = n.send_ms(r, ur, mc, dr_rc)
		if iv != 0 {
			return iv
		}
	}
	return n.do_dis(mc, dr_rc)
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
	} else {
		log.E("send message err(%v) for:%v", err.Error(), mc.RemoteAddr().String())
		return -1
	}
}

//
func (n *NIM_Rh) send_ms(r string, ur []string, mc *Msg, dr_rc map[string][]*pb.RC) int {
	if len(ur) < 1 {
		return 0
	}
	cons, err := n.Db.ListCon(ur)
	if err != nil {
		log.E("list Con by R(%v) err:%v", ur, err.Error())
		return -1
	}
	log_d("found %v online user for RS(%v) in S(%v)", len(cons), ur, n.SS.Id())
	c_sid := n.SS.Id()         //current server id.
	sr_ed := map[string]byte{} //already exec
	sender := mc.GetS()
	for _, con := range cons { //do online user
		if con.R == sender {
			continue
		}
		sr_ed[con.R] = 1
		if con.Sid == c_sid { //in current server
			log_d("sending message(%v) to con(%v)", mc.ImMsg, con)
			mc.D = &con.R
			mc.A = &r                           //setting current receive user R.
			err = n.SS.Send(con.Cid, &mc.ImMsg) //send message to client.
			if err != nil {
				log.E("sending message(%v) to R(%v) err:%v", mc.ImMsg, con.R, err.Error())
				// atomic.AddUint64(&n.DC, 1)
				// mc.Ms[con.R] = MS_DONE //mark done
				// mc.Ms[con.R] = MS_PENDING + MS_SEQ + r //mark done
			}
			mc.Ms[con.R] = append(mc.Ms[con.R], &MSS{R: r, S: MS_PENDING})
		} else { //in other distribution server
			mc.Ms[con.R] = append(mc.Ms[con.R], &MSS{R: r, S: MS_PENDING})
			tr, tc := con.R, con.Cid
			if _, ok := dr_rc[con.Sid]; ok {
				dr_rc[con.Sid] = append(dr_rc[con.Sid],
					&pb.RC{
						A: &r,
						R: &tr,
						C: &tc,
					})
			} else {
				dr_rc[con.Sid] = []*pb.RC{
					&pb.RC{
						A: &r,
						R: &tr,
						C: &tc,
					},
				}
			}
		}
	}
	// log_d("sr_ed---->%v in S(%v)", sr_ed, c_sid)
	for _, tr := range ur { //do offline user
		if _, ok := sr_ed[tr]; ok {
			continue
		}
		if tr == sender {
			continue
		}
		mc.Ms[tr] = append(mc.Ms[tr], &MSS{R: r, S: MS_PENDING})
	}
	return 0
}
func (n *NIM_Rh) do_dis(mc *Msg, dr_rc map[string][]*pb.RC) int {
	err := n.Db.Store(mc) //store mesage.
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
			log.E("sending message(%v) to distribution server(%v) err:%v", mc.GetI(), dr, err.Error())
		}
	}
	return 0
}
func (n *NIM_Rh) H(obdh *impl.OBDH) {
	obdh.AddF(MK_NRC_HB, n.HB)
	obdh.AddF(MK_NRC_LI, n.LI)
	obdh.AddF(MK_NRC_LO, n.LO)
	obdh.AddF(MK_NRC_UR, n.UR)
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
	defer r.Done()
	if r.Closed() {
		return -1
	}
	var args util.Map
	_, err := r.V(&args)
	if err != nil {
		log.W("LI V fail:%v", err.Error())
		return n.writev_ce(r, err.Error())
	}
	rv, token, ct, err := n.Db.OnLogin(r, &args)
	if err != nil {
		log.W("LI OnLogin fail:%v", err.Error())
		return n.writev_ce(r, err.Error())
	}
	con := &Con{
		Sid:   n.SS.Id(),
		Cid:   r.Id(),
		R:     rv,
		S:     "N",
		T:     CT_TCP,
		C:     ct,
		Token: token,
	}
	err = n.Db.AddCon(con)
	if err != nil {
		log.W("LI AddCon fail:%v", err.Error())
		return n.writev_ce(r, err.Error())
	}
	r.SetWait(true)
	// r.Kvs().SetVal("R", rv)
	// con.Sid = ""
	res := n.writev_c(r, con)
	// go SendUnread(n.SS, n.Db, r, rv, ct)
	log.D("LI success by R(%v),CT(%v) for:%v", rv, ct, r.RemoteAddr().String())
	return res
}
func (n *NIM_Rh) LO(r netw.Cmd) int {
	defer r.Done()
	if r.Closed() {
		return -1
	}
	var args util.Map
	_, err := r.V(&args)
	if err != nil {
		log.W("LO V fail:%v", err.Error())
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
	log.D("LO success by wait(%v) for:%v", w, r.RemoteAddr().String())
	return n.writev_c(r, con)
}
func (n *NIM_Rh) UR(r netw.Cmd) int {
	defer r.Done()
	if r.Closed() {
		return -1
	}
	tr := n.Db.FUsrR(r)
	if len(tr) < 1 {
		return n.writev_ce(r, "not login")
	}
	SendUnread(n.SS, n.Db, r, tr, 0)
	return n.writev_c(r, "OK")
}
func (n *NIM_Rh) HB(r netw.Cmd) int {
	defer r.Done()
	r.Writeb(r.Data())
	return 0
}
func (n *NIM_Rh) Push(mid string) {
	n.PushChan <- mid
}
func (n *NIM_Rh) StartPushTask() {
	go n.LoopPush()
}
func (n *NIM_Rh) LoopPush() {
	n.Running = true
	log.I("starting push task-->")
	for n.Running {
		select {
		case mid := <-n.PushChan:
			if len(mid) < 1 {
				break
			}
			sc, total, err := n.DoPush_(mid)
			if err == nil {
				log_d("doing push sc(%v),total(%v)->OK", sc, total)
			} else {
				log.W("doing push sc(%v),total(%v)->ERR:%v", sc, total, err.Error())
			}
		}
	}
	log.I("stopping push task-->")
}
func (n *NIM_Rh) DoPush_(mid string) (int, int, error) {
	msg, cons, err := n.Db.ListPushTask(n.SS.Id(), mid)
	if err != nil {
		return 0, 0, err
	}
	if len(cons) < 1 {
		return 0, 0, nil
	}
	sc := 0
	// mv := map[string]string{}
	for _, con := range cons {
		msg.D = &con.R
		for _, mss := range msg.Ms[con.R] {
			if mss.S == MS_DONE {
				continue
			}
			msg.A = &mss.S
			err = n.SS.Send(con.Cid, &msg.ImMsg)
			if err != nil {
				log.W("sending push message(%v) err:%v", msg, err.Error())
			}
		}
	}
	return sc, len(cons), nil
}

// func (n *NIM_Rh) onlo(con netw.Con) error {
// 	err := n.Db.DelCon(n.SS.Id(), con.Id(), con.Kvs().StrVal("R"), CT_TCP)
// 	con.SetWait(false)
// 	return err
// }
