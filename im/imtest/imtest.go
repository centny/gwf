package main

import (
	"fmt"
	"github.com/Centny/gwf/im"
	"github.com/Centny/gwf/log"
	"github.com/Centny/gwf/netw"
	"github.com/Centny/gwf/pool"
	"github.com/Centny/gwf/routing"
	"github.com/Centny/gwf/routing/httptest"
	"github.com/Centny/gwf/tutil"
	"github.com/Centny/gwf/util"
	"net/http"
	_ "net/http/pprof"
	"runtime"
	"time"
)

func main() {
	// var bys = []byte{50, 50, 55, 18, 5, 85, 45, 49, 45, 52, 26, 5, 85, 45, 49, 45, 49, 26, 5, 85, 45, 49, 45, 50, 26, 5, 85, 45, 49, 45, 51, 32, 0, 42, 5, 85, 45, 49, 45, 51, 50, 6, 80, 117, 115, 104, 45, 62, 58, 5, 85, 45, 49, 45, 52, 64, 187, 253, 241, 244, 186, 42}
	// fmt.Println("xxx")
	// fmt.Println(string(bys), "--->")
	log.Redirect("logs/out_%v.log", "logs/err_%v.log")
	runtime.GOMAXPROCS(util.CPU())
	go http.ListenAndServe(":2345", nil)
	run_do_imc_(100000, 1)
	// time.Sleep(100000 * time.Second)
}

func run_do_imc_(total, tc int) {
	im.ShowLog = true
	runtime.GOMAXPROCS(runtime.NumCPU() - 1)
	db := im.NewMemDbH()
	p := pool.NewBytePool(8, 102400)
	psrv := im.NewPushSrv(p, ":5498", "Push", netw.NewDoNotH(), db)
	psrv.TickLog = false
	err := psrv.Run()
	if err != nil {
		panic(err)
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
	l := im.NewListner3(db, fmt.Sprintf("S-vv-%v", 0), p, 9780, 1000000)
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
		panic(err)
	}
	<-rc
	time.Sleep(time.Second)
	//
	purl := ts.URL + "?s=%v&r=%v&c=%v&t=%v"
	// var idx = 0
	http.HandleFunc("/abc", func(w http.ResponseWriter, r *http.Request) {
		// run_do_imc_c(idx, db, purl)
		// idx++
	})
	// for i := 0; i < total; i++ {
	tutil.DoPerfV(total, tc, "", func(idx int) {
		run_do_imc_c(idx, db, purl)
	})
	fmt.Println("one done...\n\n\n\n\n\n")
	// }
	fmt.Println("all done...")
}

func run_do_imc_c(i int, db *im.MemDbH, purl string) {
	ga := fmt.Sprintf("G-%v", i)
	ua, ub, uc, ud := fmt.Sprintf("U-%v-%v", i, 1), fmt.Sprintf("U-%v-%v", i, 2),
		fmt.Sprintf("U-%v-%v", i, 3), fmt.Sprintf("U-%v-%v", i, 4)
	ta, tb, tc := fmt.Sprintf("T-%v-%v", i, 1), fmt.Sprintf("T-%v-%v", i, 2), fmt.Sprintf("T-%v-%v", i, 3)
	db.AddGrp(ga, []string{
		ua, ub, uc,
		"U-abc",
	})
	db.AddTokens(map[string]string{
		ta: ua,
		tb: ub,
		tc: uc,
	})
	di := im.NewDoImc(pool.BP, ":9780", false, []string{ta, tb, tc}, []string{ga}, 8, purl, ud)
	di.Name = fmt.Sprintf("I%v", i)
	err := di.Do()
	if err != nil {
		panic(err)
	}
	err = di.Check2(100, 100000)
	if err != nil {
		panic(err)
	}
	di.Release()
	db.DelTokens([]string{ta, tb, tc})
	db.DelGrp(ga)
	db.ClearMsg([]string{ua, ub, uc, ud})
	fmt.Printf("\n\nclient %v done...\n\n\n", i)
}
