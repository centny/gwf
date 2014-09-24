package routing

import (
	"fmt"
	"github.com/Centny/gwf/log"
	"net/http"
)

func ExampleSessionMux() {
	sb := NewSrvSessionBuilder("", "/", "example", 60*60*1000, 10000)
	mux := NewSessionMux("/example", sb)
	mux.HFilterFunc("^.*$", func(hs *HTTPSession) HResult {
		log.D("filt 001")
		return HRES_CONTINUE
	})
	//http://localhost:8080/example/ok
	mux.HFunc("^/ok(\\?.*)?$", func(hs *HTTPSession) HResult {
		hs.MsgRes("OK")
		return HRES_RETURN
	})
	//http://localhost:8080/example/data
	mux.HFunc("^/data(\\?.*)?$", func(hs *HTTPSession) HResult {
		var tid int64
		var name string
		err := hs.ValidRVal(`
			tid,R|I,R:0
			name,R|S,L:0`, //valid the argument
			&tid, &name)
		if err != nil {
			return hs.MsgResE(1, err.Error())
		}
		return hs.MsgRes(fmt.Sprintf("%v:%v", tid, name))
	})
	mux.HFunc("^/mdata(\\?.*)?$", func(hs *HTTPSession) HResult {
		hs.W.Write([]byte("some data\n"))
		return HRES_RETURN
	})
	s := http.Server{Addr: ":8080", Handler: mux}
	err := s.ListenAndServe()
	if err != nil {
		fmt.Println(err)
	}
}
