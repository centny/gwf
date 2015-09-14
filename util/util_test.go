package util

import (
	"bufio"
	"bytes"
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
func TestFexist(t *testing.T) {
	fmt.Println(Fexists("/usr/local"))
	fmt.Println(Fexists("/usr/locals"))
	fmt.Println(Fexists("/usr/local/s"))
}

func TestFile(t *testing.T) {
	fmt.Println(os.Open("/tmp/kkgg"))
}

func TestFWrite(t *testing.T) {
	err := FWrite("/tmp/test.txt", "data")
	if err != nil {
		t.Error(err.Error())
		return
	}
	err = FAppend("/tmp/test.txt", "data")
	if err != nil {
		t.Error(err.Error())
		return
	}
	err = FWrite("/kk/kkfd/d", "data")
	if err == nil {
		t.Error("not error")
	}
	err = FAppend("/kk/kkfd/d", "data")
	if err == nil {
		t.Error("not error")
	}
	err = FCopy("/tmp/test.txt", "/tmp/test2.txt")
	if err != nil {
		t.Error(err.Error())
		return
	}
	err = FCopy("/kk/kkfd/d", "data")
	if err == nil {
		t.Error("not error")
	}
	err = FCopy("/tmp/test.txt", "/dsss/dd.txt")
	if err == nil {
		t.Error("not error")
	}
	os.Remove("/tmp/test.txt")
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
func TestFTouch2(t *testing.T) {
	os.RemoveAll("/tmp/tkk")
	fmt.Println(FTouch2("/tmp/tkk", os.ModePerm))
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

func TestExec(t *testing.T) {
	fmt.Println(Exec("echo", "abc", "kk"))
}

func TestDLoad(t *testing.T) {
	DLoad("/tmp/index.html", "http/www.baidu.com")
	DLoad("/tmp/index.html", "http://www.baidu.com")
	os.Remove("/tmp/index.html")
	DLoad("/tmp/s.html", "")
}

func TestIsType(t *testing.T) {
	if !IsType(t, "T") {
		t.Error("not right")
	}
	fmt.Println(IsType(nil, "A"))
	fmt.Println(IsType(t, ""))
	fmt.Println(IsType(t, " "))
	fmt.Println(IsType(t, "\t"))
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

func TestFileProtocol(t *testing.T) {
	fmt.Println(FileProtocolPath("~"))
	fmt.Println(FileProtocolPath("sfdsf"))
	fmt.Println(FileProtocolPath("/sdfs/sfdsf"))
	fmt.Println(FileProtocolPath("C:\\s\\sdfs"))
	fmt.Println(FileProtocolPath("file://C:/s/sdfs"))
}
func TestHome(t *testing.T) {
	fmt.Println(os.Getenv("HOME"))
}

func TestStr2Int(t *testing.T) {
	fmt.Println(Str2Int("abc"))
	fmt.Println(Str2Int("11"))
	fmt.Println(Str2Int("11,22"))
	fmt.Println(Str2Int("11,"))
	fmt.Println(Str2Int("11,ssd"))
}

func TestIs2Ss(t *testing.T) {
	fmt.Println(Int2Str([]int64{1, 2}))
	fmt.Println(Is2Ss([]int64{1, 2}))
}

func TestReadW(t *testing.T) {
	r := bufio.NewReader(&Sw{})
	buf := make([]byte, 3)
	var las int64
	ReadW(r, buf, &las)
	fmt.Println(string(buf))
	fmt.Println(Now())
	ReadW(r, buf, &las)
}

type Sw struct {
	i int
}

func (s *Sw) Read(p []byte) (n int, err error) {
	if s.i < 1 {
		s.i = 1
		p[0] = 'A'
		return 1, nil
	} else if s.i < 2 {
		s.i = 2
		p[0] = 'B'
		p[1] = 'C'
		return 2, nil
	} else {
		return 0, Err("ssdsd")
	}
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

func TestCrc32(t *testing.T) {
	fmt.Println(Crc32([]byte("fwnfiuchvwhrwiuv cs")))
	fmt.Println(Crc32([]byte("fwnfiuchvwhrwiuv cs.png")))
}

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
