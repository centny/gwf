package main

import (
	"fmt"
	"github.com/Centny/gwf/im"
	"github.com/Centny/gwf/netw"
	"github.com/Centny/gwf/pool"
	"github.com/Centny/gwf/routing"
	"github.com/Centny/gwf/routing/httptest"
	"github.com/Centny/gwf/util"
	"os"
	"runtime"
	"testing"
	"time"
)

func TestImr(t *testing.T) {
	runtime.GOMAXPROCS(util.CPU())
	os.Remove("t.txt")
	os.Remove("t.log")
	// ShowLog = true
	im.ShowLog = true
	// impl.ShowLog = true
	netw.ShowLog = true
	db := im.NewMemDbH()
	db.Grp["G-xx"] = []string{"U-1", "U-2", "U-3"}
	p := pool.NewBytePool(8, 102400)
	psrv := im.NewPushSrv(p, ":5498", "Push", netw.NewDoNotH(), db)
	psrv.TickLog = false
	err := psrv.Run()
	if err != nil {
		t.Error(err.Error())
		return
	}
	l := im.NewListner3(db, fmt.Sprintf("S-vx-%v", 0), p, 9790, 1000000)
	l.WsAddr = fmt.Sprintf(":%v", 9770)
	l.PushSrvAddr = "127.0.0.1:5498"
	l.PushSrvTickLog = false
	go func() {
		err := l.Run()
		if err != nil {
			panic(err.Error())
		}
	}()
	time.Sleep(time.Second)
	ts := httptest.NewServer(func(hs *routing.HTTPSession) routing.HResult {
		srvs, _ := db.ListSrv("")
		return hs.MsgRes(srvs)
	})
	srvs, _ := db.ListSrv("")
	timc := im.NewIMC3(srvs, "token")
	pts := httptest.NewServer(func(hs *routing.HTTPSession) routing.HResult {
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
	// imc.ShowLog = true
	timc.Start()
	timc.LC.Wait()
	util.FWrite("t.txt", fmt.Sprintf("sskkdd\n%v abcc\nabkkk", timc.IC.R))
	//
	var wwc chan int = make(chan int)
	//
	os.Args = []string{"imr", "-t", "x1,x2,x3", "-l", ts.URL, "-m", "T", "-g", "G-xx",
		"-P", pts.URL + "?s=%v&r=%v&c=%v&t=%v", "-p", "U-ss", "-c", "9", "-T", "5000"}
	go func() {
		if run() != 0 {
			t.Error("error")
		}
		wwc <- 1
	}()
	<-wwc
	fmt.Println("imr->T done")
	// if true {
	// 	return
	// }
	//
	os.Args = []string{"imr", "-t", "xxx", "-l", ts.URL, "-m", "R", "-L", "N"}
	go func() {
		if run() != 0 {
			t.Error("error")
		}
		wwc <- 1
	}()
	fmt.Println("xxx-->")
	time.Sleep(time.Second)
	imc.LC.Wait()
	timc.SMS(imc.IC.R, 0, "abcc")
	time.Sleep(time.Second)
	imc.Close()
	fmt.Println("---->>>")
	<-wwc
	fmt.Println("imr->R done")
	//
	os.Args = []string{"imr", "-t", "xxx", "-s", ":9790", "-m", "C", "-L", "t.log"}
	os.Stdin, _ = os.Open("t.txt")
	go func() {
		if run() != 0 {
			t.Error("error")
		}
		wwc <- 1
	}()
	time.Sleep(time.Second)
	imc.LC.Wait()
	timc.SMS(imc.IC.R, 0, "abcc")
	fmt.Println("sending stop1")
	time.Sleep(time.Second)
	fmt.Println("sending stop2")
	imc.Close()
	fmt.Println("waiting stop")
	<-wwc
	fmt.Println("imr->C done")

	fmt.Println("other")
	os.Args = []string{"ss"}
	run()
	os.Args = []string{"imr", "-t", "xxx", "-s", ":9790", "-m", "C", "-L", "/sd/ds.log"}
	run()
	os.Args = []string{"imr", "-t", "xxx", "-l", ":979x", "-m", "X", "-L", "N"}
	run()
	os.Args = []string{"imr", "-t", "xxx", "-s", ":9790", "-m", "X", "-L", "N"}
	run()
	os.Args = []string{"imr", "-t", "xxx", "-s", ":9790", "-m", "X", "-L", "E"}
	run()
	os.Args = []string{"imr", "-t", "xxx", "-s", ":9790", "-m", "X", "-L", "W"}
	run()
	os.Args = []string{"imr", "-t", "xxx", "-s", ":9790", "-m", "X", "-L", "I"}
	run()
	os.Args = []string{"imr", "-t", "xxx", "-s", ":9790", "-m", "X", "-L", "D"}
	run()
	os.Args = []string{"imr", "-h", "ss"}
	run()
	os.Args = []string{"imr", "-t", "sss", "-m", "R", "-L", "SS"}
	run()
	os.Args = []string{"imr", "-t", "sss", "-m", "T", "-c", "xx"}
	run()
	os.Args = []string{"imr", "-t", "sss", "-m", "T", "-T", "xx"}
	run()
	os.Args = []string{"imr", "-t", "sss", "-m", "T"}
	run()
	os.Args = []string{"imr", "-t", "sss", "-m", "T", "-l", "hssfs"}
	run()
	os.Args = []string{"imr", "-t", "sss", "-L", "D", "-m", "T", "-s", "12.23"}
	go run()
	time.Sleep(3 * time.Second)
	os.Args = []string{"imr", "-h"}
	ef = func(c int) {}
	main()
	time.Sleep(time.Second)
}
