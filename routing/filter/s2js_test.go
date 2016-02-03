package filter

import (
	"fmt"
	"github.com/Centny/gwf/routing"
	"github.com/Centny/gwf/routing/httptest"
	"github.com/Centny/gwf/util"
	"testing"
)

func TestS2js(t *testing.T) {
	var ts = httptest.NewMuxServer()
	ts.Mux.HFilterFunc("^/having.*$", func(hs *routing.HTTPSession) routing.HResult {
		hs.SetVal("user", util.Map{
			"xx": "abc",
			"aa": 123,
		})
		hs.SetVal("token", "token--->")
		fmt.Println("---->")
		return routing.HRES_CONTINUE
	})
	ts.Mux.H("^/having.*$", NewS2js("abc", []string{"user", "token"}))
	ts.Mux.H("^/not.*$", NewS2js("abc", []string{"user", "token"}))
	fmt.Println(ts.G("/not"))
	fmt.Println(ts.G("/having"))
}
