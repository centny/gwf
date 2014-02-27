package routing

import (
	"fmt"
	"net/http"
	"reflect"
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
	W   http.ResponseWriter
	R   *http.Request
	S   Session
	Mux *SessionMux
}

func (h *HTTPSession) Redirect(url string) {
	http.Redirect(h.W, h.R, url, http.StatusMovedPermanently)
}
func (h *HTTPSession) SetVal(key string, val interface{}) {
	h.S.Set(key, val)
}
func (h *HTTPSession) UintVal(key string) uint64 {
	v := h.S.Val(key)
	if v == nil {
		return 0
	}
	ty := reflect.TypeOf(v)
	switch ty.Kind() {
	case reflect.Uint:
		return uint64(v.(uint))
	case reflect.Uint8:
		return uint64(v.(uint8))
	case reflect.Uint16:
		return uint64(v.(uint16))
	case reflect.Uint32:
		return uint64(v.(uint32))
	case reflect.Uint64:
		return v.(uint64)
	default:
		return 0
	}
}
func (h *HTTPSession) IntVal(key string) int64 {
	v := h.S.Val(key)
	if v == nil {
		return 0
	}
	ty := reflect.TypeOf(v)
	switch ty.Kind() {
	case reflect.Int:
		return int64(v.(int))
	case reflect.Int8:
		return int64(v.(int8))
	case reflect.Int16:
		return int64(v.(int16))
	case reflect.Int32:
		return int64(v.(int32))
	case reflect.Int64:
		return v.(int64)
	default:
		return 0
	}
}
func (h *HTTPSession) FloatVal(key string) float64 {
	v := h.S.Val(key)
	if v == nil {
		return 0
	}
	ty := reflect.TypeOf(v)
	switch ty.Kind() {
	case reflect.Float32:
		return float64(v.(float32))
	case reflect.Float64:
		return v.(float64)
	default:
		return 0
	}
}
func (h *HTTPSession) StrVal(key string) string {
	v := h.S.Val(key)
	if v == nil {
		return ""
	}
	ty := reflect.TypeOf(v)
	switch ty.Kind() {
	case reflect.String:
		return v.(string)
	default:
		return ""
	}
}
func (h *HTTPSession) CheckVal(key string) string {
	v := h.StrVal(key)
	if len(v) > 0 {
		return v
	}
	v = h.R.FormValue(key)
	if len(v) > 0 {
		return v
	}
	v = h.R.PostFormValue(key)
	return v
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
	//
	Kvs map[string]interface{}
}

func NewSessionMux2(pre string) *SessionMux {
	return NewSessionMux(pre, NewDefaultSessionBuilder())
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
	mux.Kvs = map[string]interface{}{}
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
	r.URL.Path = strings.TrimPrefix(r.URL.Path, s.Pre)
	url := r.URL.Path
	session := s.Sb.FindSession(w, r)
	hs := &HTTPSession{
		W:   w,
		R:   r,
		S:   session,
		Mux: s,
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
	if !matched { //if not matched
		http.NotFound(w, r)
	}
}
