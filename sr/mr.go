package sr

// import (
// 	"fmt"
// 	"github.com/Centny/gwf/routing"
// 	"github.com/Centny/gwf/util"
// 	"strconv"
// 	"strings"
// 	"sync"
// )

// type MR_S struct {
// 	Path string
// 	Data *util.Map
// }

// type MR struct {
// 	Pre string
// }

// func NewMR(pre string) *MR {
// 	if len(strings.Trim(pre, "\t ")) < 1 {
// 		pre = "/"
// 	}
// 	return &MR{
// 		Pre:  pre,
// 		Data: map[string]*util.Map{},
// 	}
// }
// func (m *MR) SrvHTTP(hs *routing.HTTPSession) routing.HResult {
// 	path := strings.TrimPrefix(hs.R.URL.Path, m.Pre)
// 	path = strings.TrimRight(path, "/ \t")
// 	var data, typ, exec string = "", "J", "get"
// 	err := hs.ValidF(`
// 		data,O|S,L:0;
// 		type,O|S,O:I~F~S~J;
// 		exec,O|S,O:plus~del~set~get;
// 		`, &data, &typ, &exec)
// 	if err != nil {
// 		return hs.MsgResErr2(1, "arg-err", err)
// 	}
// 	ps := strings.Split(path, "/")
// 	if len(ps) < 2 {
// 		path = ""
// 	} else {
// 		path = strings.Join(ps[:len(ps)-1], "/")
// 	}
// 	key := ps[len(ps)-1]
// 	var tmv *util.Map
// 	if mv, ok := m.Data[path]; ok {
// 		tmv = mv
// 	} else {
// 		tmv = &util.Map{}
// 		m.Data[path] = tmv
// 	}
// 	switch exec {
// 	case "set":
// 		return m.set(hs, tmv, path, key, data, typ, exec)
// 	case "del":
// 		return m.del(hs, tmv, path, key, data, typ, exec)
// 	case "plus":
// 		return m.plus(hs, tmv, path, key, data, typ, exec)
// 	default:
// 		return m.get(hs, tmv, path, key, data, typ, exec)
// 	}
// }
// func (m *MR) set(hs *routing.HTTPSession, tmv *util.Map, path, key, data, typ, exec string) routing.HResult {
// 	var val interface{}
// 	var err error
// 	switch typ {
// 	case "I":
// 		val, err = strconv.ParseInt(data, 10, 64)
// 	case "F":
// 		val, err = strconv.ParseFloat(data, 64)
// 	case "S":
// 		val = data
// 	default:
// 		val, err = util.Json2Map(data)
// 	}
// 	if err != nil {
// 		return hs.MsgResErr2(1, "arg-err", err)
// 	}
// 	tmv.SetVal(key, val)
// 	return hs.MsgRes(val)
// }
// func (m *MR) plus(hs *routing.HTTPSession, tmv *util.Map, path, key, data, typ, exec string) routing.HResult {
// 	switch typ {
// 	case "I":
// 		val, err := strconv.ParseInt(data, 10, 64)
// 		if err != nil {
// 			return hs.MsgResErr2(1, "arg-err", err)
// 		}
// 		tval := tmv.Val(key)
// 		if tval == nil {
// 			tval = int64(0)
// 		}
// 		if ti, ok := tval.(int64); ok {
// 			ti += val
// 			tmv.SetVal(key, ti)
// 			return hs.MsgRes("OK")
// 		} else {
// 			return hs.MsgResE3(1, "arg-err", "target value is not int")
// 		}
// 	case "F":
// 		val, err := strconv.ParseFloat(data, 64)
// 		if err != nil {
// 			return hs.MsgResErr2(1, "arg-err", err)
// 		}
// 		tval := tmv.Val(key)
// 		if tval == nil {
// 			tval = float64(0)
// 		}
// 		if ti, ok := tval.(float64); ok {
// 			ti += val
// 			tmv.SetVal(key, ti)
// 			return hs.MsgRes("OK")
// 		} else {
// 			return hs.MsgResE3(1, "arg-err", "target value is not float")
// 		}
// 	default:
// 		return hs.MsgResE3(1, "arg-err", "plus not support for type:"+typ)
// 	}
// }
// func (m *MR) del(hs *routing.HTTPSession, tmv *util.Map, path, key, data, typ, exec string) routing.HResult {
// 	if tmv.Exist(key) {
// 		tmv.SetVal(key, nil)
// 		return hs.MsgRes("OK")
// 	} else {
// 		return hs.MsgResE(1, fmt.Sprintf("not exist on path(%v) by key(%v)", path, key))
// 	}
// }
// func (m *MR) get(hs *routing.HTTPSession, tmv *util.Map, path, key, data, typ, exec string) routing.HResult {
// 	if key == "*" {
// 		return hs.MsgRes(tmv)
// 	}
// 	if tmv.Exist(key) {
// 		return hs.MsgRes(tmv.Val(key))
// 	} else {
// 		return hs.MsgResE(1, "not found")
// 	}
// }
