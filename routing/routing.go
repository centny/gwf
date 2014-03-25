package routing

import (
	"fmt"
	"github.com/Centny/Cny4go/log"
	"github.com/Centny/Cny4go/util"
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
	Kvs map[string]interface{}
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
		return fmt.Sprintf("%v", v)
	}
}
func (h *HTTPSession) CheckVal(key string) string {
	v := h.RVal(key)
	if len(v) > 0 {
		return v
	}
	return h.StrVal(key)
}
func (h *HTTPSession) RVal(key string) string {
	v := h.R.FormValue(key)
	if len(v) > 0 {
		return v
	}
	v = h.R.PostFormValue(key)
	return v
}

//valid all value by format,limit require.
func (h *HTTPSession) ValidVal(f string, args ...interface{}) error {
	return util.ValidAttrF(f, h.CheckVal, true, args...)
}

//valid all value by format,not limit require.
func (h *HTTPSession) ValidValN(f string, args ...interface{}) error {
	return util.ValidAttrF(f, h.CheckVal, false, args...)
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
	regex_f_ary  []*regexp.Regexp
	regex_f      map[*regexp.Regexp]int
	regex_h_ary  []*regexp.Regexp
	regex_h      map[*regexp.Regexp]int
	rs           map[*http.Request]*HTTPSession //request to session
	Kvs          map[string]interface{}
	FilterEnable bool
	HandleEnable bool
	ShowLog      bool
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
	mux.regex_f = map[*regexp.Regexp]int{}
	mux.regex_f_ary = []*regexp.Regexp{}
	mux.regex_h = map[*regexp.Regexp]int{}
	mux.regex_h_ary = []*regexp.Regexp{}
	mux.rs = map[*http.Request]*HTTPSession{}
	mux.Kvs = map[string]interface{}{}
	mux.FilterEnable = true
	mux.HandleEnable = true
	mux.ShowLog = false
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
	s.regex_f[reg] = 1
	s.regex_f_ary = append(s.regex_f_ary, reg)
}
func (s *SessionMux) HFilterFunc(pattern string, h HandleFunc) {
	reg := regexp.MustCompile(pattern)
	s.FilterFunc[reg] = h
	s.regex_f[reg] = 2
	s.regex_f_ary = append(s.regex_f_ary, reg)
}
func (s *SessionMux) H(pattern string, h Handler) {
	reg := regexp.MustCompile(pattern)
	s.Handlers[reg] = h
	s.regex_h[reg] = 1
	s.regex_h_ary = append(s.regex_h_ary, reg)
}
func (s *SessionMux) HFunc(pattern string, h HandleFunc) {
	reg := regexp.MustCompile(pattern)
	s.HandlerFunc[reg] = h
	s.regex_h[reg] = 2
	s.regex_h_ary = append(s.regex_h_ary, reg)
}
func (s *SessionMux) Handler(pattern string, h http.Handler) {
	reg := regexp.MustCompile(pattern)
	s.NHandlers[reg] = h
	s.regex_h[reg] = 3
	s.regex_h_ary = append(s.regex_h_ary, reg)
}
func (s *SessionMux) HandleFunc(pattern string, h http.HandlerFunc) {
	reg := regexp.MustCompile(pattern)
	s.NHandlerFunc[reg] = h
	s.regex_h[reg] = 4
	s.regex_h_ary = append(s.regex_h_ary, reg)
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
		Kvs: map[string]interface{}{},
	}
	s.rs[r] = hs
	defer delete(s.rs, r) //remove the http session object.
	//
	var matched bool = false
	//
	defer func() {
		if !matched { //if not matched
			http.NotFound(w, r)
		}
		if s.ShowLog {
			log.D("URL(%s),found(%v)", r.URL.String(), matched)
		}
	}()
	//match filter.
	if s.FilterEnable {
		for _, k := range s.regex_f_ary {
			if k.MatchString(url) {
				matched = true
				switch s.regex_f[k] {
				case 1:
					rv := s.Filters[k]
					if rv.SrvHTTP(hs) == HRES_RETURN {
						return
					}
				case 2:
					rv := s.FilterFunc[k]
					if rv(hs) == HRES_RETURN {
						return
					}
				}
			}
		}
	}
	//match handle
	if s.HandleEnable {
		for _, k := range s.regex_h_ary {
			if k.MatchString(url) {
				matched = true
				switch s.regex_h[k] {
				case 1:
					rv := s.Handlers[k]
					if rv.SrvHTTP(hs) == HRES_RETURN {
						return
					}
				case 2:
					rv := s.HandlerFunc[k]
					if rv(hs) == HRES_RETURN {
						return
					}
				case 3:
					rv := s.NHandlers[k]
					rv.ServeHTTP(w, r)
				case 4:
					rv := s.NHandlerFunc[k]
					rv(w, r)
				}
			}
		}
	}
}
