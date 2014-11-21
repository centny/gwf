//Package for building handler api document in routing.SessionMux.
//@Author:Centny.
//
package doc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/Centny/gwf/routing"
	"github.com/Centny/gwf/util"
	"html/template"
	"reflect"
	"regexp"
	"runtime"
	"strings"
)

const (
	FMT_HTML = "html"
	FMT_JSON = "json"
)

//marked the func document.
var Marked map[string]*Desc = map[string]*Desc{}

//mark the func document.
func ApiV(f interface{}, doc *Desc) bool {
	fnc := runtime.FuncForPC(reflect.ValueOf(f).Pointer())
	name := fnc.Name()
	if _, ok := Marked[name]; ok {
		panic(name + " already registered")
	}
	Marked[name] = doc
	return true
}
func pkgv(pkg string) string {
	pkg = strings.Replace(pkg, "/", "_", -1)
	pkg = strings.Replace(pkg, ".", "_", -1)
	return pkg
}

//api describe
type Desc struct {
	Title   string
	Url     []string                          //example URL.
	ArgsR   map[string]interface{}            //required arguments.
	ArgsO   map[string]interface{}            //option arguments.
	Option  map[string]map[string]interface{} //the argument
	ResV    interface{}                       //result
	ResJSON string                            `json:"-"` //json result
	ResHTML template.HTML                     `json:"-"` //html result
	SeeV    []interface{}                     `json:"-"` //see link
	Detail  string                            //detail
	See     []map[string]interface{}          `json:"SeeV"` //see link
}

//register api.
func (d Desc) Api(f interface{}) int {
	if d.ArgsR == nil {
		d.ArgsR = map[string]interface{}{}
	}
	if d.ArgsO == nil {
		d.ArgsO = map[string]interface{}{}
	}
	if d.Option == nil {
		d.Option = map[string]map[string]interface{}{}
	}
	if d.ResV == nil {
		d.ResV = []map[string]string{}
	}
	if d.SeeV == nil {
		d.SeeV = []interface{}{}
	}
	ApiV(f, &d)
	return 0
}

func func_pkgn(f interface{}) (string, string, string) {
	fnc := runtime.FuncForPC(reflect.ValueOf(f).Pointer())
	fnc_n := fnc.Name()
	return fnc_n, regexp.MustCompile("^.*\\.").ReplaceAllString(fnc.Name(), ""), regexp.MustCompile("\\.[^\\.]*$").ReplaceAllString(fnc.Name(), "")
}
func handle_pkgn(f interface{}) (string, string, string) {
	typ := reflect.TypeOf(reflect.Indirect(reflect.ValueOf(f)).Interface())
	return typ.PkgPath() + "/" + typ.Name(), typ.Name(), typ.PkgPath()
}

//doc visiable interface.
type Docable interface {
	Doc() *Desc
}

//the routing.SessionMux info.
type Mux struct {
	Pre          string
	Items        []string //comment items.
	Filters      []DocV
	FilterFunc   []DocV
	Handlers     []DocV
	HandlerFunc  []DocV
	NHandlers    []DocV
	NHandlerFunc []DocV
}

//the doc
type DocV struct {
	Name    string
	Pattern string
	Abs     string
	Doc     *Desc
	Marked  bool
	Pkg     string
}

//the doc viewer for building the SessionMux info.
type DocViewer struct {
	Incs         []*regexp.Regexp
	Excs         []*regexp.Regexp
	Items        []string //comment items.
	HTML         string
	Filters      bool
	FilterFunc   bool
	Handlers     bool
	HandlerFunc  bool
	NHandlers    bool
	NHandlerFunc bool
}

//new default viewer.
func NewDocViewer() *DocViewer {
	return &DocViewer{
		Items:        []string{},
		HTML:         HTML,
		Filters:      true,
		FilterFunc:   true,
		Handlers:     true,
		HandlerFunc:  true,
		NHandlers:    true,
		NHandlerFunc: true,
	}
}

//new default viewer by include.
func NewDocViewerInc(inc string) *DocViewer {
	dv := NewDocViewer()
	dv.Incs = []*regexp.Regexp{
		regexp.MustCompile(inc),
	}
	return dv
}

//new default viewer by exclude.
func NewDocViewerExc(exc string) *DocViewer {
	dv := NewDocViewer()
	dv.Excs = []*regexp.Regexp{
		regexp.MustCompile(exc),
	}
	return dv
}
func (d *DocViewer) handler_doc(reg *regexp.Regexp, f interface{}) DocV {
	_, name, pkgp := handle_pkgn(f)
	doc := DocV{}
	doc.Name = name
	doc.Pattern = reg.String()
	doc.Pkg = pkgp
	doc.Abs = fmt.Sprintf("%s_%s", doc.Name, pkgv(doc.Pkg))
	if dd, ok := f.(Docable); ok {
		doc.Doc = d.build_see(dd.Doc())
		doc.Marked = true
	} else {
		doc.Doc = nil
		doc.Marked = false
	}
	return doc
}
func (d *DocViewer) func_doc(reg *regexp.Regexp, f interface{}) DocV {
	uname, name, pkgp := func_pkgn(f)
	doc := DocV{}
	doc.Name = name
	doc.Pattern = reg.String()
	doc.Pkg = pkgp
	doc.Abs = fmt.Sprintf("%s_%s", doc.Name, pkgv(doc.Pkg))
	if dd, ok := Marked[uname]; ok {
		doc.Doc = d.build_see(dd)
		doc.Marked = true
	} else {
		doc.Doc = nil
		doc.Marked = false
	}
	return doc
}

