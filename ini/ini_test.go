package ini

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"testing"
)

type Abc struct {
}

func (a *Abc) Read(p []byte) (n int, err error) {
	return 0, errors.New("abc")
}
func TestIni(t *testing.T) {
	az, err := ParseIni("t.ini")
	if err != nil {
		t.Error(err.Error())
		return
	}
	az.Add("abc", "a", "1", "1.1")
	az.Add("abc", "a", "2", "2.2")
	az.Add("abc", "a", "3", "3.3")
	az.Add("abc", "b", "3", "3.3")
	az.Add("ab3", "a", "3", "3.3")
	az.Del("abc", "sfsdfsd", "1")
	az.Del("asfsdfsbc", "sfsdfsd", "1")
	az.Store_(bufio.NewWriter(os.Stdout))
	fmt.Println("end...")
	fmt.Println(ParseIni("terr.ini"))
	NewIni().Parse_(bufio.NewReader(&Abc{}))
}

func TestCmds(t *testing.T) {
	fmt.Println(Cmds([]string{"t.ini", "add", "a", "1", "2"}))
	fmt.Println(Cmds([]string{"/sd/t.ini", "add", "a", "1", "2"}))
	fmt.Println(Cmds([]string{"t.ini", "del", "a", "1", "2"}))
	fmt.Println(Cmds([]string{"/sd/t.ini", "del", "a", "1", "2"}))
	fmt.Println(Cmds(nil))
}
