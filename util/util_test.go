package util

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"reflect"
	"runtime"
	"strconv"
	"syscall"
	"testing"
	"time"
)

func init() {
	ShowLog = true
	go DoWeb(":65432", "./")
}
func TestDoWeb(t *testing.T) {
	go DoWeb(":65432", "./")
	time.Sleep(200 * time.Millisecond)
}

func TestBytePtrFromString(t *testing.T) {
	bys, err := syscall.BytePtrFromString(string([]byte{'/', 't', 'm', 'p', 0, '/', 'm'}))
	fmt.Println(bys, err)
	fmt.Println(os.MkdirAll(string([]byte{'/', 't', 'm', 'p', 0, '/', 'm'}), os.ModePerm))
}

func TestExec(t *testing.T) {
	fmt.Println(Exec("echo", "abc", "kk"))
}

func TestDLoad(t *testing.T) {
	DLoad("/tmp/index.html", "http/www.baidu.com")
	DLoad("/tmp/index.html", "http://www.baidu.com")
	os.Remove("/tmp/index.html")
	DLoad("/tmp/s.html", "")
}

func TestAppend(t *testing.T) {
	args := []interface{}{}
	args = Append(args, 1, nil)
	fmt.Println(args)
}

func TestList(t *testing.T) {
	ts := List("./", "^.*\\.txt$")
	if len(ts) != 1 {
		t.Error("error")
	}
	fmt.Println(ts)
}

func TestOs(t *testing.T) {
	fmt.Println(runtime.GOOS)
}

func TestHome(t *testing.T) {
	fmt.Println(os.Getenv("HOME"))
}

// func TestFsize(t *testing.T) {
// 	f, _ := os.Open("/Users/cny/Downloads/abc.mkv")
// 	fmt.Println(FormFSzie(f))
// 	f.Close()
// }

// func TestAAANet(t *testing.T) {
// 	c, err := net.Dial("tcp", "192.168.1.100:9100")
// 	if err != nil {
// 		t.Error(err.Error())
// 		return
// 	}
// 	w := bufio.NewWriter(c)

// }

type RW_ struct {
	r int
}

func (r *RW_) Read(p []byte) (n int, err error) {
	// fmt.Println("reading....")
	r.r += 5
	if r.r > 30 {
		return 0, Err("sss")
	} else if r.r > 20 {
		return 0, io.EOF
	} else if r.r == 10 {
		return 0, io.EOF
	}
	return 5, nil
}
func (r *RW_) Write(p []byte) (n int, err error) {
	if r.r > 10 {
		return 0, Err("ssff")
	}
	return len(p), nil
}

func TestCopy(t *testing.T) {
	rw := &RW_{}

	f, _ := os.Open("util.go")
	defer f.Close()
	fmt.Println(Copy2(rw, f))
	fmt.Println(Copy2(rw, rw))
	fmt.Println(Copy2(rw, rw))
	fmt.Println(Copy2(rw, rw))
	fmt.Println(Copy2(rw, rw))
	// fmt.Println("sss", base64.StdEncoding.EncodeToString(sha_))
	// fmt.Println("sss", base64.StdEncoding.EncodeToString(md5_))
	// fmt.Println(Copy(rw, rw))
	fmt.Println("--->")
}

func TestSS(t *testing.T) {
	fmt.Printf(Sha1("/Users/cny/Downloads/f.sql"))
}

func r_vv(n, o string, bv bool, t *testing.T) {
	iv, err := ChkVer(n, o)
	if err != nil {
		t.Error(err.Error())
		return
	}
	fmt.Println(iv)
	if bv {
		if iv < 1 {
			t.Errorf("errr val:%v", iv)
			return
		}
	} else {
		if iv > -1 {
			t.Errorf("errr val:%v", iv)
			return
		}
	}
}
func TestChkVer(t *testing.T) {
	r_vv("0.0.0", "", true, t)
	r_vv("0.0.0", "0.0", true, t)
	r_vv("0.0", "0.0.0", false, t)
	r_vv("1.0", "0.0.0", true, t)
	r_vv("1.0.0", "0.0.0", true, t)
	r_vv("0.0.0", "1.0.0", false, t)
	r_vv("1.0.1", "1.0.0", true, t)
	fmt.Println(ChkVer("n", "o"))
	fmt.Println(ChkVer("0.0.0", "o"))
	fmt.Println(ChkVer("ss", "0.0.0"))
	fmt.Println(ChkVer("", "0.0.0"))
}

