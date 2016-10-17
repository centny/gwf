package plugin

// import (
// 	"fmt"
// 	"github.com/Centny/gwf/log"
// 	"github.com/Centny/gwf/netw"
// 	"github.com/Centny/gwf/netw/impl"
// 	nrc "github.com/Centny/gwf/netw/rc"
// 	"github.com/Centny/gwf/pool"
// 	"github.com/Centny/gwf/util"
// 	"strings"
// 	"sync"
// 	"sync/atomic"
// )

// const (
// 	RC_M_J = "J"
// )

// type RC_S_H struct {
// 	P *pool.BytePool
// 	L *nrc.RC_Listener_m
// 	//
// 	PH_ *Ping_S_H
// 	//
// 	Conned map[string]int
// 	Closed map[string]int
// 	//
// 	cids int64

// 	//task lock
// 	ts_lck sync.RWMutex
// }

// func NewRC_S_H(p *pool.BytePool, r *nrc.RC_Runner_m) *RC_S_H {
// 	sh := &RC_S_H{
// 		P:      p,
// 		Conned: map[string]int{},
// 		Closed: map[string]int{},
// 	}
// 	r.AddHFunc("rcs_start", sh.StartH)
// 	r.AddHFunc("rcs_stop", sh.StopH)
// 	r.AddHFunc("rcs_token", sh.TokenH)
// 	r.AddHFunc("rcs_status", sh.StatusH)
// 	return sh
// }

// func (r *RC_S_H) OnConn(c netw.Con) bool {
// 	return true
// }

// func (r *RC_S_H) OnClose(c netw.Con) {
// 	token := c.Kvs().StrVal("token")
// 	if len(token) > 0 {
// 		r.Closed[token] += 1
// 	}
// }

// func (r *RC_S_H) OnCmd(c netw.Cmd) int {
// 	return 0
// }

// func (r *RC_S_H) OnLogin(rc *impl.RCM_Cmd, token string) (string, error) {
// 	r.Conned[token] += 1
// 	cids := atomic.AddInt64(&r.cids, 1)
// 	return fmt.Sprintf("N-%v", cids), nil
// }

// func (r *RC_S_H) StartH(rc *impl.RCM_Cmd) (interface{}, error) {
// 	r.ts_lck.Lock()
// 	defer r.ts_lck.Unlock()
// 	if r.L != nil {
// 		return util.Map{"code": -1, "err": "service is running"}, nil
// 	}
// 	var m, addr string = RC_M_J, ""
// 	err := rc.ValidF(`
// 		m,O|S,O:J;
// 		addr,R|S,L:0;
// 		`, &m, &addr)
// 	if err != nil {
// 		return util.Map{"code": -2, "err": err.Error()}, nil
// 	}
// 	log.I("RC_S_H will start by addr(%v),m(%v)", addr, m)
// 	switch m {
// 	case RC_M_J:
// 		r.L = nrc.NewRC_Listener_m_j(r.P, addr, r)
// 		r.L.LCH = r
// 	}
// 	r.PH_ = RegPing_S_H(r.P, r.L)
// 	err = r.L.Run()
// 	if err == nil {
// 		log.I("RC_S_H start by addr(%v),m(%v) success", addr, m)
// 		return util.Map{"code": 0}, nil
// 	} else {
// 		r.L = nil
// 		log.E("RC_S_H start by addr(%v),m(%v) error->%v", addr, m, err.Error())
// 		return util.Map{"code": -3, "err": err.Error()}, nil
// 	}
// }

// func (r *RC_S_H) StopH(rc *impl.RCM_Cmd) (interface{}, error) {
// 	r.ts_lck.Lock()
// 	defer r.ts_lck.Unlock()
// 	if r.L == nil {
// 		return util.Map{"code": -1, "err": "service is not running"}, nil
// 	}
// 	r.L.Close()
// 	r.L.Wait()
// 	r.L = nil
// 	log.I("RC_S_H stop success...")
// 	return util.Map{"code": 0}, nil
// }

