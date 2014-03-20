package session

import (
	"code.google.com/p/go-uuid/uuid"
	"github.com/Centny/Cny4go/routing"
	"net/http"
)

type RSession struct {
	W   http.ResponseWriter
	R   *http.Request
	Sb  *RSessionBuilder
	CId string
	kvs map[string]interface{}
}

func (c *RSession) Val(key string) interface{} {
	if v, ok := c.kvs[key]; ok {
		return v
	} else {
		return nil
	}
}
func (c *RSession) Set(key string, val interface{}) {
	if val == nil {
		delete(c.kvs, key)
	} else {
		c.kvs[key] = val
	}
}
func (c *RSession) Flush() error {
	return nil
}

//
type RSessionBuilder struct {
	//
	Domain string
	Path   string
	CName  string
	Cached map[string]map[string]interface{}
}

func NewRSessionBuilder(domain string, path string) *RSessionBuilder {
	sb := RSessionBuilder{}
	sb.Domain = domain
	sb.Path = path
	sb.CName = "S"
	sb.Cached = map[string]map[string]interface{}{}
	return &sb
}
func (s *RSessionBuilder) FindSession(w http.ResponseWriter, r *http.Request) routing.Session {
	c, err := r.Cookie(s.CName)
	cs := &RSession{
		W:   w,
		R:   r,
		Sb:  s,
		kvs: map[string]interface{}{},
	}
	if err == nil {
		cs.CId = c.Value
	}
	if len(cs.CId) < 1 {
		cs.CId = uuid.New()
	}
	if v, ok := s.Cached[cs.CId]; ok {
		cs.kvs = v
	} else {
		s.Cached[cs.CId] = cs.kvs
		cookie := &http.Cookie{}
		cookie.Name = s.CName
		cookie.Domain = s.Domain
		cookie.Path = s.Path
		cookie.Value = cs.CId
		cookie.MaxAge = 0
		http.SetCookie(cs.W, cookie)
	}
	return cs
}
