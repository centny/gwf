package main

import (
	"fmt"
	"github.com/Centny/gwf/routing"
	"github.com/Centny/gwf/util"
	"net/http"
	"os"
)

var TCFG = util.Fcfg{}
var W_DIR string = ""

func main() {
	if len(os.Args) < 4 {
		fmt.Println("Usage:cfgs <listen address> <token configure> <configure file base directory>")
		return
	}
	err := TCFG.InitWithFilePath(os.Args[2])
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	W_DIR = os.Args[3]
	mux := routing.NewSessionMux2("")
	mux.HFunc("^/cfg(\\?.*)?$", cfg)
	fmt.Println(http.ListenAndServe(os.Args[1], mux))
}

func cfg(hs *routing.HTTPSession) routing.HResult {
	var token string
	err := hs.ValidCheckVal(`
		token,R|S,L:0;
		`, &token)
	if err != nil {
		hs.W.WriteHeader(400)
		hs.SendT2("token is required\n")
		return routing.HRES_RETURN
	}
	if !TCFG.Exist(token) {
		hs.W.WriteHeader(400)
		hs.SendT2(fmt.Sprintf("token(%v) not found\n", token))
		return routing.HRES_RETURN
	}
	tcfg, err := util.NewFcfg(fmt.Sprintf("%v/%v", W_DIR, TCFG.Val(token)))
	if err == nil {
		hs.SendT2(tcfg.String())
	} else {
		hs.W.WriteHeader(500)
		hs.SendT2(fmt.Sprintf("server err(%v)\n", err.Error()))
	}
	return routing.HRES_RETURN
}
