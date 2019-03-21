package routing

import (
	"net"
	"net/http"
	"os"
	"strings"
)

var Shared = NewSessionMux2("")

func HFilter(pattern string, h Handler) {
	Shared.HFilter(pattern, h)
}
func HFilterFunc(pattern string, h HandleFunc) {
	Shared.HFilterFunc(pattern, h)
}
func H(pattern string, h Handler) {
	Shared.H(pattern, h)
}
func HFunc(pattern string, h HandleFunc) {
	Shared.HFunc(pattern, h)
}

func ListenAndServe(addr string) error {
	if strings.HasPrefix(addr, "/") {
		unixListener, err := net.Listen("unix", os.Args[1])
		if err != nil {
			return err
		}
		server := &http.Server{Handler: Shared}
		return server.Serve(unixListener)
	}
	return http.ListenAndServe(addr, Shared)
}
