//Package wdco provider Parser to parse golang doc to web api documnet.
//it support multi command to special api items.
package wdoc

import (
	"github.com/Centny/gwf/log"
	"github.com/Centny/gwf/routing"
	"github.com/Centny/gwf/util"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path"
	"regexp"
	"sort"
	"strings"
	"time"
)

//the @arg command regex
var ARG_REG = regexp.MustCompile("^[^\\t]*\\t(required|optional|R|O)\\t.*$")

//the @ret command regex
var RET_REG = regexp.MustCompile("^[^\\t]*\\t(S|I|F|A|O|V|string|int|float|array|object|void)\\t.*$")

//multi tab regex
var multi_t = regexp.MustCompile("\t+")
var json_m = regexp.MustCompile("(?s)^\\{.*\\}$")
var cmd_m = regexp.MustCompile("\\@[a-z]*\\,")

//Parser handler.
type Handler interface {
	//is handler checker
	ISH(dir string, decl *ast.FuncDecl) bool
}

//normal parser handler impl
type NormalH struct {
}

//create normal handler
func NewNormalH() *NormalH {
	return &NormalH{}
}

//impl is handler checker
func (n *NormalH) ISH(dir string, decl *ast.FuncDecl) bool {
	if !decl.Name.IsExported() {
		return false
	}
	var ft = decl.Type
	if ft.Results.NumFields() == 0 {
		// if decl.Recv != nil && decl.Name.Name != "ServeHTTP" {
		// 	return false
		// }
		if ft.Params.NumFields() != 2 {
			return false
		}
		sel, ok := ft.Params.List[0].Type.(*ast.SelectorExpr)
		if !ok {
			return false
		}
		if sel.Sel.Name != "ResponseWriter" {
			return false
		}
		star, ok := ft.Params.List[1].Type.(*ast.StarExpr)
		if !ok {
			return false
		}
		sel = star.X.(*ast.SelectorExpr)
		return sel.Sel.Name == "Request"

	} else if ft.Results.NumFields() == 1 {
		// if decl.Recv != nil && decl.Name.Name != "SrvHTTP" {
		// 	return false
		// }
		sel, ok := ft.Results.List[0].Type.(*ast.SelectorExpr)
		if !ok {
			return false
		}
		if sel.Sel.Name != "HResult" {
			return false
		}
		if ft.Params.NumFields() != 1 {
			return false
		}
		star, ok := ft.Params.List[0].Type.(*ast.StarExpr)
		if !ok {
			return false
		}
		sel = star.X.(*ast.SelectorExpr)
		return sel.Sel.Name == "HTTPSession"
	} else {
		return false
	}
}

//the web api parser
type Parser struct {
	Running bool
	PkgPre  string
	WebPre  string
	H       Handler
	PS      map[string]*ast.Package
	FS      map[string]map[string]*ast.FuncDecl
	Web     *Webs

	doc *Wdoc //temp
}

//create parser
func NewParser(pkg_pre, web_pre, cmdf string) *Parser {
	return &Parser{
		PkgPre: pkg_pre,
		WebPre: web_pre,
		H:      NewNormalH(),
		PS:     map[string]*ast.Package{},
		FS:     map[string]map[string]*ast.FuncDecl{},
		Web:    NewWebs(web_pre+"/html", cmdf),
	}
}

//loop parse root directory by delay and include/exclude
func (p *Parser) LoopParse(root string, inc, exc []string, delay time.Duration) {
	p.LoopParseL(strings.Split(root, ","), inc, exc, delay)
}
func (p *Parser) LoopParseL(root []string, inc, exc []string, delay time.Duration) {
	p.Running = true
	for p.Running {
		err := p.ParseDirL(root, inc, exc)
		if err != nil {
			log.E("loop parse dir(%v),inc(%v),exc(%v) error->%v", root, inc, exc, err)
		}
		p.ParseWdoc(p.PkgPre)
		time.Sleep(delay * time.Millisecond)
	}
}
func (p *Parser) ParseDirL(roots []string, inc, exc []string) error {
	for _, root := range roots {
		err := p.ParseDir(root, inc, exc)
		if err != nil {
			return err
		}
	}
	return nil
}

//parser root and child directory by include/exclude.
func (p *Parser) ParseDir(root string, inc, exc []string) error {
	dirs := util.FilterDir(root, inc, exc)
	if len(dirs) < 1 {
		return util.Err("filter root dir(%v) is empty", root)
	}
	return p.Parse(root, dirs...)
}

