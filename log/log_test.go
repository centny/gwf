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

func TestSetSLevel(t *testing.T) {
	SetSLevel("DEBUG")
	if dlev != DEBUG {
		t.Error("level is not DEBUG")
	}
	SetSLevel("INFO")
	if dlev != INFO {
		t.Error("level is not INFO")
	}
	SetSLevel("WARNING")
	if dlev != WARNING {
		t.Error("level is not WARNING")
	}
	SetSLevel("ERROR")
	if dlev != ERROR {
		t.Error("level is not ERROR")
	}
	SetSLevel("ERRORR")
	if dlev != INFO {
		t.Error("level is not INFO")
	}
	fmt.Println("test set level end...")
}

func TestNewLog(t *testing.T) {
	l := NewLog(os.Stderr, "NL ", D_LOG_FLAGS)
	l.D("show debug")
}

func TestSetWriter(t *testing.T) {
	SetWriterFlag(os.Stderr, "LL ", D_LOG_FLAGS)
	D("show debug")
	SetWriter(os.Stdout)
	D("show debug 2")
}

type MWriter struct {
}

func (m *MWriter) Write(p []byte) (n int, err error) {
	fmt.Println("mw:" + string(p))
	return len(p), nil
}

func TestCWriter(t *testing.T) {
	SetWriterFlag(&MWriter{}, "LL ", D_LOG_FLAGS)
	D("show debug")
	SetWriter(os.Stdout)
	D("show debug 2")
}
