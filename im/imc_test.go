package im

import (
	"fmt"
	"github.com/Centny/gwf/routing"
	"github.com/Centny/gwf/routing/httptest"
	// "github.com/Centny/gwf/netw"
	// "github.com/Centny/gwf/netw/impl"
	"github.com/Centny/gwf/pool"
	"runtime"
	"testing"
	"time"
)

func TestIMC(t *testing.T) {
	runtime.GOMAXPROCS(runtime.NumCPU() - 1)
	// ShowLog = true
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
	imc := NewIMC3(srvs, "token")
	// imc.ShowLog = true
	imc.TickData = []byte{}
	imc.StartRunner()
	imc.StartHB()
	imc2, err := NewIMC4(ts.URL, "token")
	if err != nil {
		t.Error(err.Error())
		return
	}
	// imc2.ShowLog = true
	imc2.TickData = []byte{}
	imc2.StartRunner()
	<-imc.LC
	<-imc2.LC
	fmt.Println(imc.IC)
	fmt.Println(imc2.IC)
	for i := 0; i < 10; i++ {
		imc.SMS(imc2.IC.R, 0, "imc1-00--->")
		imc2.SMS(imc.IC.R, 0, "imc2-00--->")
	}
	for imc.RC < 10 || imc2.RC < 10 {
		time.Sleep(500 * time.Millisecond)
	}
	imc.MCon.Close()
	<-imc.LC
	for i := 0; i < 10; i++ {
		imc.SMS(imc2.IC.R, 0, "imc1-00--->")
		imc2.SMS(imc.IC.R, 0, "imc2-00--->")
	}
	for imc.RC < 20 || imc2.RC < 20 {
		time.Sleep(500 * time.Millisecond)
	}
	imc.Close()
	imc2.Close()
	fmt.Println("all done ....")
}
