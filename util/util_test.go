package util

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"testing"
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
}

func TestReadLine(t *testing.T) {
	bf := bytes.NewBufferString("abc\ndef\nghi\n")
	r := bufio.NewReader(bf)
	for {
		bys, err := ReadLine(r, 10000, false)
		// bys, isp, err := r.ReadLine()
		fmt.Println(string(bys), err)
		if err != nil {
			break
		}
	}
}
