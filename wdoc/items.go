package wdoc

import (
	"regexp"
	"strings"
)

type Wdoc struct {
	Pkgs []Pkg `json:"pkgs,omitempty"`
}
type Pkg struct {
	Name  string `json:"name,omitempty"`
	Funcs []Func `json:"funcs,omitempty"`
}
type Func struct {
	Name  string   `json:"name,omitempty"`
	Title string   `json:"title,omitempty"`
	Desc  string   `json:"desc,omitempty"`
	Tags  []string `json:"tags,omitempty"`
	Url   *Url     `json:"url,omitempty"`
	Arg   *Arg     `json:"arg,omitempty"`
	Ret   *Arg     `json:"ret,omitempty"`
}

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

type Url struct {
	Path   string `json:"path,omitempty"`
	Method string `json:"method,omitempty"`
	Ctype  string `json:"ctype,omitempty"`
	Desc   string `json:"desc,omitempty"`
}
type Arg struct {
	Items   []Item      `json:"items,omitempty"`
	Desc    string      `json:"desc,omitempty"`
	Ctype   string      `json:"ctype,omitempty"`
	Example interface{} `json:"example,omitempty"`
}
type Item struct {
	Name string `json:"name,omitempty"`
	Type string `json:"type,omitempty"`
	Desc string `json:"desc,omitempty"`
}

type Pkgs []Pkg

func (p Pkgs) Len() int {
	return len(p)
}
func (p Pkgs) Less(i, j int) bool {
	return p[i].Name < p[j].Name
}
func (p Pkgs) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

type Funcs []Func

func (f Funcs) Len() int {
	return len(f)
}
func (f Funcs) Less(i, j int) bool {
	return f[i].Name < f[j].Name
}
func (f Funcs) Swap(i, j int) {
	f[i], f[j] = f[j], f[i]
}

type Items []Item

func (i Items) Len() int {
	return len(i)
}
func (iv Items) Less(i, j int) bool {
	return iv[i].Name < iv[j].Name
}
func (iv Items) Swap(i, j int) {
	iv[i], iv[j] = iv[j], iv[i]
}
