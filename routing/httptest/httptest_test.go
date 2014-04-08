package httptest

import (
	"fmt"
	"github.com/Centny/Cny4go/routing"
	"net/http"
	"testing"
)

var c int64 = 0

func T(hs *routing.HTTPSession) routing.HResult {
	fmt.Println(hs.CheckVal("a"))
	fmt.Println(hs.CheckVal("b"))
	fmt.Println(c)
	c = c + 1
	hs.W.Write([]byte("{\"OK\":1}"))
	return routing.HRES_RETURN
}
func NT(w http.ResponseWriter, r *http.Request) {
	fmt.Println(c)
	c = c + 1
}

type T2 struct {
}

func (t *T2) SrvHTTP(hs *routing.HTTPSession) routing.HResult {
	fmt.Println(hs.CheckVal("a"))
	fmt.Println(hs.CheckVal("b"))
	fmt.Println(c)
	c = c + 1
	hs.W.Write([]byte("{\"OK\":1}"))
	return routing.HRES_RETURN
}

func TestServer(t *testing.T) {
	ts := NewServer(T)
	defer ts.Close()
	_, err := ts.G("?a=%v", "testing")
	if err != nil {
		t.Error(err.Error())
		return
	}
	mv, err := ts.G2("?b=%v", "kkkk")
	if err != nil {
		t.Error(err.Error())
		return
	}
	if mv.IntVal("OK") != 1 {
		t.Error("not error")
	}
	_, err = ts.P("/", map[string]string{
		"a": "testing",
	})
	if err != nil {
		t.Error(err.Error())
		return
	}
	_, err = ts.P2("/", map[string]string{
		"a": "testing",
	})
	if err != nil {
		t.Error(err.Error())
		return
	}
	ts.PostF("/", "filt", "/tmp/test.txt", nil)
	ts.PostF2("/", "filt", "/tmp/test.txt", nil)
}
func TestServer2(t *testing.T) {
	ts := NewServer2(&T2{})
	ts.G("")
	Tnf(NT, "?a=%v", "testing")
	Tf(T, "?a=%v", "testing")
	Th(&T2{}, "?a=%v", "testing")
	NewMuxServer()
}
