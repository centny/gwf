package wdoc

import (
	"fmt"
	"path"

	"sort"

	"github.com/Centny/gwf/routing"
	"github.com/Centny/gwf/util"
)

type CaseL struct {
	TS []*Text `json:"text"`
	WS []*Web  `json:"webs"`
	FS []*Func `json:"func"`

	ws_m map[string]bool
	fs_m map[string]bool
}

type Cases struct {
	Data map[string]*CaseL
}

func NewCases() *Cases {
	return &Cases{
		Data: map[string]*CaseL{},
	}
}

func (c *Cases) AddCase(cs *Case, info *Func) {
	for _, key := range cs.Keys {
		var cl = c.Data[key]
		if cl == nil {
			cl = &CaseL{
				ws_m: map[string]bool{},
				fs_m: map[string]bool{},
			}
			c.Data[key] = cl
		}
		for _, web := range cs.WS {
			if cl.ws_m[web.Key] {
				continue
			}
			cl.WS = append(cl.WS, web)
			cl.ws_m[web.Key] = true
		}
		if cs.Text != nil {
			cl.TS = append(cl.TS, cs.Text)
		}
		var fk = info.Pkg + "/" + info.Name
		if cl.fs_m[fk] {
			continue
		}
		cl.FS = append(cl.FS, info)
		cl.fs_m[fk] = true
	}
}
func (c *Cases) ListCases() []util.Map {
	var res = []util.Map{}
	for key, cl := range c.Data {
		var desc = ""
		for _, text := range cl.TS {
			desc += text.Title + ";"
		}
		res = append(res, util.Map{
			"title":     key,
			"api_count": len(cl.FS),
			"desc":      desc,
		})
	}
	sort.Sort(util.NewMapSorter(res, "title", 2))
	return res
}
func (c *Cases) FindCase(key string) *CaseL {
	return c.Data[key]
}
func (c *Cases) SrvHTTP(hs *routing.HTTPSession) routing.HResult {
	var _, fn = path.Split(hs.R.URL.Path)
	switch fn {
	case "list":
		return hs.JRes(c.ListCases())
	case "data":
		var data = c.FindCase(hs.RVal("key"))
		if data == nil {
			data = &CaseL{}
		}
		util.NewFieldStringSorter("Name", data.FS).Sort(false)
		util.NewFieldStringSorter("Title", data.TS).Sort(false)
		util.NewFieldStringSorter("Index", data.TS).Sort(false)
		return hs.JRes(data)
	default:
		return hs.JRes(util.Map{
			"err": fmt.Sprintf("%v api not found", fn),
		})
	}
}
func (c *Cases) Clear() {
	c.Data = map[string]*CaseL{}
}
