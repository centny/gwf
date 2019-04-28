package routing

import (
	"net/http"
	"sync"

	"github.com/Centny/gwf/util"
)

type DefaultSession struct {
	id      string
	kvs     map[string]interface{}
	kvs_lck sync.RWMutex
}

func (d *DefaultSession) ID() string {
	return d.id
}

func (s *DefaultSession) Clear() error {
	s.kvs_lck.Lock()
	defer s.kvs_lck.Unlock()
	s.kvs = map[string]interface{}{}
	return nil
}

func (s *DefaultSession) Val(key string) interface{} {
	s.kvs_lck.RLock()
	defer s.kvs_lck.RUnlock()
	if v, ok := s.kvs[key]; ok {
		return v
	} else {
		return nil
	}
}
func (s *DefaultSession) Set(key string, val interface{}) {
	s.kvs_lck.Lock()
	defer s.kvs_lck.Unlock()
	if val == nil {
		delete(s.kvs, key)
	} else {
		s.kvs[key] = val
	}
}
func (s *DefaultSession) Flush() error {
	return nil
}

//
type DefaultSessionBuilder struct {
}

func NewDefaultSessionBuilder() *DefaultSessionBuilder {
	return &DefaultSessionBuilder{}
}
func (s *DefaultSessionBuilder) FindSession(w http.ResponseWriter, r *http.Request) Session {
	return &DefaultSession{id: util.UUID(), kvs: map[string]interface{}{}}
}
func (s *DefaultSessionBuilder) SetEvH(h SessionEvHandler) {
}
