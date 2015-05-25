package sr

// import (
// 	"fmt"
// 	"github.com/Centny/gwf/routing/httptest"
// 	"github.com/Centny/gwf/util"
// 	"testing"
// )

// func TestMr(t *testing.T) {
// 	ts := httptest.NewServer2(NewMR(""))
// 	mv, err := ts.G2("/abc/notexist")
// 	if err != nil {
// 		t.Error(err.Error())
// 		return
// 	}
// 	if mv.IntVal("code") == 0 {
// 		t.Error("not error")
// 		return
// 	}
// 	//path
// 	mv, err = ts.G2("/a1/s?exec=set&type=S&data=%v", "sval1")
// 	if mv.IntVal("code") != 0 {
// 		fmt.Println(mv)
// 		t.Error("error")
// 		return
// 	}
// 	mv, err = ts.G2("/a2/s?exec=set&type=S&data=%v", "sval2")
// 	if mv.IntVal("code") != 0 {
// 		fmt.Println(mv)
// 		t.Error("error")
// 		return
// 	}
// 	mv, err = ts.G2("/s?exec=set&type=S&data=%v", "sval3")
// 	if mv.IntVal("code") != 0 {
// 		fmt.Println(mv)
// 		t.Error("error")
// 		return
// 	}
// 	mv, err = ts.G2("/a1/*")
// 	if mv.StrValP("/data/s") != "sval1" {
// 		fmt.Println(mv)
// 		t.Error("error")
// 		return
// 	}
// 	mv, err = ts.G2("/a2/*")
// 	if mv.StrValP("/data/s") != "sval2" {
// 		fmt.Println(mv)
// 		t.Error("error")
// 		return
// 	}
// 	mv, err = ts.G2("/*")
// 	if mv.StrValP("/data/s") != "sval3" {
// 		fmt.Println(mv)
// 		t.Error("error")
// 		return
// 	}
// 	//
// 	mv, err = ts.G2("/abc/a?exec=set&data=%v", util.S2Json(map[string]interface{}{
// 		"x1": 1,
// 		"x2": "xval",
// 	}))
// 	if mv.IntVal("code") != 0 {
// 		fmt.Println(mv)
// 		t.Error("error")
// 		return
// 	}
// 	mv, err = ts.G2("/abc/s?exec=set&type=S&data=%v", "sval")
// 	if mv.IntVal("code") != 0 {
// 		fmt.Println(mv)
// 		t.Error("error")
// 		return
// 	}
// 	mv, err = ts.G2("/abc/i?exec=set&type=I&data=%v", "1")
// 	if mv.IntVal("code") != 0 {
// 		fmt.Println(mv)
// 		t.Error("error")
// 		return
// 	}
// 	mv, err = ts.G2("/abc/f?exec=set&type=F&data=%v", "122.112")
// 	if mv.IntVal("code") != 0 {
// 		fmt.Println(mv)
// 		t.Error("error")
// 		return
// 	}
// 	mv, err = ts.G2("/abc/*")
// 	if mv.IntValP("/data/a/x1") != 1 {
// 		fmt.Println(mv)
// 		t.Error("error")
// 		return
// 	}
// 	if mv.StrValP("/data/a/x2") != "xval" {
// 		fmt.Println(mv)
// 		t.Error("error")
// 		return
// 	}
// 	if mv.FloatValP("/data/f") != 122.112 {
// 		fmt.Println(mv)
// 		t.Error("error")
// 		return
// 	}
// 	if mv.IntValP("/data/i") != 1 {
// 		fmt.Println(mv)
// 		t.Error("error")
// 		return
// 	}
// 	if mv.StrValP("/data/s") != "sval" {
// 		fmt.Println(mv)
// 		t.Error("error")
// 		return
// 	}
// 	//
// 	mv, err = ts.G2("/abc/f?exec=plus&type=F&data=%v", "1")
// 	if mv.IntVal("code") != 0 {
// 		fmt.Println(mv)
// 		t.Error("error")
// 		return
// 	}
// 	mv, err = ts.G2("/abc/i?exec=plus&type=I&data=%v", "1")
// 	if mv.IntVal("code") != 0 {
// 		fmt.Println(mv)
// 		t.Error("error")
// 		return
// 	}
// 	mv, err = ts.G2("/abc/*")
// 	if mv.FloatValP("/data/f") != 123.112 {
// 		fmt.Println(mv)
// 		t.Error("error")
// 		return
// 	}
// 	if mv.IntValP("/data/i") != 2 {
// 		fmt.Println(mv)
// 		t.Error("error")
// 		return
// 	}
// 	mv, err = ts.G2("/abc/i")
// 	if mv.IntValP("/data") != 2 {
// 		fmt.Println(mv)
// 		t.Error("error")
// 		return
// 	}
// 	//
// 	mv, err = ts.G2("/abc/NI?exec=plus&type=I&data=%v", "1")
// 	mv, err = ts.G2("/abc/NI?exec=plus&type=I&data=%v", "1")
// 	mv, err = ts.G2("/abc/NI?exec=plus&type=I&data=%v", "1")
// 	if mv.IntVal("code") != 0 {
// 		fmt.Println(mv)
// 		t.Error("error")
// 		return
// 	}
// 	mv, err = ts.G2("/abc/NF?exec=plus&type=F&data=%v", "1")
// 	mv, err = ts.G2("/abc/NF?exec=plus&type=F&data=%v", "1")
// 	mv, err = ts.G2("/abc/NF?exec=plus&type=F&data=%v", "1")
// 	if mv.IntVal("code") != 0 {
// 		fmt.Println(mv)
// 		t.Error("error")
// 		return
// 	}
// 	mv, err = ts.G2("/abc/*")
// 	if mv.FloatValP("/data/NF") != 3 {
// 		fmt.Println(mv)
// 		t.Error("error")
// 		return
// 	}
// 	if mv.IntValP("/data/NI") != 3 {
// 		fmt.Println(mv)
// 		t.Error("error")
// 		return
// 	}

