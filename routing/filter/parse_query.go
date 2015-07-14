package filter

import (
	"github.com/Centny/gwf/routing"
)

func ParseQuery(hs *routing.HTTPSession) routing.HResult {
	err := hs.ParseQuery()
	if err == nil {
		return routing.HRES_CONTINUE
	} else {
		return hs.MsgResErr2(1, "arg-err", err)
	}
}
