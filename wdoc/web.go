package wdoc

import (
	"fmt"
	"github.com/Centny/gwf/log"
	"github.com/Centny/gwf/routing"
	"github.com/Centny/gwf/util"
	"net/http"
	"path/filepath"
	"regexp"
	"strings"
)

type WebH struct {
	Base  string
	CmdF  string
	Index string
	Reg   *regexp.Regexp
	H     http.Handler
}

func NewWebH(base, cmdf string) *WebH {
	return &WebH{
		Base: base,
		CmdF: cmdf,
		Reg:  regexp.MustCompile(".*\\.[(md)(MD)]+$"),
		H:    http.FileServer(http.Dir(base)),
	}
}

func (m *WebH) SrvHTTP(hs *routing.HTTPSession) routing.HResult {
	var path = hs.R.URL.Path
	if path == "" || path == "/" {
		path += m.Index
		hs.R.URL.Path = path
	}
	path = filepath.Join(m.Base, path)
	log.D("MD_H doing path %v", path)
	if !m.Reg.MatchString(path) {
		m.H.ServeHTTP(hs.W, hs.R)
		return routing.HRES_RETURN
	}
	bys, err := util.Exec(fmt.Sprintf(m.CmdF, path))
	if err != nil {
		log.W("parsing md file(%v) error->%v", path, err)
		hs.W.WriteHeader(404)
		fmt.Fprintf(hs.W, "parsing md file(%v) error->%v", path, err)
		return routing.HRES_RETURN
	}
	hs.SendT(string(bys), "text/html;charset=utf8")
	return routing.HRES_RETURN
}

type Webs struct {
	Pre  string
	CmdF string
	HS   map[string]*WebH
	Exc  []*regexp.Regexp
}

func NewWebs(pre, cmdf string) *Webs {
	return &Webs{
		Pre:  pre,
		CmdF: cmdf,
		HS:   map[string]*WebH{},
		Exc: []*regexp.Regexp{
			regexp.MustCompile(".*\\.go(\\?.*)?"),
		},
	}
}
func (w *Webs) AddWeb(name string, h *WebH) {
	w.HS[name] = h
}
func (w *Webs) AddWeb2(name, base, idx string) {
	var md = NewWebH(base, w.CmdF)
	md.Index = idx
	w.HS[name] = md
}
func (w *Webs) IfExc(path string) bool {
	for _, exc := range w.Exc {
		if exc.MatchString(path) {
			return true
		}
	}
	return false
}
func (w *Webs) SrvHTTP(hs *routing.HTTPSession) routing.HResult {
	var path = hs.R.URL.Path
	if w.IfExc(path) {
		hs.W.WriteHeader(404)
		return routing.HRES_RETURN
	}
	path = strings.TrimPrefix(path, w.Pre)
	path = strings.TrimPrefix(path, "/")
	var paths = strings.SplitN(path, "/", 2)
	var name = paths[0]
	var spath = "/"
	if len(paths) > 1 {
		spath = paths[1]
	}
	hs.R.URL.Path = spath
	if h, ok := w.HS[name]; ok {
		return h.SrvHTTP(hs)
	} else {
		hs.W.WriteHeader(404)
		fmt.Fprintf(hs.W, "the web by name(%v) error->not found", name)
		return routing.HRES_RETURN
	}
}
func (w *Webs) Clear() {
}
