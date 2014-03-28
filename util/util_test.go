package util

import (
	"bufio"
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"syscall"
	"testing"
	"time"
)

func TestFexist(t *testing.T) {
	fmt.Println(Fexists("/usr/local"))
	fmt.Println(Fexists("/usr/locals"))
	fmt.Println(Fexists("/usr/local/s"))
}

func TestFile(t *testing.T) {
	fmt.Println(os.Open("/tmp/kkgg"))
}

func TestFTouch(t *testing.T) {
	os.RemoveAll("/tmp/kkk")
	os.RemoveAll("/tmp/abc.log")
	fmt.Println(FTouch("/tmp/abc.log"))
	fmt.Println(FTouch("/tmp/kkk/abc.log"))
	fmt.Println(FTouch("/tmp/kkk/abc.log"))
	fmt.Println(FTouch("/tmp/kkk"))
	fmt.Println(FTouch("/var/libbb"))
	fmt.Println(Fexists(string([]byte{'/', 't', 'm', 'p', 0, '/', 'm', '/', 'a'})))
	fmt.Println(FTouch(string([]byte{'/', 't', 'm', 'p', 0, '/', 'm', '/', 'a'})))
	//
}
func TestBytePtrFromString(t *testing.T) {
	bys, err := syscall.BytePtrFromString(string([]byte{'/', 't', 'm', 'p', 0, '/', 'm'}))
	fmt.Println(bys, err)
	fmt.Println(os.MkdirAll(string([]byte{'/', 't', 'm', 'p', 0, '/', 'm'}), os.ModePerm))
}

// func TestFTouch2(t *testing.T) {
// 	fmt.Println(exec.Command("mkdir", "/tmp/fcg_dir").Run())
// 	fmt.Println(exec.Command("chmod", "000", "/tmp/fcg_dir").Run())
// 	fmt.Println(FTouch("/tmp/fcg_dir/aaa/a.log"))
// 	fmt.Println(exec.Command("rm", "-rf", "/tmp/fcg_dir").Run())
// }

func TestReadLine(t *testing.T) {
	f := func(end bool) {
		bf := bytes.NewBufferString("abc\ndef\nghi\n")
		r := bufio.NewReader(bf)
		for {
			bys, err := ReadLine(r, 10000, end)
			// bys, isp, err := r.ReadLine()
			fmt.Println(string(bys), err)
			if err != nil {
				break
			}
		}
	}
	f(true)
	f(false)
}

func TestTimestamp(t *testing.T) {
	tt := Timestamp(time.Now())
	bt := Time(tt)
	t2 := Timestamp(bt)
	fmt.Println(1392636938688)
	fmt.Println(tt)
	fmt.Println(t2)
	if tt != t2 {
		t.Error("convert invalid")
		return
	}
}

func TestAryExist(t *testing.T) {
	iary := []int{1, 2, 3, 4, 5, 6}
	if !AryExist(iary, 2) {
		t.Error("value exis in array.")
		return
	}
	if AryExist(iary, 8) {
		t.Error("value not exis in array.")
		return
	}
	//
	fary := []float32{1.0, 2.0, 3.0, 4.0, 5.0}
	if !AryExist(fary, float32(1.0)) {
		t.Error("value exis in array.")
		return
	}
	if AryExist(fary, float32(8.0)) {
		t.Error("value not exis in array.")
		return
	}
	//
	sary := []string{"a", "b", "c", "d", "e", "f"}
	if !AryExist(sary, "c") {
		t.Error("value exis in array.")
		return
	}
	if AryExist(sary, "g") {
		t.Error("value not exis in array.")
		return
	}
	ab := ""
	if AryExist(ab, 8) {
		t.Error("value exis in array.")
		return
	}
}

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
func TestMap2Query(t *testing.T) {
	mv := map[string]interface{}{}
	mv["abc"] = "123"
	mv["dd"] = "ee"
	fmt.Println(Map2Query(mv))
}
