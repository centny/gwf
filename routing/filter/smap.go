package filter

import (
	"github.com/Centny/gwf/routing"
	"github.com/Centny/gwf/routing/doc"
)

type SMap struct {
	SET func(hs *routing.HTTPSession, key, val string) error
	GET func(hs *routing.HTTPSession, key string) (string, error)
}

func (s *SMap) SrvHTTP(hs *routing.HTTPSession) routing.HResult {
	var key, val string
	err := hs.ValidCheckVal(`
		key,R|S,L:0;
		val,O|S,L:0;
		`, &key, &val)
	if err != nil {
		return hs.MsgResErr2(1, "arg-err", err)
	}
	if len(val) > 0 {
		err = s.SET(hs, key, val)
	} else {
		val, err = s.GET(hs, key)
	}
	if err == nil {
		return hs.MsgRes(val)
	} else {
		return hs.MsgResErr2(1, "srv-err", err)
	}
}
func (s *SMap) Doc() *doc.Desc {
	return &doc.Desc{
		Title: "Server Map",
		ArgsR: map[string]interface{}{
			"key": "the value key",
		},
		ArgsO: map[string]interface{}{
			"val": "the value",
		},
		ResV: []map[string]interface{}{
			map[string]interface{}{
				"code": "0 is success,or not",
				"data": "the value",
				"msg":  "the error message",
			},
		},
		Detail: "provide map set/get on server,val is not set,will return the target value",
	}
}
func NewSMap() *SMap {
	return &SMap{
		SET: func(hs *routing.HTTPSession, key, val string) error {
			hs.SetVal(key, val)
			return nil
		},
		GET: func(hs *routing.HTTPSession, key string) (string, error) {
			return hs.StrVal(key), nil
		},
	}
}
