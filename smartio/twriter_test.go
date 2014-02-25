package smartio

import (
	"bytes"
	"fmt"
	"os"
	"testing"
	"time"
)

func TestTwAutoFlush(t *testing.T) {
	SetFMode(0755)
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
	dw := NewDateSwitchWriter("/tmp")
	if dw.FilePath() != "" {
		t.Error("file path error")
	}
	dw.Write([]byte{'1', '1', '1', '\n'})
	dw.Write([]byte{'1', '1', '1', '\n'})
	fmt.Println(dw.FilePath())
	dw.cfn = "lll.log"
	dw.Write([]byte{'1', '1', '1', '\n'})
	dw.Write([]byte{'1', '1', '1', '\n'})
	dw.Close()
	//
	dw = NewDateSwitchWriter(string([]byte{'/', 't', 'm', 'p', 0, '/', 'm', '/', 'a'}))
	dw.Write([]byte{'1', '1', '1', '\n'})
	NewDateSwitchWriter("/tmp").Close()
}

func TestNtw(t *testing.T) {
	fw := NewTWriter(os.Stderr)
	fw.WriteString("loging \n")
	fw.WriteString("loging \n")
	fw.Flush()
	fw.Stop()
	fmt.Println("test new TWriter end ...")
}
