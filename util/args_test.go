package util

import (
	"fmt"
	"os"
	"testing"
)

func TestArgs(t *testing.T) {
	targs := os.Args
	//case 1
	os.Args = []string{"abc"}
	name, arg1, arg2 := Args()
	if name != "abc" || len(arg1) > 0 || len(arg2) > 0 {
		t.Error("error")
		return
	}
	//case 2
	os.Args = []string{"abc", "a"}
	name, arg1, arg2 = Args()
	if name != "abc" || len(arg1) > 0 || len(arg2) != 1 || arg2[0] != "a" {
		t.Error("error")
		return
	}
	//case 3
	os.Args = []string{"abc", "a", "b"}
	name, arg1, arg2 = Args()
	if name != "abc" || len(arg1) > 0 || len(arg2) != 2 || arg2[0] != "a" || arg2[1] != "b" {
		t.Error("error")
		return
	}
	//case 3-1
	os.Args = []string{"abc", "a", "b", "c"}
	name, arg1, arg2 = Args()
	if name != "abc" || len(arg1) > 0 || len(arg2) != 3 || arg2[0] != "a" || arg2[1] != "b" || arg2[2] != "c" {
		t.Error("error")
		return
	}
	//case 4
	os.Args = []string{"abc", "-a", "b"}
	name, arg1, arg2 = Args()
	if name != "abc" || len(arg1) != 1 || len(arg2) != 0 || arg1["a"].([]string)[0] != "b" {
		fmt.Println(S2Json(arg1))
		t.Error("error")
		return
	}
	//case 5
	os.Args = []string{"abc", "-a", "b", "c"}
	name, arg1, arg2 = Args()
	if name != "abc" || len(arg1) != 1 || len(arg2) != 1 || arg1["a"].([]string)[0] != "b" || arg2[0] != "c" {
		t.Error("error")
		return
	}
	//case 6
	os.Args = []string{"abc", "-a", "b", "-x"}
	name, arg1, arg2 = Args()
	if name != "abc" || len(arg1) != 2 || len(arg2) != 0 || arg1["a"].([]string)[0] != "b" || arg1["x"] == nil {
		t.Error("error")
		return
	}
	//case 7
	os.Args = []string{"abc", "-a", "b", "c"}
	name, arg1, arg2 = Args()
	if name != "abc" || len(arg1) != 1 || len(arg2) != 1 || arg1["a"].([]string)[0] != "b" || arg2[0] != "c" {
		t.Error("error")
		return
	}
	//case 8
	os.Args = []string{"abc", "-a", "b", "c", "-x"}
	name, arg1, arg2 = Args()
	if name != "abc" || len(arg1) != 2 || len(arg2) != 1 || arg1["a"].([]string)[0] != "b" || arg1["x"] == nil || arg2[0] != "c" {
		t.Error("error")
		return
	}
	//case 9
	os.Args = []string{"abc", "-a", "b", "c", "-x", "-l"}
	name, arg1, arg2 = Args()
	if name != "abc" || len(arg1) != 3 || len(arg2) != 1 || arg1["a"].([]string)[0] != "b" || arg1["x"] == nil || arg1["l"] == nil || arg2[0] != "c" {
		t.Error("error")
		return
	}

	//
	os.Args = targs
	fmt.Println(os.Args)
}

func TestParseArgs(t *testing.T) {
	args := ParseArgs("")
	if len(args) != 0 {
		t.Error("error len")
		return
	}
	fmt.Println(args)
	//
	args = ParseArgs("a")
	if len(args) != 1 {
		t.Error("error len")
		return
	}
	fmt.Println(args)
	//
	args = ParseArgs("a b")
	if len(args) != 2 {
		t.Error("error len")
		return
	}
	fmt.Println(args)
	//
	args = ParseArgs("a b c")
	if len(args) != 3 {
		t.Error("error len")
		return
	}
	fmt.Println(args)
	//
	args = ParseArgs("a \"b   x\" c")
	if len(args) != 3 {
		t.Error("error len")
		return
	}
	fmt.Println(args)
	//
	args = ParseArgs("a 'b   x' c")
	if len(args) != 3 {
		t.Error("error len")
		return
	}
	fmt.Println(args)
	//
	args = ParseArgs("a 'b \"s\"  x' c")
	if len(args) != 3 {
		t.Error("error len")
		return
	}
	fmt.Println(args)
	//
	args = ParseArgs("a \"b  'xd' x\" c")
	if len(args) != 3 {
		t.Error("error len")
		return
	}
	fmt.Println(args)
	//
	args = ParseArgs("\"b  'xd' x\"")
	if len(args) != 1 {
		t.Error("error len")
		return
	}
	fmt.Println(args)
	//
	args = ParseArgs("'b   x'")
	if len(args) != 1 {
		t.Error("error len")
		return
	}
	fmt.Println(args)
	//
	args = ParseArgs("a\tb\tc")
	if len(args) != 3 {
		t.Error("error len")
		return
	}
	fmt.Println(args)
}
