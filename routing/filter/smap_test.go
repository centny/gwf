package filter

import (
	"fmt"
	"github.com/Centny/gwf/routing"
	"github.com/Centny/gwf/routing/httptest"
	"github.com/Centny/gwf/util"
	"testing"
)

func TestSMap(t *testing.T) {
	smap := NewSMap()
	ts := httptest.NewServer2(smap)
	ts.G2("")
	fmt.Println(ts.G2("?key=abc&val=123"))
	tv, _ := ts.G2("?key=abc")
	if "123" == tv.StrVal("data") {
		t.Error("sdfss")
	}
	ts.G2("?key=abc&val=")
	fmt.Println(ts.G2("?key=abc"))
	smap.GET = func(hs *routing.HTTPSession, key string) (string, error) {
		return "", util.Err("sdfsfdsf->")
	}
	fmt.Println(ts.G2("?key=abc"))
}