// func (r *RC_S_H) TokenH(rc *impl.RCM_Cmd) (interface{}, error) {
// 	if r.L == nil {
// 		return util.Map{"code": -1, "err": "service is not running"}, nil
// 	}
// 	var tokens string
// 	err := rc.ValidF(`
// 		tokens,R|S,L:0;
// 		`, &tokens)
// 	if err != nil {
// 		return util.Map{"code": -2, "err": err.Error()}, nil
// 	}
// 	r.L.AddToken2(strings.Split(tokens, ","))
// 	log.I("RC_S_H adding token(%v)", tokens)
// 	return util.Map{"code": 0}, nil
// }

// func (r *RC_S_H) Status() util.Map {
// 	if r.L == nil {
// 		return util.Map{
// 			"code":   0,
// 			"status": TS_NOT_START,
// 		}
// 	} else {
// 		return util.Map{
// 			"code":   0,
// 			"status": TS_RUNNING,
// 			"ping":   r.PH_.Status(),
// 			"conned": r.Conned,
// 			"closed": r.Closed,
// 		}
// 	}
// }
// func (r *RC_S_H) StatusH(rc *impl.RCM_Cmd) (interface{}, error) {
// 	return r.Status(), nil
// }

// const (
// 	S_CONNECT   = "CONNECT"
// 	S_CLOSED    = "CLOSED"
// 	S_LOGINED   = "LOGINED"
// 	S_LOGIN_ERR = "LOGIN_ERR"
// )

// type RC_C_H struct {
// 	P *pool.BytePool
// 	R *nrc.RC_Runner_m
// 	//
// 	PH_ *Ping_C_H
// 	//
// 	Conned    int
// 	Closed    int
// 	Logined   int
// 	Token     string
// 	ConStatus string
// 	//
// 	wc_l chan string
// 	//task lock
// 	ts_lck sync.RWMutex
// }

// func NewRC_C_H(p *pool.BytePool, r *nrc.RC_Runner_m) *RC_C_H {
// 	sh := &RC_C_H{
// 		P:    p,
// 		wc_l: make(chan string, 3),
// 	}
// 	r.AddHFunc("rcc_start", sh.StartH)
// 	r.AddHFunc("rcc_stop", sh.StopH)
// 	r.AddHFunc("rcc_start_ping", sh.StartPingH)
// 	r.AddHFunc("rcc_stop_ping", sh.StopPingH)
// 	r.AddHFunc("rcc_status", sh.StatusH)
// 	return sh
// }

// func (r *RC_C_H) OnConn(c netw.Con) bool {
// 	r.Conned += 1
// 	go r.DoLogin()
// 	r.ConStatus = S_CONNECT
// 	return true
// }

// func (r *RC_C_H) OnClose(c netw.Con) {
// 	r.Closed += 1
// 	r.ConStatus = S_CLOSED
// }

// func (r *RC_C_H) OnCmd(c netw.Cmd) int {
// 	return 0
// }

// func (r *RC_C_H) DoLogin() {
// 	err := r.R.Login_(r.Token)
// 	if err == nil {
// 		r.Logined += 1
// 		r.ConStatus = S_LOGINED
// 		r.wc_l <- ""
// 		log.I("RC_C_H login by token(%v) success", r.Token)
// 	} else {
// 		r.ConStatus = S_LOGIN_ERR
// 		r.wc_l <- err.Error()
// 		log.E("RC_C_H login by token(%v) err->%v", r.Token, err.Error())
// 	}
// }

