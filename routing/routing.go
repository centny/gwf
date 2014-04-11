package routing

import (
	"errors"
	"fmt"
	"github.com/Centny/Cny4go/log"
	"github.com/Centny/Cny4go/util"
	"io"
	"net/http"
	"os"
	"reflect"
	"regexp"
	"strings"
)

type HResult int

const (
	HRES_CONTINUE HResult = iota
	HRES_RETURN
)

func (h HResult) String() string {
	if h == HRES_CONTINUE {
		return "CONTINUE"
	} else {
		return "RETURN"
	}
}

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

func (h *HTTPSession) SetCookie(key string, val string) {
	cookie := &http.Cookie{}
	cookie.Name = key
	cookie.Domain = h.Mux.Domain
	cookie.Path = h.Mux.Path
	cookie.Value = val
	cookie.MaxAge = 0
	http.SetCookie(h.W, cookie)
}

func (h *HTTPSession) Cookie(key string) string {
	c, err := h.R.Cookie(key)
	if c == nil || err != nil {
		return ""
	}
	return c.Value
}

func (h *HTTPSession) Redirect(url string) {
	http.Redirect(h.W, h.R, url, http.StatusMovedPermanently)
}

func (h *HTTPSession) SetVal(key string, val interface{}) {
	h.S.Set(key, val)
}

