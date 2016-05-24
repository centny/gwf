package filter

import (
	"github.com/Centny/gwf/log"
	"github.com/Centny/gwf/routing"
	"github.com/Centny/gwf/util"
	"html/template"
	"io/ioutil"
	"net/url"
	"path/filepath"
	"regexp"
	"strings"
)

var REG_EXT_LINE = regexp.MustCompile("^<!--[\\s]*R:.*-->$")

type RenderH interface {
	LoadData(r *Render, hs *routing.HTTPSession) (tmpl *TmplF, data interface{}, err error)
}

type RenderDataH interface {
	LoadData(r *Render, hs *routing.HTTPSession, tmpl *TmplF, args url.Values, info interface{}) (data interface{}, err error)
}

type RENDER_DATA_F func(r *Render, hs *routing.HTTPSession, tmpl *TmplF, args url.Values, info interface{}) (data interface{}, err error)

func (f RENDER_DATA_F) LoadData(r *Render, hs *routing.HTTPSession, tmpl *TmplF, args url.Values, info interface{}) (data interface{}, err error) {
	return f(r, hs, tmpl, args, info)
}

type RenderWebData struct {
	Url  string
	Path string
}

func NewRenderWebData(url string) *RenderWebData {
	return &RenderWebData{Url: url}
}

func (r *RenderWebData) LoadData(rd *Render, hs *routing.HTTPSession, tmpl *TmplF, args url.Values, info interface{}) (data interface{}, err error) {
	var url string
	if strings.Contains(r.Url, "?") {
		url = r.Url + "&" + args.Encode()
	} else {
		url = r.Url + "?" + args.Encode()
	}
	res, err := util.HGet2("%v", url)
	if err == nil && len(r.Path) > 0 {
		data, err = res.ValP(r.Path)
	}
	if err != nil {
		err = util.Err("RenderWebData do request by url(%v) fail with error(%v)->%v", url, err, util.S2Json(res))
	}
	return data, err
}

type RenderNamedF struct {
	DataF map[string]RenderDataH
}

func NewRenderNamedF() *RenderNamedF {
	return &RenderNamedF{
		DataF: map[string]RenderDataH{},
	}
}

func (r *RenderNamedF) AddDataF(key string, f RENDER_DATA_F) {
	r.DataF[key] = f
}

func (r *RenderNamedF) AddDataH(key string, h RenderDataH) {
	r.DataF[key] = h
}
func (r *RenderNamedF) LoadData(rd *Render, hs *routing.HTTPSession) (tmpl *TmplF, data interface{}, err error) {
	var args url.Values
	tmpl, args, err = rd.LoadTmpF(hs)
	if err != nil {
		return
	}
	dataf, ok := r.DataF[tmpl.Key]
	if !ok {
		err = util.Err("the data provider by key(%v) is not found", tmpl.Key)
		return
	}
	data, err = dataf.LoadData(rd, hs, tmpl, args, nil)
	if err != nil {
		err = util.Err("load provider(%v) data by args(%v) fail with error->%v", tmpl.Key, args.Encode(), err)
	}
	return
}

type TmplF struct {
	Path string             `json:"path"`
	Text string             `json:"text"`
	Key  string             `json:"key"`
	U    *url.URL           `json:"-"`
	T    *template.Template `json:"-"`
}

type Render struct {
	Dir string
	H   RenderH
}

func NewRender(dir string, h RenderH) *Render {
	return &Render{
		Dir: dir,
		H:   h,
	}
}

func (r *Render) NewTmplFv(path string) (tmpl *TmplF, err error) {
	tmpl = &TmplF{}
	tmpl.Path = path
	fpath := filepath.Join(r.Dir, path)
	bys, err := ioutil.ReadFile(fpath)
	if err != nil {
		err = util.Err("read template file(%v) fail with error->%v", fpath, err)
		return nil, err
	}
	tmpl.Text = string(bys)
	ext := strings.SplitN(tmpl.Text, "\n", 2)[0]
	ext = strings.TrimSpace(ext)
	if REG_EXT_LINE.MatchString(ext) {
		ext = strings.TrimPrefix(ext, "<!--")
		ext = strings.TrimSuffix(ext, "-->")
		ext = strings.TrimSpace(ext)
		ext = strings.TrimPrefix(ext, "R:")
		tmpl.U, err = url.Parse(ext)
		if err != nil {
			err = util.Err("parsing extern line(%v) on file(%v) fail with error->%v", fpath, err)
			return nil, err
		}
		tmpl.Key = tmpl.U.Path
	}
	tmpl.T, err = template.New(tmpl.Path).Parse(tmpl.Text)
	return
}

func (r *Render) LoadTmpFv(path string) (*TmplF, error) {
	return r.NewTmplFv(path)
}

func (r *Render) LoadTmpF(hs *routing.HTTPSession) (*TmplF, url.Values, error) {
	path := strings.TrimSpace(hs.R.URL.Path)
	path = strings.Trim(path, "/ \t")
	if len(path) < 1 {
		path = "index.html"
	}
	tmpl, err := r.LoadTmpFv(path)
	if err != nil {
		return nil, nil, util.Err("loading template fail->%v", err)
	}
	targs := hs.R.URL.Query()
	if tmpl.U != nil {
		for key, vals := range tmpl.U.Query() {
			targs[key] = vals
		}
	}
	return tmpl, targs, nil
}

func (r *Render) SrvHTTP(hs *routing.HTTPSession) routing.HResult {
	log.D("Render doing %v", hs.R.URL.Path)
	tmpl, data, err := r.H.LoadData(r, hs)
	if err != nil {
		hs.W.WriteHeader(500)
		return hs.Printf("loading data fail with error->%v", err)
	}
	if hs.RVal("_data_") == "1" {
		return hs.JRes(data)
	}
	err = tmpl.T.Execute(hs.W, data)
	if err != nil {
		return hs.Printf("execute template fail with error(%v) by data:\n%v", err, util.S2Json(data))
	}
	return routing.HRES_RETURN
}
