package log

import (
	"fmt"
	"os"
	"testing"
)

func TestLog(t *testing.T) {
	l := NewDefaultLog(os.Stderr)
	l.D("DEBUG")
	l.I("INFO")
	l.W("WARNING")
	l.E("ERROR")
	l.Println("................")
	l.SetLevel(DEBUG)
	l.D("DEBUG")
	l.I("INFO")
	l.W("WARNING")
	l.E("ERROR")
	l.Println("................")
	l.SetLevel(INFO)
	l.D("DEBUG")
	l.I("INFO")
	l.W("WARNING")
	l.E("ERROR")
	l.Println("................")
	l.SetLevel(WARNING)
	l.D("DEBUG")
	l.I("INFO")
	l.W("WARNING")
	l.E("ERROR")
	l.Println("................")
	l.SetLevel(ERROR)
	l.D("DEBUG")
	l.I("INFO")
	l.W("WARNING")
	l.E("ERROR")
}
func TestDLog(t *testing.T) {
	D("DEBUG")
	I("INFO")
	W("WARNING")
	E("ERROR")
	fmt.Println("................")
	SetLevel(DEBUG)
	D("DEBUG")
	I("INFO")
	W("WARNING")
	E("ERROR")
	fmt.Println("................")
	SetLevel(INFO)
	D("DEBUG")
	I("INFO")
	W("WARNING")
	E("ERROR")
	fmt.Println("................")
	SetLevel(WARNING)
	D("DEBUG")
	I("INFO")
	W("WARNING")
	E("ERROR")
	fmt.Println("................")
	SetLevel(ERROR)
	D("DEBUG")
	I("INFO")
	W("WARNING")
	E("ERROR")
}

func TestFmt(t *testing.T) {
	SetLevel(DEBUG)
	D("%s---%d", "abc", 1000)
}
