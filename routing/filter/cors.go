package filter

import (
	"fmt"
	"strings"
	// "github.com/Centny/gwf/log"
	"github.com/Centny/gwf/routing"
	"net/http"
)

type CORS struct {
	Sites   map[string]int //sites for access
	Headers []string
	Methods []string
}

func (c *CORS) exec(w http.ResponseWriter, r *http.Request) routing.HResult {
	origin := r.Header.Get("Origin")
	found := func(origin string) routing.HResult {
		// log.D("sending CORS to %s", origin)
		w.Header().Set("Access-Control-Allow-Origin", origin)
		if len(c.Headers) > 0 {
			w.Header().Set("Access-Control-Allow-Headers", strings.Join(c.Headers, ", "))
		}
		if len(c.Methods) > 0 {
			w.Header().Set("Access-Control-Allow-Methods", strings.Join(c.Methods, ", "))
		}
		if r.Method == "OPTIONS" {
			return routing.HRES_RETURN
		} else {
			return routing.HRES_CONTINUE
		}
	}
	if len(origin) > 0 {
		if v, ok := c.Sites["*"]; ok && v > 0 {
			return found("*")
		} else if v, ok := c.Sites[origin]; ok && v > 0 {
			return found(origin)
		} else {
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			return routing.HRES_RETURN
		}
	} else {
		return routing.HRES_CONTINUE
	}
}
func (c *CORS) SrvHTTP(hs *routing.HTTPSession) routing.HResult {
	return c.exec(hs.W, hs.R)
}
func (c *CORS) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c.exec(w, r)
}
func (c *CORS) AddSite(site string) {
	c.Sites[site] = 1
}
func (c *CORS) DelSite(site string) {
	delete(c.Sites, site)
}
func NewCORS() *CORS {
	cors := &CORS{}
	cors.Sites = map[string]int{}
	return cors
}
func NewCORS2(site string) *CORS {
	cors := &CORS{}
	cors.Sites = map[string]int{
		site: 1,
	}
	return cors
}

func NewCORS_GET_POST(site string) *CORS {
	cors := NewCORS()
	cors.AddSite(site)
	cors.Methods = []string{"GET", "POST"}
	cors.Headers = []string{"Origin", "X-Requested-With", "Content-Type", "Accept"}
	return cors
}

func NewCORS_All() *CORS {
	return NewCORS_GET_POST("*")
}

type P3P struct {
	val string
}

func NewP3P(options []string) *P3P {
	var p3p = &P3P{}
	p3p.SetOptions(options)
	return p3p
}

func (p *P3P) SrvHTTP(hs *routing.HTTPSession) routing.HResult {
	hs.W.Header().Set("P3P", p.val)
	return routing.HRES_CONTINUE
}

func (p *P3P) SetOptions(options []string) {
	p.val = fmt.Sprintf("CP=\" %v \"", strings.Join(options, " "))
}
