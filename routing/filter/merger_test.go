package filter

import (
	"testing"

	"fmt"

	"bytes"

	"github.com/Centny/gwf/routing"
	"github.com/Centny/gwf/routing/httptest"
	"github.com/Centny/gwf/util"
)

func TestSimpleMerger(t *testing.T) {
	var ts = httptest.NewMuxServer()
	fcfg, err := util.NewFcfg("merger.properties?TSRV=" + ts.URL)
	if err != nil {
		t.Error(err)
		return
	}
	ts.Mux.HFunc("^/api/a(\\?.*)?$", func(hs *routing.HTTPSession) routing.HResult {
		if len(hs.CheckVal("merger")) < 1 {
			return hs.MsgRes2(1, "error")
		}
		return hs.MsgRes(hs.RVal("keya"))
	})
	ts.Mux.HFunc("^/api/b(\\?.*)?$", func(hs *routing.HTTPSession) routing.HResult {
		if len(hs.CheckVal("merger")) < 1 {
			return hs.MsgRes2(1, "error")
		}
		return hs.MsgRes(util.Map{
			"val": hs.RVal("keyb"),
		})
	})
	ts.Mux.HFunc("^/api/c(\\?.*)?$", func(hs *routing.HTTPSession) routing.HResult {
		if len(hs.CheckVal("merger")) < 1 {
			return hs.MsgRes2(1, "error")
		}
		var res = util.Map{}
		var err = hs.UnmarshalJ(&res)
		if err != nil {
			fmt.Println(err)
			return hs.MsgResErr(1, "arg-err", err)
		}
		return hs.MsgRes(res)
	})
	ts.Mux.HFunc("^/api/err(\\?.*)?$", func(hs *routing.HTTPSession) routing.HResult {
		if len(hs.CheckVal("merger")) < 1 {
			return hs.MsgRes2(1, "error")
		}
		return hs.MsgResErr(1, "arg-err", fmt.Errorf("%v", "error"))
	})
	ts.Mux.HFunc("^/api/err2(\\?.*)?$", func(hs *routing.HTTPSession) routing.HResult {
		panic("error")
	})
	HandMerger(ts.Mux, fcfg)
	var rdata = util.Map{
		"c": util.Map{
			"ival": 123,
			"sval": "abc",
		},
	}
	mres, err := ts.PostN2(
		"/merge/t1?merger=%v&a.keya=%v&b.keyb=%v&",
		"application/json", bytes.NewBufferString(util.S2Json(rdata)),
		"a,b,c", "123", "abc")
	if err != nil {
		t.Error(err)
		return
	}
	if mres.IntVal("code") != 0 {
		fmt.Println(util.S2Json(mres))
		t.Error("error")
		return
	}
	if mres.StrValP("/data/a") != "123" {
		fmt.Println(util.S2Json(mres))
		t.Error("error a")
		return
	}
	if mres.StrValP("/data/b/val") != "abc" {
		fmt.Println(util.S2Json(mres))
		t.Error("error b")
		return
	}
	if mres.StrValP("/data/c/sval") != "abc" {
		fmt.Println(util.S2Json(mres))
		t.Error("error c")
		return
	}
	if mres.StrValP("/data/c/ival") != "123" {
		fmt.Println(util.S2Json(mres))
		t.Error("error c")
		return
	}

	//test reverse error
	mres, err = ts.PostN2(
		"/merge/t1?merger=%v&a.keya=%v&b.keyb=%v&",
		"application/json", bytes.NewBufferString(util.S2Json(rdata)),
		"a,b,c,err", "123", "abc")
	if err != nil {
		t.Error(err)
		return
	}
	if mres.IntVal("code") == 0 {
		fmt.Println(util.S2Json(mres))
		t.Error("error")
		return
	}
	//test reverse error2
	mres, err = ts.PostN2(
		"/merge/t1?merger=%v&a.keya=%v&b.keyb=%v&",
		"application/json", bytes.NewBufferString(util.S2Json(rdata)),
		"a,b,c,err2", "123", "abc")
	if err != nil {
		t.Error(err)
		return
	}
	if mres.IntVal("code") == 0 {
		fmt.Println(util.S2Json(mres))
		t.Error("error")
		return
	}
	//test reverse not merger key
	mres, err = ts.PostN2(
		"/merge/t1?merger=%v&a.keya=%v&b.keyb=%v&",
		"application/json", bytes.NewBufferString(util.S2Json(rdata)),
		"", "123", "abc")
	if err != nil {
		t.Error(err)
		return
	}
	if mres.IntVal("code") == 0 {
		fmt.Println(util.S2Json(mres))
		t.Error("error")
		return
	}
	//test reverse merger key not found
	mres, err = ts.PostN2(
		"/merge/t1?merger=%v&a.keya=%v&b.keyb=%v&",
		"application/json", bytes.NewBufferString(util.S2Json(rdata)),
		"xddd", "123", "abc")
	if err != nil {
		t.Error(err)
		return
	}
	if mres.IntVal("code") == 0 {
		fmt.Println(util.S2Json(mres))
		t.Error("error")
		return
	}
	//test hand not configure
	HandMerger(ts.Mux, util.NewFcfg3())
}
