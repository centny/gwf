package routing

import (
	"net/http"
)

type DefaultSession map[string]interface{}

func (s DefaultSession) Val(key string) interface{} {
	if v, ok := s[key]; ok {
		return v
	} else {
		return nil
	}
}
func (s DefaultSession) Set(key string, val interface{}) {
	if val == nil {
		delete(s, key)
	} else {
		s[key] = val
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
	return &DefaultSession{}
}
