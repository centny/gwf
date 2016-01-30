package wdoc

import (
	"encoding/xml"
	"regexp"
	"strings"
)

const DESC_L = 2

//the web api doc
type Wdoc struct {
	XMLName xml.Name       `xml:"coverage"`
	Tags    map[string]int `json:"tags,omitempty" xml:"-"`                //all tags
	Pkgs    []*Pkg         `json:"pkgs,omitempty" xml:"packages>package"` //packages
	Rate    float32        `json:"rate,omitempty" xml:"line-rate,attr"`   //rate
}

func (w *Wdoc) Marshal() ([]byte, error) {
	return xml.MarshalIndent(w, " ", "  ")
}

func (w *Wdoc) RateV() {
	for _, p := range w.Pkgs {
		p.RateV()
	}
}

//the package
type Pkg struct {
	Name  string  `json:"name,omitempty" xml:"name,attr"`      //the package name
	Funcs []*Func `json:"funcs,omitempty" xml:"classes>class"` //the functions
}

func (p *Pkg) RateV() {
	for _, f := range p.Funcs {
		f.RateV()
	}
}

//the func
type Func struct {
	Name     string    `json:"name,omitempty" xml:"name,attr"` //the func name
	Title    string    `json:"title,omitempty" xml:"-"`        //the func title
	Desc     string    `json:"desc,omitempty" xml:"-"`         //the func desc
	Tags     []string  `json:"tags,omitempty" xml:"-"`         //the func tags
	Url      *Url      `json:"url,omitempty" xml:"-"`          //the func url
	Arg      *Arg      `json:"arg,omitempty" xml:"-"`          //the func argument
	Ret      *Arg      `json:"ret,omitempty" xml:"-"`          //the func return
	Author   *Author   `json:"author,omitempty" xml:"-"`       //the func author
	Methods  []*Method `json:"-" xml:"methods>method"`         //the methods
	Filename string    `json:"-" xml:"filename,attr"`          //the filename
}

func (f *Func) RateV() {
	f.Filename = f.Name
	f.Methods = []*Method{}
	//
	f.Methods = append(f.Methods, NewMethod("name", desc_hits(f.Name)))
	f.Methods = append(f.Methods, NewMethod("title", desc_hits(f.Title)))
	f.Methods = append(f.Methods, NewMethod("desc", desc_hits(f.Desc)))
	if len(f.Tags) > 0 {
		f.Methods = append(f.Methods, NewMethod("tags", 1))
	} else {
		f.Methods = append(f.Methods, NewMethod("tags", 0))
	}
	if f.Url == nil {
		f.Methods = append(f.Methods, NewMethod("url", 0))
	} else {
		f.Methods = append(f.Methods, f.Url.RateV()...)
	}
	if f.Arg == nil {
		f.Methods = append(f.Methods, NewMethod("arg", 0))
	} else {
		f.Methods = append(f.Methods, f.Arg.RateV()...)
	}
	if f.Ret == nil {
		f.Methods = append(f.Methods, NewMethod("ret", 0))
	} else {
		f.Methods = append(f.Methods, f.Ret.RateV()...)
	}
	if f.Author == nil {
		f.Methods = append(f.Methods, NewMethod("author", 0))
	} else {
		f.Methods = append(f.Methods, f.Author.RateV()...)
	}
}

//chekc if matched by key,tags
func (f *Func) Matched(key, tags string) bool {
	if len(tags) > 0 {
		var tags_ = strings.Split(tags, ",")
		for _, tag := range tags_ {
			var matched = false
			for _, t := range f.Tags {
				if tag == t {
					matched = true
					break
				}
			}
			if !matched {
				return false
			}
		}
	}
	if len(key) > 0 {
		var reg = regexp.MustCompile(key)
		if reg.MatchString(f.Title) ||
			reg.MatchString(f.Desc) ||
			reg.MatchString(f.Name) {
			return true
		} else {
			return false
		}
	}
	return true
}

//the author
type Author struct {
	Name string `json:"name,omitempty"` //the author name
	Date int64  `json:"date,omitempty"` //the create date
	Desc string `json:"desc,omitempty"` //the auth desc
}

func (a *Author) RateV() []*Method {
	return []*Method{
		NewMethod("author.name", desc_hits(a.Name)),
		NewMethod("author.desc", desc_hits(a.Desc)),
	}
}

//the url
type Url struct {
	Path   string `json:"path,omitempty"`   //the url path
	Method string `json:"method,omitempty"` //the request method
	Ctype  string `json:"ctype,omitempty"`  //the content type
	Desc   string `json:"desc,omitempty"`   //the url dec
}

func (u *Url) RateV() []*Method {
	return []*Method{
		NewMethod("url.path", desc_hits(u.Path)),
		NewMethod("url.desc", desc_hits(u.Desc)),
	}
}

//the argument
type Arg struct {
	Items   []Item      `json:"items,omitempty"`   //the item list
	Desc    string      `json:"desc,omitempty"`    //the argument desc
	Ctype   string      `json:"ctype,omitempty"`   //the request content type
	Example interface{} `json:"example,omitempty"` //the example
}

func (a *Arg) RateV() []*Method {
	var ms = []*Method{}
	for _, i := range a.Items {
		ms = append(ms, NewMethod("arg."+i.Name, i.Hits()))
	}
	if a.Example == nil {
		ms = append(ms, NewMethod("arg.example", 0))
	} else {
		ms = append(ms, NewMethod("arg.example", 1))
	}
	ms = append(ms, NewMethod("arg.desc", desc_hits(a.Desc)))
	return ms
}

//the item
type Item struct {
	Name string `json:"name,omitempty"` //the item name
	Type string `json:"type,omitempty"` //the item type
	Desc string `json:"desc,omitempty"` //the item desc
}

func (i *Item) Hits() int {
	return desc_hits(i.Desc)
}

func desc_hits(desc string) int {
	if len(desc) > DESC_L {
		return 1
	} else {
		return 0
	}
}

type pkgs_l []*Pkg

func (p pkgs_l) Len() int {
	return len(p)
}
func (p pkgs_l) Less(i, j int) bool {
	return p[i].Name < p[j].Name
}
func (p pkgs_l) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

type funcs_l []*Func

func (f funcs_l) Len() int {
	return len(f)
}
func (f funcs_l) Less(i, j int) bool {
	return f[i].Name < f[j].Name
}
func (f funcs_l) Swap(i, j int) {
	f[i], f[j] = f[j], f[i]
}

type items_l []Item

func (i items_l) Len() int {
	return len(i)
}
func (iv items_l) Less(i, j int) bool {
	return iv[i].Name < iv[j].Name
}
func (iv items_l) Swap(i, j int) {
	iv[i], iv[j] = iv[j], iv[i]
}

type Method struct {
	Name      string  `xml:"name,attr"`
	Signature string  `xml:"signature,attr"`
	Lines     []*Line `xml:"lines>line"`
}

func NewMethod(name string, hits int) *Method {
	return &Method{
		Name: name,
		Lines: []*Line{
			&Line{
				Hits: hits,
			},
		},
	}
}

type Line struct {
	Number int `xml:"number,attr"`
	Hits   int `xml:"hits,attr"`
}