func TestASecM(t *testing.T) {
	m := map[string]interface{}{
		"v": 123456789999,
	}
	fmt.Println(S2Json(m))
	mv, _ := Json2Map(S2Json(m))
	var vv int64
	fmt.Println(reflect.TypeOf(mv["v"]))
	fmt.Println(mv.IntVal("v"))
	fmt.Println(strconv.ParseFloat(fmt.Sprintf("%v", mv["v"]), 64))
	err := mv.ValidF(`
		v,R|I,R:0
		`, &vv)
	if err != nil {
		t.Error(err.Error())
		return
	}
	fmt.Println(vv)
}

func TestBuffer(t *testing.T) {
	pr, pw := io.Pipe()
	go func() {
		buf := make([]byte, 1024)
		for {
			bl, err := pr.Read(buf)
			if err != nil {
				fmt.Println("Err->", err.Error())
				break
			}
			fmt.Println("R->", string(buf[0:bl]))
		}
		fmt.Println("R->done...")
	}()
	go func() {
		for i := 0; i < 5; i++ {
			pw.Write([]byte(fmt.Sprintf("d-%v", i)))
		}
		fmt.Println("W->done...")
	}()
	time.Sleep(1 * time.Second)
}

func TestCopyp(t *testing.T) {
	FWrite("/tmp/kskds.txt", "data")
	f, _ := os.Open("/tmp/kskds.txt")
	Copyp("/tmp/xxkj.txt", f)
	f.Close()
	f, _ = os.Open("/tmp/kskds.txt")
	Copyp2("/tmp/xxkj2.txt", f)
	f.Close()
	f, _ = os.Open("/tmp/kskds.txt")
	Copyp2("/tmp/x/xkj2.txt", f)
	f.Close()
}

type xxk struct {
	Kk int64 `json:"kk"`
}

func TestJsonInt(t *testing.T) {
	var xk xxk
	xk.Kk = Now()
	bys, err := json.Marshal(&xk)
	if err != nil {
		t.Error(err.Error())
		return
	}
	fmt.Println(string(bys))
	err = json.Unmarshal(bys, &xk)
	if err != nil {
		t.Error(err.Error())
		return
	}
	fmt.Println(xk.Kk)
}

func TestBysSize(t *testing.T) {
	fmt.Println(BysSize(10))
	fmt.Println(BysSize(1024))
	fmt.Println(BysSize(1034))
	fmt.Println(BysSize(1024 * 1024))
	fmt.Println(BysSize(1024*1024 + 1))
	fmt.Println(BysSize(1024 * 1024 * 1024))
	fmt.Println(BysSize(1024*1024*1024 + 1))
	fmt.Println(BysSize(1024*1024*1024*1024 + 1))
	fmt.Println(BysSize(1024 * 1024 * 1024 * 1024 * 1024))
	fmt.Println(BysSize(1024*1024*1024*1024*1024 + 1))
}

func TestATime(t *testing.T) {
	sv := "2016-08-10 22:48:47"
	xx, _ := time.ParseInLocation("2006-01-02 15:04:05", sv, time.Local)
	if xx.Format("2006-01-02 15:04:05") != sv {
		t.Error("error")
		return
	}
	if Time(1470840527000).Format("2006-01-02 15:04:05") != sv {
		t.Error("error")
		return
	}
	if Time(Timestamp(xx)).Format("2006-01-02 15:04:05") != sv {
		t.Error("error")
		return
	}
	if Timestamp(xx) != 1470840527000 {
		t.Error("error")
		return
	}
	// fmt.Println(util.Timestamp(xx))
	// fmt.Println(util.Time(1470840527000).Format("2006-01-02 15:04:05"))
	// fmt.Println(util.Timestamp(xx), 1470869327000, 1470840527000)
}
