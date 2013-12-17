package util

import (
	"fmt"
	"log"
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
	log.New(nil, "", 0).Println("ss")
}

func TestFTouch(t *testing.T) {
	os.RemoveAll("/tmp/kkk")
	os.RemoveAll("/tmp/abc.log")
	fmt.Println(FTouch("/tmp/abc.log"))
	fmt.Println(FTouch("/tmp/kkk/abc.log"))
	fmt.Println(FTouch("/tmp/kkk/abc.log"))
	fmt.Println(FTouch("/tmp/kkk"))
}
