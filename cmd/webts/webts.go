package main

import (
	"github.com/Centny/gwf/routing"
	"net/http"
)

var WWW string = "."

func Hand(pre string, mux *routing.SessionMux) {
	mux.HFunc("^/g_args(\\?.*)?$", GetArgs)
	mux.HFunc("^/p_args(\\?.*)?$", PostArgs)
	mux.HFunc("^/m_args(\\?.*)?$", MultipartArgs)
	mux.HFunc("^/s_ss(\\?.*)?$", SetSs)
	mux.HFunc("^/g_ss(\\?.*)?$", GetSs)
	mux.HFunc("^/upload(\\?.*)?$", Upload)
	mux.HFunc("^/body(\\?.*)?$", Body)
	mux.HFunc("^/req_ctype(\\?.*)?$", ReqCType)
	mux.HFunc("^/res_ctype(\\?.*)?$", ResCType)
	mux.HFunc("^/echo(\\?.*)?$", Echo)
	mux.Handler("^.*$", http.FileServer(http.Dir(WWW)))
}
