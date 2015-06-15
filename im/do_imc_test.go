package im

import (
	"fmt"
	"github.com/Centny/gwf/netw"
	"github.com/Centny/gwf/pool"
	"github.com/Centny/gwf/routing"
	"github.com/Centny/gwf/routing/httptest"
	"net/http"
	"runtime"
	"testing"
	"time"
)

func TestDoImc(t *testing.T) {
	ShowLog = true
	runtime.GOMAXPROCS(runtime.NumCPU() - 1)
	db := NewMemDbH()
	p := pool.NewBytePool(8, 102400)
	psrv := NewPushSrv(p, ":5498", "Push", netw.NewDoNotH(), db)
	psrv.TickLog = false
	err := psrv.Run()
	if err != nil {
		t.Error(err.Error())
		return
	}
	ts := httptest.NewServer(func(hs *routing.HTTPSession) routing.HResult {
		var r, c, s string
		var t int64
		err := hs.ValidCheckVal(`
		s,R|S,L:0;
		r,R|S,L:0;
		c,R|S,L:0;
		t,R|I,O:0~1~101;
		`, &s, &r, &c, &t)
		if err != nil {
			return hs.MsgResErr2(1, "arg-err", err)
		}
		_, err = psrv.PushN(s, r, c, uint32(t))
		if err == nil {
			return hs.MsgRes("OK")
		} else {
			return hs.MsgResErr2(1, "srv-err", err)
		}
	})
	l := NewListner3(db, fmt.Sprintf("S-vv-%v", 0), p, 9780, 1000000)
	l.WsAddr = fmt.Sprintf(":%v", 9770)
	l.PushSrvAddr = "127.0.0.1:5498"
	l.PushSrvTickLog = false
	rc := make(chan int)
	go func() {
		rc <- 1
		hs := &http.Server{
			Handler: l.WIM_L.WsS(),
			Addr:    l.WsAddr,
		}
		hs.ListenAndServe()
	}()
	err = l.Run()
	if err != nil {
		t.Error(err.Error())
		return
	}
	<-rc
	time.Sleep(time.Second)
	//
	purl := ts.URL + "?s=%v&r=%v&c=%v&t=%v"
	db.Grp["G-x"] = []string{"U-1", "U-2", "U-3", "U-abc"}
	di := NewDoImc(":9780", false, []string{"a", "b", "c"}, []string{"G-x"}, 8, purl, "U-4")
	err = di.Do()
	if err != nil {
		t.Error(err.Error())
		return
	}
	err = di.Check2(1000, 10000)
	if err != nil {
		t.Error(err.Error())
		fmt.Println(di.Res)
		return
	}
	fmt.Println(di.Res)
}
