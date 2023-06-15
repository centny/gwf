package filter

import (
	"bytes"
	"html/template"
	"io"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"

	"github.com/Centny/gwf/log"
	"github.com/Centny/gwf/routing"
	"github.com/Centny/gwf/util"
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
	Dir      string
	H        RenderH
	Err      string
	Funcs    template.FuncMap
	CacheErr bool
	CacheDir string
	latest   map[string][]byte
	cacheLck sync.RWMutex
}

func NewRender(dir string, h RenderH) *Render {
	return &Render{
		Dir:      dir,
		H:        h,
		Err:      "error.html",
		CacheErr: true,
		CacheDir: os.TempDir(),
		latest:   map[string][]byte{},
		cacheLck: sync.RWMutex{},
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
			err = util.Err("parsing extern line(%v) on file(%v) fail with error->%v", ext, fpath, err)
			return nil, err
		}
		tmpl.Key = tmpl.U.Path
	}
	stdtmpl := template.New(tmpl.Path)
	if r.Funcs != nil {
		stdtmpl = stdtmpl.Funcs(r.Funcs)
	}
	tmpl.T, err = stdtmpl.Parse(tmpl.Text)
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
	if hs.RVal("_data_") == "1" {
		_, data, err := r.H.LoadData(r, hs)
		if err == nil {
			return hs.JRes(data)
		}
		return hs.Printf("load data fail with %v", err)
	}
	buffer := bytes.NewBuffer(nil)
	_, _, err := r.prepareResponseData(buffer, hs)
	if err == nil {
		cache := buffer.Bytes()
		hs.W.Write(cache)
		err = r.storeCacheData(hs, cache)
		if err != nil {
			log.E("Render store cache data fail with %v", err)
		}
	} else {
		log.E("Render prepare response data fail with %v", err)
		cache, lerr := r.loadCacheData(hs)
		if lerr == nil && len(cache) > 0 {
			log.E("Render prepare response data fail with %v, and using cache(%v)", err, len(cache))
			hs.W.Write(cache)
		} else {
			log.E("Render prepare response data fail with %v, and load cache fail with len(%v),%v", err, len(cache), lerr)
			hs.Printf("Render prepare response data fail with %v, and load cache fail with len(%v),%v", err, len(cache), lerr)
		}
	}
	return routing.HRES_RETURN
}

func (r *Render) prepareResponseData(w io.Writer, hs *routing.HTTPSession) (tmpl *TmplF, data interface{}, err error) {
	tmpl, data, err = r.H.LoadData(r, hs)
	if err == nil {
		err = tmpl.T.Execute(w, data)
	} else {
		var bys []byte
		bys, err = ioutil.ReadFile(filepath.Join(r.Dir, r.Err))
		if err == nil {
			_, err = w.Write(bys)
		}
	}
	return
}

func (r *Render) cacheFilename(hs *routing.HTTPSession) (name string) {
	path := strings.TrimSpace(hs.R.URL.Path)
	path = strings.Trim(path, "/ \t")
	if len(path) < 1 {
		path = "index.html"
	}
	name = strings.Replace(strings.Replace(path, "/", "_", -1), "\\", "_", -1) + ".cache"
	return
}

func (r *Render) loadCacheData(hs *routing.HTTPSession) (cache []byte, err error) {
	if !r.CacheErr {
		return
	}
	filename := r.cacheFilename(hs)
	r.cacheLck.Lock()
	defer r.cacheLck.Unlock()
	cache = r.latest[filename]
	if len(cache) > 0 {
		return
	}
	if len(r.CacheDir) > 0 {
		cacheFile := filepath.Join(r.CacheDir, r.cacheFilename(hs))
		cache, err = ioutil.ReadFile(cacheFile)
		if err == nil {
			r.latest[filename] = cache
			log.D("Render read cache from %v success", cacheFile)
		}
	}
	return
}

func (r *Render) storeCacheData(hs *routing.HTTPSession, cache []byte) (err error) {
	if !r.CacheErr {
		return
	}
	filename := r.cacheFilename(hs)
	r.cacheLck.Lock()
	defer r.cacheLck.Unlock()
	having := r.latest[filename]
	if len(having) > 0 && bytes.Equal(having, cache) {
		return
	}
	r.latest[filename] = cache
	if len(r.CacheDir) > 0 {
		cacheFile := filepath.Join(r.CacheDir, r.cacheFilename(hs))
		log.D("Redner saving cache to file %v", cacheFile)
		err = os.Remove(cacheFile)
		if err == nil || os.IsNotExist(err) {
			err = ioutil.WriteFile(cacheFile, cache, os.ModePerm)
		}
	}
	return
}
