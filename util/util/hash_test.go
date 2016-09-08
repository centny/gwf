package util

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

func TestCrc32(t *testing.T) {
	fmt.Println(Crc32([]byte("fwnfiuchvwhrwiuv cs")))
	fmt.Println(Crc32([]byte("fwnfiuchvwhrwiuv cs.png")))
}

func TestMd5(t *testing.T) {
	var tf = "/tmp/a.txt"
	defer os.Remove(tf)
	var bys = []byte("abcxxx")
	ioutil.WriteFile(tf, bys, os.ModePerm)
	var hash, err = Md5(tf)
	if err != nil {
		t.Error(err)
		return
	}
	if Md5Byte(bys) != hash {
		t.Error("error")
		return
	}
	_, err = Md5("/sd/sd.txt")
	if err == nil {
		t.Error("error")
		return
	}
}

func TestSha1(t *testing.T) {
	var tf = "/tmp/a.txt"
	defer os.Remove(tf)
	var bys = []byte("abcxxx")
	ioutil.WriteFile(tf, bys, os.ModePerm)
	var hash, err = Sha1(tf)
	if err != nil {
		t.Error(err)
		return
	}
	if Sha1Byte(bys) != hash {
		t.Error("error")
		return
	}
	_, err = Sha1("/sd/sd.txt")
	if err == nil {
		t.Error("error")
		return
	}
}

func TestShortLink(t *testing.T) {
	fmt.Println(ShortLink("http://google.com"))
}
