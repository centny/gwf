package filter

import (
	"net/url"
	"strings"

	"fmt"

	"bytes"

	"github.com/Centny/gwf/log"
	"github.com/Centny/gwf/routing"
	"github.com/Centny/gwf/util"
)

type SimpleMerger struct {
	Keys  []string
	Route map[string]string
}

func NewSimpleMerger() *SimpleMerger {
	return &SimpleMerger{
		Route: map[string]string{},
	}
}

func (s *SimpleMerger) SrvHTTP(hs *routing.HTTPSession) routing.HResult {
	var merger = hs.CheckValA("merger")
	if len(merger) < 1 {
		err := fmt.Errorf("the merger key is not found")
		return hs.MsgResErr2(1, "arg-err", err)
	}
	var body = util.Map{}
	var ctype = hs.R.Header.Get("Content-Type")
	if strings.HasPrefix(ctype, "application/json") {
		err := hs.UnmarshalJ(&body)
		if err != nil {
			err = fmt.Errorf("parse json body fail with error(%v)", err)
			return hs.MsgResErr2(1, "arg-err", err)
		}
	}
	var res = util.Map{}
	var mkeys = strings.Split(merger, ",")
	for _, mkey := range mkeys {
		rurl, ok := s.Route[mkey]
		if !ok {
			err := fmt.Errorf("found invalid key(%v) is not in list(%v)", mkey, s.Keys)
			return hs.MsgResErr2(2, "arg-err", err)
		}
		data, err := s.ReverseRoute(mkey, rurl, body.MapVal(mkey), hs)
		if err != nil {
			log.W("SimpleMerger reverse url(%v) fail with error(%v)", rurl, err)
			return hs.MsgResErr2(1, "arg-err", err)
		}
		if data.IntVal("code") != 0 {
			err = fmt.Errorf("SimpleMerger reverse url(%v) fail with %v", rurl, util.S2Json(data))
			return hs.MsgResErr2(2, "arg-err", err)
		}
		res[mkey] = data.Val("data")
	}
	return hs.MsgRes(res)
}

func (s *SimpleMerger) ReverseRoute(key, rurl string, body util.Map, hs *routing.HTTPSession) (res util.Map, err error) {
	var pre = key + "."
	var kvs = url.Values{}
	for k, v := range hs.R.Form {
		if !strings.Contains(k, ".") {
			kvs[k] = v
			continue
		}
		if strings.HasPrefix(k, pre) {
			kvs[strings.TrimPrefix(k, pre)] = v
		}
	}
	turl := rurl + "?" + kvs.Encode()
	if body == nil {
		log.D("SimpleMerger send get with %v", turl)
		res, err = util.HGet2("%v", turl)
	} else {
		log.D("SimpleMerger send post with %v", turl)
		_, res, err = util.HPostN2(turl, "application/json", bytes.NewBufferString(util.S2Json(body)))
	}
	return
}

func HandMerger(mux *routing.SessionMux, fcfg *util.Fcfg) {
	var merger = fcfg.Val2("loc/merger", "")
	if len(merger) < 1 {
		log.I("HandMerger parse merger done wiht loc/merger configure not found")
		return
	}
	var mkeys = strings.Split(merger, ",")
	for _, mkey := range mkeys {
		mtype := fcfg.Val2(mkey+"/type", "")
		if len(mtype) < 1 {
			log.E("HandMerger parse merger(%v) fail with (%v/type) not found", mkey, mkey)
			continue
		}
		switch mtype {
		case "simple":
			parseSimpeMerger(mux, fcfg, mkey)
		}
	}
}

func parseSimpeMerger(mux *routing.SessionMux, fcfg *util.Fcfg, mkey string) {
	route := fcfg.Val2(mkey+"/route", "")
	if len(route) < 1 {
		log.E("HandMerger parse simple merger(%v) fail with (%v/route) not found", mkey, mkey)
		return
	}
	allkeys := fcfg.Val2(mkey+"/keys", "")
	if len(allkeys) < 1 {
		log.E("HandMerger parse simple merger(%v) fail with (%v/keys) not found", mkey, mkey)
		return
	}
	simple := NewSimpleMerger()
	simple.Keys = strings.Split(allkeys, ",")
	for _, ckey := range simple.Keys {
		rurl := fcfg.Val2(mkey+"/"+ckey, "")
		if len(rurl) < 1 {
			log.E("HandMerger parse simple merger(%v) fail with (%v/%v) not found", mkey, mkey, ckey)
			return
		}
		simple.Route[ckey] = rurl
	}
	mux.H(route, simple)
	log.I("HandMerger parse simple merger(%v) on route(%v) success", mkey, route)
}
