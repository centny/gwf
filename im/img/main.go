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
	netw.ShowLog_C = true
	impl.ShowLog = true
	p := pool.NewBytePool(8, 1024)
	go db.GrpBuilder()
	l := im.NewListner2(db, "S-vv-1", p, 9891)
	l.WsAddr = ":9892/ws"
	err := l.Run()
	if err != nil {
		panic(err.Error())
	}
	mux := routing.NewSessionMux2("")
	mux.HFunc("/listSrv", ListSrv)
	mux.HFunc("/listRs", ListRs)
	mux.Handler("^/ws$", l.WIM_L.WsS())
	mux.Handler("^.*$", http.FileServer(http.Dir("www")))
	http.Handle("/", mux)
	http.ListenAndServe(":9892", nil)
	l.Wait()
}

func ListSrv(hs *routing.HTTPSession) routing.HResult {
	srv, _ := db.ListSrv("")
	return hs.MsgRes(srv)
}
func ListRs(hs *routing.HTTPSession) routing.HResult {
	usr, _ := db.ListR()
	return hs.MsgRes(usr)
}
