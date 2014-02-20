package srvmux

import (
	"net/http"
	"regexp"
	"strings"
)

type SrvMux struct {
	hdls  map[string]http.Handler
	hfuns map[string]http.HandlerFunc
	pre   string
}

func NewSrvMux(pre string) *SrvMux {
	mux := SrvMux{}
	mux.pre = pre
	mux.hdls = make(map[string]http.Handler)
	mux.hfuns = make(map[string]http.HandlerFunc)
	return &mux
}
func (s *SrvMux) Handler(pattern string, handler http.Handler) {
	s.hdls[pattern] = handler
}

func (s *SrvMux) HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	s.hfuns[pattern] = handler
}

func (s *SrvMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	url := strings.TrimPrefix(r.URL.String(), s.pre)
	if len(url) == len(r.URL.String()) {
		http.NotFound(w, r)
		return
	}
	if v, ok := s.hdls[url]; ok {
		v.ServeHTTP(w, r)
		return
	}
	if v, ok := s.hfuns[url]; ok {
		v(w, r)
		return
	}
	http.NotFound(w, r)
}

type RegMux struct {
	hdls  map[*regexp.Regexp]http.Handler
	hfuns map[*regexp.Regexp]http.HandlerFunc
	pre   string
}

func NewRegMux(pre string) *RegMux {
	mux := RegMux{}
	mux.pre = pre
	mux.hdls = make(map[*regexp.Regexp]http.Handler)
	mux.hfuns = make(map[*regexp.Regexp]http.HandlerFunc)
	return &mux
}
func (s *RegMux) Handler(pattern string, handler http.Handler) {
	reg := regexp.MustCompile(pattern)
	s.hdls[reg] = handler
}

func (s *RegMux) HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	reg := regexp.MustCompile(pattern)
	s.hfuns[reg] = handler
}

func (s *RegMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	url := strings.TrimPrefix(r.URL.String(), s.pre)
	if len(url) == len(r.URL.String()) {
		http.NotFound(w, r)
		return
	}
	for k, v := range s.hdls {
		if k.MatchString(url) {
			v.ServeHTTP(w, r)
			return
		}
	}
	for k, v := range s.hfuns {
		if k.MatchString(url) {
			v(w, r)
			return
		}
	}
	http.NotFound(w, r)
}
