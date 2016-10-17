package plugin

// import (
// 	"fmt"
// 	"github.com/Centny/gwf/netw"
// 	"github.com/Centny/gwf/netw/rc"
// 	"github.com/Centny/gwf/pool"
// 	"github.com/Centny/gwf/util"
// 	"runtime"
// 	"testing"
// 	"time"
// )

// func TestPing(t *testing.T) {
// 	runtime.GOMAXPROCS(util.CPU())
// 	bp := pool.NewBytePool(8, 1024000)
// 	rcs := rc.NewRC_Listener_m_j(bp, ":28741", netw.NewDoNotH())
// 	err := rcs.Run()
// 	if err != nil {
// 		t.Error(err.Error())
// 		return
// 	}
// 	ping3 := RegPing_S_H(bp, rcs)
// 	ping3.ShowLog = true
// 	fmt.Println(ping3.Status())
// 	rcc := rc.NewRC_Runner_m_j(bp, "127.0.0.1:28741", netw.NewDoNotH())
// 	err = rcc.Run()
// 	if err != nil {
// 		t.Error(err.Error())
// 		return
// 	}
// 	ping := NewPing_C_H(bp, rcc.RCM_Con)
// 	ping.ShowLog = true
// 	fmt.Println(ping.Ping(8, 1024))
// 	fmt.Println(ping.Ping(0, 0))
// 	fmt.Println(ping.PingC(8, 1024, 10))
// 	rcc.VExec_m("ping", util.Map{
// 		"len": "sss",
// 	})
// 	//
// 	ping.StartPingS(500, 8, 102400)
// 	time.Sleep(time.Second)
// 	res := ping.Status()
// 	if res.IntVal("ping_c") < 1 {
// 		t.Error("error")
// 		return
// 	}
// 	if res.IntVal("ping_s/running") < 1 {
// 		t.Error("error")
// 		return
// 	}
// 	ping.StopPingS()
// 	time.Sleep(time.Second)
// 	ping2 := NewPing_C_H2(bp, rcc)
// 	ping2.ShowLog = true
// 	fmt.Println(ping2.Ping(8, 1024))
// 	rcs.Close()
// 	rcs.Wait()
// 	fmt.Println("rcs--->")
// 	time.Sleep(100 * time.Millisecond)
// 	go func() {
// 		ping.PingC(8, 8, 1)
// 	}()
// 	time.Sleep(100 * time.Millisecond)
// 	rcc.Stop()
// 	rcc.Wait()
// 	fmt.Println("xxx->")
// 	time.Sleep(1000 * time.Millisecond)
// 	// ping.StopPingS()
// 	fmt.Println("xxx->1")
// 	ping3.S = nil
// 	ping2.VExec_m = nil
// 	func() {
// 		defer func() {
// 			fmt.Println(recover())
// 		}()
// 		fmt.Println(ping2.Ping(8, 1024))
// 	}()
// 	func() {
// 		defer func() {
// 			fmt.Println(recover())
// 		}()
// 		fmt.Println(ping2.PingC(8, 1024, 10))
// 	}()
// 	func() {
// 		defer func() {
// 			fmt.Println(recover())
// 		}()
// 		fmt.Println(ping3.PingH(nil))
// 	}()
// 	func() {
// 		defer func() {
// 			fmt.Println(recover())
// 		}()
// 		fmt.Println(ping2.PingS(8, 9))
// 	}()
// 	func() {
// 		defer func() {
// 			fmt.Println(recover())
// 		}()
// 		ping2.DoPingS(500, 8, 9)
// 	}()
// 	func() {
// 		defer func() {
// 			fmt.Println(recover())
// 		}()
// 		ping2.StartPingS(500, 8, 9)
// 	}()
// 	// var tx int
// 	ping.VExec_m = func(name string, args interface{}) (util.Map, error) {
// 		// if tx%2 == 1 {
// 		return nil, util.Err("error")
// 		// } else {
// 		// 	return util.Map{}, nil
// 		// }
// 	}
// 	ping.StartPingS(500, 8, 1024)
// 	time.Sleep(time.Second)
// 	ping.StopPingS()
// 	ping.PingC(8, 8, 2)
// 	ping.VExec_m = func(name string, args interface{}) (util.Map, error) {
// 		return util.Map{}, nil
// 	}
// 	ping.PingS(8, 8)
// 	fmt.Println("Ping done...")
// }
