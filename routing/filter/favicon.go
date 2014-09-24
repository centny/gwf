package filter

import (
	"fmt"
	"github.com/Centny/gwf/routing"
	"github.com/Centny/gwf/util"
	"net/http"
)

type Favicon struct {
	Path string
}

func (f *Favicon) SrvHTTP(hs *routing.HTTPSession) routing.HResult {
	fmt.Println(hs, f)
	hs.SendF("favicon.ico", f.Path, "image/x-icon", false)
	return routing.HRES_RETURN
}

func (f *Favicon) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	routing.SendF(w, r, "favicon.ico", f.Path, "image/x-icon", false)
}

func NewFavicon(path string) *Favicon {
	if !util.Fexists(path) {
		fmt.Println(fmt.Sprintf("%s not found", path))
		return nil
	}
	return &Favicon{Path: path}
}
