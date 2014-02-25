package routing

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"
)

type HResult int

const (
	HRES_CONTINUE HResult = iota
	HRES_RETURN
)

type SessionBuilder interface {
	FindSession(w http.ResponseWriter, r *http.Request) Session
}
type Session interface {
	Val(key string) interface{}
	Set(key string, val interface{})
	Flush() error
}
type HandleFunc func(*HTTPSession) HResult
type Handler interface {
	SrvHTTP(*HTTPSession) HResult
}

type HTTPSession struct {
	W http.ResponseWriter
	R *http.Request
	S Session
}

type SessionMux struct {
	Pre string
	//
	Sb SessionBuilder
	//
	Filters      map[*regexp.Regexp]Handler
	FilterFunc   map[*regexp.Regexp]HandleFunc
	Handlers     map[*regexp.Regexp]Handler
	HandlerFunc  map[*regexp.Regexp]HandleFunc
	NHandlers    map[*regexp.Regexp]http.Handler
	NHandlerFunc map[*regexp.Regexp]http.HandlerFunc
	rs           map[*http.Request]*HTTPSession //request to session
}

func NewSessionMux(pre string, sb SessionBuilder) *SessionMux {
	if sb == nil {
		fmt.Println("session builder is nil")
		return nil
	}
	mux := SessionMux{}
	mux.Pre = pre
	mux.Sb = sb
	mux.Filters = map[*regexp.Regexp]Handler{}
	mux.Handlers = map[*regexp.Regexp]Handler{}
	mux.NHandlers = map[*regexp.Regexp]http.Handler{}
	mux.FilterFunc = map[*regexp.Regexp]HandleFunc{}
	mux.HandlerFunc = map[*regexp.Regexp]HandleFunc{}
	mux.NHandlerFunc = map[*regexp.Regexp]http.HandlerFunc{}
	mux.rs = map[*http.Request]*HTTPSession{}
	return &mux
}

func (s *SessionMux) RSession(r *http.Request) *HTTPSession {
	if v, ok := s.rs[r]; ok {
		return v
	} else {
		return nil
	}
}
func (s *SessionMux) HFilter(pattern string, h Handler) {
	reg := regexp.MustCompile(pattern)
	s.Filters[reg] = h
}
func (s *SessionMux) HFilterFunc(pattern string, h HandleFunc) {
	reg := regexp.MustCompile(pattern)
	s.FilterFunc[reg] = h
}
func (s *SessionMux) H(pattern string, h Handler) {
	reg := regexp.MustCompile(pattern)
	s.Handlers[reg] = h
}
func (s *SessionMux) HFunc(pattern string, h HandleFunc) {
	reg := regexp.MustCompile(pattern)
	s.HandlerFunc[reg] = h
}
func (s *SessionMux) Handler(pattern string, h http.Handler) {
	reg := regexp.MustCompile(pattern)
	s.NHandlers[reg] = h
}
func (s *SessionMux) HandleFunc(pattern string, h http.HandlerFunc) {
	reg := regexp.MustCompile(pattern)
	s.NHandlerFunc[reg] = h
}

//
func (s *SessionMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var url string = r.URL.String()
	if len(s.Pre) > 0 {
		url = strings.TrimPrefix(r.URL.String(), s.Pre)
		if len(url) == len(r.URL.String()) {
			http.NotFound(w, r)
			return
		}
	}
	session := s.Sb.FindSession(w, r)

	hs := &HTTPSession{
		W: w,
		R: r,
		S: session,
	}
	s.rs[r] = hs
	defer delete(s.rs, r) //remove the http session object.
	//
	var matched bool = false
	//match filter.
	for k, v := range s.Filters {
		if k.MatchString(url) {
			matched = true
			if v.SrvHTTP(hs) == HRES_RETURN {
				return
			}
		}
	}
	for k, v := range s.FilterFunc {
		if k.MatchString(url) {
			matched = true
			if v(hs) == HRES_RETURN {
				return
			}
		}
	}
	//match handle
	for k, v := range s.Handlers {
		if k.MatchString(url) {
			matched = true
			if v.SrvHTTP(hs) == HRES_RETURN {
				return
			}
		}
	}
	for k, v := range s.HandlerFunc {
		if k.MatchString(url) {
			matched = true
			if v(hs) == HRES_RETURN {
				return
			}
		}
	}
	//match normal handle
	for k, v := range s.NHandlers {
		if k.MatchString(url) {
			matched = true
			v.ServeHTTP(w, r)
		}
	}
	for k, v := range s.NHandlerFunc {
		if k.MatchString(url) {
			matched = true
			v(w, r)
		}
	}
	//
	if matched { //if not matched
		session.Flush()
	} else {
		http.NotFound(w, r)
	}
}
