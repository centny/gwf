package util

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

func TestHTTPGet(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// for i := 0; i < 5; i++ {
		// 	if i == 3 {
		// 		panic("ss")
		// 	}
		// 	w.Write([]byte("kkkkkkkk"))
		// 	fmt.Println("writing ...")
		// 	time.Sleep(1000 * time.Millisecond)
		// 	fmt.Println(reflect.TypeOf(w))
		// }
	}))
	ts2 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("kkkkkkkk"))
	}))
	res := HTTPGet("kkkk")
	if len(res) > 0 {
		t.Error("testing data not empty")
		return
	}
	res = HTTPGet(ts2.URL)
	if len(res) < 1 {
		t.Error("testing data is empty")
		return
	}
	go func() {
		res = HTTPGet(ts.URL)
		if len(res) > 0 {
			t.Error("testing data not empty", res)
			return
		}
		fmt.Println("..........", res)
	}()
	time.Sleep(5000 * time.Millisecond)
}

func TestReadAllStr(t *testing.T) {
	res, _ := readAllStr(nil)
	if len(res) > 0 {
		t.Error("not empty")
		return
	}
	r, _ := os.Open("name")
	res, _ = readAllStr(r)
	if len(res) > 0 {
		t.Error("not empty")
		return
	}
}

func TestHTTP2(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("{\"code\":1}"))
	}))
	ts2 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("{\"code:1}"))
	}))
	ts3 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	}))
	res := HTTPGet2(ts.URL)
	fmt.Println(res)
	res = HTTPGet2(ts2.URL)
	fmt.Println(res)
	res = HTTPGet2(ts3.URL)
	fmt.Println(res)
	_, err := HPostF(ts.URL, map[string]string{"ma": "123"}, "abc", "")
	if err == nil {
		t.Error("not error")
		return
	}
	_, err = HPostF("hhh", map[string]string{"ma": "123"}, "abc", "test.txt")
	if err == nil {
		t.Error("not error")
		return
	}
	_, err = HPostF(ts.URL, map[string]string{"ma": "123"}, "abc", "test.txt")
	if err != nil {
		t.Error(err.Error())
		return
	}
	_, err = HPostF(ts.URL, map[string]string{"ma": "123"}, "abc", "/tmp")
	if err == nil {
		t.Error("not error")
		return
	}
	_, err = HPostF2(ts.URL, map[string]string{"ma": "123"}, "abc", "test.txt")
	if err != nil {
		t.Error(err.Error())
		return
	}
	HPostF2("kkk", map[string]string{"ma": "123"}, "abc", "test.txt")
	HTTPPost(ts.URL, map[string]string{"ma": "123"})
	HTTPPost2(ts.URL, map[string]string{"ma": "123"})
	HTTPPost2("jhj", map[string]string{"ma": "123"})
}

//
type osize struct {
}

func (o *osize) Size() int64 {
	return 100
}

type ostat struct {
	F *os.File
}

func (o *ostat) Stat() (os.FileInfo, error) {
	return o.F.Stat()
}
func TestFormFSzie(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("{\"code\":1}"))
		src, _, err := r.FormFile("abc")
		if err != nil {
			t.Error(err.Error())
			return
		}
		fsize := FormFSzie(src)
		if fsize < 1 {
			t.Error("not size")
		}
	}))
	_, err := HPostF(ts.URL, map[string]string{"ma": "123"}, "abc", "test.txt")
	if err != nil {
		t.Error(err.Error())
	}
	f, _ := os.Open("test.txt")
	defer f.Close()
	fsize := FormFSzie(f)
	if fsize < 1 {
		t.Error("not right")
	}
	fsize = FormFSzie(&osize{})
	if fsize < 1 {
		t.Error("not right")
	}
}
func TestHpp(t *testing.T) {
	HGet("kk")
	HGet2("kk")
	HPost("jjjj", nil)
	HPost2("kkk", nil)
	HGet2("kkk")
}
func TestMap2Query(t *testing.T) {
	mv := map[string]interface{}{}
	mv["abc"] = "123"
	mv["dd"] = "ee"
	fmt.Println(Map2Query(mv))
}

func TestAHttpPost(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			r.ParseMultipartForm(10000000)
			// r.PostFormValue(key)
			fmt.Println(r.PostFormValue("kkk"))
		}))
	HPostF2s(ts.URL, map[string]string{
		"ab": "233",
	}, "", "")
}

func HPostF2s(url string, fields map[string]string, fkey string, fp string) (string, error) {
	ctype, bodyBuf, err := CreateFormBody2(fields, fkey, fp)
	if err != nil {
		return "", err
	}
	res, err := http.Post(url, ctype, bodyBuf)
	if err != nil {
		return "", err
	}
	return readAllStr(res.Body)
}

func CreateFormBody2(fields map[string]string, fkey string, fp string) (string, *bytes.Buffer, error) {
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)
	for k, v := range fields {
		bodyWriter.WriteField(k, v)
	}
	w, _ := bodyWriter.CreateFormField("kkk")
	w.Write([]byte("kkkkkkk"))
	if len(fkey) > 0 {
		fileWriter, err := bodyWriter.CreateFormFile(fkey, fp)
		if err != nil {
			return "", nil, err
		}
		fh, err := os.Open(fp)
		if err != nil {
			return "", nil, err
		}
		defer fh.Close()
		_, err = io.Copy(fileWriter, fh)
		if err != nil {
			return "", nil, err
		}
	}
	ctype := bodyWriter.FormDataContentType()
	bodyWriter.Close()
	return ctype, bodyBuf, nil
}

func TestJson2Ary(t *testing.T) {
	ary, err := Json2Ary(`
		[1,2,"ss"]
		`)
	if err != nil {
		t.Error(err.Error())
		return
	}
	fmt.Println(ary)
	_, err = Json2Ary(`
		[1,2,ss"]
		`)
	if err == nil {
		t.Error("not error")
		return
	}
}

func TestPostN(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Ok"))
	}))
	_, data, err := HPostN(ts.URL, "text/plain", bytes.NewBuffer([]byte("WWW")))
	if err != nil {
		t.Error(err.Error())
		return
	}
	fmt.Println(data)
	fmt.Println(HPostN("kkk://sssss", "text/plain", bytes.NewBuffer([]byte("WWW"))))
	fmt.Println(HPostN("http:///kkkfjdfsfsd", "text/plain", bytes.NewBuffer([]byte("WWW"))))
}
