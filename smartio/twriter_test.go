package smartio

import (
	"bufio"
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
	f, err := os.OpenFile(path, os.O_RDWR|os.O_APPEND, os.ModePerm)
	fmt.Println(f, err)
	tw := NewTWriter(f)
	tw.WriteString("sfsdfsdfsdfsdfsdfs\n")
	tw.WriteString("sfsdfsdfsdfsdfsdfs\n")
	tw.WriteString("sfsdfsdfsdfsdfsdfs\n")
	tw.WriteString("sfsdfsdfsdfsdfsdfs\n")
	tw.WriteString("sfsdfsdfsdfsdfsdfs\n")
	time.Sleep(1000 * time.Millisecond)
	tw.Stop()
	TimeWriterWait()
	fmt.Println("exit")
}
func TestNewTimeWriter2(t *testing.T) {
	TimeWriterWait()
}
func TestWriteFile(t *testing.T) {
	path := "/tmp/kka/a.log"
	// util.FTouch(path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_APPEND, os.ModePerm)
	fmt.Println(f, err)
	w := bufio.NewWriterSize(f, 1000)
	w.WriteString("9999999999999")
	w.Flush()
	f.Close()
}
func TestMemoryWriter(t *testing.T) {

}

// func TestFile(t *testing.T) {
// f, err := os.Create("/tmp/tt")
// }

// func TestMkdir(t *testing.T) {
// os.MkdirAll("/tmp/lll", os.ModePerm)
// }
