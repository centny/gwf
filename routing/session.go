package routing

import (
	"net/http"
	"sync"
	"time"

	"github.com/Centny/gwf/log"
	"github.com/Centny/gwf/util"
)

type SrvSession struct {
	token string
	begin int64
	lck   sync.RWMutex
	kvs   map[string]interface{}
}

func NewSrvSession() *SrvSession {
	return &SrvSession{
		lck: sync.RWMutex{},
		kvs: map[string]interface{}{},
	}
}

func (s *SrvSession) Val(key string) interface{} {
	s.lck.RLock()
	defer s.lck.RUnlock()
	if v, ok := s.kvs[key]; ok {
		return v
	} else {
		return nil
	}
}
func (s *SrvSession) Set(key string, val interface{}) {
	s.lck.Lock()
	defer s.lck.Unlock()
	if val == nil {
		delete(s.kvs, key)
	} else {
		s.kvs[key] = val
	}
}
func (s *SrvSession) Token() string {
	return s.token
}
func (s *SrvSession) Flush() error {
	s.begin = util.Timestamp(time.Now())
	return nil
}

//
type SrvSessionBuilder struct {
	//
	Domain    string
	Path      string
	Timeout   int64
	CDelay    time.Duration
	CookieKey string //cookie key
	ShowLog   bool
	//
	evh     SessionEvHandler
	looping bool
	ks      map[string]*SrvSession //key session
	ks_lck  sync.RWMutex
}

func NewSrvSessionBuilder(domain string, path string, ckey string, timeout int64, cdelay time.Duration) *SrvSessionBuilder {
	sb := SrvSessionBuilder{}
	sb.Domain = domain
	sb.Path = path
	sb.Timeout = timeout
	sb.CDelay = cdelay
	sb.CookieKey = ckey
	sb.ks = map[string]*SrvSession{}
	sb.ShowLog = false
	sb.SetEvH(SessionEvHFunc(func(t string, s Session) {
	}))
	return &sb
}
func (s *SrvSessionBuilder) log(f string, args ...interface{}) {
	if s.ShowLog {
		log.D(f, args...)
	}
}
func (s *SrvSessionBuilder) SetEvH(h SessionEvHandler) {
	s.evh = h
}
func (s *SrvSessionBuilder) FindSession(w http.ResponseWriter, r *http.Request) Session {
	c, err := r.Cookie(s.CookieKey)
	ncookie := func() {
		c = &http.Cookie{}
		c.Name = s.CookieKey
		c.Value = util.UUID()
		c.Path = s.Path
		c.Domain = s.Domain
		c.MaxAge = 10 * 24 * 60 * 60
		//
		session := NewSrvSession()
		session.token = c.Value
		session.Flush()
		//
		// s.ks_lck.Lock()
		s.ks[c.Value] = session
		// s.ks_lck.Unlock()
		http.SetCookie(w, c)
		s.evh.OnCreate(session)
		// s.log("setting cookie %v=%v to %v", c.Name, c.Value, r.Host)
	}
	s.ks_lck.Lock()
	defer s.ks_lck.Unlock()
	if err != nil {
		ncookie()
	}
	if _, ok := s.ks[c.Value]; !ok { //if not found,reset cookie
		ncookie()
	}
	ss := s.ks[c.Value]
	ss.Flush()
	return ss
}

func (s *SrvSessionBuilder) Session(token string) Session {
	s.ks_lck.RLock()
	defer s.ks_lck.RUnlock()
	if v, ok := s.ks[token]; ok {
		return v
	} else {
		return nil
	}
}

//
func (s *SrvSessionBuilder) StartLoop() {
	s.looping = true
	go s.Loop()
}
func (s *SrvSessionBuilder) StopLoop() {
	s.looping = false
}

//
func (s *SrvSessionBuilder) Loop() {
	for s.looping {
		ary := []string{}
		now := util.Timestamp(time.Now())
		s.ks_lck.RLock()
		for k, v := range s.ks {
			delay := now - v.begin
			if delay > s.Timeout {
				s.evh.OnTimeout(v)
				ary = append(ary, k)
			}
		}
		s.ks_lck.RUnlock()
		if len(ary) > 0 {
			s.log("looping session time out,removing (%v)", ary)
		}
		s.ks_lck.Lock()
		for _, v := range ary {
			delete(s.ks, v)
		}
		s.ks_lck.Unlock()
		time.Sleep(s.CDelay * time.Millisecond)
	}
}

func (s *SrvSessionBuilder) Clear() {
	s.ks_lck.Lock()
	for k, v := range s.ks {
		s.evh.OnTimeout(v)
		delete(s.ks, k)
	}
	s.ks_lck.Unlock()
}
