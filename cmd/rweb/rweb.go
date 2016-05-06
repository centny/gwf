package main

import (
	"fmt"
	"github.com/Centny/gwf/routing"
	"github.com/Centny/gwf/routing/filter"
	"github.com/Centny/gwf/util"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
)

func main() {
	var addr = ":80"
	_, args, _ := util.Args()
	var paddr string
	var tpl, preg []string
	if args.Exist("h") {
		fmt.Printf(`Usage: rweb <options>
	-h		show help
	-addr <addr>				the listen addr
	-paddr <proxy address>	 	the proxy address
	-preg <proxy path regex>	the proxy path regex
	-tpl <template path regex>	the template html path regex
	-T<key> http://xxxx/		the tpl data url`)
		os.Exit(1)
		return
	}
	var err = args.ValidF(`
		addr,O|S,L:0;
		paddr,O|S,L:0;
		preg,O|S,L:0;
		tpl,O|S,L:0;
		`, &addr, &paddr, &preg, &tpl)
	if err != nil {
		fmt.Println("check value fail->", err)
		os.Exit(1)
		return
	}
	fmt.Println("running on", addr)
	mux := routing.NewSessionMux2("")
	if len(paddr) > 0 {
		var burl, err = url.Parse(paddr)
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
		for _, reg := range preg {
			mux.Handler(reg, proxy)
		}
	}
	if len(tpl) > 0 {
		var rn = filter.NewRenderNamedF()
		for key, _ := range args {
			if !strings.HasPrefix(key, "T") {
				continue
			}
			url := args.StrVal2(key)
			if len(url) < 1 {
				continue
			}
			web := filter.NewRenderWebData(url)
			keys := strings.SplitN(strings.TrimPrefix(key, "T"), "=", 2)
			if len(keys) > 1 {
				web.Path = keys[1]
			}
			rn.AddDataH(keys[0], web)
		}
		var rd = filter.NewRender(".", rn)
		for _, t := range tpl {
			mux.H(t, rd)
		}
	}
	mux.HFilter("^.*$", filter.NewP3P2())
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
