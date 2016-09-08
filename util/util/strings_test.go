package util

import (
	"fmt"
	"strings"
	"testing"
)

func TestSplit(t *testing.T) {
	vals := strings.Split("", ",")
	if len(vals) != 1 {
		t.Error("error")
		return
	}
	vals = Split("", ",")
	if len(vals) != 0 {
		t.Error("error")
		return
	}
	vals = Split("1", ",")
	if len(vals) != 1 {
		t.Error("error")
		return
	}
}

func TestTrimStrs(t *testing.T) {
	vals := TrimStrs([]string{"xa", " ", "xb", "\t", "xc", "", "xa"}, " \t")
	if len(vals) != 4 {
		t.Error("error")
		return
	}
	if vals[0] != "xa" || vals[1] != "xb" || vals[2] != "xc" || vals[3] != "xa" {
		t.Error("error")
		return
	}
	vals = TrimStrsRepeat([]string{"xa", " ", "xb", "\t", "xc", "xb", "", "xa"}, " \t", true)
	if len(vals) != 3 {
		t.Error("error")
		return
	}
	if vals[0] != "xa" || vals[1] != "xb" || vals[2] != "xc" {
		t.Error("error")
		return
	}
	if Trim("x ") != "x" {
		t.Error("error")
		return
	}
	fmt.Println("done...")
}
