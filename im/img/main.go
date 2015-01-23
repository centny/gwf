package main

import (
	"github.com/Centny/gwf/im"
	"github.com/Centny/gwf/netw"
	"github.com/Centny/gwf/netw/impl"
	"github.com/Centny/gwf/pool"
	"github.com/Centny/gwf/routing"
	"net/http"
)

var db = im.NewMemDbH()

func main() {
	im.ShowLog = true
	netw.ShowLog = true
	p := pool.NewBytePool(8, 1024)
	go db.GrpBuilder()
	l := im.NewListner(db, "S-vv-1", p, 9891,
		impl.Json_V2B, impl.Json_B2V, impl.Json_ND, impl.Json_NAV, impl.Json_VNA)
	err := l.Run()
	if err != nil {
		panic(err.Error())
	}
	mux := routing.NewSessionMux2("")
	mux.HFunc("/listSrv", ListSrv)
	mux.Handler("/ws", l.WsH())
	mux.Handler("^.*$", http.FileServer(http.Dir("www")))
	http.Handle("/", mux)
	http.ListenAndServe(":9892", nil)
	l.Wait()
}

func ListSrv(hs *routing.HTTPSession) routing.HResult {
	srv, _ := db.ListSrv("")
	return hs.MsgRes(srv)
}