func (d *DocViewer) match(t *regexp.Regexp) bool {
	if len(d.Incs) > 0 {
		for _, inc := range d.Incs {
			if inc.MatchString(t.String()) {
				return true
			}
		}
		return false
	}
	if len(d.Excs) > 0 {
		for _, exc := range d.Excs {
			if exc.MatchString(t.String()) {
				return false
			}
		}
		return true
	}
	return true
}
func (d *DocViewer) build_see(desc *Desc) *Desc {
	if desc.SeeV == nil {
		return desc
	}
	desc.See = []map[string]interface{}{}
	if desc.ResV != nil {
		dst := &bytes.Buffer{}
		json.Indent(dst, []byte(util.S2Json(desc.ResV)), "", "  ")
		desc.ResJSON = dst.String()
		htmlv := desc.ResJSON
		htmlv = strings.Replace(htmlv, "\"", "&quot;", -1)
		htmlv = strings.Replace(htmlv, " ", "&nbsp;", -1)
		htmlv = strings.Replace(htmlv, "\n", "<br/>", -1)
		desc.ResHTML = template.HTML(htmlv)
	}
	var name, pkgn string
	for _, v := range desc.SeeV {
		typ := reflect.TypeOf(v)
		if typ.Kind() == reflect.Func {
			_, name, pkgn = func_pkgn(v)
		} else {
			_, name, pkgn = handle_pkgn(v)
		}
		mv := map[string]interface{}{}
		mv["Name"] = name
		mv["Pkg"] = pkgn
		mv["Abs"] = name + "_" + pkgv(pkgn)
		desc.See = append(desc.See, mv)
	}
	return desc
}

//build Mux by SessionMux
func (d *DocViewer) BuildMux(smux *routing.SessionMux) *Mux {
	tmux := &Mux{
		Items:        d.Items,
		Pre:          smux.Pre,
		Filters:      []DocV{},
		FilterFunc:   []DocV{},
		Handlers:     []DocV{},
		HandlerFunc:  []DocV{},
		NHandlers:    []DocV{},
		NHandlerFunc: []DocV{},
	}
	if d.Filters {
		for reg, f := range smux.Filters {
			if d.match(reg) {
				tmux.Filters = append(tmux.Filters, d.handler_doc(reg, f))
			}
		}
	}
	if d.FilterFunc {
		for reg, f := range smux.FilterFunc {
			if d.match(reg) {
				tmux.FilterFunc = append(tmux.FilterFunc, d.func_doc(reg, f))
			}
		}
	}
	if d.Handlers {
		for reg, f := range smux.Handlers {
			if d.match(reg) {
				tmux.Handlers = append(tmux.Handlers, d.handler_doc(reg, f))
			}
		}
	}
	if d.HandlerFunc {
		for reg, f := range smux.HandlerFunc {
			if d.match(reg) {
				tmux.HandlerFunc = append(tmux.HandlerFunc, d.func_doc(reg, f))
			}
		}
	}
	if d.NHandlers {
		for reg, f := range smux.NHandlers {
			if d.match(reg) {
				tmux.NHandlers = append(tmux.NHandlers, d.handler_doc(reg, f))
			}
		}
	}
	if d.NHandlerFunc {
		for reg, f := range smux.NHandlerFunc {
			if d.match(reg) {
				tmux.NHandlerFunc = append(tmux.NHandlerFunc, d.func_doc(reg, f))
			}
		}
	}
	return tmux
}

//srv
func (d *DocViewer) SrvHTTP(hs *routing.HTTPSession) routing.HResult {
	format := hs.CheckVal("fmt")
	if len(format) < 1 {
		format = FMT_HTML
	}
	if format == FMT_HTML {
		t, err := template.New("DocViewer").Parse(d.HTML)
		if err == nil {
			mux := d.BuildMux(hs.Mux)
			t.Execute(hs.W, map[string]interface{}{
				"Filters":      mux.Filters,
				"FilterFunc":   mux.FilterFunc,
				"Handlers":     mux.Handlers,
				"HandlerFunc":  mux.HandlerFunc,
				"NHandlers":    mux.NHandlers,
				"NHandlerFunc": mux.NHandlerFunc,
			})
		} else {
			hs.MsgResE(1, err.Error())
		}
		return routing.HRES_RETURN
	} else if format == FMT_JSON {
		return hs.MsgRes(d.BuildMux(hs.Mux))
	} else {
		return hs.MsgRes("error format")
	}
}

//doc
func (d *DocViewer) Doc() *Desc {
	return &Desc{
		ResV: &Mux{
			Filters:      []DocV{DocV{}},
			FilterFunc:   []DocV{DocV{}},
			Handlers:     []DocV{DocV{}},
			HandlerFunc:  []DocV{DocV{}},
			NHandlers:    []DocV{DocV{}},
			NHandlerFunc: []DocV{DocV{}},
		},
	}
}
