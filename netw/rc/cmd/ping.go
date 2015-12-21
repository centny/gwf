package cmd

import (
	"github.com/Centny/gwf/log"
	"github.com/Centny/gwf/netw/impl"
	"github.com/Centny/gwf/netw/rc"
	"github.com/Centny/gwf/pool"
	"github.com/Centny/gwf/util"
)

type Ping struct {
	ShowLog bool
	P       *pool.BytePool
	S       *rc.RC_Listener_m
	C       *rc.RC_Runner_m
}

func RegPingH(p *pool.BytePool, s *rc.RC_Listener_m) *Ping {
	ph := &Ping{
		P: p,
		S: s,
	}
	s.AddHFunc("ping", ph.PingH)
	return ph
}

func NewPing(p *pool.BytePool, c *rc.RC_Runner_m) *Ping {
	return &Ping{
		P: p,
		C: c,
	}
}

func (p *Ping) slog(f string, args ...interface{}) {
	if p.ShowLog {
		log.D_(1, f, args...)
	}
}

//ping command handler
func (p *Ping) PingH(rc *impl.RCM_Cmd) (interface{}, error) {
	var l int64 = 8
	var data string = ""
	err := rc.ValidF(`
		len,O|I,R:0;
		data,O|S,L:0;
		`, &l, &data)
	if err != nil {
		return nil, err
	}
	p.slog("Ping->receive ping len(%v)", l)
	bys := p.P.Alloc(int(l))
	for i := 0; i < int(l); i++ {
		bys[i] = byte('A' + i%26)
	}
	defer p.S.P.Free(bys)
	return util.Map{
		"data": string(bys),
	}, nil
}

//ping to server
func (p *Ping) tping(ul, dl int) (int64, error) {
	if ul < 1 {
		ul = 8
	}
	if dl < 1 {
		dl = 8
	}
	bys := p.P.Alloc(ul)
	for i := 0; i < ul; i++ {
		bys[i] = byte('A' + i%26)
	}
	defer p.P.Free(bys)
	beg := util.Now()
	_, err := p.C.VExec_m("ping", util.Map{
		"len":  dl,
		"data": string(bys),
	})
	return util.Now() - beg, err
}

func (p *Ping) Ping(ul, dl int) (int64, error) {
	p.slog("Ping->ping to server by u_len(%v),d_len(%v)", ul, dl)
	used, err := p.tping(ul, dl)
	p.slog("Ping->response from server by used(%v),err(%v)", used, err)
	return used, err
}

//ping to server multi count
func (p *Ping) PingC(ul, dl, c int) (int64, error) {
	p.slog("Ping->ping to server by u_len(%v),d_len(%v),c(%v)", ul, dl, c)
	var total int64 = 0
	for i := 0; i < c; i++ {
		used, err := p.tping(ul, dl)
		if err != nil {
			p.slog("Ping->response from server by used(%v),err(%v)", total, err)
			return total, err
		}
		total += used
	}
	p.slog("Ping->response from server by used(%v),err(%v)", total, nil)
	return total, nil
}
