package io

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestFileSize(t *testing.T) {
	var val string
	val = FileSize(0).String()
	if val != "0B" {
		t.Error(val)
		return
	}
	val = FileSize(1).String()
	if val != "1B" {
		t.Error(val)
		return
	}
	val = FileSize(1024).String()
	if val != "1KB" {
		t.Error(val)
		return
	}
	val = FileSize(1025).String()
	if val != "1KB" {
		t.Error(val)
		return
	}
	val = FileSize(1024 * 1024).String()
	if val != "1MB" {
		t.Error(val)
		return
	}
	val = FileSize(1024 * 1024 * 1024).String()
	if val != "1GB" {
		t.Error(val)
		return
	}
	val = FileSize(1024 * 1024 * 1024 * 1024).String()
	if val != "1TB" {
		t.Error(val)
		return
	}
	val = FileSize(1024 * 1024 * 1024 * 1024 * 2).String()
	if val != "2TB" {
		t.Error(val)
		return
	}
	val = FileSize(1024 * 1024 * 1024 * 1024 * 1024).String()
	if val != "1024TB" {
		t.Error(val)
		return
	}
	fmt.Println("test file size done...")
}

func TestWalk(t *testing.T) {
	var t_call = func(v func(string) string) {
		var dirs = Walk("../", true, []string{".*util"}, []string{".*io"}, v)
		var util_exist = false
		for _, dir := range dirs {
			if strings.HasSuffix(dir, "io") {
				t.Error("io found")
				return
			} else if strings.HasSuffix(dir, "util") {
				util_exist = true
			}
		}
		if !util_exist {
			t.Error("util not exist")
			return
		}
	}
	t_call(nil)
	t_call(func(v string) string {
		return v
	})
	fmt.Println("test walk done...")
}

func TestFileExist(t *testing.T) {
	dir, exist := FileExists("../")
	if !exist || !dir {
		t.Error("error")
		return
	}
	dir, exist = FileExists("file.go")
	if !exist || dir {
		t.Error("error")
		return
	}
	dir, exist = FileExists("filexxx.go")
	if exist || dir {
		t.Error("error")
		return
	}
	fmt.Println("test file exist done...")
}

func TestTouchFile(t *testing.T) {
	defer os.RemoveAll("tmp")
	defer os.RemoveAll("tmp.dat")
	os.RemoveAll("tmp")
	os.RemoveAll("tmp.dat")
	test_touch := func(path string) {
		err := TouchFile(path)
		if err != nil {
			t.Errorf("touch file error->%v", err)
			return
		}
		if dir, exist := FileExists(path); dir || !exist {
			t.Error("file not exist")
		}
		err = TouchFile(path)
		if err != nil {
			t.Errorf("touch file error->%v", err)
			return
		}
	}
	test_touch(filepath.Join("tmp.dat"))           //file
	test_touch(filepath.Join("tmp", "t.dat"))      //folder file
	test_touch(filepath.Join("tmp", "a", "t.dat")) //multi folder file
	//
	//test error
	err := TouchFile("/usr/localx/t.dat")
	if err == nil {
		t.Error("error")
		return
	}
	err = TouchFile("/usr/local/t.dat")
	if err == nil {
		t.Error("error")
		return
	}
	err = TouchFile("/usr/local")
	if err == nil {
		t.Error("error")
		return
	}
	fmt.Println("test touch file done...")
}

func TestFileReadWrite(t *testing.T) {
	defer os.RemoveAll("tmp")
	defer os.RemoveAll("tmp.txt")
	os.RemoveAll("tmp")
	os.RemoveAll("tmp.txt")
	os.MkdirAll("tmp", os.ModePerm)
	//
	//test write
	err := WriteFileString("tmp/a.txt", "abcd")
	if err != nil {
		t.Error(err)
		return
	}
	err = WriteFileReader("tmp/b.txt", bytes.NewBufferString("abcd"))
	if err != nil {
		t.Error(err)
		return
	}
	err = AppendFileString("tmp/c.txt", "ab")
	if err != nil {
		t.Error(err)
		return
	}
	err = AppendFileString("tmp/c.txt", "c")
	if err != nil {
		t.Error(err)
		return
	}
	err = AppendFileReader("tmp/c.txt", bytes.NewBufferString("d"))
	if err != nil {
		t.Error(err)
		return
	}
	bys, err := ioutil.ReadFile("tmp/a.txt")
	if err != nil {
		t.Error(err)
		return
	}
	if string(bys) != "abcd" {
		t.Error("error")
		return
	}
	bys, err = ioutil.ReadFile("tmp/b.txt")
	if err != nil {
		t.Error(err)
		return
	}
	if string(bys) != "abcd" {
		t.Error("error")
		return
	}
	bys, err = ioutil.ReadFile("tmp/c.txt")
	if err != nil {
		t.Error(err)
		return
	}
	if string(bys) != "abcd" {
		t.Error("error")
		return
	}
	//test copy and check read
	bys, err = CheckReadFile("tmp.txt", "file.go")
	if err != nil {
		t.Error(err)
		return
	}
	if string(bys) == "abcd" {
		t.Error("error")
		return
	}
	_, err = CopyFile("tmp/a.txt", "tmp.txt")
	if err != nil {
		t.Error(err)
		return
	}
	bys, err = CheckReadFile("tmp.txt", "file.go")
	if err != nil {
		t.Error(err)
		return
	}
	if string(bys) != "abcd" {
		t.Error("error")
		return
	}
	//
	//test error
	err = WriteFileString("xx/a.tx", "data")
	if err == nil {
		t.Error("error")
		return
	}
	err = WriteFileReader("xx/a.tx", bytes.NewBufferString("data"))
	if err == nil {
		t.Error("error")
		return
	}
	err = AppendFileString("xx/a.tx", "data")
	if err == nil {
		t.Error("error")
		return
	}
	err = AppendFileReader("xx/a.tx", bytes.NewBufferString("data"))
	if err == nil {
		t.Error("error")
		return
	}
	_, err = CopyFile("tmp/a.txt", "xx/b.txt")
	if err == nil {
		t.Error("error")
		return
	}
	_, err = CopyFile("xx/c.txt", "xx/b.txt")
	if err == nil {
		t.Error("error")
		return
	}
	_, err = CheckReadFile("xx/a.txt", "xx/b.txt")
	if err == nil {
		t.Error("error")
		return
	}
	fmt.Println("test read/write file done...")
}

func TestFileProtocol(t *testing.T) {
	fmt.Println(FileProtocolPath("~"))
	fmt.Println(FileProtocolPath("sfdsf"))
	fmt.Println(FileProtocolPath("/sdfs/sfdsf"))
	fmt.Println(FileProtocolPath("C:\\s\\sdfs"))
	fmt.Println(FileProtocolPath("file://C:/s/sdfs"))
}
