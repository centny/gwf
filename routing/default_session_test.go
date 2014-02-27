package routing

import (
	"github.com/Centny/Cny4go/util"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDefaultSessionMux(t *testing.T) {
	mux := NewSessionMux2("/t")
	mux.HFunc("^.*$", func(hs *HTTPSession) HResult {
		hs.SetVal("abc", "123")
		hs.StrVal("abc")
		hs.SetVal("abc", nil)
		hs.StrVal("abc")
		hs.S.Flush()
		return HRES_RETURN
	})
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mux.ServeHTTP(w, r)
	}))
	util.HTTPGet("%s/t", ts.URL)
}