// func (r *RC_C_H) StartH(rc *impl.RCM_Cmd) (interface{}, error) {
// 	r.ts_lck.Lock()
// 	defer r.ts_lck.Unlock()
// 	if r.R != nil {
// 		return util.Map{"code": -1, "err": "runner is running"}, nil
// 	}
// 	var m, addr, token string = RC_M_J, "", ""
// 	err := rc.ValidF(`
// 		m,O|S,O:J;
// 		addr,R|S,L:0;
// 		token,R|S,L:0;
// 		`, &m, &addr, &token)
// 	if err != nil {
// 		return util.Map{"code": -2, "err": err.Error()}, nil
// 	}
// 	log.I("RC_C_H will start by addr(%v),token(%v),m(%v)", addr, token, m)
// 	r.Token = token
// 	switch m {
// 	case RC_M_J:
// 		r.R = nrc.NewRC_Runner_m_j(r.P, addr, r)
// 	}
// 	r.PH_ = NewPing_C_H2(r.P, r.R)
// 	err = r.R.Run()
// 	if err != nil {
// 		r.R = nil
// 		log.E("RC_C_H start by addr(%v),token(%v),m(%v) run err->%v", addr, token, m, err.Error())
// 		return util.Map{"code": -3, "err": err.Error()}, nil
// 	}
// 	msg := <-r.wc_l
// 	if len(msg) < 1 {
// 		log.I("RC_C_H start by addr(%v),token(%v),m(%v) success", addr, token, m)
// 		return util.Map{"code": 0}, nil
// 	} else {
// 		r.R = nil
// 		log.E("RC_C_H start by addr(%v),token(%v),m(%v) login err->%v", addr, token, m, msg)
// 		return util.Map{"code": -4, "err": msg}, nil
// 	}

// }

// func (r *RC_C_H) StopH(rc *impl.RCM_Cmd) (interface{}, error) {
// 	r.ts_lck.Lock()
// 	defer r.ts_lck.Unlock()
// 	if r.R == nil {
// 		return util.Map{"code": -1, "err": "runner is not running"}, nil
// 	}
// 	r.R.Stop()
// 	r.R.Wait()
// 	r.R = nil
// 	log.I("RC_C_H stop success...")
// 	return util.Map{"code": 0}, nil
// }

// func (r *RC_C_H) StartPingH(rc *impl.RCM_Cmd) (interface{}, error) {
// 	r.ts_lck.Lock()
// 	defer r.ts_lck.Unlock()
// 	if r.R == nil {
// 		return util.Map{"code": -1, "err": "runner is not running"}, nil
// 	}
// 	if r.PH_.Running > 0 {
// 		return util.Map{"code": -2, "err": "PingS is running"}, nil
// 	}
// 	var delay int64 = 3 * 60 * 1000
// 	var min, max int64 = 8, 1024000
// 	err := rc.ValidF(`
// 		delay,O|I,R:0;
// 		min,O|I,R:0;
// 		max,O|I,R:0;
// 		`, &delay, &min, &max)
// 	if err != nil {
// 		return util.Map{"code": -3, "err": err.Error()}, nil
// 	}
// 	r.PH_.StartPingS(delay, int(min), int(max))
// 	log.I("RC_C_H start PingS by delay(%v),min(%v),max(%v) success", delay, min, max)
// 	return util.Map{"code": 0}, nil
// }

// func (r *RC_C_H) StopPingH(rc *impl.RCM_Cmd) (interface{}, error) {
// 	r.ts_lck.Lock()
// 	defer r.ts_lck.Unlock()
// 	if r.R == nil {
// 		return util.Map{"code": -1, "err": "runner is not running"}, nil
// 	}
// 	if r.PH_.Running < 1 {
// 		return util.Map{"code": -2, "err": "PingS is not running"}, nil
// 	}
// 	r.PH_.StopPingS()
// 	log.I("RC_C_H stop PingS success")
// 	return util.Map{"code": 0}, nil
// }

// func (r *RC_C_H) Status() util.Map {
// 	var res util.Map
// 	if r.R == nil {
// 		res = util.Map{
// 			"code":   0,
// 			"status": TS_NOT_START,
// 		}
// 	} else {
// 		res = util.Map{
// 			"code":       0,
// 			"status":     TS_RUNNING,
// 			"con_status": r.ConStatus,
// 			"conned":     r.Conned,
// 			"closed":     r.Closed,
// 			"logined":    r.Logined,
// 			"token":      r.Token,
// 			"ping":       r.PH_.Status(),
// 		}
// 	}
// 	return res
// }

// func (r *RC_C_H) StatusH(rc *impl.RCM_Cmd) (interface{}, error) {
// 	return r.Status(), nil
// }

