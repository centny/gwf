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
	Db DbH
}

func (p *PushSrv) Notify(mid string) int {
	return p.Writev(&util.Map{
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
	msg.Ms = map[string]string{}
	gr, ur, err := p.Db.Sift(msg.R)
	if err != nil {
		return nil, err
	}
	for _, r := range ur {
		if len(msg.Ms[r]) < 1 {
			msg.Ms[r] = MS_PENDING + MS_SEQ + r
		} else {
			msg.Ms[r] += MS_SEQ + r
		}
	}
	gur, err := p.Db.ListUsrR(gr)
	if err != nil {
		return nil, err
	}
	for g, ur := range gur {
		for _, r := range ur {
			if len(msg.Ms[r]) < 1 {
				msg.Ms[r] = MS_PENDING + MS_SEQ + g
			} else {
				msg.Ms[r] += MS_SEQ + g
			}
		}
	}
	err = p.Db.Store(msg)
	if err != nil {
		return nil, err
	}
	sc := p.Notify(msg.GetI())
	log.D("push message(%v) to server(%v)", msg, sc)
	return msg, nil

}
func NewPushSrv(p *pool.BytePool, port string, n string, h netw.CCHandler, db DbH) *PushSrv {
	return NewPushSrvN(p, port, n, h, impl.Json_NewCon, db)
}

func NewPushSrvN(p *pool.BytePool, port string, n string, h netw.CCHandler, ncf netw.NewConF, db DbH) *PushSrv {
	return &PushSrv{
		Listener: netw.NewListenerN(p, port, n, h, ncf),
		Db:       db,
	}
}
