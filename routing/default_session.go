package routing

import (
	"net/http"
	"sync"
)

type DefaultSession struct {
	kvs     map[string]interface{}
	kvs_lck sync.RWMutex
}

func (s DefaultSession) Val(key string) interface{} {
	s.kvs_lck.RLock()
	defer s.kvs_lck.RUnlock()
	if v, ok := s.kvs[key]; ok {
		return v
	} else {
		return nil
	}
}
func (s DefaultSession) Set(key string, val interface{}) {
	s.kvs_lck.Lock()
	defer s.kvs_lck.Unlock()
	if val == nil {
		delete(s.kvs, key)
	} else {
		s.kvs[key] = val
	}
}
func (s DefaultSession) Flush() error {
	return nil
}

//
type DefaultSessionBuilder struct {
}

func NewDefaultSessionBuilder() *DefaultSessionBuilder {
	return &DefaultSessionBuilder{}
}
func (s *DefaultSessionBuilder) FindSession(w http.ResponseWriter, r *http.Request) Session {
	return &DefaultSession{kvs: map[string]interface{}{}}
}
func (s *DefaultSessionBuilder) SetEvH(h SessionEvHandler) {
}
