package im

import (
	"github.com/Centny/gwf/log"
	"github.com/Centny/gwf/netw"
	"github.com/Centny/gwf/netw/impl"
	"github.com/Centny/gwf/util"
	"sync/atomic"
)

//
type NIM_Rh struct {
	Db DbH
	SS Sender
	DS Sender
	DC uint64
}

func (n *NIM_Rh) OnConn(c netw.Con) bool {
	return true
}
func (n *NIM_Rh) OnClose(c netw.Con) {
	n.Db.OnUsrLogout(c.Kvs().StrVal("R"), nil)
	n.onlo(c)
}

func (n *NIM_Rh) OnCmd(c netw.Cmd) int {
	defer c.Done()
	// log_d("NIM_Rh receive data:%v", string(c.Data()))
	var mc Msg
	_, err := c.V(&mc)
	if err != nil {
		log.E("convert valus to IM msg error:%v", err.Error())
		return -1
	}
	mc.S = c.Id()
	mc.Cmd = c
	mc.Ms = map[string]string{}
	return n.OnMsg(&mc)
}
func (n *NIM_Rh) OnMsg(mc *Msg) int {
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
	log_d("sending message(%v) to RS(%v)", mc, ur)
	//
	cons, err := n.Db.ListCon(ur)
	if err != nil {
		log.E("list Con by R(%v) err:%v", ur, err.Error())
		return -1
	}
	log_d("found %v online user for RS(%v)", len(cons), ur)
	c_sid := n.SS.Id()                      //current server id.
	sr_ed := map[string]byte{}              //already exec
	dr_rc := map[string]map[string]string{} //
	for _, con := range cons {              //do online user
		sr_ed[con.R] = 1
		if con.Sid == c_sid { //in current server
			mc.D = con.R                 //setting current receive user R.
			err = n.SS.Send(con.Cid, mc) //send message to client.
			if err == nil {
				atomic.AddUint64(&n.DC, 1)
				mc.Ms[con.R] = MS_DONE //mark done
			} else {
				log_d("sending message to R(%v) err:%v", con.R, err.Error())
				mc.Ms[con.R] = MS_ERR + err.Error() //mark send error.
			}
		} else { //in other distribution server
			mc.Ms[con.R] = MS_PENDING //mark to pending.
			if _, ok := dr_rc[con.Sid]; ok {
				dr_rc[con.Sid][con.R] = con.Cid
			} else {
				dr_rc[con.Sid] = map[string]string{
					con.R: con.Cid,
				}
			}
		}
	}

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
	for dr, rc := range dr_rc { //if having distribution message.
		dmc := &DsMsg{
			M:  *mc,
			RC: rc,
		}
		err = n.DS.Send(dr, dmc)
		if err == nil { //if not err,the other distribution server will makr result.
			continue
		} else {
			log.E("sending message(%v) to distribution server(%v) err:%v", mc, dr, err.Error())
		}
	}
	return 0
}

func (n *NIM_Rh) Exec(r *impl.RCM_Cmd) (interface{}, error) {
	log_d("call action(%v)", r.Name)
	switch r.Name {
	case "LI":
		return n.LI(r)
	case "LO":
		return n.LO(r)
	}
	return nil, util.Err("action not found by name(%v)", r.Name)
}

func (n *NIM_Rh) LI(r *impl.RCM_Cmd) (interface{}, error) {
	rv, err := n.Db.OnUsrLogin(r, r.Map)
	if err != nil {
		return r.CRes(1, err.Error())
	}
	con := &Con{
		Sid: n.SS.Id(),
		Cid: r.Id(),
		R:   rv,
		S:   "N",
		T:   CT_TCP,
	}
	err = n.Db.AddCon(con)
	if err != nil {
		return r.CRes(1, err.Error())
	}
	r.SetWait(true)
	r.Kvs().SetVal("R", rv)
	// con.Sid = ""
	return r.CRes(0, con)
}
func (n *NIM_Rh) LO(r *impl.RCM_Cmd) (interface{}, error) {
	err := n.Db.OnUsrLogout(r.Kvs().StrVal("R"), r.Map)
	err = n.onlo(r)
	if err == nil {
		return r.CRes(0, "OK")
	} else {
		return r.CRes(1, err.Error())
	}
}
func (n *NIM_Rh) onlo(con netw.Con) error {
	err := n.Db.DelCon(n.SS.Id(), con.Id(), con.Kvs().StrVal("R"), CT_TCP)
	con.SetWait(false)
	return err
}
