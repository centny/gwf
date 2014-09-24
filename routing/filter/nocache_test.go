package filter

import (
	"github.com/Centny/gwf/routing"
	"github.com/Centny/gwf/util"
	"net/http/httptest"
	"testing"
)

func TestNoCache(t *testing.T) {
	mux := routing.NewSessionMux2("")
	mux.HFilterFunc("^.*$", NoCacheFilter)
	mux.HFunc("^.*$", func(hs *routing.HTTPSession) routing.HResult {
		return routing.HRES_RETURN
	})
	ts := httptest.NewServer(mux)
	util.HGet("%v", ts.URL)
}
