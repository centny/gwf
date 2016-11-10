package im

import (
	"fmt"

	"github.com/Centny/gwf/util"
	// "github.com/Centny/gwf/netw"
	"github.com/Centny/gwf/routing"
	"github.com/Centny/gwf/routing/httptest"
	// "github.com/Centny/gwf/netw/impl"
	"runtime"
	"testing"
	"time"

	"github.com/Centny/gwf/pool"
)

func TestIMC(t *testing.T) {
	runtime.GOMAXPROCS(util.CPU())
	ShowLog = true
	// impl.ShowLog = true
	// netw.ShowLog = true
	db := NewMemDbH()
	p := pool.NewBytePool(8, 102400)
	l := NewListner3(db, fmt.Sprintf("S-vx-%v", 0), p, 9790, 1000000)
	go func() {
		err := l.Run()
		if err != nil {
			panic(err.Error())
		}
	}()
	ts := httptest.NewServer(func(hs *routing.HTTPSession) routing.HResult {
		srvs, _ := db.ListSrv("")
		return hs.MsgRes(srvs)
	})
	time.Sleep(time.Second)
	srvs, err := db.ListSrv("")
	if err != nil {
		t.Error(err.Error())
		return
	}
	fmt.Println("imc--->01-00")
	imc := NewIMC3(pool.BP, srvs, "token")
	// imc.ShowLog = true
	imc.TickData = []byte{}
	imc.Start()
	imc2, err := NewIMC4(pool.BP, ts.URL, "token")
	if err != nil {
		t.Error(err.Error())
		return
	}
	// imc2.ShowLog = true
	imc2.TickData = []byte{}
	imc2.Start()
	fmt.Println("imc--->01-01")
	imc.LC.Wait()
	imc2.LC.Wait()
	imc.StartHB()
	imc2.StartHB()
	fmt.Println("imc--->01-02")
	fmt.Println(imc.IC)
	fmt.Println(imc2.IC)
	imc.UR()
	imc2.UR()
	fmt.Println(imc.Logined())
	imc.SMS("S-Robot-X", 0, "Robot")
	fmt.Println("\n\n\n")
	time.Sleep(time.Second)
	for i := 0; i < 10; i++ {
		imc.SMS(imc2.IC.Uid, 0, "imc1-00--->")
		imc2.SMS(imc.IC.Uid, 0, "imc2-00--->")
	}
	fmt.Println("imc--->01-03")
	for imc.RC < 10 || imc2.RC < 10 {
		fmt.Println("-->", imc.RC)
		time.Sleep(300 * time.Millisecond)
	}
	fmt.Println("\n\n\n")
	imc.MCon.Close()
	time.Sleep(time.Second)
	imc.LC.Wait()
	for i := 0; i < 10; i++ {
		imc.SMS(imc2.IC.Uid, 0, "imc1-00--->")
		imc2.SMS(imc.IC.Uid, 0, "imc2-00--->")
	}
	for imc.RC < 20 || imc2.RC < 20 {
		time.Sleep(300 * time.Millisecond)
	}
	db.Grp["G-abc"] = []string{"U-a", "U-b"}
	ur, err := imc.GR([]string{"G-abc"})
	if err != nil {
		t.Error(err.Error())
		return
	}
	urs := ur["G-abc"]
	if len(urs) < 2 || urs[0] != "U-a" || urs[1] != "U-b" {
		t.Error("error")
		return
	}
	fmt.Println(ur)
	fmt.Println("\n\n\n")
	fmt.Println(db.Show())
	imc.Close()
	fmt.Println("all done ....")
}

func TestMessageToSelf(t *testing.T) {
	runtime.GOMAXPROCS(util.CPU())
	ShowLog = true
	// impl.ShowLog = true
	// netw.ShowLog = true
	db := NewMemDbH()
	db.Tokens["token"] = "u1"
	p := pool.NewBytePool(8, 102400)
	l := NewListner3(db, fmt.Sprintf("S-vx-%v", 0), p, 9790, 1000000)
	go func() {
		err := l.Run()
		if err != nil {
			panic(err.Error())
		}
	}()
	time.Sleep(time.Second)
	srvs, err := db.ListSrv("")
	if err != nil {
		t.Error(err.Error())
		return
	}
	imc := NewIMC3(pool.BP, srvs, "token")
	// imc.ShowLog = true
	imc.TickData = []byte{}
	imc.Start()
	imc2 := NewIMC3(pool.BP, srvs, "token")
	// imc2.ShowLog = true
	imc2.TickData = []byte{}
	imc2.Start()
	imc.LC.Wait()
	imc2.LC.Wait()
	imc.StartHB()
	imc2.StartHB()
	fmt.Println("imc--->01-02")
	fmt.Println(imc.IC)
	fmt.Println(imc2.IC)
	imc.UR()
	imc2.UR()
	time.Sleep(time.Second)
	for i := 0; i < 10; i++ {
		imc.SMS("xx1", 0, "imc1-00--->")
	}
	fmt.Println("imc--->01-03")
	for i := 0; i < 10 && imc2.RC < 10; i++ {
		fmt.Println("-->", imc2.RC)
		time.Sleep(300 * time.Millisecond)
	}
	if imc2.RC != 10 {
		t.Error("time out")
		return
	}
	fmt.Println("\n\n\n")
	db.Grp["G-abc"] = []string{"U-a", "U-b", "u1"}
	for i := 0; i < 10; i++ {
		imc2.SMS("G-abc", 0, "imc1-00--->")
	}
	for i := 0; i < 10 && imc.RC < 10; i++ {
		fmt.Println("-->", imc.RC)
		time.Sleep(300 * time.Millisecond)
	}
	if imc.RC != 10 {
		t.Error("time out")
		return
	}
	fmt.Println("all done ....")
}
