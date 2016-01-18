package dtm

import (
	"fmt"
	"github.com/Centny/gwf/netw"
	"github.com/Centny/gwf/pool"
	"github.com/Centny/gwf/util"
	"runtime"
	"testing"
	"time"
)

func TestDtmBase(t *testing.T) {
	runtime.GOMAXPROCS(util.CPU())
	bp := pool.NewBytePool(8, 10240000)
	netw.ShowLog = true
	// impl.ShowLog = true
	dtmh := NewDTM_S_Proc()
	dtms := NewDTM_S_j(bp, ":23244", dtmh)
	dtms.AddToken2([]string{"abc"})
	err := dtms.Run()
	if err != nil {
		t.Error(err.Error())
		return
	}
	dtmc := NewDTM_C_j(bp, "127.0.0.1:23244")
	dtmc.Start()
	go dtmc.RunProcH(23245)
	dtmc.Login_("abc")
	var cid string
	for cid_, _ := range dtms.CmdCs() {
		cid = cid_
	}
	fmt.Println("---->")
	if len(cid) < 1 {
		t.Error("clinet not found")
		return
	}
	do_c := func(wait bool) {
		tid, err := dtms.StartTask(cid, "./dtm_test.sh 1 http://127.0.0.1:${proc_port}/proc?tid=${tid}")
		if err != nil {
			t.Error(err.Error())
			return
		}
		if _, ok := dtmh.Rates[tid]; !ok {
			t.Error("tid not found")
			return
		}
		if wait {
			err = dtms.WaitTask(cid, tid)
			if err != nil {
				t.Error(err.Error())
				return
			}
		} else {
			time.Sleep(time.Second)
			dtms.StopTask(cid, tid)
			err = dtms.WaitTask(cid, tid)
			if err == nil {
				t.Error("error")
				return
			}
		}
	}
	go do_c(false)
	do_c(true)
	//
	dtms.StartTask(cid, "")
	dtms.StartTask(cid, "xx")
	dtms.StopTask(cid, "")
	dtms.StopTask(cid, "tid")
	dtms.WaitTask(cid, "")
	dtms.WaitTask(cid, "tid")
	dtms.StartTask("cid", "cmds")
	dtms.StopTask("cid", "tid")
	dtms.WaitTask("cid", "tid")
	dtmc.OnCmd(nil)
	dtmc.Writev2([]byte{CMD_M_PROC}, "{s}")
	dtmc.Writev2([]byte{CMD_M_PROC}, util.Map{})
	util.HGet("http://127.0.0.1:%v/proc", dtmc.ProcPort)
	//
	fmt.Println("--------->>-a")
	dtms.Close()
	time.Sleep(time.Second)
	go func() {
		fmt.Println(util.HGet("http://127.0.0.1:%v/proc?tid=%v&process=0.5", dtmc.ProcPort, "tid"))
		fmt.Println("--------->>-a-0")
	}()
	time.Sleep(time.Second)
	fmt.Println("--------->>-b")
	dtmc.Stop()
	time.Sleep(time.Second)
	//
	//
	fmt.Println("done...")
}
