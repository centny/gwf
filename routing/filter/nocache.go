package filter

import (
	"github.com/Centny/gwf/routing"
	"time"
)

func NoCacheFilter(hs *routing.HTTPSession) routing.HResult {
	hs.R.Header.Set("Expires", "Tue, 01 Jan 1980 1:00:00 GMT")
	hs.R.Header.Set("Last-Modified", time.Now().Format(time.RFC1123))
	hs.R.Header.Set("Cache-Control", "no-stroe,no-cache,must-revalidate,post-check=0,pre-check=0")
	hs.R.Header.Set("Pragma", "no-cache")
	return routing.HRES_CONTINUE
}
