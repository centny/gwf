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

func (p *Parser) do_arg_ret(cmd, text string, valid *regexp.Regexp, info util.Map) {
	lines := strings.Split(text, "\n")
	info.SetValP(cmd+"/desc", strings.Trim(lines[0], " \t"))
	if len(lines) < 2 {
		return
	}
	var items = []util.Map{}
	lines = lines[1:]
	var sidx = -1
	for idx, line := range lines {
		line = strings.Trim(line, " \t")
		if !valid.MatchString(line) {
			sidx = idx
			break
		}
		vals := strings.SplitN(line, "\t", 3)
		items = append(items, util.Map{
			"name": vals[0],
			"type": vals[1],
			"desc": vals[2],
		})
	}
	info.SetValP(cmd+"/items", items)
	if sidx > -1 {
		var ctext = ""
		for i := sidx; i < len(lines); i++ {
			ctext += "\n" + lines[i]
		}
		ctext = strings.Trim(ctext, " \t")
		cm, err := util.Json2Map(ctext)
		if err == nil {
			info.SetValP(cmd+"/example", cm)
		} else {
			info.SetValP(cmd+"/example", ctext)
		}
	}
}
func (p *Parser) do_url(text string, info util.Map) {
	lines := strings.Split(text, "\n")
	info.SetValP("/url/desc", strings.Trim(lines[0], " \t"))
	if len(lines) < 2 {
		return
	}
	vals := strings.SplitN(strings.Trim(lines[1], " \t"), "\t", 3)
	info.SetValP("/url/path", vals[0])
	if len(vals) > 1 {
		info.SetValP("/url/method", vals[1])
	}
	if len(vals) > 2 {
		info.SetValP("/url/ctype", vals[2])
	}
}
func (p *Parser) Func2Map(fn string, f *ast.FuncDecl) util.Map {
	var info = util.Map{
		"name":  fn,
		"title": "",
		"desc":  "",
		"tags":  []string{},
		"url": util.Map{
			"desc":   "",
			"path":   "",
			"method": "",
			"ctype":  "",
		},
		"arg": util.Map{
			"desc":    "",
			"items":   nil,
			"example": nil,
		},
		"ret": util.Map{
			"desc":    "",
			"items":   nil,
			"ctype":   "",
			"example": nil,
		},
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
	info.SetValP("/title", desces[0])
	if len(desces) > 1 {
		info.SetValP("/desc", desces[1])
	}
	if len(cmds) < 1 || len(cmds) != len(dataes)-1 {
		return info
	}
	for idx, cmd := range cmds {
		var text = strings.Trim(dataes[idx+1], " \t\n")
		switch cmd {
		case "@url,":
			p.do_url(text, info)
		case "@arg,":
			p.do_arg_ret("/arg", text, ARG_REG, info)
		case "@ret,":
			p.do_arg_ret("/ret", text, RET_REG, info)
		case "@tag,":
			info.SetValP("/tags", strings.Split(text, ","))
		default:
			log.E("unknow command(%v) for data(%v)", cmd, text)
		}
	}
	return info
}
func (p *Parser) ToM(prefix string) util.Map {
	var res = util.Map{}
	var pkgs = []util.Map{}
	for name, fs := range p.FS {
		var tfs = []util.Map{}
		for fn, f := range fs {
			tfs = append(tfs, p.Func2Map(fn, f))
		}
		sort.Sort(&util.MapSorter{
			Maps: tfs,
			Key:  "/name",
			Type: 2,
		})
		names := strings.SplitN(name, "src/", 2)
		if len(names) == 2 {
			name = names[1]
		}
		name = strings.TrimPrefix(name, prefix)
		pkgs = append(pkgs, util.Map{
			"name":  name,
			"items": tfs,
		})
	}
	sort.Sort(&util.MapSorter{
		Maps: pkgs,
		Key:  "/name",
		Type: 2,
	})
	res["pkgs"] = pkgs
	return res
}

func (p *Parser) SrvHTTP(hs *routing.HTTPSession) routing.HResult {
	// t, err := template.New("Parser").Parse(HTML)
	// if err == nil {
	// 	err = t.Execute(hs.W, p.ToM(p.Pre))
	// } else {
	// 	fmt.Fprintf(hs.W, "%v", err.Error())
	// }
	// return routing.HRES_RETURN
	return hs.MsgRes(p.ToM(p.Pre))
}
