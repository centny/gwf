package rctest

import (
	"fmt"
	"github.com/Centny/gwf/netw/impl"
	"github.com/Centny/gwf/util"
	"testing"
)

func rc_run(rc *impl.RCM_Cmd) (interface{}, error) {
	return rc.Map, nil
}
func TestRCTest(t *testing.T) {
	rct := NewRCTest_j2(":12334")
	rct.ShowLog(true)
	rct.L.AddHFunc("abc", rc_run)
	res, err := rct.R.VExec_m("abc", util.Map{
		"a": 1,
		"b": "x",
		"c": "kkk",
	})
	if err != nil {
		t.Error(err.Error())
		return
	}
	if res.IntVal("a") != 1 || res.StrVal("b") != "x" || res.StrVal("c") != "kkk" {
		t.Error("error")
		return
	}
	rct.Runner()
	rct.Listener()
	fmt.Print("TestRCTest...")
}