//parser all directory
func (p *Parser) Parse(root string, dirs ...string) error {
	var fs = token.NewFileSet()
	for _, dir := range dirs {
		pkgs, err := parser.ParseDir(fs, dir, func(f os.FileInfo) bool {
			return !strings.HasSuffix(f.Name(), "_test.go")
		}, parser.ParseComments)
		if err == nil {
			p.parse_pkgs(dir, pkgs)
		} else {
			log.E("parse error with dirs(%v)->%v", len(dirs), err)
			return err
		}
	}
	log.D("parse success with dirs(%v) on root(%v)", len(dirs), root)
	return nil
}

//parser packages on directory and return the FuncDecl map
func (p *Parser) parse_pkgs(dir string, pkgs map[string]*ast.Package) {
	for _, pkg := range pkgs {
		doutf := p.parse_pkg_(dir, pkg)
		if len(doutf) > 0 {
			p.FS[dir] = doutf
			p.PS[dir] = pkg
		}
	}
}

//parser package and return the FuncDecl map
func (p *Parser) parse_pkg_(dir string, pkg *ast.Package) map[string]*ast.FuncDecl {
	outf := map[string]*ast.FuncDecl{}
	for _, f := range pkg.Files {
		p.parse_file(dir, f, outf)
	}
	return outf
}

//parser file and return the FuncDecl map
func (p *Parser) parse_file(dir string, f *ast.File, outf map[string]*ast.FuncDecl) {
	for _, decl := range f.Decls {
		fdecl, ok := decl.(*ast.FuncDecl)
		if !ok {
			continue
		}
		if !p.ISH(dir, fdecl) {
			continue
		}
		if fdecl.Recv == nil {
			outf[fdecl.Name.Name] = fdecl
			continue
		}
		rtype := fdecl.Recv.List[0].Type
		sel, ok := rtype.(*ast.Ident)
		if !ok {
			sel = rtype.(*ast.StarExpr).X.(*ast.Ident)
		}
		outf[sel.Name+"."+fdecl.Name.Name] = fdecl
	}
}

//is web api handler H
func (p *Parser) ISH(dir string, decl *ast.FuncDecl) bool {
	return p.H.ISH(dir, decl)
}

//do arg/ret command
func (p *Parser) do_arg_ret(cmd, text string, valid *regexp.Regexp, arg *Arg) {
	lines := strings.Split(text, "\n")
	arg.Desc = strings.Trim(lines[0], " \t")
	if len(lines) < 2 {
		return
	}
	lines = lines[1:]
	var sidx = -1
	for idx, line := range lines {
		line = strings.Trim(line, " \t")
		line = multi_t.ReplaceAllString(line, "\t")
		if !valid.MatchString(line) {
			sidx = idx
			break
		}
		vals := strings.SplitN(line, "\t", 3)
		arg.Items = append(arg.Items, Item{
			Name: vals[0],
			Type: vals[1],
			Desc: vals[2],
		})
	}
	sort.Sort(items_l(arg.Items))
	if sidx > -1 {
		var ctext = ""
		var matched = false
		for i := sidx; i < len(lines); i++ {
			var line = strings.Trim(lines[i], " \t\n")
			if len(ctext) < 1 && len(line) < 1 {
				continue
			}
			if !(matched || strings.HasPrefix(line, "[") || strings.HasPrefix(line, "{")) {
				//remove fist line is not json
				matched = true
				continue
			}
			ctext += "\n" + line
			matched = true
		}
		ctext = strings.Trim(ctext, " \t\n")
		ctext = strings.TrimPrefix(ctext, "样例")
		ctext = strings.TrimPrefix(ctext, "example")
		ctext = strings.Trim(ctext, " \t\n")
		ctexts := strings.Split(ctext, "\n")
		var cmt, json string
		var cmt_reg = regexp.MustCompile("^//.*$")
		var do_append = func() {
			if len(json) < 1 {
				return
			}
			cm, err := util.Json2Map(json)
			if err == nil {
				arg.Example = append(arg.Example, &Example{
					Cmt:  cmt,
					Data: cm,
				})
			} else {
				if json_m.MatchString(ctext) {
					log.W("parsing liked json data error(%v) by \n%v\n", ctext, err)
				}
				arg.Example = append(arg.Example, &Example{
					Cmt:  cmt,
					Data: json,
				})
			}
		}
		for _, txt := range ctexts {
			txt = strings.TrimSpace(txt)
			if len(txt) < 1 {
				continue
			}
			if !cmt_reg.MatchString(txt) {
				json += txt + "\n"
				continue
			}
			if len(json) < 1 {
				cmt += txt + "\n"
				continue
			}
			do_append()
			json = ""
			cmt = txt
		}
		do_append()
	}
}

