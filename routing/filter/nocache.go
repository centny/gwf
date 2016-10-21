package filter

import (
	"time"

	"github.com/Centny/gwf/routing"
)

func NoCacheFilter(hs *routing.HTTPSession) routing.HResult {
	hs.W.Header().Set("Expires", "Tue, 01 Jan 1980 1:00:00 GMT")
	hs.W.Header().Set("Last-Modified", time.Now().UTC().Format("Mon, 02 Jan 2006 15:04:05 GMT"))
	hs.W.Header().Set("Cache-Control", "no-stroe,no-cache,must-revalidate,post-check=0,pre-check=0")
	hs.W.Header().Set("Pragma", "no-cache")
	return routing.HRES_CONTINUE
}