// type RC_CTL struct {
// 	Exec_m func(name string, args interface{}) (util.Map, error)
// }

// func NewRC_CTL(c *impl.RCM_Con) *RC_CTL {
// 	return &RC_CTL{
// 		Exec_m: c.Exec_m,
// 	}
// }
// func (r *RC_CTL) RC_S_Start(addr, m string) error {
// 	res, err := r.Exec_m("rcs_start", util.Map{
// 		"m":    m,
// 		"addr": addr,
// 	})
// 	if err != nil {
// 		return err
// 	}
// 	if res.IntVal("code") == 0 {
// 		return nil
// 	} else {
// 		return util.Err("RC_S_Start err->%v", res.StrVal("err"))
// 	}
// }
// func (r *RC_CTL) RC_S_Stop() error {
// 	res, err := r.Exec_m("rcs_stop", util.Map{})
// 	if err != nil {
// 		return err
// 	}
// 	if res.IntVal("code") == 0 {
// 		return nil
// 	} else {
// 		return util.Err("RC_S_Stop err->%v", res.StrVal("err"))
// 	}
// }
// func (r *RC_CTL) RC_S_Token(ts []string) error {
// 	res, err := r.Exec_m("rcs_token", util.Map{
// 		"tokens": strings.Join(ts, ","),
// 	})
// 	if err != nil {
// 		return err
// 	}
// 	if res.IntVal("code") == 0 {
// 		return nil
// 	} else {
// 		return util.Err("RC_S_Token err->%v", res.StrVal("err"))
// 	}
// }
// func (r *RC_CTL) RC_S_Status() (util.Map, error) {
// 	return r.Exec_m("rcs_status", util.Map{})
// }
// func (r *RC_CTL) RC_C_Start(addr, token, m string) error {
// 	res, err := r.Exec_m("rcc_start", util.Map{
// 		"m":     m,
// 		"addr":  addr,
// 		"token": token,
// 	})
// 	if err != nil {
// 		return err
// 	}
// 	if res.IntVal("code") == 0 {
// 		return nil
// 	} else {
// 		return util.Err("RC_C_Start err->%v", res.StrVal("err"))
// 	}
// }
// func (r *RC_CTL) RC_C_Stop() error {
// 	res, err := r.Exec_m("rcc_stop", util.Map{})
// 	if err != nil {
// 		return err
// 	}
// 	if res.IntVal("code") == 0 {
// 		return nil
// 	} else {
// 		return util.Err("RC_C_Stop err->%v", res.StrVal("err"))
// 	}
// }
// func (r *RC_CTL) RC_C_StartPing(delay int64, min, max int) error {
// 	res, err := r.Exec_m("rcc_start_ping", util.Map{
// 		"delay": delay,
// 		"min":   min,
// 		"max":   max,
// 	})
// 	if err != nil {
// 		return err
// 	}
// 	if res.IntVal("code") == 0 {
// 		return nil
// 	} else {
// 		return util.Err("RC_C_StartPing err->%v", res.StrVal("err"))
// 	}
// }
// func (r *RC_CTL) RC_C_StopPing() error {
// 	res, err := r.Exec_m("rcc_stop_ping", util.Map{})
// 	if err != nil {
// 		return err
// 	}
// 	if res.IntVal("code") == 0 {
// 		return nil
// 	} else {
// 		return util.Err("RC_C_StopPing err->%v", res.StrVal("err"))
// 	}
// }
// func (r *RC_CTL) RC_C_Status() (util.Map, error) {
// 	return r.Exec_m("rcc_status", util.Map{})
// }

// type RC_CTL_S_H struct {
// 	L *nrc.RC_Listener_m
// }

// func NewRC_CTL_S_H(l *nrc.RC_Listener_m) *RC_CTL_S_H {
// 	return &RC_CTL_S_H{
// 		L: l,
// 	}
// }

// func (r *RC_CTL_S_H) CTL(cid string) *RC_CTL {
// 	con := r.L.CmdC(cid)
// 	if con == nil {
// 		return nil
// 	} else {
// 		return NewRC_CTL(con)
// 	}
// }
