package main

import (
	"bytes"
	"code.google.com/p/go.net/publicsuffix"
	"fmt"
	"github.com/Centny/gwf/routing/httptest"
	"github.com/Centny/gwf/util"
	"net/http/cookiejar"
	"net/url"
	"os"
	"testing"
)

func TestWebts(t *testing.T) {
	var ts = httptest.NewMuxServer()
	Hand("", ts.Mux)
	os.Mkdir("test", os.ModePerm)
	WWW = "test/"
	//
	var res, err = ts.G2("/g_args?a=%v&b=%v&c=%v", 1, 2, "xx")
	if err != nil {
		t.Error(err.Error())
		return
	}
	if res.IntVal("a") != 1 || res.IntVal("b") != 2 || res.StrVal("c") != "xx" {
		t.Error("error")
	}
	//
	var uargs = url.Values{}
	uargs.Set("a", "1")
	uargs.Set("b", "2")
	uargs.Set("c", "xx")
	_, res, _, err = ts.PostFormV2("/p_args", nil, uargs)
	if err != nil {
		t.Error(err.Error())
		return
	}
	if res.IntVal("a") != 1 || res.IntVal("b") != 2 || res.StrVal("c") != "xx" {
		t.Error("error")
	}
	//
	var margs = map[string]string{}
	margs["a"] = "1"
	margs["b"] = "2"
	margs["c"] = "xx"
	res, err = ts.PostF2("/m_args", "", "", margs)
	if err != nil {
		t.Error(err.Error())
		return
	}
	if res.IntVal("a") != 1 || res.IntVal("b") != 2 || res.StrVal("c") != "xx" {
		t.Error("error")
	}
	//
	sha1, err := util.Sha1("api.go")
	if err != nil {
		t.Error("error")
		return
	}
	res, err = ts.PostF2("/upload", "file", "api.go", nil)
	if err != nil {
		t.Error("error")
		return
	}
	if res.StrVal("fn") != "api.go" || res.StrVal("sha") != sha1 {
		t.Error("error")
		return
	}
	sha1_, err := util.Sha1("api.go")
	if err != nil {
		t.Error("error")
		return
	}
	if sha1 != sha1_ {
		t.Error("error")
		return
	}
	//
	res, err = ts.PostN2("/body", "application/json", bytes.NewBufferString(util.S2Json(util.Map{
		"a": 1,
		"b": 2,
		"c": "xx",
	})))
	if err != nil {
		t.Error(err.Error())
		return
	}
	if res.IntVal("a") != 1 || res.IntVal("b") != 2 || res.StrVal("c") != "xx" {
		t.Error("error")
	}
	//
	_, res, _, err = ts.PostFormV2("/req_ctype", map[string]string{
		"A": "1",
		"B": "2",
		"C": "xx",
	}, nil)
	if err != nil {
		t.Error(err.Error())
		return
	}
	if res.StrVal("Content-Type") != "application/x-www-form-urlencoded" {
		t.Error("error")
	}
	//
	_, _, rh, err := util.HTTPClient.DoGet2(nil, "%v/res_ctype?a=%v&b=%v&c=%v", ts.URL, 1, 2, "xx")
	if err != nil {
		t.Error(err.Error())
		return
	}
	if rh["A"] != "1" || rh["B"] != "2" || rh["C"] != "xx" {
		t.Error("error")
		return
	}
	//
	//
	options := cookiejar.Options{
		PublicSuffixList: publicsuffix.List,
	}
	jar, err := cookiejar.New(&options)
	if err != nil {
		t.Error(err.Error())
		return
	}
	util.HTTPClient.Jar = jar
	res, err = ts.G2("/s_ss?xa=%v&xb=%v&xc=%v", 1, 2, "xx")
	if err != nil {
		t.Error(err.Error())
		return
	}
	if res.IntVal("xa") != 1 || res.IntVal("xb") != 2 || res.StrVal("xc") != "xx" {
		t.Error("error")
	}
	//
	res, err = ts.G2("/g_ss?keys=%v", "xa,xb,xc")
	if err != nil {
		t.Error(err.Error())
		return
	}
	if res.IntVal("xa") != 1 || res.IntVal("xb") != 2 || res.StrVal("xc") != "xx" {
		t.Error("error")
	}
	fmt.Println("webts done...")
}
