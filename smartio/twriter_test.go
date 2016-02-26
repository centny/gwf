package smartio

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/Centny/gwf/util"
	"os"
	"sync"
	"testing"
	"time"
)

func TestTwAutoFlush(t *testing.T) {
	b := bytes.NewBuffer([]byte{})
	tw := NewTimeWriter(b, 100000, 1000)
	tw.WriteString("12345\n")
	//wait for auto flush.
	time.Sleep(1500 * time.Millisecond)
	if b.Len() != 6 {
		t.Error("Auto flush error")
	}
	tw.Stop()
	TimeWriterWait()
	fmt.Println("Test auto flush success")
}

func TestTwBuffer(t *testing.T) {
	b := bytes.NewBuffer([]byte{})
	tw := NewTimeWriter(b, 100, 10000)
	for i := 0; i < 11; i++ {
		tw.WriteString("123456789\n")
	}
	//wait for auto buffer flush.
	time.Sleep(300 * time.Millisecond)
	if b.Len() != 100 {
		t.Error("Auto buffer flush error")
	}
	tw.Stop()
	if b.Len() != 110 {
		t.Error("Auto buffer flush error")
	}
	TimeWriterWait()
	fmt.Println("Test auto buffer flush success")
}

func TestDwNormal(t *testing.T) {
	dw := NewDateSwitchWriter2("/tmp")
	if dw.FilePath() != "" {
		t.Error("file path error")
	}
	dw.Write([]byte{'1', '1', '1', '\n'})
	dw.Write([]byte{'1', '1', '1', '\n'})
	fmt.Println(dw.FilePath())
	dw.cfn = "lll.log"
	dw.Write([]byte{'1', '1', '1', '\n'})
	dw.Write([]byte{'1', '1', '1', '\n'})
	dw.F.Close()
	dw.Write([]byte{'1', '1', '1', '\n'})
	dw.Close()
	//
	dw = NewDateSwitchWriter2(string([]byte{'/', 't', 'm', 'p', 0, '/', 'm', '/', 'a'}))
	dw.Write([]byte{'1', '1', '1', '\n'})
	NewDateSwitchWriter2("/tmp").Close()
}
func TestDwNormal2(t *testing.T) {
	dw := NewDateSwitchWriter2("/tmp")
	for i := 0; i < 1000; i++ {
		var ij = i
		go func() {
			_, err := dw.Write([]byte(fmt.Sprintf("ksjfksdfjksdfjskfjskfsfs:%v\n", ij)))
			if err != nil {
				t.Error(err.Error())
			}
		}()
	}
	time.Sleep(time.Second)
	dw.Close()
}
func TestDwNormal3(t *testing.T) {
	dw := NewDateSwitchWriter2("/tmp")
	dw.Write([]byte("jjjsfs\n"))
	dw.F.Close()
	dw.Write([]byte("jjjsfs\n"))
	dw.Close()
}
func TestNtw(t *testing.T) {
	fw := NewTWriter(os.Stderr)
	fw.WriteString("loging \n")
	fw.WriteString("loging \n")
	fw.buf.Flush()
	fw.Stop()
	fmt.Println("test new TWriter end ...")
}

func TestDTW(t *testing.T) {
	dw := NewDateSwitchWriter2("/tmp/kkjj/kk")
	tw := NewTimeWriter(dw, 1024, 100)
	wg := sync.WaitGroup{}
	for i := 0; i < 1000; i++ {
		var ij = i
		go func() {
			wg.Add(1)
			_, err := tw.Write([]byte(fmt.Sprintf("ksjfksdfjksdfjskfjskfsfs:%v\n", ij)))
			if err != nil {
				t.Error(err.Error())
			}
			wg.Done()
		}()
	}
	time.Sleep(time.Second)
	wg.Wait()
	dw.Close()
}

type ErrWriter struct {
}

func (e *ErrWriter) Write(p []byte) (int, error) {
	return 0, errors.New("error")
}
func TestDTW2(t *testing.T) {
	tw := NewTimeWriter(&ErrWriter{}, 1024, 100)
	wg := sync.WaitGroup{}
	tw.Write([]byte("ksjfksdfjksdfjskfjskfsfs:%v\n"))
	time.Sleep(time.Second)
	wg.Wait()
}

func TestRedirect(t *testing.T) {
	RedirectStdout3("test/out%v.log")
	RedirectStderr3("test/err%v.log")
	fmt.Printf("%v", "abc")
	fmt.Fprintf(os.Stderr, "%v", "abc")
	time.Sleep(4 * time.Second)
	if !util.Fexists(fmt.Sprintf("out%v.log", time.Now().Format("2006-01-02"))) {
		t.Error("error")
		return
	}
	if !util.Fexists(fmt.Sprintf("err%v.log", time.Now().Format("2006-01-02"))) {
		t.Error("error")
		return
	}
	if RedirectStdout3("/xsss/sdfss%v") == nil {
		t.Error("error")
		return
	}
}
