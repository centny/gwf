package httptest

import (
	"fmt"
	"github.com/Centny/Cny4go/routing"
	"github.com/Centny/Cny4go/util"
	"net/http"
	"net/http/httptest"
)

type Server struct {
	URL string
	S   *httptest.Server
	mux *routing.SessionMux
}

func (s *Server) Close() {
	s.S.Close()
}

func (s *Server) G(f string, args ...interface{}) (string, error) {
	return util.HGet(fmt.Sprintf("%v%v", s.URL, f), args...)
}

func (s *Server) G2(f string, args ...interface{}) (util.Map, error) {
	return util.HGet2(fmt.Sprintf("%v%v", s.URL, f), args...)
}

func (s *Server) P(url string, fields map[string]string) (string, error) {
	return util.HPost(fmt.Sprintf("%v%v", s.URL, url), fields)
}

func (s *Server) P2(url string, fields map[string]string) (util.Map, error) {
	return util.HPost2(fmt.Sprintf("%v%v", s.URL, url), fields)
}

func NewServer(f routing.HandleFunc) *Server {
	sb := routing.NewSrvSessionBuilder("", "/", "tsrv", 60000, 200)
	mux := routing.NewSessionMux("", sb)
	mux.HFunc("^.*$", f)
	return NewMuxServer(mux)
}
func NewServer2(h routing.Handler) *Server {
	sb := routing.NewSrvSessionBuilder("", "/", "tsrv", 60000, 200)
	mux := routing.NewSessionMux("", sb)
	mux.H("^.*$", h)
	return NewMuxServer(mux)
}
func NewMuxServer(mux *routing.SessionMux) *Server {
	srv := &Server{mux: mux}
	srv.S = httptest.NewServer(mux)
	srv.URL = srv.S.URL
	return srv
}

//test normal handler
func Tnh(h http.Handler, f string, args ...interface{}) error {
	ts := httptest.NewServer(h)
	_, err := util.HGet(fmt.Sprintf("%v%v", ts.URL, f), args...)
	return err
}

//test normal handler function
func Tnf(h func(http.ResponseWriter, *http.Request), f string, args ...interface{}) error {
	return Tnh(http.HandlerFunc(h), f, args...)
}

func Th(h routing.Handler, f string, args ...interface{}) error {
	ts := NewServer2(h)
	_, err := ts.G(f, args...)
	return err
}

func Tf(h routing.HandleFunc, f string, args ...interface{}) error {
	ts := NewServer(h)
	_, err := ts.G(f, args...)
	return err
}
