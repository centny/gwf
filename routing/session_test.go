package routing

import (
	"code.google.com/p/go.net/publicsuffix"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"testing"
	"time"
)

type CSrv struct {
	Count int
	Res   HResult
}

func (s *CSrv) SrvHTTP(hs *HTTPSession) HResult {
	s.Count = s.Count + 1
	hs.S.Set("abc", "123456789")
	fmt.Println(hs.S.Val("abc"))
	fmt.Println(hs.S.Val("abc1"))
	hs.S.Set("kkk", nil)
	return s.Res
}

func TestSessionMux(t *testing.T) {
	sb := NewSrvSessionBuilder("", "/", 2000, 500)
	mux := NewSessionMux("/t", sb)
	// mux.CDelay = 500
	ssrv1 := Ssrv{Count: 0}
	csrv1 := CSrv{Count: 0, Res: HRES_CONTINUE}
	csrv2 := CSrv{Count: 0, Res: HRES_RETURN}
	count := 0
	mux.Handler("^/a(\\?.*)?$", &ssrv1)
	mux.HandleFunc("^/a(\\?.*)?$", func(w http.ResponseWriter, r *http.Request) {
		count = count + 1
		w.Write([]byte("abc"))
		fmt.Println(mux.RSession(r))
	})
	mux.HFilter("^/a(\\?.*)?$", &csrv1)
	mux.HFilterFunc("^/a(\\?.*)?$", func(hs *HTTPSession) HResult {
		count = count + 1
		fmt.Println(sb.Session(hs.S.(*SrvSession).Token()))
		return HRES_CONTINUE
	})
	mux.HFilter("^/r1(\\?.*)?$", &csrv2)
	mux.HFilterFunc("^/r2(\\?.*)?$", func(hs *HTTPSession) HResult {
		count = count + 1
		fmt.Println(sb.Session(hs.S.(*SrvSession).Token()))
		return HRES_RETURN
	})
	mux.H("^/a(\\?.*)?$", &csrv1)
	mux.HFunc("^/a(\\?.*)?$", func(hs *HTTPSession) HResult {
		count = count + 1
		return HRES_CONTINUE
	})
	mux.H("^/r3(\\?.*)?$", &csrv2)
	mux.HFunc("^/r4(\\?.*)?$", func(hs *HTTPSession) HResult {
		count = count + 1
		return HRES_RETURN
	})
	sb.StartLoop()
	//
	http.Handle("/t/", mux)
	http.Handle("/t2/", mux)
	http.Handle("/abc/", mux)
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
	c.Get("http://127.0.0.1:2789/t/b")
	c.Get("http://127.0.0.1:2789/t2/b")
	fmt.Println(ssrv1.Count, csrv1.Count, count)
	if ssrv1.Count != 1 || csrv1.Count != 2 || count != 3 {
		t.Error("count error")
	}
	time.Sleep(3 * time.Second)
	c.Get("http://127.0.0.1:2789/t/a")
	c.Get("http://127.0.0.1:2789/t/b")
	c.Get("http://127.0.0.1:2789/t2/b")
	c.Get("http://127.0.0.1:2789/t/r1")
	c.Get("http://127.0.0.1:2789/t/r2")
	c.Get("http://127.0.0.1:2789/t/r3")
	c.Get("http://127.0.0.1:2789/t/r4")
	c.Get("http://127.0.0.1:2789/abc/a")
	//
	if mux.RSession(nil) != nil {
		t.Error("not nil")
		return
	}
	if sb.Session("ss") != nil {
		t.Error("not nil")
		return
	}
	//
	sb.StopLoop()
	//
	NewSessionMux("/", nil)
	//
	fmt.Println("TestSessionMux end")
}

// func TestCookie(t *testing.T) {
// 	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
// 		c, err := r.Cookie("token3")
// 		fmt.Println(c, err)

// 		// expire := time.Now().AddDate(0, 0, 1)
// 		c = &http.Cookie{}
// 		c.Name = "token3"
// 		c.Value = "tokenvalue"
// 		c.Path = "/"
// 		// c.Domain = "127.0.0.1"
// 		fmt.Println(len(c.Domain))
// 		http.SetCookie(w, c)
// 	})
// 	go http.ListenAndServe(":2789", nil)

// 	time.Sleep(300 * time.Second)
// }
