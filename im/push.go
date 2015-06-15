package im

import (
	"github.com/Centny/gwf/log"
	"github.com/Centny/gwf/netw"
	"github.com/Centny/gwf/netw/impl"
	"github.com/Centny/gwf/pool"
	"github.com/Centny/gwf/util"
	"strings"
)

type PushSrv struct {
	*netw.Listener
	Db      DbH
	TickLog bool
}

func (p *PushSrv) Notify(mid string) int {
	return p.Writev2([]byte{MK_PUSH_N}, &util.Map{
		"MID": mid,
	})
}
func (p *PushSrv) PushN(s string, r string, c string, t uint32) (*Msg, error) {
	return p.PushV(s, strings.Split(r, ","), []byte(c), t)
}
func (p *PushSrv) PushV(s string, r []string, c []byte, t uint32) (*Msg, error) {
	s = strings.Trim(s, " \t")
	if len(s) < 1 || r == nil || len(r) < 1 || c == nil || len(c) < 1 {
		return nil, util.Err("arguments s/r/c having nil or empty")
	}
	for _, tr := range r {
		tr = strings.Trim(tr, " \t")
		if len(tr) < 1 {
			return nil, util.Err("arguments r having empty")
		}
	}
	var ii string = p.Db.NewMid()
	var time int64 = util.Now()
	msg := &Msg{}
	msg.R = r
	msg.T = &t
	msg.C = c
	msg.I = &ii
	msg.Time = &time
	msg.S = &s
	msg.Ms = map[string][]*MSS{}
	msg.added = map[string]bool{}
	gr, ur, err := p.Db.Sift(msg.R)
	if err != nil {
		return nil, err
	}
	for _, r := range ur {
		msg.Ms[r] = append(msg.Ms[r], &MSS{R: s, S: MS_PENDING})
	}
	gur, err := p.Db.ListUsrR(gr)
	if err != nil {
		return nil, err
	}
	for g, ur := range gur {
		for _, r := range ur {
			msg.Ms[r] = append(msg.Ms[r], &MSS{R: g, S: MS_PENDING})
		}
	}
	err = p.Db.Store(msg)
	if err != nil {
		return nil, err
	}
	sc := p.Notify(msg.GetI())
	log_d("push message(%v) to (%v) server", msg, sc)
	return msg, nil

}
func (p *PushSrv) OnCmd(c netw.Cmd) int {
	if p.TickLog {
		log.D("receive tick message(%v) from %v", string(c.Data()), c.RemoteAddr())
	}
	return 0
}

// func (p *PushSrv) DoPcm(hs *routing.HTTPSession) routing.HResult {
// 	var r, c, s string
// 	var t int64
// 	err := hs.ValidCheckVal(`
// 		s,O|S,L:0;
// 		r,R|S,L:0;
// 		c,R|S,L:0;
// 		t,R|I,O:0~1~101;
// 		`, &s, &r, &c, &t)
// 	if err != nil {
// 		return hs.MsgResErr2(1, "arg-err", err)
// 	}
// 	if len(s) > 0 {
// 		_, err = p.PushN(s, r, c, uint32(t))
// 	} else {
// 		uid := hs.IntVal("uid")
// 		if s, err = p.Db.FindUsrR(uid); err == nil {
// 			_, err = p.PushN(s, r, c, uint32(t))
// 		}
// 	}
// 	if err == nil {
// 		return hs.MsgRes("OK")
// 	} else {
// 		return hs.MsgResErr2(1, "srv-err", err)
// 	}
// }

func NewPushSrv(p *pool.BytePool, port string, n string, h netw.CCHandler, db DbH) *PushSrv {
	return NewPushSrvN(p, port, n, h, impl.Json_NewCon, db)
}

func NewPushSrvN(p *pool.BytePool, port string, n string, h netw.CCHandler, ncf netw.NewConF, db DbH) *PushSrv {
	obdh := impl.NewOBDH()
	psrv := &PushSrv{
		Listener: netw.NewListenerN(p, port, n, netw.NewCCH(h, obdh), ncf),
		Db:       db,
	}
	obdh.AddH(MK_PUSH_HB, psrv)
	return psrv
}
