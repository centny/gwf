package filter

import (
	"fmt"
	"github.com/Centny/gwf/netw/rc/rctest"
	"github.com/Centny/gwf/routing"
	"github.com/Centny/gwf/routing/httptest"
	"github.com/Centny/gwf/util"
	"testing"
	"time"
)

func TestAttr(t *testing.T) {
	var rt = rctest.NewRCTest_j2(":25322")
	var af = NewAttrFilter("tk")
	af.ShowLog = true
	af.Delay = 100
	af.Timeout = 1000
	af.StartTimeout()
	af.Hand(rt.L)
	tk, err := RegisterAttr(rt.R, util.Map{"abc": 1})
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(tk)
	var ts = httptest.NewMuxServer()
	ts.Mux.HFilter("^.*$", af)
	ts.Mux.HFunc("^/abc.*$", func(hs *routing.HTTPSession) routing.HResult {
		var val = hs.IntVal("abc")
		if val != 1 {
			t.Error("error")
		}
		return hs.MsgRes("OK")
	})
	res, _ := ts.G2("/abc?tk=%v", tk)
	if res.IntVal("code") != 0 {
		t.Error("error")
	}
	time.Sleep(2 * time.Second)
	res, _ = ts.G2("/abc?tk=%v", tk)
	if res.IntVal("code") == 0 {
		t.Error("error")
	}
	res, _ = ts.G2("/abc?tk=%v", "")
	if res.IntVal("code") == 0 {
		t.Error("error")
	}
	fmt.Println(res)

}

func TestSs(t *testing.T) {
}
