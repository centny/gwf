package util

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"runtime"
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

func TestFsize(t *testing.T) {
	f, _ := os.Open("/Users/cny/Downloads/abc.mkv")
	fmt.Println(FormFSzie(f))
	f.Close()
}
