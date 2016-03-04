package filter

import (
	"github.com/Centny/gwf/log"
	"github.com/Centny/gwf/netw/impl"
	"github.com/Centny/gwf/netw/rc"
	"github.com/Centny/gwf/routing"
	"github.com/Centny/gwf/util"
	"sync"
	"time"
)

type AttrFilter struct {
	Key     string
	Attrs   util.Map
	TimeL   map[string]int64
	Delay   int64
	Timeout int64
	ShowLog bool
	running bool
	lck     sync.RWMutex
}

func NewAttrFilter(key string) *AttrFilter {
	return &AttrFilter{
		Key:     key,
		Attrs:   util.Map{},
		TimeL:   map[string]int64{},
		Delay:   5000,
		Timeout: 30 * 60000,
		ShowLog: false,
		lck:     sync.RWMutex{},
	}
}
func (a *AttrFilter) SrvHTTP(hs *routing.HTTPSession) routing.HResult {
	var token = hs.CheckVal(a.Key)
	if len(token) < 1 {
		return hs.MsgResErr2(1, "arg-err", util.Err("the %v args is required", a.Key))
	}
	if !a.Attrs.Exist(token) {
		return hs.MsgResErr2(1, "arg-err", util.Err("the token(%v) by key(%v) not found", token, a.Key))
	}
	var attrs = a.Attrs.MapVal(token)
	if a.ShowLog {
		log.D("AttrFilter do filter to path(%v) with token(%v),attrs(%v)", hs.R.URL.Path, token, util.S2Json(attrs))
	}
	for key, val := range attrs {
		hs.SetVal(key, val)
	}
	return routing.HRES_CONTINUE
}
func (a *AttrFilter) RegisterH_w(hs *routing.HTTPSession) routing.HResult {
	hs.R.ParseForm()
	var args = util.Map{}
	for key, _ := range hs.R.Form {
		args[key] = hs.R.FormValue(key)
	}
	a.lck.Lock()
	defer a.lck.Unlock()
	var token = util.UUID()
	a.Attrs[token] = args
	a.TimeL[token] = util.Now()
	if a.ShowLog {
		log.D("AttrFilter(Web) register token(%v) by args(%v) success", token, util.S2Json(args))
	}
	return hs.MsgRes(token)
}

func (a *AttrFilter) Hand_w(pre string, mux *routing.SessionMux) {
	mux.HFunc("^"+pre+"/attr/reg", a.RegisterH_w)
}

func (a *AttrFilter) RegisterH_rc(rc *impl.RCM_Cmd) (interface{}, error) {
	a.lck.Lock()
	defer a.lck.Unlock()
	var token = util.UUID()
	a.Attrs[token] = *rc.Map
	a.TimeL[token] = util.Now()
	if a.ShowLog {
		log.D("AttrFilter(RC) register token(%v) by args(%v) success", token, util.S2Json(rc.Map))
	}
	return token, nil
}

func (a *AttrFilter) Hand_rc(l *rc.RC_Listener_m) {
	l.AddHFunc("attr/reg", a.RegisterH_rc)
}

func (a *AttrFilter) ChkTimeout() {
	a.lck.Lock()
	defer a.lck.Unlock()
	var now = util.Now()
	var removed int = 0
	for token, beg := range a.TimeL {
		if now-beg >= a.Timeout {
			delete(a.TimeL, token)
			delete(a.Attrs, token)
			removed += 1
		}
	}
	if removed > 0 {
		log.D("AttrFilter do timeout with %v token be removed", removed)
	}
}

func (a *AttrFilter) loop_timeout(delay int64) {
	a.running = true
	log.D("AttrFilter start loop timeout with delay(%v),timeout(%v)", delay, a.Timeout)
	for a.running {
		a.ChkTimeout()
		time.Sleep(time.Duration(delay) * time.Millisecond)
	}
	log.D("AttrFilter loop timeout is stopped")
}

func (a *AttrFilter) StartTimeout() {
	go a.loop_timeout(a.Delay)
}

func (a *AttrFilter) StopTimeout() {
	a.running = false
}

func RegisterAttr(runner *rc.RC_Runner_m, attrs util.Map) (string, error) {
	return runner.VExec_s("attr/reg", attrs)
}
