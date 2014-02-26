package routing

import (
	"code.google.com/p/go-uuid/uuid"
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
	Domain  string
	Path    string
	Timeout int64
	CDelay  time.Duration
	//
	looping bool
	ks      map[string]*SrvSession //key session
	ks_lck  sync.RWMutex
}

func NewSrvSessionBuilder(domain string, path string, timeout int64, cdelay time.Duration) *SrvSessionBuilder {
	sb := SrvSessionBuilder{}
	sb.Domain = domain
	sb.Path = path
	sb.Timeout = timeout
	sb.CDelay = cdelay
	sb.ks = map[string]*SrvSession{}
	return &sb
}
func (s *SrvSessionBuilder) FindSession(w http.ResponseWriter, r *http.Request) Session {
	c, err := r.Cookie("token")
	ncookie := func() {
		c = &http.Cookie{}
		c.Name = "token"
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
				ary = append(ary, k)
			}
		}
		s.ks_lck.RLock()
		for _, v := range ary {
			delete(s.ks, v)
		}
		s.ks_lck.RUnlock()
		time.Sleep(s.CDelay * time.Millisecond)
	}
}