package wdoc

import (
	"time"
	// "fmt"
	"github.com/Centny/gwf/log"
	"github.com/Centny/gwf/routing"
	"github.com/Centny/gwf/util"
	"go/ast"
	"go/parser"
	"go/token"
	// "html/template"
	"os"
	"regexp"
	"sort"
	"strings"
)

var ARG_REG = regexp.MustCompile("^[^\\t]*\\t(required|optional|R|O)\\t.*$")

var RET_REG = regexp.MustCompile("^[^\\t]*\\t(S|I|F|A|O|string|int|float|array|object)\\t.*$")

type Handler interface {
	ISH(dir string, decl *ast.FuncDecl) bool
	Filter(f os.FileInfo) bool
}

type NormalH struct {
	Inc []*regexp.Regexp
	Exc []*regexp.Regexp
}

func NewNormalH() *NormalH {
	return &NormalH{
		Exc: []*regexp.Regexp{regexp.MustCompile(".*_test.go")},
	}
}

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
func (n *NormalH) Filter(f os.FileInfo) bool {
	var name = f.Name()
	if len(n.Inc) > 0 {
		for _, inc := range n.Inc {
			if inc.MatchString(name) {
				return true
			}
		}
		return false
	}
	if len(n.Exc) > 0 {
		for _, exc := range n.Exc {
			if exc.MatchString(name) {
				return false
			}
		}
		return true
	}
	return true
}

//xxx
//
//	sds	abc <a>skks</a> fdsf dsf kskd fkd sfjds kf djslfj sdl fks jdfk dsjfsd fjdksd lfjs dkfjs df jsd kf  jskdf abc skks fdsf dsf kskd fkd sfjds kf djslfj sdl fks jdfk dsjfsd fjdksd lfjs dkfjs df jsd kf  jskdf
//	fs	required	xxsds
//
type Parser struct {
	Pre string
	H   Handler
	PS  map[string]*ast.Package
	FS  map[string]map[string]*ast.FuncDecl
}

func NewParser() *Parser {
	return &Parser{
		H:  &NormalH{},
		PS: map[string]*ast.Package{},
		FS: map[string]map[string]*ast.FuncDecl{},
	}
}

func (p *Parser) LoopParse(root string, inc, exc []string, delay time.Duration) {
	for {
		err := p.ParseDir(root, inc, exc)
		if err != nil {
			log.E("loop parse dir(%v) error->%v", err)
		}
		time.Sleep(delay * time.Millisecond)
	}
}

func (p *Parser) ParseDir(root string, inc, exc []string) error {
	dirs := util.FilterDir(root, inc, exc)
	if len(dirs) < 1 {
		return util.Err("filter root dir(%v) is empty", root)
	}
	err := p.Parse(dirs...)
	if err == nil {
		log.D("parse dir(%v) success with dirs(%v),inc(%v),exc(%v)", root, len(dirs), inc, exc)
	} else {
		log.E("parse dir(%v) error with dirs(%v),inc(%v),exc(%v)->%v", root, len(dirs), inc, exc, err)
	}
	return err
}
func (p *Parser) Parse(dirs ...string) error {
	var fs = token.NewFileSet()
	for _, dir := range dirs {
		pkgs, err := parser.ParseDir(fs, dir, p.Filter, parser.ParseComments)
		if err == nil {
			p.parse_pkgs(dir, pkgs)
		} else {
			return err
		}
	}
	return nil
}

func (p *Parser) parse_pkgs(dir string, pkgs map[string]*ast.Package) {
	for _, pkg := range pkgs {
		doutf := p.parse_pkg_(dir, pkg)
		if len(doutf) > 0 {
			p.FS[dir] = doutf
			p.PS[dir] = pkg
		}
	}
}

func (p *Parser) parse_pkg_(dir string, pkg *ast.Package) map[string]*ast.FuncDecl {
	outf := map[string]*ast.FuncDecl{}
	for _, f := range pkg.Files {
		p.parse_file(dir, f, outf)
	}
	return outf
}

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

func (p *Parser) Filter(f os.FileInfo) bool {
	return p.H.Filter(f)
}

func (p *Parser) ISH(dir string, decl *ast.FuncDecl) bool {
	return p.H.ISH(dir, decl)
}

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
	if sidx > -1 {
		var ctext = ""
		for i := sidx; i < len(lines); i++ {
			ctext += "\n" + lines[i]
		}
		ctext = strings.Trim(ctext, " \t")
		cm, err := util.Json2Map(ctext)
		if err == nil {
			arg.Example = cm
		} else {
			arg.Example = strings.Trim(ctext, " \t\n")
		}
	}
}
func (p *Parser) do_url(text string, url *Url) {
	lines := strings.Split(text, "\n")
	url.Desc = strings.Trim(lines[0], " \t")
	if len(lines) < 2 {
		return
	}
	vals := strings.SplitN(strings.Trim(lines[1], " \t"), "\t", 3)
	url.Path = vals[0]
	if len(vals) > 1 {
		url.Method = vals[1]
	}
	if len(vals) > 2 {
		url.Ctype = vals[2]
	}
}
func (p *Parser) Func2Map(fn string, f *ast.FuncDecl) Func {
	var info = Func{
		Name: fn,
	}
	if f.Doc == nil {
		return info
	}
	var doc = f.Doc.Text()
	doc = strings.Trim(doc, " \t")
	if len(doc) < 1 {
		return info
	}
	reg := regexp.MustCompile("\\@[^\\,]*\\,")
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
			info.Tags = strings.Split(text, ",")
		default:
			log.E("unknow command(%v) for data(%v)", cmd, text)
		}
	}
	return info
}
func (p *Parser) ToMv(prefix, key, tags string) *Wdoc {
	var res = &Wdoc{}
	var pkgs = []Pkg{}
	for name, fs := range p.FS {
		var tfs = []Func{}
		for fn, f := range fs {
			ff := p.Func2Map(fn, f)
			if ff.Matched(key, tags) {
				tfs = append(tfs, ff)
			}
		}
		if len(tfs) < 1 {
			continue
		}
		sort.Sort(Funcs(tfs))
		names := strings.SplitN(name, "src/", 2)
		if len(names) == 2 {
			name = names[1]
		}
		name = strings.TrimPrefix(name, prefix)
		pkgs = append(pkgs, Pkg{
			Name:  name,
			Funcs: tfs,
		})
	}
	sort.Sort(Pkgs(pkgs))
	res.Pkgs = pkgs
	return res
}
func (p *Parser) ToM(prefix string) *Wdoc {
	return p.ToMv(prefix, "", "")
}

func (p *Parser) SrvHTTP(hs *routing.HTTPSession) routing.HResult {
	var key string
	var tags string
	err := hs.ValidCheckVal(`
		key,O|S,L:0;
		tags,O|S,L:0;
		`, &key, &tags)
	if err != nil {
		return hs.MsgResErr2(1, "arg-err", err)
	}
	// t, err := template.New("Parser").Parse(HTML)
	// if err == nil {
	// 	err = t.Execute(hs.W, p.ToM(p.Pre))
	// } else {
	// 	fmt.Fprintf(hs.W, "%v", err.Error())
	// }
	// return routing.HRES_RETURN
	return hs.MsgRes(p.ToMv(p.Pre, ".*"+key+".*", tags))
}
