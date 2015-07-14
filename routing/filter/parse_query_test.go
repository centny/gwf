package filter

import (
	"fmt"
	"github.com/Centny/gwf/routing"
	"github.com/Centny/gwf/util"
	"net/http/httptest"
	"testing"
)

func trec_f(hs *routing.HTTPSession) routing.HResult {
	var a string
	err := hs.ValidCheckVal(`
		a,R|S,L:0;
		`, &a)
	if err != nil {
		return hs.MsgResErr2(1, "arg-err", err)
	}
	fmt.Println("A->", a)
	_, err = hs.RecF("file", "abc2.txt")
	if err == nil {
		return hs.MsgRes("OK")
	} else {
		return hs.MsgResErr2(1, "srv-err", err)
	}
}

func TestParseQuery(t *testing.T) {
	mux := routing.NewSessionMux2("")
	mux.HFilterFunc("^.*$", ParseQuery)
	mux.HFunc("^.*$", trec_f)
	ts := httptest.NewServer(mux)
	util.FWrite("abc.txt", "abc-123")
	fmt.Println(util.HPostF2(ts.URL+"?a=abc", nil, "file", "abc.txt"))
}
