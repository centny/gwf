package srvmux

import (
	"fmt"
	"net/http"
	"testing"
)

type Ssrv struct {
	Count int
}

func (s *Ssrv) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.Count = s.Count + 1
}
func TestSrvMux(t *testing.T) {
	ssrv1 := Ssrv{Count: 0}
	smx1 := NewSrvMux("/srv")
	smx1.Handler("/a", &ssrv1)
	smx1.HandleFunc("/b", func(w http.ResponseWriter, r *http.Request) {
		ssrv1.Count = ssrv1.Count + 1
	})
	http.Handle("/srv/", smx1)
	//
	ssrv2 := Ssrv{}
	smx2 := NewSrvMux("/srv")
	smx2.Handler("/a", &ssrv2)
	smx2.HandleFunc("/b", func(w http.ResponseWriter, r *http.Request) {
		ssrv2.Count = ssrv2.Count + 1
	})
	http.Handle("/srv2/", smx2)
	http.Handle("/sr2/", smx2)
	//
	ssrv3 := Ssrv{Count: 0}
	rgx1 := NewRegMux("/reg")
	rgx1.Handler("^\\/a.*$", &ssrv3)
	rgx1.HandleFunc("^\\/b.*$", func(w http.ResponseWriter, r *http.Request) {
		ssrv3.Count = ssrv3.Count + 1
	})
	http.Handle("/reg/", rgx1)
	//
	ssrv4 := Ssrv{Count: 0}
	rgx2 := NewRegMux("/reg")
	rgx2.Handler("^\\/a.*$", &ssrv4)
	rgx2.HandleFunc("^\\/b.*$", func(w http.ResponseWriter, r *http.Request) {
		ssrv4.Count = ssrv4.Count + 1
	})
	http.Handle("/reg2/", rgx2)
	http.Handle("/re2/", rgx2)
	//
	go http.ListenAndServe(":6789", nil)
	http.Get("http://127.0.0.1:6789/srv/a")
	http.Get("http://127.0.0.1:6789/srv/b")
	http.Get("http://127.0.0.1:6789/srv2/a")
	http.Get("http://127.0.0.1:6789/srv2/b")
	http.Get("http://127.0.0.1:6789/sr2/a")
	http.Get("http://127.0.0.1:6789/sr2/b")
	http.Get("http://127.0.0.1:6789/reg/a1")
	http.Get("http://127.0.0.1:6789/reg/b2")
	http.Get("http://127.0.0.1:6789/reg2/a3")
	http.Get("http://127.0.0.1:6789/reg2/b4")
	http.Get("http://127.0.0.1:6789/re2/a3")
	http.Get("http://127.0.0.1:6789/re2/b4")
	// time.Sleep(time.Second)
	if ssrv1.Count != 2 || ssrv2.Count > 0 || ssrv3.Count != 2 || ssrv4.Count > 0 {
		t.Error("testing error")
	}
	fmt.Println("TestSrvMux end")
}

func TestRegex(t *testing.T) {
	// fmt.Println(regexp.MustCompile("\\idi"))
}
