package routing

import (
	"fmt"
	"github.com/Centny/gwf/util"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRouting(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r.URL.Path)
		fmt.Println(r.URL.String())
		r.URL.Path = strings.TrimPrefix(r.URL.Path, "/t")
		fmt.Println(r.URL.Path)
		fmt.Println(r.URL.String())
	}))
	util.HTTPGet("%s/t/b/c", ts.URL)

}

// import (
// 	"code.google.com/p/go-uuid/uuid"
// 	"code.google.com/p/go.net/publicsuffix"
// 	"fmt"
// 	"github.com/Centny/gwf/util"
// 	"net/http"
// 	"net/http/cookiejar"
// 	"sync"
// 	"testing"
// 	"time"
// )

// type CSrv struct {
// 	Count int
// 	Res   HResult
// }

// func (s *CSrv) SrvHTTP(hs *HTTPSession) HResult {
// 	s.Count = s.Count + 1
// 	hs.S.Set("abc", "123456789")
// 	fmt.Println(hs.S.Val("abc"))
// 	fmt.Println(hs.S.Val("abc1"))
// 	return s.Res
// }

// func TestSessionMux(t *testing.T) {
// 	sb := NewSrvSessionBuilder("", "/", 2000, 500)
// 	mux := NewSessionMux("/t", sb)
// 	// mux.CDelay = 500
// 	ssrv1 := Ssrv{Count: 0}
// 	csrv1 := CSrv{Count: 0, Res: HRES_CONTINUE}
// 	csrv2 := CSrv{Count: 0, Res: HRES_RETURN}
// 	count := 0
// 	mux.Handler("^/a(\\?.*)?$", &ssrv1)
// 	mux.HandleFunc("^/a(\\?.*)?$", func(w http.ResponseWriter, r *http.Request) {
// 		count = count + 1
// 		w.Write([]byte("abc"))
// 		fmt.Println(mux.RSession(r))
// 	})
// 	mux.HFilter("^/a(\\?.*)?$", &csrv1)
// 	mux.HFilterFunc("^/a(\\?.*)?$", func(hs *HTTPSession) HResult {
// 		count = count + 1
// 		fmt.Println(sb.Session(hs.S.(*SrvSession).Token()))
// 		return HRES_CONTINUE
// 	})
// 	mux.HFilter("^/r1(\\?.*)?$", &csrv2)
// 	mux.HFilterFunc("^/r2(\\?.*)?$", func(hs *HTTPSession) HResult {
// 		count = count + 1
// 		fmt.Println(sb.Session(hs.S.(*SrvSession).Token()))
// 		return HRES_RETURN
// 	})
// 	mux.H("^/a(\\?.*)?$", &csrv1)
// 	mux.HFunc("^/a(\\?.*)?$", func(hs *HTTPSession) HResult {
// 		count = count + 1
// 		return HRES_CONTINUE
// 	})
// 	mux.H("^/r3(\\?.*)?$", &csrv2)
// 	mux.HFunc("^/r4(\\?.*)?$", func(hs *HTTPSession) HResult {
// 		count = count + 1
// 		return HRES_RETURN
// 	})
// 	sb.StartLoop()
// 	//
// 	http.Handle("/t/", mux)
// 	http.Handle("/t2/", mux)
// 	go http.ListenAndServe(":2789", nil)
// 	options := cookiejar.Options{
// 		PublicSuffixList: publicsuffix.List,
// 	}
// 	jar, err := cookiejar.New(&options)
// 	if err != nil {
// 		t.Error(err.Error())
// 		return
// 	}
// 	c := http.Client{Jar: jar}
// 	c.Get("http://127.0.0.1:2789/t/a")
// 	c.Get("http://127.0.0.1:2789/t/b")
// 	c.Get("http://127.0.0.1:2789/t2/b")
// 	fmt.Println(ssrv1.Count, csrv1.Count, count)
// 	if ssrv1.Count != 1 || csrv1.Count != 2 || count != 3 {
// 		t.Error("count error")
// 	}
// 	time.Sleep(3 * time.Second)
// 	c.Get("http://127.0.0.1:2789/t/a")
// 	c.Get("http://127.0.0.1:2789/t/b")
// 	c.Get("http://127.0.0.1:2789/t2/b")
// 	c.Get("http://127.0.0.1:2789/t/r1")
// 	c.Get("http://127.0.0.1:2789/t/r2")
// 	c.Get("http://127.0.0.1:2789/t/r3")
// 	c.Get("http://127.0.0.1:2789/t/r4")
// 	//
// 	if mux.RSession(nil) != nil {
// 		t.Error("not nil")
// 		return
// 	}
// 	if sb.Session("ss") != nil {
// 		t.Error("not nil")
// 		return
// 	}
// 	//
// 	sb.StopLoop()
// 	//
// 	//
// 	fmt.Println("TestSessionMux end")
// }

