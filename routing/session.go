package routing

import (
	"code.google.com/p/go-uuid/uuid"
	"github.com/Centny/Cny4go/log"
	"github.com/Centny/Cny4go/util"
	"net/http"
	"sync"
	"time"
)

type SrvSession struct {
	token string
	begin int64
	kvs   map[string]interface{}
}

func (s *SrvSession) Val(key string) interface{} {
	if v, ok := s.kvs[key]; ok {
		return v
	} else {
		return nil
	}
}
func (s *SrvSession) Set(key string, val interface{}) {
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
		c.Value = uuid.New()
		c.Path = s.Path
		c.Domain = s.Domain
		c.MaxAge = 0
		//
		session := &SrvSession{}
		session.token = c.Value
		session.kvs = map[string]interface{}{}
		session.Flush()
		//
		s.ks_lck.RLock()
		s.ks[c.Value] = session
		s.ks_lck.RUnlock()
		http.SetCookie(w, c)
		s.evh.OnCreate(session)
		// s.log("setting cookie %v=%v to %v", c.Name, c.Value, r.Host)
	}
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
		for k, v := range s.ks {
			delay := now - v.begin
			if delay > s.Timeout {
				s.evh.OnTimeout(v)
				ary = append(ary, k)
			}
		}
		if len(ary) > 0 {
			s.log("looping session time out,removing (%v)", ary)
		}
		s.ks_lck.RLock()
		for _, v := range ary {
			delete(s.ks, v)
		}
		s.ks_lck.RUnlock()
		time.Sleep(s.CDelay * time.Millisecond)
	}
}

func (s *SrvSessionBuilder) Clear() {
	s.ks_lck.RLock()
	for k, _ := range s.ks {
		delete(s.ks, k)
	}
	s.ks_lck.RUnlock()
}
