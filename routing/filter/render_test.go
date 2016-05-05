package filter

import (
	"fmt"
	"github.com/Centny/gwf/routing"
	"github.com/Centny/gwf/routing/httptest"
	"github.com/Centny/gwf/util"
	"net/url"
	"testing"
)

func TestReander(t *testing.T) {
	var rn = NewRenderNamedF()
	var r = NewRender(".", rn)
	var ts = httptest.NewServer2(r)
	rn.AddDataF("/abc", func(r *Render, hs *routing.HTTPSession, tmpl *TmplF, args url.Values) (interface{}, error) {
		return util.Map{"name": "abc"}, nil
	})
	rn.AddDataF("", func(r *Render, hs *routing.HTTPSession, tmpl *TmplF, args url.Values) (interface{}, error) {
		return util.Map{"name": "default"}, nil
	})
	fmt.Println(ts.G("/render_test1.html"))
	fmt.Println(ts.G("/render_test2.html"))
	fmt.Println(ts.G("/render_test3.html"))
	fmt.Println(ts.G("/render_test4.html"))
	fmt.Println(ts.G("/render_test1.html?_data_=1"))
	//
	r = NewRender(".", rn)
	ts = httptest.NewServer2(r)
	rn.AddDataF("/abc", func(r *Render, hs *routing.HTTPSession, tmpl *TmplF, args url.Values) (interface{}, error) {
		return nil, util.Err("error")
	})
	fmt.Println(ts.G("/render_test1.html"))
	fmt.Println(ts.G("/render_test2.html"))
	fmt.Println(ts.G(""))
	//
}
