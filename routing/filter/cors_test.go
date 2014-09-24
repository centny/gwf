package filter

import (
	"github.com/Centny/gwf/routing/httptest"
	"net/http"
	"testing"
)

func TestCors(t *testing.T) {
	cors := NewCORS()
	ts := httptest.NewServer2(cors)
	client := &http.Client{}
	//not origin
	req, _ := http.NewRequest("GET", ts.URL, nil)
	client.Do(req)
	//specified origin not access.
	req, _ = http.NewRequest("GET", ts.URL, nil)
	req.Header.Set("Origin", ts.URL)
	res, _ := client.Do(req)
	if res.StatusCode != http.StatusForbidden {
		t.Error("not right")
		return
	}
	//specified origin access
	cors.AddSite(ts.URL)
	res, _ = client.Do(req)
	if res.StatusCode != http.StatusOK {
		t.Error("not right")
		return
	}
	if res.Header.Get("Access-Control-Allow-Origin") != ts.URL {
		t.Error("not right")
		return
	}
	cors.DelSite(ts.URL)
	//
	//all access
	cors.AddSite("*")
	res, _ = client.Do(req)
	if res.StatusCode != http.StatusOK {
		t.Error("not right")
		return
	}
	if res.Header.Get("Access-Control-Allow-Origin") != "*" {
		t.Error("not right")
		return
	}
	//option require
	req, _ = http.NewRequest("OPTIONS", ts.URL, nil)
	req.Header.Set("Origin", ts.URL)
	res, _ = client.Do(req)
	if res.StatusCode != http.StatusOK {
		t.Error("not right")
		return
	}
}

func TestNormal(t *testing.T) {
	httptest.Tnh(NewCORS(), "")
}