//do url command
func (p *Parser) do_url(text string, url *Url) {
	lines := strings.Split(text, "\n")
	url.Desc = strings.Trim(lines[0], " \t")
	if len(lines) < 2 {
		return
	}
	line := strings.Trim(lines[1], " \t")
	line = multi_t.ReplaceAllString(line, "\t")
	vals := strings.SplitN(line, "\t", 3)
	url.Path = vals[0]
	if len(vals) > 1 {
		url.Method = vals[1]
	}
	if len(vals) > 2 {
		url.Ctype = vals[2]
	}
}

//do auth command
func (p *Parser) do_author(text string, author *Author) {
	line := strings.Trim(text, " \t")
	line = multi_t.ReplaceAllString(line, "\t")
	vals := strings.SplitN(line, ",", 3)
	author.Name = vals[0]
	if len(vals) > 1 {
		date, err := time.Parse("2006-01-02", vals[1])
		if err == nil {
			author.Date = util.Timestamp(date)
		} else {
			log.W("parsing date(%v) on line(%v) error->%v", vals[1], text, err)
		}
	}
	if len(vals) > 2 {
		author.Desc = vals[2]
	}
}

func (p *Parser) do_web(pkg_path, text string, web *Web) {
	line := strings.Trim(text, " \t")
	line = multi_t.ReplaceAllString(line, "\t")
	vals := strings.SplitN(line, ",", 3)
	if len(vals) < 2 {
		log.W("parsing web line(%v) error->%v", line, "must having key and index name")
		return
	}
	web.Key = vals[0]
	web.Index = vals[1]
	if len(vals) > 2 {
		web.Desc = vals[2]
	}
	log.D("Parser adding web by path(%v),key(%v),index(%v)", pkg_path, web.Key, web.Index)
	p.Web.AddMD2(web.Key, pkg_path, web.Index)
}

func (p *Parser) do_see(pkg_path, text string, see *See) {
	line := strings.Trim(text, " \t")
	line = multi_t.ReplaceAllString(line, "\t")
	vals := strings.SplitN(line, ",", 2)
	if len(vals) < 1 {
		log.W("parsing see line(%v) error->%v", line, "must having key and index name")
		return
	}
	if strings.HasPrefix(vals[0], ".") {
		vals[0] = pkg_path + strings.TrimPrefix(vals[0], ".")
	}
	see.Pkg, see.Name = path.Split(vals[0])
	if len(see.Pkg) < 1 {
		see.Pkg = pkg_path
	}
	var pkgs = strings.SplitN(see.Pkg, "src/", 2)
	if len(pkgs) > 1 {
		see.Pkg = pkgs[1]
	}
	see.Pkg = strings.Trim(see.Pkg, "/ ")
	if len(vals) > 1 {
		see.Desc = vals[1]
	}
}

//parse matched func to Func
func (p *Parser) Func2Map(path, fn string, f *ast.FuncDecl) *Func {
	var info = &Func{
		Name: fn,
	}
	if f.Doc == nil {
		return info
	}
	var doc = f.Doc.Text()
	doc = strings.Trim(doc, " \t")
	if strings.Contains(doc, "@ignore") {
		return nil
	}
	if len(doc) < 1 {
		return info
	}
	var reg = cmd_m
	// fmt.Println(doc)
	cmds := reg.FindAllString(doc, -1)
	dataes := reg.Split(doc, -1)
	desces := strings.SplitN(dataes[0], "\n", 2)
	info.Title = desces[0]
	if len(desces) > 1 {
		info.Desc = desces[1]
	}
	if len(cmds) < 1 || len(cmds) != len(dataes)-1 {
		return info
	}
	for idx, cmd := range cmds {
		var text = strings.Trim(dataes[idx+1], " \t\n")
		// log.D("parsing command(%v) by text(%v)", cmd, text)
		switch cmd {
		case "@url,":
			info.Url = &Url{}
			p.do_url(text, info.Url)
		case "@arg,":
			info.Arg = &Arg{}
			p.do_arg_ret("/arg", text, ARG_REG, info.Arg)
		case "@ret,":
			info.Ret = &Arg{}
			p.do_arg_ret("/ret", text, RET_REG, info.Ret)
		case "@tag,":
			info.Tags = []string{}
			info.Tags = strings.Split(text, ",")
		case "@author,":
			info.Author = &Author{}
			p.do_author(text, info.Author)
		case "@web,":
			var web = &Web{}
			p.do_web(path, text, web)
			info.WS = append(info.WS, web)
		case "@see,":
			var see = &See{}
			p.do_see(path, text, see)
			info.See = append(info.See, see)
		default:
			log.E("unknow command(%v) for data(%v)", cmd, text)
		}
	}
	return info
}

