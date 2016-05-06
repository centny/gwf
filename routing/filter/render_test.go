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
	rn.AddDataF("/abc", func(r *Render, hs *routing.HTTPSession, tmpl *TmplF, args url.Values, info interface{}) (interface{}, error) {
		return util.Map{"name": "abc"}, nil
	})
	rn.AddDataF("", func(r *Render, hs *routing.HTTPSession, tmpl *TmplF, args url.Values, info interface{}) (interface{}, error) {
		return util.Map{"name": "default"}, nil
	})
	web := NewRenderWebData("http://pes.dev.gdy.io/pub/api/page/GetPage")
	web.Path = "/data"
	rn.AddDataH("/web", web)
	fmt.Println(ts.G("/render_test1.html"))
	fmt.Println(ts.G("/render_test2.html"))
	fmt.Println(ts.G("/render_test3.html"))
	fmt.Println(ts.G("/render_test4.html"))
	fmt.Println(ts.G("/render_test5.html"))
	fmt.Println(ts.G("/render_test1.html?_data_=1"))
	//
	r = NewRender(".", rn)
	ts = httptest.NewServer2(r)
	rn.AddDataF("/abc", func(r *Render, hs *routing.HTTPSession, tmpl *TmplF, args url.Values, info interface{}) (interface{}, error) {
		return nil, util.Err("error")
	})
	fmt.Println(ts.G("/render_test1.html"))
	fmt.Println(ts.G("/render_test2.html"))
	fmt.Println(ts.G(""))
	//
}

// func TestHGet(t *testing.T) {
// 	// uu, err := url.Parse(`abc?keys=[{"key":"courselist","param":{"limit":10,"page":1}}]`)
// 	// if err != nil {
// 	// 	t.Error(err)
// 	// 	return
// 	// }
// 	fmt.Println(util.HGet2("http://pes.dev.gdy.io/pub/api/page/GetPage?keys=%v", "%5B%7B%22key%22%3A%22courselist%22%2C%22param%22%3A%7B%22limit%22%3A10%2C%22page%22%3A1%7D%7D%5D"))
// }
