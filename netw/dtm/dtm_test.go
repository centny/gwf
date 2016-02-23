package dtm

import (
	"fmt"
	"github.com/Centny/gwf/netw"
	"github.com/Centny/gwf/pool"
	"github.com/Centny/gwf/routing/httptest"
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
	//
	//test not client case
	_, _, err = dtms.StartTask3("cmds")
	if err == nil {
		t.Error("error")
		return
	}

	//test client connected case
	var dtmc *DTM_C

	dtmc = NewDTM_C_j(bp, "127.0.0.1:23244")
	dtmc.RunProcH() //not configure
	dtmc.Cfg.SetVal("proc_addr", ":23245")
	dtmc.Start()
	go dtmc.RunProcH()
	dtmc.Login_("abc")
	fmt.Println("---->")
	if len(dtms.CmdCs()) < 1 {
		t.Error("clinet not found")
		return
	}
	cid := dtmh.MinUsedCid(dtms)
	if len(cid) < 1 {
		t.Error("error")
		return
	}
	do_c := func(wait bool) {
		cid, tid, err := dtms.StartTask3("./dtm_test.sh 1 http://127.0.0.1${proc_addr}/proc?tid=${proc_tid}")
		if err != nil {
			t.Error(err.Error())
			return
		}
		if _, ok := dtmh.Rates[cid][tid]; !ok {
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
	var args = util.Map{}
	if dtmc.cmd_do_res(args, "xxxx", `
		sfs
		fsdf

		`) != 0 {
		t.Error("error")
		return
	}
	if dtmc.cmd_do_res(args, "xxxx", `
	----------------result----------------

		`) != 0 {
		t.Error("error")
		return
	}
	if dtmc.cmd_do_res(args, "xxxx", `
	----------------result----------------
	[json]
	{}
	[/json]
		`) != 0 {
		t.Error("error")
		return
	}
	if dtmc.cmd_do_res(args, "xxxx", `
	----------------result----------------
	[json]
	{xx}
	[/json]
		`) == 0 {
		t.Error("error")
		return
	}
	//
	httptest.NewServer2(dtms.H.(*DTM_S_Proc)).G("?xx=1")
	//
	//test error
	fmt.Println("xxxx->")
	err = dtms.StartTask(cid, "", "")
	if err == nil {
		t.Error("error")
		return
	}
	fmt.Println("\n\n\n\n\n\n\n\n", err, "\n\n\n\n\n\n\n\n")
	// if true {
	// 	return
	// }
	dtms.StartTask(cid, "", "xx")
	dtms.StopTask(cid, "")
	dtms.StopTask(cid, "tid")
	dtms.WaitTask(cid, "")
	dtms.WaitTask(cid, "tid")
	dtms.StartTask("cid", "ss", "cmds")
	dtms.StopTask("cid", "tid")
	dtms.WaitTask("cid", "tid")
	dtmc.OnCmd(nil)
	dtmc.Writev2([]byte{CMD_M_PROC}, "{s}")
	dtmc.Writev2([]byte{CMD_M_PROC}, util.Map{})
	dtmc.Writev2([]byte{CMD_M_DONE}, "{s}")
	dtmc.Writev2([]byte{CMD_M_DONE}, util.Map{})
	util.HGet("http://127.0.0.1%v/proc", dtmc.Cfg.Val("proc_addr"))
	//
	//test stop
	fmt.Println("--------->>-a")
	dtms.Close()
	time.Sleep(time.Second)
	go func() {
		fmt.Println(util.HGet("http://127.0.0.1%v/proc?tid=%v&process=0.5", dtmc.Cfg.Val("proc_addr"), "tid"))
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