func (p *Parser) ParseWdoc(prefix string) *Wdoc {
	var res = &Wdoc{}
	var pkgs_ = []*Pkg{}
	var tags_ = map[string]int{}
	for name, fs := range p.FS {
		var tfs = []*Func{}
		for fn, f := range fs {
			ff := p.Func2Map(name, fn, f)
			if ff == nil {
				continue
			}
			tfs = append(tfs, ff)
			for _, tag := range ff.Tags {
				tags_[tag] += 1
			}
		}
		if len(tfs) < 1 {
			continue
		}
		sort.Sort(funcs_l(tfs))
		names := strings.SplitN(name, "src/", 2)
		if len(names) == 2 {
			name = names[1]
		}
		name = strings.TrimPrefix(name, prefix)
		pkgs_ = append(pkgs_, &Pkg{
			Name:  name,
			Funcs: tfs,
		})
	}
	sort.Sort(pkgs_l(pkgs_))
	res.Pkgs = pkgs_
	res.Tags = tags_
	p.doc = res
	return res
}
func (p *Parser) get_wdoc() *Wdoc {
	if p.doc == nil {
		p.ParseWdoc(p.PkgPre)
	}
	return p.doc
}

//parse and search all matched func to Wdoc
func (p *Parser) ToMv(key, tags string) *Wdoc {
	return p.get_wdoc().Filter(key, tags)
}

//parse all matched func to Wdoc
func (p *Parser) ToM() *Wdoc {
	return p.ToMv("", "")
}

//list all web api doc
//list all web api doc by key,tags
//@url,the normal GET request
//	~/wdoc	GET
//@arg,the normal query arguments
//	key		O	the search key for seaching doc
//	tags	O	the filter tag for filter Api
//	~/wdoc?key=xx&tags=wdoc,godoc
//@ret,the json result
//	pkgs		A	the package list
//	name		S	the unique name for package/function
//	title		S	the title for package/function
//	desc		S	the description for package/function/item.
//	tags		A	the function tags for filter or seach.
//	funcs		A	the function list on package.
//	url			O	the web api url info object, contain method/path/desc field
//	method		S	the web api request method on GET/POST
//	path		S	the web api releative path
//	arg/type	S	the argument item type, options:R/required,O/optional
//	ret/type	S	the return item value type, options:S/string,I/int,F/float,A/array,O/object,V/void.
//	example		V	the example value. exammple is object type when return value is json, other case is string type
/*	example
	{
		"tags": {
			"godoc": 1,
			"wdoc": 1
		},
		"pkgs": [{
			"name": "github.com/Centny/gwf/wdoc/",
			"funcs": [{
				"name": "Parser.SrvHTTP",
				"title": "list all web api doc",
				"tags": ["wdoc", "godoc"],
				"url": {
					"path": "~/wdoc",
					"method": "GET",
					"desc": "the normal GET request"
				},
				"arg": {
					"items": [{
						"name": "key",
						"type": "O",
						"desc": "the search key for seaching doc"
					}],
					"desc": "the normal query arguments",
					"example": "~/wdoc?key=xx\u0026tags=wdoc,godoc"
				},
				"ret": {
					"items": [{
						"name": "pkgs",
						"type": "A",
						"desc": "the package list"
					}],
					"desc": "the json result"
				},
				"author": {
					"name": "Centny",
					"date": 1453939200000
				}
			}]
		}]
	}
*/
//@tag,wdoc,godoc
//
//@author,Centny,2016-01-28
//@web,readme_cn,README_cn.md,the chinese doc
//@see,Webs.SrvHTTP,the webs
//@see,./Webs.SrvHTTP
func (p *Parser) SrvHTTP(hs *routing.HTTPSession) routing.HResult {
	var path = hs.R.URL.Path
	path = strings.TrimPrefix(path, p.WebPre)
	path = strings.TrimPrefix(path, "/")
	if strings.HasPrefix(path, "html") {
		return p.Web.SrvHTTP(hs)
	}
	var key string = hs.CheckVal("key")
	var tags string = hs.CheckVal("tags")
	hs.JsonRes(p.ToMv(".*"+key+".*", tags))
	return routing.HRES_RETURN
}

// func (p *Parser) LoadHtml(hs *routing.HTTPSession) routing.HResult {
// 	var _, html = path.Split(strings.TrimSuffix(hs.R.URL.Path, "/"))
// 	html = strings.Trim(html, "/ \t")
// 	html = strings.TrimSuffix(html, ".html")
// 	if len(html) < 1 {
// 		hs.W.WriteHeader(404)
// 		return routing.HRES_RETURN
// 	}
// 	if content, ok := p.HTML[html]; ok {
// 		log.D("Parser load html by key(%v) success", html)
// 		hs.SendT(content, "text/html")
// 	} else {
// 		log.D("Parser load html by key(%v) fail with key is not found", html)
// 		hs.W.WriteHeader(404)
// 	}
// 	return routing.HRES_RETURN
// }
