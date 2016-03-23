package cookie

import (
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/Centny/gwf/routing"
	"golang.org/x/net/publicsuffix"
	"net/http"
	"net/http/cookiejar"
	"testing"
)

type Abc struct {
	A string
}

func TestSessionMux(t *testing.T) {
	sb := NewCookieSessionBuilder("", "/t")
	mux := routing.NewSessionMux("/t", sb)
	//
	mux.HFilterFunc("^/a(\\?.*)?$", func(hs *routing.HTTPSession) routing.HResult {
		hs.S.Set("testing", "abc")
		hs.S.Flush()
		hs.S.Flush()
		return routing.HRES_CONTINUE
	})
	mux.HFilterFunc("^/c(\\?.*)?$", func(hs *routing.HTTPSession) routing.HResult {
		return routing.HRES_CONTINUE
	})
	mux.HFilterFunc("^/b/[^\\?/]*(\\?.*)?$", func(hs *routing.HTTPSession) routing.HResult {
		fmt.Println(hs.R.URL.String())
		if hs.S.Val("testing") != "abc" {
			t.Error("not error testing")
		}
		hs.S.Set("testing", nil)
		if hs.S.Val("testing") != nil {
			t.Error("not nil")
		}
		return routing.HRES_CONTINUE
	})
	mux.HFunc("^/a(\\?.*)?$", func(hs *routing.HTTPSession) routing.HResult {
		if hs.S.Val("testing") != "abc" {
			t.Error("not error testing")
		}
		return routing.HRES_CONTINUE
	})
	http.Handle("/t/", mux)
	//
	esb := NewCookieSessionBuilder("", "/e")
	emux := routing.NewSessionMux("/e", esb)
	//
	emux.HFilterFunc("^/a(\\?.*)?$", func(hs *routing.HTTPSession) routing.HResult {
		hs.S.Set("testing", "abc")
		fmt.Println(hs.S.Flush())
		return routing.HRES_CONTINUE
	})
	emux.HFilterFunc("^/b(\\?.*)?$", func(hs *routing.HTTPSession) routing.HResult {
		if hs.S.Val("testing") != nil {
			t.Error("not nil")
		}
		return routing.HRES_CONTINUE
	})
	emux.HFilterFunc("^/c(\\?.*)?$", func(hs *routing.HTTPSession) routing.HResult {
		hs.S.Set("testing", "abc")
		hs.S.Flush()
		return routing.HRES_CONTINUE
	})
	emux.HFilterFunc("^/d/[^\\?/]*(\\?.*)?$", func(hs *routing.HTTPSession) routing.HResult {
		fmt.Println(hs.R.URL.String())
		if hs.S.Val("testing") != nil {
			t.Error("not nil")
		}
		return routing.HRES_CONTINUE
	})
	http.Handle("/e/", emux)
	//
	go http.ListenAndServe(":2789", nil)
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
	c.Get("http://127.0.0.1:2789/t/b/1")
	c.Get("http://127.0.0.1:2789/t/c")
	//
	esb.Crypto = func(bys []byte) ([]byte, error) {
		return bys, errors.New("crypto error")
	}
	c.Get("http://127.0.0.1:2789/e/a")
	c.Get("http://127.0.0.1:2789/e/b")
	esb.Crypto = func(bys []byte) ([]byte, error) {
		return bys, nil
	}
	esb.UnCrypto = func(bys []byte) ([]byte, error) {
		return bys, errors.New("uncrypto error")
	}
	c.Get("http://127.0.0.1:2789/e/c")
	c.Get("http://127.0.0.1:2789/e/d/1")
	//
	esb.Crypto = func(bys []byte) ([]byte, error) {
		return bys, errors.New("uncrypto error")
	}
	esb.UnCrypto = func(bys []byte) ([]byte, error) {
		return bys, errors.New("uncrypto error")
	}
	cs := &CookieSession{
		Sb:      esb,
		kvs:     map[string]interface{}{},
		updated: false,
	}
	_, err = cs.Crypto()
	if err == nil {
		t.Error("not error")
	}
	cs.Set("abc", &Abc{})
	_, err = cs.Crypto()
	if err == nil {
		t.Error("not error")
	}
	cs.UnCrypto("")
	cs.UnCrypto("测试")
	cs.UnCrypto("e1f")
	cs.UnCrypto(hex.EncodeToString([]byte("kkjjjjj")))
	esb.UnCrypto = func(bys []byte) ([]byte, error) {
		return bys, nil
	}
	cs.UnCrypto(hex.EncodeToString([]byte("kkjjjjj")))
	// cs.kvs["abc"] = 1235
	// cs.kvs["dg"] = "kkkkk"

	//
	fmt.Println("TestSessionMux end")
}
