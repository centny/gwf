package jcr

import (
	"fmt"
	"github.com/Centny/Cny4go/log"
	"github.com/Centny/Cny4go/routing"
	"github.com/Centny/Cny4go/routing/filter"
	"github.com/Centny/Cny4go/util"
	"net/http"
	"path/filepath"
	"sync"
)

var _jcr_js_ string = `
function jcr(){
	if(window.__coverage__){
		$.ajax({
	        url: jcr_store,
	        data: "cover=" + JSON.stringify(window.__coverage__),
	        type: "POST",
	        async: false,
	        success: function(data) {
	        	if(data.code){
	            	alter(data.msg)
	        	}
	        }
		});
	}
}
window.onunload=jcr;
`

type Config struct {
	Name   string `json:"name"`
	Dir    string `json:"dir"`
	Count  int    `json:"count"`
	Listen string `json:"listen"`
}

func (c *Config) SavePath() string {
	defer func() { c.Count++ }()
	return fmt.Sprintf("%s/%s_%.3d.json", c.Dir, c.Name, c.Count)
}

// func (c *Config) JcrJs() string {
// 	return fmt.Sprintf("var jcr_store='%s';\n%s", c.Store, _jcr_js_)
// }

var _conf_ *Config = &Config{}

func JcrConf(hs *routing.HTTPSession) routing.HResult {
	hs.JsonRes(_conf_)
	return routing.HRES_RETURN
}
func JcrStore(hs *routing.HTTPSession) routing.HResult {
	var cover string
	err := hs.ValidRVal(`
		cover,R|S,L:0
		`, &cover)
	if err != nil {
		return hs.MsgResE(1, err.Error())
	}
	spath := _conf_.SavePath()
	log.D("saving coverage report to %s", spath)
	err = util.FWrite(spath, cover)
	if err != nil {
		log.E("saving coverage report to %s error:%s", spath, err.Error())
		return hs.MsgResE(1, err.Error())
	} else {
		return hs.MsgRes("SUCCESS")
	}
}
func JcrJs(hs *routing.HTTPSession) routing.HResult {
	hs.SendT(fmt.Sprintf("var jcr_store='http://%s%s/store';%s",
		hs.R.Host, filepath.Dir(hs.R.RequestURI), _jcr_js_),
		"text/javascript;charset=utf-8")
	return routing.HRES_RETURN
}
func JcrExit(hs *routing.HTTPSession) routing.HResult {
	log.D("Jcr receiving exit command...")
	StopSrv()
	return hs.MsgRes("SUCCESS")
}
func NewJcrMux() *routing.SessionMux {
	mux := routing.NewSessionMux2("/jcr")
	cors := filter.NewCORS()
	cors.AddSite("*")
	mux.HFilter("^/(conf|store|jcr)(\\?.*)?$", cors)
	mux.HFunc("^/conf(\\?.*)?$", JcrConf)
	mux.HFunc("^/store(\\?.*)?$", JcrStore)
	mux.HFunc("^/jcr\\.js(\\?.*)?$", JcrJs)
	mux.HFunc("^/exit(\\?.*)?$", JcrExit)
	return mux
}

var lock sync.WaitGroup
var s_running bool
var s http.Server

func StartSrv(name, dir, port string) {
	_conf_.Name = name
	_conf_.Dir = dir
	_conf_.Listen = port
	mux := http.NewServeMux()
	mux.Handle("/jcr/", NewJcrMux())
	log.D("running server on %v", _conf_.Listen)
	s = http.Server{Addr: _conf_.Listen, Handler: mux}
	err := s.ListenAndServe()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}

//run the server.
func RunSrv(name, dir, port string) {
	s_running = true
	lock.Add(1)
	go StartSrv(name, dir, port)
	lock.Wait()
	s_running = false
	log.D("jcr server stopped...")
}

//stop the server.
func StopSrv() {
	if s_running {
		lock.Done()
	}
}
