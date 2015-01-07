package hutil

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"testing"
)

func TestRh(t *testing.T) {
	rr()
}

func rr() {
	remote, err := url.Parse("http://192.168.1.2:8080")
	if err != nil {
		panic(err)
	}

	proxy := httputil.NewSingleHostReverseProxy(remote)
	http.HandleFunc("/app2", handler(proxy))
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}

func handler(p *httputil.ReverseProxy) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.URL)
		r.URL.Path = "/"
		p.ServeHTTP(w, r)
	}
}
