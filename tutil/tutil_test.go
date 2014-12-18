package tutil

import (
	"fmt"
	"github.com/Centny/gwf/util"
	"testing"
	"time"
)

func TestRun3(t *testing.T) {
	NewTSK("", nil).Run()
	NewTSK("sfdsfsdf", nil).Run()
	NewTSK("sfdsfsdf", nil).Conn(nil)
	NewTSk_C("2323")
	tsk_ := NewTSK2(":23423", func(tc *TSK_C, msg string) error {
		fmt.Println("---->", msg)
		tc.Write("MSG\n")
		return nil
	})
	go tsk_.Run()
	time.Sleep(200 * time.Millisecond)
	tc, _ := tsk_.Conn(func(tc *TSK_C, msg string) error {
		fmt.Println("---->", msg)
		return nil
	})
	tc2, _ := tsk_.Conn(func(tc *TSK_C, msg string) error {
		fmt.Println("---->", msg)
		return util.NOT_FOUND
	})
	tc.Write("II\n")
	tc2.Write("AA\n")
	time.Sleep(200 * time.Millisecond)
	tc.Close()
	time.Sleep(200 * time.Millisecond)
	tsk_.Stop()
	time.Sleep(time.Second)
}
