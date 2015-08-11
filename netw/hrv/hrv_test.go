package hrv

import (
	"fmt"
	"github.com/Centny/gwf/netw"
	"github.com/Centny/gwf/netw/impl"
	"github.com/Centny/gwf/pool"
	"github.com/Centny/gwf/routing"
	"github.com/Centny/gwf/routing/httptest"
	"github.com/Centny/gwf/util"
	"runtime"
	"testing"
	"time"
)

type hrv_h struct {
}

func (h *hrv_h) OnConn(c netw.Con) bool {
	return true
}
func (h *hrv_h) OnClose(c netw.Con) {
}
func (h *hrv_h) OnCmd(c netw.Cmd) int {
	return -1
}
func (h *hrv_h) OnLogin(token, name, alias string) error {
	if token == "tacc" {
		return nil
	} else {
		return util.Err("invalid token")
	}
}
func TestHrv(t *testing.T) {
	runtime.GOMAXPROCS(runtime.NumCPU())
	ts := httptest.NewMuxServer()
	ts.Mux.HFunc("^/ax/a1(\\?.*)?$", func(hs *routing.HTTPSession) routing.HResult {
		hs.R.ParseForm()
		kvs := util.Map{}
		for k, v := range hs.R.Form {
			kvs[k] = v[0]
		}
		for k, v := range hs.R.Header {
			kvs[k] = v[0]
		}
		return hs.MsgRes(kvs)
	})
	ts.Mux.HFunc("^/ax/a2(\\?.*)?$", func(hs *routing.HTTPSession) routing.HResult {
		hs.R.ParseForm()
		hs.R.ParseMultipartForm(10240)
		kvs := util.Map{}
		for k, v := range hs.R.Form {
			kvs[k] = v[0]
		}
		for k, v := range hs.R.Header {
			kvs[k] = v[0]
		}
		return hs.MsgRes(kvs)
	})
	bp := pool.NewBytePool(8, 10240)
	hs := NewHrvS_j(bp, ":8234")
	hs.H = &hrv_h{}
	hs.Pre = "/hh"
	hs.ShowLog = true
	hs.Headers = map[string]bool{
		"HABC1": true,
		"HABC2": true,
	}
	hs.Args["aa"] = 123
	hs.Args["bb"] = 321
	hs.Head["Abc1"] = 123
	hs.Head["Abc2"] = 321
	hs.SetWww(".")
	hs.AddPattern("^a.*$")
	tss := httptest.NewMuxServer()
	tss.Mux.HFunc("^/.*$", hs.Doh)
	err := hs.Run()
	if err != nil {
		t.Error(err.Error())
		return
	}
	hc := NewHrvC_j(bp, "127.0.0.1:8234", fmt.Sprintf("%v/ax", ts.URL))
	hc.ShowLog = true
	hc.Token = "tacc"
	hc.Name = "ax"
	hc.H = hs.H
	hc.Start()
	time.Sleep(time.Second)
	tss.G("/hh?xx=1&yy=2&h:Hx=abc&h:Hy=123")
	res, err := tss.G2("/hh/ax/a1?x=1&y=2")
	if err != nil || res.IntVal("code") != 0 ||
		res.StrValP("/data/x") != "1" || res.StrValP("/data/y") != "2" ||
		res.StrValP("/data/xx") != "1" || res.StrValP("/data/yy") != "2" ||
		res.StrValP("/data/Hx") != "abc" || res.StrValP("/data/Hy") != "123" ||
		res.StrValP("/data/aa") != "123" || res.StrValP("/data/bb") != "321" ||
		res.StrValP("/data/Abc1") != "123" || res.StrValP("/data/Abc2") != "321" {
		t.Error(fmt.Sprintf("%v/%v", res, err))
		return
	}
	res, err = util.HPostFv2(tss.URL+"/hh/ax/a2", map[string]string{
		"x": "11",
		"y": "12",
	}, map[string]string{
		"Habc1": "val1",
		"Habc2": "val2",
	}, "", "")
	if err != nil || res.IntVal("code") != 0 ||
		res.StrValP("/data/x") != "11" || res.StrValP("/data/y") != "12" ||
		res.StrValP("/data/Habc1") != "val1" || res.StrValP("/data/Habc2") != "val2" ||
		res.StrValP("/data/xx") != "1" || res.StrValP("/data/yy") != "2" ||
		res.StrValP("/data/Hx") != "abc" || res.StrValP("/data/Hy") != "123" ||
		res.StrValP("/data/aa") != "123" || res.StrValP("/data/bb") != "321" ||
		res.StrValP("/data/Abc1") != "123" || res.StrValP("/data/Abc2") != "321" {
		t.Error(fmt.Sprintf("%v/%v", res, err))
		return
	}
	res_, err := tss.G("/hh/ax")
	if err != nil {
		t.Error(err.Error())
		return
	}
	fmt.Println(res_)
	res, err = tss.G2("/hh/notsss/a1?x=1&y=2")
	if err != nil || res.IntVal("code") == 0 {
		t.Error("not error")
		return
	}
	hc.HB()
	hc.Login(hc.Token, hc.Name, hc.Alias)
	//
	hs.MsgC("ax").Writeb([]byte("server->"))
	hc.MC.Writeb([]byte("client->"))

	//
	//test login
	lg := NewHrvC_j(bp, "127.0.0.1:8234", fmt.Sprintf("%v/ax", ts.URL))
	lg.ShowLog = true
	go func() {
		time.Sleep(200 * time.Millisecond)
		lg.Timeout()
	}()
	go lg.HB()
	err = lg.Login("", "", "")
	if err == nil {
		t.Error("not error")
		return
	}
	lg.Start()
	err = lg.Login("token", "name", "alias")
	if err == nil {
		t.Error("not error")
		return
	}
	err = lg.Login("token", "", "alias")
	if err == nil {
		t.Error("not error")
		return
	}
	lg.Doh(&impl.RCM_Cmd{
		Map: &util.Map{
			"M": "PUT",
		},
	})
	//
	//base empty error
	hc.Base = ""
	tss.G2("/hh/ax/a1?x=1&y=2")
	hs.F = nil
	fmt.Println(tss.G("/hh/ax"))
	//
	//close
	hc.Close()
	lg.Close()
	hs.Close()
	//
	//test error for coverage.
	tmap := util.Map{}
	tres := &Res{}
	HRV_B2V(nil, tres)
	HRV_B2V(nil, tmap)
	HRV_V2B(tres)
	HRV_V2B(HRV_V2B)
	tres.String()
	tres = nil
	tres.GetCode()
	tres.GetData()
	//
	hs.H = nil
	hs.OnCmd(nil)
	hs.OnConn(nil)
	hs.OnClose(nil)
	hs.onlogin("token", "name", "alias")
	hc.H = nil
	hc.OnCmd(nil)
	hc.OnConn(&Con_{Con_: &netw.Con_{}})
	hs.Parse("{{html}")
}

type Con_ struct {
	*netw.Con_
}

func (c *Con_) SetWait(t bool) {
}
