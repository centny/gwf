package util

import (
	"bufio"
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
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
		for i := 0; i < 5; i++ {
			if i == 3 {
				panic("ss")
			}
			w.Write([]byte("kkkkkkkk"))
			fmt.Println("writing ...")
			time.Sleep(1000 * time.Millisecond)
			fmt.Println(reflect.TypeOf(w))
		}
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