// ////////////////////////////////////

// type SrvSession struct {
// 	token string
// 	begin int64
// 	kvs   map[string]interface{}
// }

// func (s *SrvSession) Val(key string) interface{} {
// 	if v, ok := s.kvs[key]; ok {
// 		return v
// 	} else {
// 		return nil
// 	}
// }
// func (s *SrvSession) Set(key string, val interface{}) {
// 	s.kvs[key] = val
// }
// func (s *SrvSession) Token() string {
// 	return s.token
// }
// func (s *SrvSession) Flush() {
// 	s.begin = util.Timestamp(time.Now())
// }

// //
// type SrvSessionBuilder struct {
// 	//
// 	Domain  string
// 	Path    string
// 	Timeout int64
// 	CDelay  time.Duration
// 	//
// 	looping bool
// 	ks      map[string]*SrvSession //key session
// 	ks_lck  sync.RWMutex
// }

// func NewSrvSessionBuilder(domain string, path string, timeout int64, cdelay time.Duration) *SrvSessionBuilder {
// 	sb := SrvSessionBuilder{}
// 	sb.Domain = domain
// 	sb.Path = path
// 	sb.Timeout = timeout
// 	sb.CDelay = cdelay
// 	sb.ks = map[string]*SrvSession{}
// 	return &sb
// }
// func (s *SrvSessionBuilder) FindSession(w http.ResponseWriter, r *http.Request) Session {
// 	c, err := r.Cookie("token")
// 	ncookie := func() {
// 		c = &http.Cookie{}
// 		c.Name = "token"
// 		c.Value = uuid.New()
// 		c.Path = s.Path
// 		c.Domain = s.Domain
// 		c.MaxAge = 0
// 		//
// 		session := &SrvSession{}
// 		session.token = c.Value
// 		session.kvs = map[string]interface{}{}
// 		session.Flush()
// 		//
// 		s.ks_lck.RLock()
// 		s.ks[c.Value] = session
// 		s.ks_lck.RUnlock()
// 		http.SetCookie(w, c)
// 	}
// 	if err != nil {
// 		ncookie()
// 	}
// 	if _, ok := s.ks[c.Value]; !ok { //if not found,reset cookie
// 		ncookie()
// 	}
// 	session := s.ks[c.Value]
// 	session.Flush()
// 	return session
// }

// func (s *SrvSessionBuilder) Session(token string) Session {
// 	if v, ok := s.ks[token]; ok {
// 		return v
// 	} else {
// 		return nil
// 	}
// }

// //
// func (s *SrvSessionBuilder) StartLoop() {
// 	s.looping = true
// 	go s.Loop()
// }
// func (s *SrvSessionBuilder) StopLoop() {
// 	s.looping = false
// }

// //
// func (s *SrvSessionBuilder) Loop() {
// 	for s.looping {
// 		ary := []string{}
// 		now := util.Timestamp(time.Now())
// 		for k, v := range s.ks {
// 			delay := now - v.begin
// 			if delay > s.Timeout {
// 				ary = append(ary, k)
// 			}
// 		}
// 		s.ks_lck.RLock()
// 		for _, v := range ary {
// 			delete(s.ks, v)
// 		}
// 		s.ks_lck.RUnlock()
// 		time.Sleep(s.CDelay * time.Millisecond)
// 	}
// }

// func TestTesting(t *testing.T) {
// 	sb := NewSrvSessionBuilder("", "/", 2000, 500)
// 	mux := NewSessionMux("/t", sb)
// 	ts, _ := httptest.NewServer(func(w http.ResponseWriter, r *http.Request) {
// 		mux.ServeHTTP(w, r)
// 	})
// 	http.Get(fmt.Sprintf("%s", ...))
// }
type S3 struct {
	A string
}

func TestOrder(t *testing.T) {
	mv := map[*S3]string{}
	mv[&S3{A: "3"}] = "abc"
	mv[&S3{A: "1"}] = "abc"
	mv[&S3{A: "2"}] = "abc"
	for k, v := range mv {
		fmt.Println(k.A, v)
	}
}

// func TestPrintRef(t *testing.T) {
// 	// var val = reflect.Indirect()
// 	var val = reflect.ValueOf(HFunc)
// 	fmt.Println(reflect.Method(val))
// }
