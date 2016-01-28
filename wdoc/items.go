package wdoc

import (
	"regexp"
	"strings"
)

//the web api doc
type Wdoc struct {
	Tags map[string]int `json:"tags,omitempty"` //all tags
	Pkgs []Pkg          `json:"pkgs,omitempty"` //packages
}

//the package
type Pkg struct {
	Name  string `json:"name,omitempty"`  //the package name
	Funcs []Func `json:"funcs,omitempty"` //the functions
}

//the func
type Func struct {
	Name   string   `json:"name,omitempty"`   //the func name
	Title  string   `json:"title,omitempty"`  //the func title
	Desc   string   `json:"desc,omitempty"`   //the func desc
	Tags   []string `json:"tags,omitempty"`   //the func tags
	Url    *Url     `json:"url,omitempty"`    //the func url
	Arg    *Arg     `json:"arg,omitempty"`    //the func argument
	Ret    *Arg     `json:"ret,omitempty"`    //the func return
	Author *Author  `json:"author,omitempty"` //the func author
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

//the url
type Url struct {
	Path   string `json:"path,omitempty"`   //the url path
	Method string `json:"method,omitempty"` //the request method
	Ctype  string `json:"ctype,omitempty"`  //the content type
	Desc   string `json:"desc,omitempty"`   //the url dec
}

//the argument
type Arg struct {
	Items   []Item      `json:"items,omitempty"`   //the item list
	Desc    string      `json:"desc,omitempty"`    //the argument desc
	Ctype   string      `json:"ctype,omitempty"`   //the request content type
	Example interface{} `json:"example,omitempty"` //the example
}

//the item
type Item struct {
	Name string `json:"name,omitempty"` //the item name
	Type string `json:"type,omitempty"` //the item type
	Desc string `json:"desc,omitempty"` //the item desc
}

type pkgs_l []Pkg

func (p pkgs_l) Len() int {
	return len(p)
}
func (p pkgs_l) Less(i, j int) bool {
	return p[i].Name < p[j].Name
}
func (p pkgs_l) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

type funcs_l []Func

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
