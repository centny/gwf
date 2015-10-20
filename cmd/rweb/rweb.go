package main

import (
	"fmt"
	"github.com/Centny/gwf/routing"
	"net/http"
	"os"
	"strings"
)

func main() {
	addr := ":80"
	if len(os.Args) > 1 {
		addr = os.Args[1]
	}
	fmt.Println("running on", addr)
	mux := routing.NewSessionMux2("")
	mux.HFilterFunc("^.*$", MicroMessengerFilter)
	mux.Handler("^/.*$", http.FileServer(http.Dir(".")))
	fmt.Println(http.ListenAndServe(addr, mux))
}

func MicroMessengerFilter(h *routing.HTTPSession) routing.HResult {
	uag := h.R.Header.Get("User-Agent")
	if strings.Index(uag, "MicroMessenger") == -1 {
		return routing.HRES_CONTINUE
	} else {
		h.W.WriteHeader(404)
		return routing.HRES_RETURN
	}
}
