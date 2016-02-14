package main

import (
	"fmt"
	"github.com/Centny/gwf/routing"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
)

func main() {
	addr := ":80"
	if len(os.Args) > 1 {
		if os.Args[1] == "-h" {
			fmt.Println("Usage: rweb <addr> <proxy addres> <proxy path regex>")
			return
		}
		addr = os.Args[1]
	}
	fmt.Println("running on", addr)
	mux := routing.NewSessionMux2("")
	if len(os.Args) > 2 {
		var burl, err = url.Parse(os.Args[2])
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		var proxy = httputil.NewSingleHostReverseProxy(burl)
		var proxy_d = proxy.Director
		proxy.Director = func(r *http.Request) {
			r.Host = burl.Host
			proxy_d(r)
		}
		for _, arg := range os.Args[3:] {
			mux.Handler(arg, proxy)
		}

	}
	mux.HFunc("^/_echo_.*$", func(hs *routing.HTTPSession) routing.HResult {
		hs.R.ParseForm()
		fmt.Println("---Header---")
		for k, v := range hs.R.Header {
			fmt.Println(k, "\t", v)
		}
		fmt.Println("---Form---")
		for k, v := range hs.R.Form {
			fmt.Println(k, "\t", v)
		}
		fmt.Println("---PostForm---")
		for k, v := range hs.R.PostForm {
			fmt.Println(k, "\t", v)
		}
		hs.W.Write([]byte("OK"))
		return routing.HRES_RETURN
	})
	mux.HFilterFunc("^.*$", MicroMessengerFilter)
	mux.HFilterFunc("^.*\\.apk$", func(hs *routing.HTTPSession) routing.HResult {
		hs.W.Header().Set("Content-Type", "application/vnd.android.package-archive")
		return routing.HRES_CONTINUE
	})
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
