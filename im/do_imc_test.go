package im

import (
	"fmt"
	"github.com/Centny/gwf/netw/impl"
	"sync"
	// "github.com/Centny/gwf/log"
	"github.com/Centny/gwf/netw"
	"github.com/Centny/gwf/pool"
	"github.com/Centny/gwf/routing"
	"github.com/Centny/gwf/routing/filter"
	"github.com/Centny/gwf/tutil"
	"github.com/Centny/gwf/util"
	"net/http"
	"runtime"
	"sync/atomic"
	"testing"
	"time"
)

func TestDoImc(t *testing.T) {
	run_do_imc_(t, 1000, 100)
}

func run_do_imc_(t *testing.T, total, tc int) {
	// log.SetLevel(log.ERROR)
	ShowLog = true
	netw.ShowLog_C = true
	netw.ShowLog = true
	impl.ShowLog = true
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
	mux := routing.Shared
	mux.HFunc(".*", func(hs *routing.HTTPSession) routing.HResult {
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
	monitor := filter.NewMonitorH()
	l := NewListner3(db, fmt.Sprintf("S-vv-%v", 0), p, 9780, 1000000)
	l.StartMonitor()
	l.WsAddr = fmt.Sprintf(":%v", 9770)
	l.PushSrvAddr = "127.0.0.1:5498"
	l.PushSrvTickLog = false
	monitor.AddMonitor("L", l)
	mux.HFilter("^/adm/status$", monitor)
	mux.HFilterFunc("^/adm/show$", func(hs *routing.HTTPSession) routing.HResult {
		var res = util.Map{}
		dolck.Lock()
		for idx, imc := range doing {
			if len(imc.Res) < 1 {
				continue
			}
			var tval = util.Map{}
			tval["res"] = imc.Res
			for key, xx := range imc.imcs {
				cons := []string{}
				for _, con := range xx.Cons() {
					cons = append(cons, con.LocalAddr().String())
				}
				tval[key] = cons
			}
			res[fmt.Sprintf("I-%v", idx)] = tval
		}
		dolck.Unlock()
		var msss = map[string]*Msg{}
		for mid, msg := range db.Ms {
			for _, mss := range msg.Ms {
				for _, ms := range mss {
					if ms.S == MS_DONE {
						continue
					}
					msss[mid] = msg
				}
			}
		}
		return hs.MsgRes(util.Map{
			"res":   res,
			"mdb":   msss,
			"res_l": len(res),
		})
	})
	go routing.ListenAndServe(":2388")
	// fmt.Println(ts.URL)
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
	pool.BP = pool.NewBytePool(8, 10240000)
	//
	purl := "http://127.0.0.1:2388/" + "?s=%v&r=%v&c=%v&t=%v"
	var done uint32
	tutil.DoPerfV(total, tc, "a.log", func(i int) {
		run_do_imc_c(i, db, purl, t, monitor)
		atomic.AddUint32(&done, 1)
		fmt.Printf("%v/%v done(%v)...\n", i, total, done)
	})
	if len(db.Ms) > 0 || len(db.Cons) > 0 || len(db.Usr) > 0 ||
		len(db.Grp) > 0 || len(db.Tokens) > 0 || len(db.U2M) > 0 {
		t.Error("error")
		return
	}
	fmt.Println("all done...")
}

var doing = map[int]*DoImc{}
var dolck = sync.RWMutex{}

func run_do_imc_c(i int, db *MemDbH, purl string, t *testing.T, m *filter.MonitorH) {
	ga := fmt.Sprintf("G-%v", i)
	ua, ub, uc, ud := fmt.Sprintf("U-%v-%v", i, 1), fmt.Sprintf("U-%v-%v", i, 2),
		fmt.Sprintf("U-%v-%v", i, 3), fmt.Sprintf("U-%v-%v", i, 4)
	ta, tb, tc := fmt.Sprintf("T-%v-%v", i, 1), fmt.Sprintf("T-%v-%v", i, 2), fmt.Sprintf("T-%v-%v", i, 3)
	fmt.Println("---->", ua, ub, uc, ud, "->", ta, tb, tc)
	db.AddGrp(ga, []string{
		ua, ub, uc,
	})
	db.AddTokens(map[string]string{
		ta: ua,
		tb: ub,
		tc: uc,
	})
	di := NewDoImc(pool.BP, ":9780", false, []string{ta, tb, tc}, []string{ga}, 8, purl, ud)
	dolck.Lock()
	doing[i] = di
	dolck.Unlock()
	// m.AddMonitor(fmt.Sprintf("M-%v", i), di)
	err := di.Do()
	if err != nil {
		t.Error(err.Error())
		return
	}
	err = di.Check2(100, 1000000)
	if err != nil {
		t.Error(err.Error())
	}
	di.Release()
	db.DelTokens([]string{ta, tb, tc})
	db.DelGrp(ga)
	db.ClearMsg([]string{ua, ub, uc, ud})
	dolck.Lock()
	delete(doing, i)
	dolck.Unlock()
}

func TestByte(t *testing.T) {
	var bys = []byte{77, 45, 54, 52, 18, 5, 85, 45, 48, 45, 50, 26, 3, 71, 45, 48, 32, 0, 42, 5, 85, 45, 48, 45, 51, 50, 6, 71, 45, 48, 45, 62, 51, 58, 3, 71, 45, 48, 64, 196, 222, 145, 174, 218, 42}
	fmt.Println(string(bys))
}
