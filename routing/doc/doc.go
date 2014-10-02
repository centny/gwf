package doc

import (
	"github.com/Centny/gwf/routing"
	"reflect"
	"regexp"
	"runtime"
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

//api describe
type Desc struct {
	Title  string
	Url    string      //example URL.
	ArgsR  interface{} //required arguments.
	ArgsO  interface{} //option arguments.
	ResV   interface{} //result
	Detail string      //detail
}

//register api.
func (d Desc) Api(f interface{}) int {
	if d.ArgsR == nil {
		d.ArgsR = []map[string]string{}
	}
	if d.ArgsO == nil {
		d.ArgsO = []map[string]string{}
	}
	if d.ResV == nil {
		d.ResV = []map[string]string{}
	}
	ApiV(f, &d)
	return 0
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
	Doc     *Desc
	Marked  bool
	Pkg     string
}

//the doc viewer for building the SessionMux info.
type DocViewer struct {
	Incs         []*regexp.Regexp
	Excs         []*regexp.Regexp
	Items        []string //comment items.
	Pkg          bool
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
		Pkg:          true,
		Items:        []string{},
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
	typ := reflect.TypeOf(reflect.Indirect(reflect.ValueOf(f)).Interface())
	doc := DocV{}
	doc.Name = typ.String()
	doc.Pattern = reg.String()
	if d.Pkg {
		doc.Pkg = typ.PkgPath()
	}
	if dd, ok := f.(Docable); ok {
		doc.Doc = dd.Doc()
		doc.Marked = true
	} else {
		doc.Doc = nil
		doc.Marked = false
	}
	return doc
}
func (d *DocViewer) func_doc(reg *regexp.Regexp, f interface{}) DocV {
	fnc := runtime.FuncForPC(reflect.ValueOf(f).Pointer())
	doc := DocV{}
	doc.Name = regexp.MustCompile("^.*\\.").ReplaceAllString(fnc.Name(), "")
	doc.Pattern = reg.String()
	if d.Pkg {
		doc.Pkg = regexp.MustCompile("\\.[^\\.]*$").ReplaceAllString(fnc.Name(), "")
	}
	if dd, ok := Marked[fnc.Name()]; ok {
		doc.Doc = dd
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
	return hs.MsgRes(d.BuildMux(hs.Mux))
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