// 	//root
// 	mv, err = ts.G2("/a?exec=set&data=%v", util.S2Json(map[string]interface{}{
// 		"x1": 1,
// 		"x2": "xval",
// 	}))
// 	if mv.IntVal("code") != 0 {
// 		fmt.Println(mv)
// 		t.Error("error")
// 		return
// 	}
// 	mv, err = ts.G2("/s?exec=set&type=S&data=%v", "sval")
// 	if mv.IntVal("code") != 0 {
// 		fmt.Println(mv)
// 		t.Error("error")
// 		return
// 	}
// 	mv, err = ts.G2("/i?exec=set&type=I&data=%v", "1")
// 	if mv.IntVal("code") != 0 {
// 		fmt.Println(mv)
// 		t.Error("error")
// 		return
// 	}
// 	mv, err = ts.G2("/f?exec=set&type=F&data=%v", "122.112")
// 	if mv.IntVal("code") != 0 {
// 		fmt.Println(mv)
// 		t.Error("error")
// 		return
// 	}
// 	mv, err = ts.G2("/*")
// 	if mv.IntValP("/data/a/x1") != 1 {
// 		fmt.Println(mv)
// 		t.Error("error")
// 		return
// 	}
// 	if mv.StrValP("/data/a/x2") != "xval" {
// 		fmt.Println(mv)
// 		t.Error("error")
// 		return
// 	}
// 	if mv.FloatValP("/data/f") != 122.112 {
// 		fmt.Println(mv)
// 		t.Error("error")
// 		return
// 	}
// 	if mv.IntValP("/data/i") != 1 {
// 		fmt.Println(mv)
// 		t.Error("error")
// 		return
// 	}
// 	if mv.StrValP("/data/s") != "sval" {
// 		fmt.Println(mv)
// 		t.Error("error")
// 		return
// 	}
// 	//
// 	mv, err = ts.G2("/abc/f?exec=del")
// 	if mv.IntVal("code") != 0 {
// 		fmt.Println(mv)
// 		t.Error("error")
// 		return
// 	}
// 	mv, err = ts.G2("/abc/f?exec=del")
// 	if mv.IntVal("code") == 0 {
// 		fmt.Println(mv)
// 		t.Error("error")
// 		return
// 	}
// 	mv, err = ts.G2("/abc/*")
// 	if mv.IntVal("code") != 0 {
// 		fmt.Println(mv)
// 		t.Error("error")
// 		return
// 	}
// 	if mv.Exist("/data/f") {
// 		fmt.Println(mv)
// 		t.Error("error")
// 		return
// 	}
// 	//
// 	//test error
// 	fmt.Println(ts.G2("/xxx/a?exec=ss"))
// 	fmt.Println(ts.G2("/xxx/a?exec=set&type=Jdata=abc"))
// 	fmt.Println(ts.G2("/xxx/a?exec=set&type=I&data=abc"))
// 	fmt.Println(ts.G2("/xxx/a?exec=set&type=F&data=abc"))
// 	//
// 	fmt.Println(ts.G2("/xxx/a?exec=plus&type=S&data=abc"))
// 	ts.G2("/xxx/s?exec=set&type=S&data=abc")
// 	ts.G2("/xxx/s?exec=plus&type=I&data=abc")
// 	ts.G2("/xxx/s?exec=plus&type=F&data=abc")
// 	ts.G2("/xxx/s?exec=plus&type=I&data=1")
// 	ts.G2("/xxx/s?exec=plus&type=F&data=1")
// }
