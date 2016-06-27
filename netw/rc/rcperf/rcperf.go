package main

import (
	"fmt"
	"github.com/Centny/gwf/log"
	"github.com/Centny/gwf/netw"
	"github.com/Centny/gwf/netw/impl"
	"github.com/Centny/gwf/netw/rc"
	"github.com/Centny/gwf/pool"
	"github.com/Centny/gwf/routing"
	"github.com/Centny/gwf/routing/filter"
	"github.com/Centny/gwf/tutil"
	"github.com/Centny/gwf/util"
	"os"
	"runtime"
	"sync/atomic"
	"time"
)

func pref_exec(rc *impl.RCM_Cmd) (interface{}, error) {
	var dc int64
	err := rc.ValidF(`
		dc,R|I,R:0;
		`, &dc)
	if err != nil {
		log.E("pref_exec valid args error:%v", err.Error())
		return nil, err
	}
	time.Sleep(50 * time.Millisecond)
	return util.Map{
		"code": 0,
		"data": dc,
	}, nil
}

var CR *rc.RC_Runner_m

func server() {
	bp := pool.NewBytePool(8, 102400)
	lm := rc.NewRC_Listener_m_j(bp, ":10812", netw.NewDoNotH())
	lm.StartMonitor()
	lm.AddHFunc("exec", pref_exec)
	err := lm.Run()
	if err != nil {
		fmt.Println(err)
		return
	}
	cr := rc.NewRC_Runner_m_j(bp, "127.0.0.1:10812", netw.NewDoNotH())
	cr.Start()
	cr.StartMonitor()
	CR = cr
	// go client2()
	cr2 := rc.NewRC_Runner_m_j(bp, "127.0.0.1:10812", netw.NewDoNotH())
	cr2.Start()
	cr2.StartMonitor()
	cr3 := rc.NewRC_Runner_m_j(bp, "127.0.0.1:10812", netw.NewDoNotH())
	cr3.Start()
	cr3.StartMonitor()
	monitor := filter.NewMonitorH()
	monitor.AddMonitor("Hand", lm)
	monitor.AddMonitor("Runner", cr)
	monitor.AddMonitor("http", routing.Shared)
	// monitor.AddMonitor("chan", lm.CH)
	routing.Shared.StartMonitor()
	routing.H("^/adm/status(\\?.*)?$", monitor)
	routing.HFunc("^/adm/list(\\?.*)?$", func(hs *routing.HTTPSession) routing.HResult {
		var res = util.Map{}
		for k, v := range cr.RC_Con.Cts {
			res[fmt.Sprintf("e_%v", k)] = v
		}
		return hs.MsgRes(res)
	})
	var fail int64 = 0
	var sw_c int64 = 0
	routing.HFunc("^/test(\\?.*)?$", func(hs *routing.HTTPSession) routing.HResult {
		var res util.Map
		var err error
		tsw := atomic.AddInt64(&sw_c, 1)
		switch tsw % 3 {
		case 0:
			res, err = cr.VExec_m("exec", util.Map{"dc": hs.RVal("idx")})
		case 1:
			res, err = cr2.VExec_m("exec", util.Map{"dc": hs.RVal("idx")})
			// return hs.MsgRes("OK")
		default:
			res, err = cr3.VExec_m("exec", util.Map{"dc": hs.RVal("idx")})
			// return hs.MsgRes("OK")
		}
		if err != nil {
			fail++
			fmt.Println(err.Error())
			return hs.MsgResErr(1, "error", err)
		}
		if res.IntVal("code") != 0 {
			return hs.MsgResErr(1, "error", util.Err("%v", util.S2Json(res)))
		} else {
			return hs.MsgRes(res.Val("data"))
		}
	})
	routing.ListenAndServe(":8344")
}

func client() {
	runtime.GOMAXPROCS(util.CPU())
	all := 3000
	tc := 500
	util.ShowLog = true
	var done int64 = 0
	var errc int64 = 0
	var beg = util.Now()
	total, err := tutil.DoPerfV_(all, tc, "", func(i int) error {
		res, err := util.HGet2("http://127.0.0.1:8344/test?idx=%v", i+1)
		atomic.AddInt64(&done, 1)
		fmt.Printf("http request %v/%v done...\n", i+1, done)
		if err != nil {
			fmt.Println(err)
			atomic.AddInt64(&errc, 1)
			return nil
		}
		if res.IntVal("code") == 0 {
			if res.IntVal("data") != int64(i+1) {
				return util.Err("%v bytes found, expect %v", res.Val("data"), i+1)
			}
			return nil
		} else {
			atomic.AddInt64(&errc, 1)
			return util.Err("%v", util.S2Json(res))
		}
	})
	fmt.Println("total:", total, "avg:", total/int64(all), "errc:", errc, "err:", err, "used:", util.Now()-beg)
}

func client2() {
	all := 100000
	tc := 50000
	util.ShowLog = true
	var errc int64 = 0
	var beg = util.Now()
	total, err := tutil.DoPerfV_(all, tc, "", func(i int) error {
		CR.VExec_m("exec", util.Map{"dc": 1})
		return nil
	})
	fmt.Println("total:", total, "avg:", total/int64(all), "errc:", errc, "err:", err, "used:", util.Now()-beg)
}

func main() {
	if len(os.Args) < 2 {
		return
	}
	runtime.GOMAXPROCS(util.CPU())
	switch os.Args[1] {
	case "s":
		server()
	case "c":
		client()
	}
}
