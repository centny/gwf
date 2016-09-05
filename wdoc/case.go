package wdoc

import (
	"fmt"
	"github.com/Centny/gwf/routing"
	"github.com/Centny/gwf/util"
	"path"
)

type CaseL struct {
	TS []*Text          `json:"text"`
	WS map[string]*Web  `json:"webs"`
	FS map[string]*Func `json:"func"`
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
				WS: map[string]*Web{},
				FS: map[string]*Func{},
			}
			c.Data[key] = cl
		}
		for _, web := range cs.WS {
			cl.WS[web.Index] = web
		}
		if cs.Text != nil {
			cl.TS = append(cl.TS, cs.Text)
		}
		cl.FS[info.Name] = info
	}
}
func (c *Cases) ListCases() map[string]int {
	var res = map[string]int{}
	for key, cl := range c.Data {
		res[key] = len(cl.FS)
	}
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
