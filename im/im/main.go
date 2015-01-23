package main

import (
	"github.com/Centny/gwf/im"
	"github.com/Centny/gwf/netw"
	"github.com/Centny/gwf/netw/impl"
	"github.com/Centny/gwf/pool"
	"net/http"
)

func main() {
	im.ShowLog = true
	netw.ShowLog = true
	db := im.NewMemDbH()
	p := pool.NewBytePool(8, 1024)
	go db.GrpBuilder()
	l := im.NewListner(db, "S-vv-1", p, ":9891",
		impl.Json_V2B, impl.Json_B2V, impl.Json_ND, impl.Json_NAV, impl.Json_VNA)
	err := l.Run()
	if err != nil {
		panic(err.Error())
	}
	http.Handle("/ws", l.WsH())
	http.Handle("/", http.FileServer(http.Dir("www")))
	http.ListenAndServe(":9892", nil)
	l.Wait()
}
