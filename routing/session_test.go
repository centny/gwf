package routing

import (
	"code.google.com/p/go.net/publicsuffix"
	"errors"
	"fmt"
	"github.com/Centny/gwf/util"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"net/url"
	"regexp"
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
	hs.S.Set("int", int(123))
	hs.S.Set("int8", int8(123))
	hs.S.Set("int16", int16(123))
	hs.S.Set("int32", int32(123))
	hs.S.Set("int64", int64(123))
	hs.S.Set("uint", uint(123))
	hs.S.Set("uint8", uint8(123))
	hs.S.Set("uint16", uint16(123))
	hs.S.Set("uint32", uint32(123))
	hs.S.Set("uint64", uint64(123))
	hs.S.Set("float32", float32(123.34))
	hs.SetVal("float64", float64(123.34))
	fmt.Println(hs.S.Val("abc"))
	fmt.Println(hs.S.Val("abc1"))
	fmt.Println(hs.FloatVal("abc"))
	fmt.Println(hs.IntVal("abc"))
	fmt.Println(hs.StrVal("int"))
	//
	fmt.Println(hs.IntVal("int"))
	fmt.Println(hs.IntVal("int8"))
	fmt.Println(hs.IntVal("int16"))
	fmt.Println(hs.IntVal("int32"))
	fmt.Println(hs.IntVal("int64"))
	fmt.Println(hs.IntVal("intsss"))
	//
	fmt.Println(hs.UintVal("uint"))
	fmt.Println(hs.UintVal("uint8"))
	fmt.Println(hs.UintVal("uint16"))
	fmt.Println(hs.UintVal("uint32"))
	fmt.Println(hs.UintVal("uint64"))
	fmt.Println(hs.UintVal("float32"))
	fmt.Println(hs.UintVal("intsss"))
	//
	fmt.Println(hs.FloatVal("float32"))
	fmt.Println(hs.FloatVal("float64"))
	fmt.Println(hs.FloatVal("floss"))
	fmt.Println(hs.StrVal("abc"))
	fmt.Println(hs.CheckVal("abc"))
	fmt.Println(hs.StrVal("abcss"))
	fmt.Println(hs.CheckVal("abcss"))
	fmt.Println(hs.Host())
	hs.S.Set("kkk", nil)
	fmt.Println(hs.Val("kkk"))
	//
	var iv int64
	err := hs.ValidCheckVal("int,R|I,R:50~300", &iv)
	fmt.Println(err, iv)
	if iv != 123 {
		panic("hava error")
	}
	hs.ValidCheckVal("int,R|I,R:50!300", &iv)
	hs.ValidRVal("int,R|I,R:50~300", &iv)
	hs.ValidRValN("int,R|I,R:50~300", &iv)
	hs.ValidCheckValN("int,R|I,R:50~300", &iv)
	hs.Cookie("key")
	hs.SetCookie("kk", "sfsf")
	hs.Cookie("kk")
	hs.SetCookie("kk", "")
	return s.Res
}

