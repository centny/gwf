package session

import (
	"code.google.com/p/go.net/publicsuffix"
	"fmt"
	"github.com/Centny/Cny4go/routing"
	"net/http"
	"net/http/cookiejar"
	"testing"
)

type Abc struct {
	A string
}

func TestSessionMux(t *testing.T) {
	sb := NewRSessionBuilder("", "/t")
	mux := routing.NewSessionMux("/t", sb)
	//
	mux.HFilterFunc("^/a(\\?.*)?$", func(hs *routing.HTTPSession) routing.HResult {
		hs.S.Set("testing", "abc")
		hs.SetVal("kkk", nil)
		return routing.HRES_RETURN
	})
	mux.HFilterFunc("^/b(\\?.*)?$", func(hs *routing.HTTPSession) routing.HResult {
		if hs.StrVal("testing") != "abc" {
			t.Error("testing empty")
		}
		hs.StrVal("kkkkk")
		hs.S.Flush()
		return routing.HRES_RETURN
	})
	http.Handle("/t/", mux)
	//
	go http.ListenAndServe(":2799", nil)
	options := cookiejar.Options{
		PublicSuffixList: publicsuffix.List,
	}
	jar, err := cookiejar.New(&options)
	if err != nil {
		t.Error(err.Error())
		return
	}
	c := http.Client{Jar: jar}
	c.Get("http://127.0.0.1:2789/t/a")
	c.Get("http://127.0.0.1:2789/t/b")
	//
	//
	fmt.Println("TestSessionMux end")
}
