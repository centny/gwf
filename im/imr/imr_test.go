package main

import (
	"fmt"
	"github.com/Centny/gwf/im"
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
	runtime.GOMAXPROCS(runtime.NumCPU() - 1)
	os.Remove("t.txt")
	os.Remove("t.log")
	// ShowLog = true
	// impl.ShowLog = true
	// netw.ShowLog = true
	db := im.NewMemDbH()
	p := pool.NewBytePool(8, 102400)
	l := im.NewListner3(db, fmt.Sprintf("S-vx-%v", 0), p, 9790, 1000000)
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
	// imc.ShowLog = true
	timc.StartRunner()
	<-timc.LC
	util.FWrite("t.txt", fmt.Sprintf("sskkdd\n%v abcc\nabkkk", timc.IC.R))
	//
	var wwc chan int = make(chan int)
	//
	os.Args = []string{"imr", "-t", "xxx", "-l", ts.URL, "-m", "R", "-L", "N"}
	go func() {
		run()
		wwc <- 1
	}()
	time.Sleep(time.Second)
	timc.SMS(imc.IC.R, 0, "abcc")
	time.Sleep(time.Second)
	imc.StopRunner()
	imc.Close()
	fmt.Println("---->>>")
	<-wwc
	//
	os.Args = []string{"imr", "-t", "xxx", "-s", ":9790", "-m", "C", "-L", "t.log"}
	os.Stdin, _ = os.Open("t.txt")
	go func() {
		run()
		wwc <- 1
	}()
	time.Sleep(time.Second)
	timc.SMS(imc.IC.R, 0, "abcc")
	time.Sleep(time.Second)
	imc.StopRunner()
	<-wwc
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
	// main()
}
