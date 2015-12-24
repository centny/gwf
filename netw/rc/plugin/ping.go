package plugin

import (
	"github.com/Centny/gwf/log"
	"github.com/Centny/gwf/netw/impl"
	"github.com/Centny/gwf/netw/rc"
	"github.com/Centny/gwf/pool"
	"github.com/Centny/gwf/util"
	"time"
)

type Ping_S_H struct {
	ShowLog bool
	P       *pool.BytePool
	S       *rc.RC_Listener_m
	C       int64
}

type Ping_C_H struct {
	ShowLog bool
	P       *pool.BytePool
	VExec_m func(name string, args interface{}) (util.Map, error)
	C       int64

	//ping speed field
	Running int
	Delay   int64
	Speed   int64
	Err     string
	ErrC    int64
	ps_wc_  chan int
}

func RegPing_S_H(p *pool.BytePool, s *rc.RC_Listener_m) *Ping_S_H {
	ph := &Ping_S_H{
		P: p,
		S: s,
	}
	s.AddHFunc("ping", ph.PingH)
	return ph
}

func NewPing_C_H(p *pool.BytePool, c *impl.RCM_Con) *Ping_C_H {
	return &Ping_C_H{
		P: p,
		VExec_m: func(name string, args interface{}) (util.Map, error) {
			var dest util.Map
			_, err := c.Exec(name, args, &dest)
			return dest, err
		},
		ps_wc_: make(chan int),
	}
}

func NewPing_C_H2(p *pool.BytePool, c *rc.RC_Runner_m) *Ping_C_H {
	return &Ping_C_H{
		P: p,
		VExec_m: func(name string, args interface{}) (util.Map, error) {
			var dest util.Map
			_, err := c.Exec(name, args, &dest)
			return dest, err
		},
		ps_wc_: make(chan int),
	}
}

func (p *Ping_S_H) slog(f string, args ...interface{}) {
	if p.ShowLog {
		log.D_(1, f, args...)
	}
}

func (p *Ping_C_H) slog(f string, args ...interface{}) {
	if p.ShowLog {
		log.D_(1, f, args...)
	}
}

//ping command handler
func (p *Ping_S_H) PingH(rc *impl.RCM_Cmd) (interface{}, error) {
	if p.S == nil {
		panic("RC_Listener_m is not initial")
	}
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
	p.C += 1
	return util.Map{
		"data": string(bys),
	}, nil
}

//ping to server
func (p *Ping_C_H) tping(up, down int) (int64, error) {
	if up < 1 {
		up = 8
	}
	if down < 1 {
		down = 8
	}
	bys := p.P.Alloc(up)
	for i := 0; i < up; i++ {
		bys[i] = byte('A' + i%26)
	}
	defer p.P.Free(bys)
	beg := util.Now()
	_, err := p.VExec_m("ping", util.Map{
		"len":  down,
		"data": string(bys),
	})
	p.C += 1
	return util.Now() - beg, err
}

//ping to server
func (p *Ping_C_H) Ping(up, down int) (int64, error) {
	if p.VExec_m == nil {
		panic("VExec_m is not initial")
	}
	p.slog("Ping->ping to server by u_len(%v),d_len(%v)", up, down)
	used, err := p.tping(up, down)
	p.slog("Ping->response from server by used(%v),err(%v)", used, err)
	return used, err
}

//ping to server multi count by upload data leng,download data len and ping count
func (p *Ping_C_H) PingC(up, down, c int) (int64, error) {
	if p.VExec_m == nil {
		panic("VExec_m is not initial")
	}
	p.slog("Ping->ping to server by u_len(%v),d_len(%v),c(%v)", up, down, c)
	var total int64 = 0
	for i := 0; i < c; i++ {
		used, err := p.tping(up, down)
		if err != nil {
			p.slog("Ping->response from server by used(%v),err(%v)", total, err)
			return total, err
		}
		total += used
	}
	p.slog("Ping->response from server by used(%v),err(%v)", total, nil)
	return total, nil
}

//ping to server for testing delay and speed by min/max data length.
//return the min data ping delay and max data ping speed by second
func (p *Ping_C_H) PingS(min, max int) (delay, speed int64, err error) {
	if p.VExec_m == nil {
		panic("VExec_m is not initial")
	}
	delay, err = p.Ping(min, min)
	if err != nil {
		return 0, 0, err
	}
	if max < 10240 {
		max = 10240
	}
	max = (max / 10240) * 10240
	var used int64
	used, err = p.PingC(10240, 10240, max/10240)
	if used < 1 {
		used = 1
	}
	speed = int64(max) / used
	return
}

//loop ping to server for testing delay and speed.
func (p *Ping_C_H) DoPingS(delay int64, min, max int) {
	if p.VExec_m == nil {
		panic("VExec_m is not initial")
	}
	defer func() {
		p.ps_wc_ <- 0
		log.I("Ping_C_H DoPingS is stopping")
	}()
	log.I("Ping_C_H start PingS by delay(%v),min(%v),max(%v)", delay, min, max)
	var err error
	var waiting int64
	p.Running = 1
	p.ps_wc_ <- 1
	for p.Running > 0 {
		p.slog("Ping_C_H do PingS by min(%v),max(%v)", min, max)
		p.Delay, p.Speed, err = p.PingS(min, max)
		if err != nil {
			p.Err = err.Error()
			p.ErrC += 1
		}
		waiting = delay
		for p.Running > 0 && waiting > 0 {
			time.Sleep(500 * time.Millisecond)
			waiting -= 500
		}
	}
}

//start loop task for pinging to server for testing delay and speed.
func (p *Ping_C_H) StartPingS(delay int64, min, max int) {
	if p.VExec_m == nil {
		panic("VExec_m is not initial")
	}
	go p.DoPingS(delay, min, max)
	<-p.ps_wc_
}

//stop loop task for pinging to server for testing delay and speed.
func (p *Ping_C_H) StopPingS() {
	p.Running = 0
	<-p.ps_wc_
}

func (p *Ping_S_H) Status() util.Map {
	return util.Map{
		"code":   0,
		"ping_c": p.C,
	}
}

func (p *Ping_C_H) Status() util.Map {
	return util.Map{
		"code":   0,
		"ping_c": p.C,
		"ping_s": util.Map{
			"running": p.Running,
			"delay":   p.Delay,
			"speed":   p.Speed,
			"err":     p.Err,
			"err_c":   p.ErrC,
		},
	}
}
