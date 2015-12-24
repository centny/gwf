package plugin

import (
	"fmt"
	"github.com/Centny/gwf/netw"
	"github.com/Centny/gwf/netw/rc"
	"github.com/Centny/gwf/pool"
	"github.com/Centny/gwf/util"
	"runtime"
	"testing"
)

func TestRCH(t *testing.T) {
	runtime.GOMAXPROCS(util.CPU())
	bp := pool.NewBytePool(8, 1024000)
	rcs := rc.NewRC_Listener_m_j(bp, ":28742", netw.NewDoNotH())
	err := rcs.Run()
	if err != nil {
		t.Error(err.Error())
		return
	}
	rcc := rc.NewRC_Runner_m_j(bp, "127.0.0.1:28742", netw.NewDoNotH())
	err = rcc.Run()
	if err != nil {
		t.Error(err.Error())
		return
	}
	rcs.AddToken2([]string{"abc"})
	err = rcc.Login_("abc")
	if err != nil {
		t.Error(err.Error())
		return
	}
	var cid string
	for _, tcid := range rcs.RCH.MID {
		cid = tcid
	}
	fmt.Println(cid)
	//
	rcs_h := NewRC_S_H(bp, rcc)
	rcc_h := NewRC_C_H(bp, rcc)
	fmt.Println(rcs_h.Status(), rcc_h.Status())
	rcs_h.OnCmd(nil)
	rcc_h.OnCmd(nil)
	ctls := NewRC_CTL_S_H(rcs)
	ctl := ctls.CTL("xxxxx")
	if ctl != nil {
		t.Error("error")
	}
	ctl = ctls.CTL(cid)
	if ctl == nil {
		t.Error("error")
		return
	}
	{
		err = ctl.RC_S_Start("", RC_M_J)
		if err == nil {
			t.Error("error")
			return
		}
		err = ctl.RC_S_Start("127.0.0.28773", RC_M_J)
		if err == nil {
			t.Error("error")
			return
		}
	}
	err = ctl.RC_S_Start(":28773", RC_M_J)
	if err != nil {
		t.Error(err.Error())
		return
	}
	{
		err = ctl.RC_S_Start(":28773", RC_M_J)
		if err == nil {
			t.Error("error")
			return
		}
	}
	res, err := ctl.RC_S_Status()
	if err != nil {
		t.Error("error")
		return
	}
	if res.IntVal("status") != TS_RUNNING {
		t.Error("not running")
		return
	}
	err = ctl.RC_S_Token([]string{"abc"})
	if err != nil {
		t.Error(err.Error())
		return
	}
	{
		err = ctl.RC_S_Token([]string{""})
		if err == nil {
			t.Error("error")
			return
		}
	}
	{
		err = ctl.RC_C_Start("", "abc", RC_M_J)
		if err == nil {
			t.Error("error")
			return
		}
		err = ctl.RC_C_Start("xxx:xxx", "abc", RC_M_J)
		if err == nil {
			t.Error("error")
			return
		}
	}
	err = ctl.RC_C_Start("127.0.0.1:28773", "abc", RC_M_J)
	if err != nil {
		t.Error(err.Error())
		return
	}
	res, err = ctl.RC_C_Status()
	if err != nil {
		t.Error(err.Error())
		return
	}
	if res.IntVal("status") != TS_RUNNING {
		t.Error("not running")
		return
	}
	{
		err = ctl.RC_C_Start("127.0.0.1:28773", "abc", RC_M_J)
		if err == nil {
			t.Error("error")
			return
		}
		res, err = ctl.Exec_m("rcc_start_ping", util.Map{
			"delay": "ssdsfs",
		})
		if err != nil {
			t.Error("error")
			return
		}
		if res.IntVal("code") == 0 {
			t.Error("error")
			return
		}

	}
	err = ctl.RC_C_StartPing(500, 8, 102400)
	if err != nil {
		t.Error(err.Error())
		return
	}
	res, err = ctl.RC_C_Status()
	if err != nil {
		t.Error(err.Error())
		return
	}
	if res.IntVal("ping/ping_c") < 0 {
		t.Error("not running")
		return
	}
	if res.IntVal("ping/ping_s/running") < 0 {
		t.Error("not running")
		return
	}
	fmt.Println(util.S2Json(res), "--->")
	{
		err = ctl.RC_C_StartPing(500, 8, 102400)
		if err == nil {
			t.Error("error")
			return
		}
	}
	err = ctl.RC_C_StopPing()
	if err != nil {
		t.Error(err.Error())
		return
	}
	{
		err = ctl.RC_C_StopPing()
		if err == nil {
			t.Error("error")
			return
		}
	}
	err = ctl.RC_C_Stop()
	if err != nil {
		t.Error(err.Error())
		return
	}
	{
		err = ctl.RC_C_Stop()
		if err == nil {
			t.Error("error")
			return
		}
		err = ctl.RC_C_StartPing(500, 8, 102400)
		if err == nil {
			t.Error("error")
			return
		}
		err = ctl.RC_C_StopPing()
		if err == nil {
			t.Error("error")
			return
		}
		err = ctl.RC_C_Start("127.0.0.1:28773", "xxx", RC_M_J)
		if err == nil {
			t.Error("error")
			return
		}
	}

	err = ctl.RC_S_Stop()
	if err != nil {
		t.Error(err.Error())
		return
	}
	{
		err = ctl.RC_S_Stop()
		if err == nil {
			t.Error("error")
			return
		}
		err = ctl.RC_S_Token([]string{"abc"})
		if err == nil {
			t.Error("error")
			return
		}
		err = ctl.RC_S_Token([]string{""})
		if err == nil {
			t.Error("error")
			return
		}
	}
}

func TestRCHErr(t *testing.T) {
	ctl := &RC_CTL{
		Exec_m: func(name string, args interface{}) (util.Map, error) {
			return nil, util.Err("error")
		},
	}
	ctl.RC_S_Start("addr", "m")
	ctl.RC_S_Token([]string{"ss"})
	ctl.RC_S_Stop()
	ctl.RC_C_Start("addr", "token", "m")
	ctl.RC_C_Stop()
	ctl.RC_C_StartPing(500, 8, 8)
	ctl.RC_C_StopPing()
}