func (h *HTTPSession) Val(key string) interface{} {
	return h.S.Val(key)
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
func (h *HTTPSession) FormFInfo(name string) (int64, string, error) {
	src, fh, err := h.R.FormFile("file")
	if err != nil {
		return 0, "", err
	}
	fsize := util.FormFSzie(src)
	if fsize < 1 {
		return 0, "", errors.New("file size error")
	}
	return fsize, fh.Filename, nil
}

func (h *HTTPSession) RecF(name, tfile string) (int64, error) {
	src, _, err := h.R.FormFile("file")
	if err != nil {
		return 0, err
	}
	err = util.FTouch(tfile)
	if err != nil {
		return 0, err
	}
	dst, _ := os.OpenFile(tfile, os.O_RDWR|os.O_APPEND, 0)
	csize, err := io.Copy(dst, src)
	dst.Close()
	if err != nil {
		os.Remove(tfile)
		return 0, err
	}
	return csize, nil
}

func (h *HTTPSession) SendF(fname, tfile, ctype string, attach bool) {
	SendF(h.W, h.R, fname, tfile, ctype, attach)
}

//valid require value by format,limit require.
func (h *HTTPSession) ValidRVal(f string, args ...interface{}) error {
	return util.ValidAttrF(f, h.RVal, true, args...)
}

//valid require value by format,not limit require.
func (h *HTTPSession) ValidRValN(f string, args ...interface{}) error {
	return util.ValidAttrF(f, h.RVal, false, args...)
}

//valid all value by format,limit require.
func (h *HTTPSession) ValidCheckVal(f string, args ...interface{}) error {
	return util.ValidAttrF(f, h.CheckVal, true, args...)
}

//valid all value by format,not limit require.
func (h *HTTPSession) ValidCheckValN(f string, args ...interface{}) error {
	return util.ValidAttrF(f, h.CheckVal, false, args...)
}

type SessionMux struct {
	Pre    string
	Domain string
	Path   string
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
	regex_m      map[*regexp.Regexp]string
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
	mux.Domain = ""
	mux.Path = "/"
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
	mux.regex_m = map[*regexp.Regexp]string{}
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
	s.HFilterM(pattern, h, "*")
}
func (s *SessionMux) HFilterM(pattern string, h Handler, m string) {
	reg := regexp.MustCompile(pattern)
	s.Filters[reg] = h
	s.regex_f[reg] = 1
	s.regex_f_ary = append(s.regex_f_ary, reg)
	s.regex_m[reg] = m
}
func (s *SessionMux) HFilterFunc(pattern string, h HandleFunc) {
	s.HFilterFuncM(pattern, h, "*")
}
func (s *SessionMux) HFilterFuncM(pattern string, h HandleFunc, m string) {
	reg := regexp.MustCompile(pattern)
	s.FilterFunc[reg] = h
	s.regex_f[reg] = 2
	s.regex_f_ary = append(s.regex_f_ary, reg)
	s.regex_m[reg] = m
}
func (s *SessionMux) H(pattern string, h Handler) {
	s.HM(pattern, h, "*")
}
func (s *SessionMux) HM(pattern string, h Handler, m string) {
	reg := regexp.MustCompile(pattern)
	s.Handlers[reg] = h
	s.regex_h[reg] = 1
	s.regex_h_ary = append(s.regex_h_ary, reg)
	s.regex_m[reg] = m
}
func (s *SessionMux) HFunc(pattern string, h HandleFunc) {
	s.HFuncM(pattern, h, "*")
}
func (s *SessionMux) HFuncM(pattern string, h HandleFunc, m string) {
	reg := regexp.MustCompile(pattern)
	s.HandlerFunc[reg] = h
	s.regex_h[reg] = 2
	s.regex_h_ary = append(s.regex_h_ary, reg)
	s.regex_m[reg] = m
}
func (s *SessionMux) Handler(pattern string, h http.Handler) {
	s.HandlerM(pattern, h, "*", true)
}
func (s *SessionMux) HandlerM(pattern string, h http.Handler, m string, ret bool) {
	reg := regexp.MustCompile(pattern)
	s.NHandlers[reg] = h
	s.regex_h[reg] = 3
	s.regex_h_ary = append(s.regex_h_ary, reg)
	if ret {
		m = fmt.Sprintf("%s,:RETURN", m)
	} else {
		m = fmt.Sprintf("%s,:CONTINUE", m)
	}
	s.regex_m[reg] = m
}
func (s *SessionMux) HandleFunc(pattern string, h http.HandlerFunc) {
	s.HandleFuncM(pattern, h, "*", true)
}
func (s *SessionMux) HandleFuncM(pattern string, h http.HandlerFunc, m string, ret bool) {
	reg := regexp.MustCompile(pattern)
	s.NHandlerFunc[reg] = h
	s.regex_h[reg] = 4
	s.regex_h_ary = append(s.regex_h_ary, reg)
	if ret {
		m = fmt.Sprintf("%s,:RETURN", m)
	} else {
		m = fmt.Sprintf("%s,:CONTINUE", m)
	}
	s.regex_m[reg] = m
}

func (s *SessionMux) slog(fmt string, args ...interface{}) {
	if s.ShowLog {
		log.D(fmt, args...)
	}
}
func (s *SessionMux) check_m(reg *regexp.Regexp, m string) bool {
	if tm, ok := s.regex_m[reg]; ok {
		return strings.Contains(tm, "*") || strings.Contains(tm, m)
	}
	return false
}

func (s *SessionMux) check_continue(reg *regexp.Regexp) bool {
	if tm, ok := s.regex_m[reg]; ok {
		fmt.Println(tm)
		return strings.Contains(tm, ":CONTINUE")
	}
	return false
}

func (s *SessionMux) exec_f(hs *HTTPSession) (bool, HResult) {
	url := hs.R.URL.Path
	var matched bool = false
	for _, k := range s.regex_f_ary {
		if !k.MatchString(url) {
			continue
		}
		if !s.check_m(k, hs.R.Method) {
			s.slog("not mathced method %v to %v", hs.R.Method, s.regex_m[k])
			continue
		}
		matched = true
		switch s.regex_f[k] {
		case 1:
			rv := s.Filters[k]
			res := rv.SrvHTTP(hs)
			s.slog("mathced filter %v to %v (%v)", k, hs.R.URL.Path, res.String())
			if res == HRES_RETURN {
				return matched, res
			}
		case 2:
			rv := s.FilterFunc[k]
			res := rv(hs)
			s.slog("mathced filter func %v to %v (%v)", k, hs.R.URL.Path, res.String())
			if res == HRES_RETURN {
				return matched, res
			}
		}
	}
	return matched, HRES_CONTINUE
}
func (s *SessionMux) exec_h(hs *HTTPSession) (bool, HResult) {
	url := hs.R.URL.Path
	var matched bool = false
	for _, k := range s.regex_h_ary {
		if !k.MatchString(url) {
			continue
		}
		if !s.check_m(k, hs.R.Method) {
			s.slog("not mathced method %v to %v", hs.R.Method, s.regex_m[k])
			continue
		}
		matched = true
		switch s.regex_h[k] {
		case 1:
			rv := s.Handlers[k]
			res := rv.SrvHTTP(hs)
			s.slog("mathced handler %v to %v (%v)", k, hs.R.URL.Path, res.String())
			if res == HRES_RETURN {
				return matched, res
			}
		case 2:
			rv := s.HandlerFunc[k]
			res := rv(hs)
			s.slog("mathced handler func %v to %v (%v)", k, hs.R.URL.Path, res.String())
			if res == HRES_RETURN {
				return matched, res
			}
		case 3:
			rv := s.NHandlers[k]
			rv.ServeHTTP(hs.W, hs.R)
			if s.check_continue(k) {
				s.slog("mathced normal handler %v to %v (%v)",
					k, hs.R.URL.Path, HRES_CONTINUE.String())
				continue
			} else {
				s.slog("mathced normal handler %v to %v (%v)",
					k, hs.R.URL.Path, HRES_RETURN.String())
				return matched, HRES_RETURN
			}
		case 4:
			rv := s.NHandlerFunc[k]
			rv(hs.W, hs.R)
			if s.check_continue(k) {
				s.slog("mathced normal handler func %v to %v (%v)",
					k, hs.R.URL.Path, HRES_CONTINUE.String())
				continue
			} else {
				s.slog("mathced normal handler func %v to %v (%v)", k,
					hs.R.URL.Path, HRES_RETURN.String())
				return matched, HRES_RETURN
			}
		}
	}
	return matched, HRES_CONTINUE
}

//
func (s *SessionMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	r.URL.Path = strings.TrimPrefix(r.URL.Path, s.Pre)
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
			s.slog("not matchd any filter:%s", r.URL.Path)
			http.NotFound(w, r)
		}
	}()
	//match filter.
	if s.FilterEnable {
		mrv, res := s.exec_f(hs)
		matched = mrv
		if res == HRES_RETURN {
			return
		}
	}
	//match handle
	if s.HandleEnable {
		mrv, res := s.exec_h(hs)
		matched = matched || mrv
		if res == HRES_RETURN {
			return
		}
	}
}