func TestSessionMux(t *testing.T) {
	sb := NewSrvSessionBuilder("", "/", "rtest", 2000, 500)
	sb.ShowLog = true
	mux := NewSessionMux("/t", sb)
	mux.ShowLog = true
	// mux.CDelay = 500
	ssrv1 := Ssrv{Count: 0}
	csrv1 := CSrv{Count: 0, Res: HRES_CONTINUE}
	csrv2 := CSrv{Count: 0, Res: HRES_RETURN}
	count := 0
	mux.HandlerM("^/a(\\?.*)?$", &ssrv1, "*", false)
	mux.Handler("^/atestadd(\\?.*)?$", &ssrv1)
	mux.HandleFuncM("^/a(\\?.*)?$", func(w http.ResponseWriter, r *http.Request) {
		count = count + 1
		w.Write([]byte("abc"))
		fmt.Println(mux.RSession(r))
	}, "*", false)
	mux.HandleFunc("^/atestadd(\\?.*)?$", func(w http.ResponseWriter, r *http.Request) {
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

	mux.HFilterFunc("^/redirect(\\?.*)?$", func(hs *HTTPSession) HResult {
		// count = count + 1
		hs.Redirect("http://www.baidu.com")
		// fmt.Println(sb.Session(hs.S.(*SrvSession).Token()))
		return HRES_CONTINUE
	})
	mux.HFilterFunc("^/abc(\\?.*)?$", func(hs *HTTPSession) HResult {
		// count = count + 1
		fmt.Println(hs.CheckVal("ttv"))
		fmt.Println(hs.CheckValA("ttv"))
		fmt.Println(hs.CheckValA("ttvssss"))
		// fmt.Println(sb.Session(hs.S.(*SrvSession).Token()))
		return HRES_CONTINUE
	})
	fmt.Println(mux.regex_m)
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
	fmt.Println(ssrv1.Count, csrv1.Count, count)
	if ssrv1.Count != 1 || csrv1.Count != 2 || count != 3 {
		t.Error("count error")
		return
	}
	c.Get("http://127.0.0.1:2789/t/b")
	c.Post("http://127.0.0.1:2789/t/b", "application/x-www-form-urlencoded", nil)
	mux.ShowLog = true
	c.Get("http://127.0.0.1:2789/t/redirect")
	c.Get("http://127.0.0.1:2789/t2/b")
	time.Sleep(3 * time.Second)
	c.Get("http://127.0.0.1:2789/t/a")
	c.Get("http://127.0.0.1:2789/t/b")
	c.Get("http://127.0.0.1:2789/t2/b")
	c.Get("http://127.0.0.1:2789/t/r1")
	c.Get("http://127.0.0.1:2789/t/r2")
	c.Get("http://127.0.0.1:2789/t/r3")
	c.Get("http://127.0.0.1:2789/t/r4")
	//
	mux.FilterEnable = false
	mux.HandleEnable = false
	c.Get("http://127.0.0.1:2789/t/r1")
	c.Get("http://127.0.0.1:2789/t/r2")
	c.Get("http://127.0.0.1:2789/t/r3")
	c.Get("http://127.0.0.1:2789/t/r4")
	mux.FilterEnable = true
	mux.HandleEnable = true
	//
	c.Get("http://127.0.0.1:2789/abc/a")
	c.Get("http://127.0.0.1:2789/t/abc?ttv=1111")
	c.PostForm("http://127.0.0.1:2789/t/abc", url.Values{"ttv": {"1111"}})
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
	sb.Clear()
	sb.StopLoop()
	//
	NewSessionMux("/", nil)
	//
	fmt.Println("TestSessionMux end")
}
func TestSessionMux2(t *testing.T) {
	sb := NewSrvSessionBuilder("", "/", "rtest2", 2000, 500)
	sb.ShowLog = true
	mux := NewSessionMux("/t", sb)
	mux.ShowLog = true
	mux.HFunc("^/a$", func(hs *HTTPSession) HResult {
		return hs.MsgRes("aaaa")
	})
	mux.HFunc("^/b$", func(hs *HTTPSession) HResult {
		return hs.MsgRes2(200, "bbbb")
	})
	mux.HFunc("^/c$", func(hs *HTTPSession) HResult {
		return hs.MsgResE(1, "cccc")
	})
	mux.HFunc("^/c2$", func(hs *HTTPSession) HResult {
		fmt.Println("------->")
		err := hs.JsonRes(func() {})
		if err == nil {
			t.Error("not error")
		}
		return HRES_RETURN
	})
	mux.HFunc("^/c3$", func(hs *HTTPSession) HResult {
		return hs.MsgResE(1, "cccc")
	})
	mux.HFunc("^/c4$", func(hs *HTTPSession) HResult {
		return hs.MsgResErr2(1, "cccc", errors.New("text"))
	})
	mux.HFunc("^/c5$", func(hs *HTTPSession) HResult {
		return hs.MsgResE2(1, "cccc", "text")
	})
	mux.HFunc("^/c6$", func(hs *HTTPSession) HResult {
		return hs.MsgResErr(1, "cccc", errors.New("text"))
	})
	mux.HFunc("^/c7$", func(hs *HTTPSession) HResult {
		return hs.MsgResE3(1, "cccc", "text")
	})
	mux.HFunc("^/d$", func(hs *HTTPSession) HResult {
		hs.SendT("abc", "text/plain;charset=utf-8")
		return HRES_RETURN
	})
	ts := httptest.NewServer(mux)
	fmt.Println(util.HGet("%s/t/a", ts.URL))
	fmt.Println(util.HGet("%s/t/b", ts.URL))
	fmt.Println(util.HGet("%s/t/c", ts.URL))
	fmt.Println(util.HGet("%s/t/c2", ts.URL))
	fmt.Println(util.HGet("%s/t/c3", ts.URL))
	fmt.Println(util.HGet("%s/t/c4", ts.URL))
	fmt.Println(util.HGet("%s/t/c5", ts.URL))
	fmt.Println(util.HGet("%s/t/c6", ts.URL))
	fmt.Println(util.HGet("%s/t/c7", ts.URL))
	fmt.Println(util.HGet("%s/t/d", ts.URL))
}
func RecF(hs *HTTPSession) HResult {
	hs.FormFInfo("file")
	hs.RecF("file", "/tmp/test2.txt")
	return HRES_RETURN
}
func RecF2(hs *HTTPSession) HResult {
	hs.RecF("file", "/t/mp/test2.txt")
	return HRES_RETURN
}

func TestRecf(t *testing.T) {
	sb := NewSrvSessionBuilder("", "/", "rtest", 2000, 500)
	mux := NewSessionMux("", sb)
	mux.HFunc("^/t1.*$", RecF)
	mux.HFunc("^/t2.*$", RecF2)
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mux.ServeHTTP(w, r)
	}))
	defer ts.Close()
	util.FWrite("/tmp/test.txt", "testing")
	fmt.Println(util.HPostF(fmt.Sprintf("%v/t1", ts.URL), nil, "file", "/tmp/test.txt"))
	fmt.Println(util.HPostF(fmt.Sprintf("%v/t1", ts.URL), nil, "file2", "/tmp/test.txt"))
	fmt.Println(util.HPostF(fmt.Sprintf("%v/t2", ts.URL), nil, "file", "/tmp/test.txt"))
	fmt.Println(util.HPostF(fmt.Sprintf("%v/t1", ts.URL), nil, "file", "/tmp/test.txt2"))
}

