package filter

import (
	"fmt"
	"net/url"
	"os"
	"strings"
	"testing"

	"github.com/Centny/gwf/routing"
	"github.com/Centny/gwf/routing/httptest"
	"github.com/Centny/gwf/util"
)

func assertGet(ts *httptest.Server, expect string, trim bool, f string, args ...interface{}) {
	data, err := ts.G(f, args...)
	if err != nil {
		panic(err)
	}
	if trim {
		data = strings.Trim(data, "\r\n\t ")
	}
	if data != expect {
		panic(fmt.Sprintf("expect %v, but %v", []byte(expect), []byte(data)))
	}
}

func assertGetLike(ts *httptest.Server, expect string, f string, args ...interface{}) {
	data, err := ts.G(f, args...)
	if err != nil {
		panic(err)
	}
	if !strings.Contains(data, expect) {
		panic(fmt.Sprintf("expect %v, but %v", expect, data))
	}
}

func TestReander(t *testing.T) {
	util.Exec("rm -rf " + os.TempDir() + "/render_test*")
	var rn = NewRenderNamedF()
	var r = NewRender(".", rn)
	var ts = httptest.NewServer2(r)
	var abcVal = util.Map{"name": "abc"}
	rn.AddDataF("/abc", func(r *Render, hs *routing.HTTPSession, tmpl *TmplF, args url.Values, info interface{}) (interface{}, error) {
		return abcVal, nil
	})
	rn.AddDataF("", func(r *Render, hs *routing.HTTPSession, tmpl *TmplF, args url.Values, info interface{}) (interface{}, error) {
		return util.Map{"name": "default"}, nil
	})
	web := NewRenderWebData("http://pes.dev.gdy.io/pub/api/page/GetPage")
	web.Path = "/data"
	rn.AddDataH("/web", web)
	assertGet(ts, "abc", true, "/render_test1.html")
	assertGet(ts, "abc", true, "/render_test1.html")
	assertGet(ts, "default", true, "/render_test2.html")
	assertGetLike(ts, "render_test3.html", "/render_test3.html")
	assertGetLike(ts, "error.html", "/render_test4.html")
	assertGetLike(ts, "error.html", "/render_test5.html")
	assertGet(ts, `{"name":"abc"}`, true, "/render_test1.html?_data_=1")
	//
	//test cache error
	assertGetLike(ts, "render_test6.html", "/render_test6.html")
	abcVal = util.Map{"name": []string{"abc"}}
	assertGet(ts, "abc", true, "/render_test6.html")
	//using memory cache
	abcVal = util.Map{"name": "abc"}
	assertGet(ts, "abc", true, "/render_test6.html")
	//using file cache
	r.latest = map[string][]byte{} //clear cache
	assertGet(ts, "abc", true, "/render_test6.html")
	//
	fmt.Printf("test normal done...\n\n\n")
	//
	//
	r = NewRender(".", rn)
	ts = httptest.NewServer2(r)
	rn.AddDataF("/abc", func(r *Render, hs *routing.HTTPSession, tmpl *TmplF, args url.Values, info interface{}) (interface{}, error) {
		return nil, util.Err("error")
	})
	fmt.Println(ts.G("/render_test1.html"))
	fmt.Println(ts.G("/render_test2.html"))
	fmt.Println(ts.G(""))
	r.Err = "render_test1.html"
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
