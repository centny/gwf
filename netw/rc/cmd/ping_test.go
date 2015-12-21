package cmd

import (
	"fmt"
	"github.com/Centny/gwf/netw"
	"github.com/Centny/gwf/netw/rc"
	"github.com/Centny/gwf/pool"
	"github.com/Centny/gwf/util"
	"runtime"
	"testing"
	"time"
)

func TestPing(t *testing.T) {
	runtime.GOMAXPROCS(util.CPU())
	bp := pool.NewBytePool(8, 1024000)
	rcs := rc.NewRC_Listener_m_j(bp, ":28741", netw.NewDoNotH())
	err := rcs.Run()
	if err != nil {
		t.Error(err.Error())
		return
	}
	RegPingH(bp, rcs).ShowLog = true
	rcc := rc.NewRC_Runner_m_j(bp, "127.0.0.1:28741", netw.NewDoNotH())
	err = rcc.Run()
	if err != nil {
		t.Error(err.Error())
		return
	}
	ping := NewPing(bp, rcc)
	ping.ShowLog = true
	fmt.Println(ping.Ping(8, 1024))
	fmt.Println(ping.Ping(0, 0))
	fmt.Println(ping.PingC(8, 1024, 10))
	rcc.VExec_m("ping", util.Map{
		"len": "sss",
	})
	rcs.Close()
	rcs.Wait()
	fmt.Println("rcs--->")
	time.Sleep(100 * time.Millisecond)
	go func() {
		ping.PingC(8, 8, 1)
	}()
	time.Sleep(100 * time.Millisecond)
	rcc.Close()
	rcc.Stop()
	time.Sleep(100 * time.Millisecond)
	fmt.Println("Ping done...")
}
