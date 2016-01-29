package wdoc

import (
	"encoding/xml"
	"regexp"
	"strings"
)

const DESC_L = 5

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

func (w *Wdoc) RateV() float32 {
	var rate float32 = 0
	for _, p := range w.Pkgs {
		rate += p.RateV()
	}
	if len(w.Pkgs) > 0 {
		rate /= float32(len(w.Pkgs))
	}
	w.Rate = rate
	return w.Rate
}

//the package
type Pkg struct {
	Name  string  `json:"name,omitempty" xml:"name,attr"`      //the package name
	Funcs []*Func `json:"funcs,omitempty" xml:"classes>class"` //the functions
	Rate  float32 `json:"rate,omitempty" xml:"line-rate,attr"` //rate
}

func (p *Pkg) RateV() float32 {
	var rate float32 = 0
	for _, f := range p.Funcs {
		rate += f.RateV()
	}
	if len(p.Funcs) > 0 {
		rate /= float32(len(p.Funcs))
	}
	p.Rate = rate
	return p.Rate
}

//the func
type Func struct {
	Name     string    `json:"name,omitempty" xml:"name,attr"`      //the func name
	Title    string    `json:"title,omitempty" xml:"-"`             //the func title
	Desc     string    `json:"desc,omitempty" xml:"-"`              //the func desc
	Tags     []string  `json:"tags,omitempty" xml:"-"`              //the func tags
	Url      *Url      `json:"url,omitempty" xml:"-"`               //the func url
	Arg      *Arg      `json:"arg,omitempty" xml:"-"`               //the func argument
	Ret      *Arg      `json:"ret,omitempty" xml:"-"`               //the func return
	Author   *Author   `json:"author,omitempty" xml:"-"`            //the func author
	Methods  []*Method `json:"-" xml:"methods>method"`              //the methods
	Rate     float32   `json:"rate,omitempty" xml:"line-rate,attr"` //rate
	Filename string    `json:"-" xml:"filename,attr"`
}

func (f *Func) RateV() float32 {
	//
	var title = desc_rate(f.Title)
	f.Methods = append(f.Methods, &Method{
		Name: "title",
		Rate: title,
	})
	//
	var desc = desc_rate(f.Desc)
	f.Methods = append(f.Methods, &Method{
		Name: "desc",
		Rate: desc,
	})
	//
	var tag float32 = 0
	if len(f.Tags) < 1 {
		tag = 0
	} else {
		tag = 1
	}
	f.Methods = append(f.Methods, &Method{
		Name: "tag",
		Rate: tag,
	})
	//
	var url float32 = 0
	if f.Url != nil {
		url = f.Url.RateV()
	}
	f.Methods = append(f.Methods, &Method{
		Name: "url",
		Rate: url,
	})
	//
	var arg float32 = 0
	if f.Arg != nil {
		arg = f.Arg.RateV()
	}
	f.Methods = append(f.Methods, &Method{
		Name: "arg",
		Rate: arg,
	})
	//
	var ret float32 = 0
	if f.Ret != nil {
		ret = f.Ret.RateV()
	}
	f.Methods = append(f.Methods, &Method{
		Name: "ret",
		Rate: ret,
	})
	//
	var author float32 = 0
	if f.Author != nil {
		author = f.Author.RateV()
	}
	f.Methods = append(f.Methods, &Method{
		Name: "author",
		Rate: author,
	})
	f.Rate = (title + desc + tag + url + arg + ret + author) / 7
	return f.Rate
}

type Method struct {
	Name      string  `xml:"name,attr"`
	Rate      float32 `xml:"line-rate,attr"`
	Signature string  `xml:"signature,attr"`
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

func (a *Author) RateV() float32 {
	if a.Date > 0 {
		return 1
	} else {
		return 0
	}
}

//the url
type Url struct {
	Path   string `json:"path,omitempty"`   //the url path
	Method string `json:"method,omitempty"` //the request method
	Ctype  string `json:"ctype,omitempty"`  //the content type
	Desc   string `json:"desc,omitempty"`   //the url dec
}

func (u *Url) RateV() float32 {
	return desc_rate(u.Desc)
}

//the argument
type Arg struct {
	Items   []Item      `json:"items,omitempty"`   //the item list
	Desc    string      `json:"desc,omitempty"`    //the argument desc
	Ctype   string      `json:"ctype,omitempty"`   //the request content type
	Example interface{} `json:"example,omitempty"` //the example
}

func (a *Arg) RateV() float32 {
	var item float32 = 0
	for _, i := range a.Items {
		item += i.RateV()
	}
	if len(a.Items) > 0 {
		item /= float32(len(a.Items))
	}
	var exm float32 = 0
	if a.Example == nil {
		exm = 0
	} else {
		exm = 1
	}
	var desc = desc_rate(a.Desc)
	return 0.1*item + 0.4*desc + 0.5*exm
}

//the item
type Item struct {
	Name string `json:"name,omitempty"` //the item name
	Type string `json:"type,omitempty"` //the item type
	Desc string `json:"desc,omitempty"` //the item desc
}

func (i *Item) RateV() float32 {
	return desc_rate(i.Desc)
}

func desc_rate(desc string) float32 {
	var clen = len(desc)
	if clen >= DESC_L {
		return 1
	} else {
		return float32(clen) / float32(DESC_L)
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