func SendF1(hs *HTTPSession) HResult {
	hs.SendF("test.txt", "/tmp/test.txt", "", false)
	return HRES_RETURN
}

func SendF2(hs *HTTPSession) HResult {
	hs.SendF("test.txt", "/tmp/test.txt", "application/text", false)
	return HRES_RETURN
}

func SendF3(hs *HTTPSession) HResult {
	hs.SendF("test.txt", "/tmp/test.txt", "", true)
	return HRES_RETURN
}

func SendF4(hs *HTTPSession) HResult {
	hs.SendF("test.txt", "/tmp/jj/test.txt", "", true)
	return HRES_RETURN
}
func SendF5(hs *HTTPSession) HResult {
	hs.SendF("test.txt", "/tmp", "", true)
	return HRES_RETURN
}

func TestSendF(t *testing.T) {
	sb := NewSrvSessionBuilder("", "/", "rtest", 2000, 500)
	mux := NewSessionMux("", sb)
	mux.HFunc("^/t1.*$", SendF1)
	mux.HFunc("^/t2.*$", SendF2)
	mux.HFunc("^/t3.*$", SendF3)
	mux.HFunc("^/t4.*$", SendF4)
	mux.HFunc("^/t5.*$", SendF5)
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mux.ServeHTTP(w, r)
	}))
	defer ts.Close()
	util.FWrite("/tmp/test.txt", "testing")
	fmt.Println(util.HGet("%s/t1", ts.URL))
	fmt.Println(util.HGet("%s/t2", ts.URL))
	fmt.Println(util.HGet("%s/t3", ts.URL))
	fmt.Println(util.HGet("%s/t4", ts.URL))
	fmt.Println(util.HGet("%s/t5", ts.URL))
}

type rhtp struct {
}

func (rh *rhtp) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK\n"))
}
func TestMatch(t *testing.T) {
	sb := NewSrvSessionBuilder("", "/", "rtest", 2000, 500)
	mux := NewSessionMux("", sb)
	mux.HandleFuncM("^/a1.*$", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK\n"))
	}, "POST", true)
	mux.HandleFuncM("^/a2.*$", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK\n"))
	}, "*", true)
	mux.HandlerM("^/a3.*$", &rhtp{}, "*", true)
	mux.HFilterFuncM("^/a4.*$", func(hs *HTTPSession) HResult {
		return HRES_RETURN
	}, "POST")
	mux.HFilterFuncM("^/a5.*$", func(hs *HTTPSession) HResult { //test check m
		hs.Mux.check_continue(regexp.MustCompile(".*"))
		hs.Mux.check_m(regexp.MustCompile(".*"), "*")
		return HRES_RETURN
	}, "*")
	ts := httptest.NewServer(mux)
	util.HGet("%v/a1", ts.URL)
	util.HGet("%v/a2", ts.URL)
	util.HGet("%v/a3", ts.URL)
	util.HGet("%v/a4", ts.URL)
	util.HGet("%v/a5", ts.URL)
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
