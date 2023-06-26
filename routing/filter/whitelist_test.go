package filter

import (
	"net/http/httptest"
	"testing"

	"github.com/Centny/gwf/routing"
	"github.com/Centny/gwf/util"
)

func TestWhiteList(t *testing.T) {

	f, err := NewWhitelistFilter("172.10.0.0/16")
	if err != nil {
		t.Error("New whitelist filter failed", err)
		return
	}

	cases := map[string]bool{
		"":                false,
		"172.10.1.12":     true,
		"172.10.1.13":     true,
		"172.10.2.13":     true,
		"172.11.1.12":     false,
		"192.0.1.1":       false,
		"211.123.123.123": false,
	}

	for ip, expect := range cases {
		if f.IsAllowed(ip) != expect {
			t.Errorf("ip(%v) match expected %v but %v", ip, expect, f.IsAllowed(ip))
		}
	}

	f, err = NewWhitelistFilter("192.168.0.0/16")
	if err != nil {
		t.Error("New whitelist filter failed", err)
		return
	}

	cases = map[string]bool{
		"":             false,
		"172.10.1.12":  false,
		"192.0.1.1":    false,
		"192.168.3.13": true,
	}

	for ip, expect := range cases {
		if f.IsAllowed(ip) != expect {
			t.Errorf("ip(%v) match expected %v but %v", ip, expect, f.IsAllowed(ip))
		}
	}

	f, err = NewWhitelistFilter("192.168.1.13")
	if err != nil {
		t.Error("New whitelist filter failed", err)
		return
	}

	cases = map[string]bool{
		"":             false,
		"172.10.1.12":  false,
		"192.168.3.13": false,
		"192.168.1.13": true,
	}

	for ip, expect := range cases {
		if f.IsAllowed(ip) != expect {
			t.Errorf("ip(%v) match expected %v but %v", ip, expect, f.IsAllowed(ip))
		}
	}

	if !IsIPAddress("192.168.0.1") {
		t.Error("192.168.0.1 should be ip")
	}

	f, err = NewWhitelistFilter("192.168.0.1", "172.10.0.1/16")
	if err != nil {
		t.Error("New whitelist filter failed", err)
		return
	}
	mux := routing.NewSessionMux2("")
	mux.HFilter("^.*$", f)
	mux.HFunc("^.*$", func(hs *routing.HTTPSession) routing.HResult {
		return hs.MsgRes("ok")
	})
	ts := httptest.NewServer(mux)
	res, err := util.HGet2("%v", ts.URL)
	if err != nil {
		t.Error(err)
	}
	if res.IntVal("code") != 401 {
		t.Error("code should be 401")
	}

	f.AddIPRange("127.0.0.0/8")
	res, err = util.HGet2("%v", ts.URL)
	if err != nil {
		t.Error(err)
	}
	if res.IntVal("code") == 401 {
		t.Error("code should not be 401")
	}

	f, err = NewWhitelistFilter("192.168.0.1", "127.0.0.1/16")
	if err != nil {
		t.Error("New whitelist filter failed", err)
		return
	}
	mux = routing.NewSessionMux2("")
	mux.HFilter("^.*$", f)
	mux.HFunc("^.*$", func(hs *routing.HTTPSession) routing.HResult {
		return hs.MsgRes("ok")
	})
	ts = httptest.NewServer(mux)
	res, err = util.HGet2("%v", ts.URL)
	if err != nil {
		t.Error(err)
	}
	if res.IntVal("code") == 401 {
		t.Error("code should not be 401")
	}
}
