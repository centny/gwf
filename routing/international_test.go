package routing

import (
	"fmt"
	"github.com/Centny/Cny4go/util"
	"net/http/httptest"
	"testing"
)

func TestINT(t *testing.T) {
	sb := NewSrvSessionBuilder("", "/", "rtest", 2000, 500)
	// sb.ShowLog = true
	mux := NewSessionMux("/t", sb)
	// mux.ShowLog = true
	INT, _ := NewJsonINT(".")
	INT.Default = "en"
	mux.HFunc("^/info.*$", func(hs *HTTPSession) HResult {
		fmt.Println(hs.LocalVal("abc"))
		return hs.MsgRes("OK")
	})
	mux.HFunc("^/set.*$", func(hs *HTTPSession) HResult {
		hs.SetLocal("zh")
		return hs.MsgRes("OK")
	})
	header := map[string]string{}
	ts := httptest.NewServer(mux)
	fmt.Println(util.HTTPClient.HGet_H(header, "%s/t/info", ts.URL))
	fmt.Println(util.HTTPClient.HGet_H(header, "%s/t/set", ts.URL))
	mux.INT = INT
	fmt.Println(util.HTTPClient.HGet_H(header, "%s/t/info", ts.URL))
	header["Accept-Language"] = "zh-cn,zh;q=0.8,en;q=0.3,en-us;q=0.5"
	fmt.Println(util.HTTPClient.HGet_H(header, "%s/t/info", ts.URL))
	header["Accept-Language"] = "zh-cn,zh;q=0.8,en;q=0.9,en-us;q=0.5"
	fmt.Println(util.HTTPClient.HGet_H(header, "%s/t/info", ts.URL))
	header["Accept-Language"] = "zh-cn,zh;q=0.8,en;q=0.3,en-us;q=0.9"
	fmt.Println(util.HTTPClient.HGet_H(header, "%s/t/info", ts.URL))
	header["Accept-Language"] = "zh-cn,zh;"
	fmt.Println(util.HTTPClient.HGet_H(header, "%s/t/info", ts.URL))
	header["Accept-Language"] = "zh-cn,zh;q=0.8,en;q=0.3,abc;q=0.9"
	fmt.Println(util.HTTPClient.HGet_H(header, "%s/t/info", ts.URL))
	//
	fmt.Println(util.HTTPClient.HGet_H(header, "%s/t/set", ts.URL))
	header["Accept-Language"] = "zh-cn,zh;q=0.8,en;q=0.9"
	fmt.Println(util.HTTPClient.HGet_H(header, "%s/t/info", ts.URL))
}
