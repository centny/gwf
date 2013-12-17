package twriter

import (
	"fmt"
	"github.com/Centny/Cny4go/util"
	"os"
	"testing"
	"time"
)

func TestNewTimeWriter(t *testing.T) {
	os.RemoveAll("/tmp/kk")
	path := "/tmp/kka/a.log"
	util.FTouch(path)
	f, err := os.Open(path)
	fmt.Println(f, err)
	tw := NewTWriter(f)
	tw.Writer().WriteString("sfsdfsdfsdfsdfsdfs")
	time.Sleep(1000 * time.Millisecond)
	tw.Stop()
	Wait()
	fmt.Println("exit")
}

// func TestFile(t *testing.T) {
// f, err := os.Create("/tmp/tt")
// }

// func TestMkdir(t *testing.T) {
// os.MkdirAll("/tmp/lll", os.ModePerm)
// }
