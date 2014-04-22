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
	Mux *routing.SessionMux
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

func (s *Server) PostF(url, fkey, fp string, fields map[string]string) (string, error) {
	return util.HPostF(fmt.Sprintf("%v%v", s.URL, url), fields, fkey, fp)
}

func (s *Server) PostF2(url, fkey, fp string, fields map[string]string) (util.Map, error) {
	return util.HPostF2(fmt.Sprintf("%v%v", s.URL, url), fields, fkey, fp)
}

func NewServer(f routing.HandleFunc) *Server {
	sb := routing.NewSrvSessionBuilder("", "/", "tsrv", 60000, 200)
	mux := routing.NewSessionMux("", sb)
	mux.HFunc("^.*$", f)
	return NewMuxServer2(mux)
}
func NewServer2(h routing.Handler) *Server {
	sb := routing.NewSrvSessionBuilder("", "/", "tsrv", 60000, 200)
	mux := routing.NewSessionMux("", sb)
	mux.H("^.*$", h)
	return NewMuxServer2(mux)
}
func NewMuxServer() *Server {
	sb := routing.NewSrvSessionBuilder("", "/", "tsrv", 60000, 200)
	mux := routing.NewSessionMux("", sb)
	srv := &Server{Mux: mux}
	srv.S = httptest.NewServer(mux)
	srv.URL = srv.S.URL
	return srv
}
func NewMuxServer2(mux *routing.SessionMux) *Server {
	srv := &Server{Mux: mux}
	srv.S = httptest.NewServer(mux)
	srv.URL = srv.S.URL
	return srv
}

//test normal handler
func Tnh(h http.Handler, f string, args ...interface{}) error {
	_, err := Tnh2(h, f, args...)
	return err
}
func Tnh2(h http.Handler, f string, args ...interface{}) (string, error) {
	ts := httptest.NewServer(h)
	return util.HGet(fmt.Sprintf("%v%v", ts.URL, f), args...)
}

//test normal handler function a=%v&
func Tnf(h func(http.ResponseWriter, *http.Request), f string, args ...interface{}) error {
	return Tnh(http.HandlerFunc(h), f, args...)
}
func Tnf2(h func(http.ResponseWriter, *http.Request), f string, args ...interface{}) (string, error) {
	return Tnh2(http.HandlerFunc(h), f, args...)
}
func Th(h routing.Handler, f string, args ...interface{}) error {
	_, err := Th2(h, f, args...)
	return err
}
func Th2(h routing.Handler, f string, args ...interface{}) (string, error) {
	ts := NewServer2(h)
	return ts.G(f, args...)
}
func Tf(h routing.HandleFunc, f string, args ...interface{}) error {
	_, err := Tf2(h, f, args...)
	return err
}
func Tf2(h routing.HandleFunc, f string, args ...interface{}) (string, error) {
	ts := NewServer(h)
	return ts.G(f, args...)
}
