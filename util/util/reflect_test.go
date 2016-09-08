package util

import (
	"fmt"
	"strings"
	"testing"
)

func TestSliceExists(t *testing.T) {
	iary := []int{1, 2, 3, 4, 5, 6}
	if !SliceExists(iary, 2) {
		t.Error("value exis in array.")
		return
	}
	if SliceExists(iary, 8) {
		t.Error("value not exis in array.")
		return
	}
	//
	fary := []float32{1.0, 2.0, 3.0, 4.0, 5.0}
	if !SliceExists(fary, float32(1.0)) {
		t.Error("value exis in array.")
		return
	}
	if SliceExists(fary, float32(8.0)) {
		t.Error("value not exis in array.")
		return
	}
	//
	sary := []string{"a", "b", "c", "d", "e", "f"}
	if !SliceExists(sary, "c") {
		t.Error("value exis in array.")
		return
	}
	if SliceExists(sary, "g") {
		t.Error("value not exis in array.")
		return
	}
	ab := ""
	if SliceExists(ab, 8) {
		t.Error("value exis in array.")
		return
	}
	fmt.Println("test slice exists done...")
}

func TestIsType(t *testing.T) {
	if !IsType(t, "T") {
		t.Error("not right")
	}
	fmt.Println(IsType(nil, "A"))
	fmt.Println(IsType(t, ""))
	fmt.Println(IsType(t, " "))
	fmt.Println(IsType(t, "\t"))
}

func TestJoin(t *testing.T) {
	fmt.Println("...")
	if Join([]int{1, 2, 3}, ",") != "1,2,3" {
		t.Error("error")
		return
	}
	if Join([]float64{1.1, 2.2, 3.3}, ",") != "1.1,2.2,3.3" {
		t.Error("error")
		return
	}
	if Join([]string{"1", "2", "3"}, ",") != "1,2,3" {
		t.Error("error")
		return
	}
	if Join(nil, ",") != "" {
		t.Error("error")
		return
	}
	if Join("xx", ",") != "" {
		t.Error("error")
		return
	}
	if Join([]string{}, ",") != "" {
		t.Error("error")
		return
	}
}

func TestReflectName(t *testing.T) {
	var abc = []string{}
	var abc2 = []*testing.T{}
	var abc3 = []testing.T{}
	if StructName(t) != "testing.T" {
		t.Error("error")
		return
	}
	if StructName(TestReflectName) != "func(*testing.T)" {
		t.Error("error")
		return
	}
	if StructName(abc) != "[]string" {
		t.Error("error")
		return
	}
	if StructName(abc2) != "[]*testing.T" {
		t.Error("error")
		return
	}
	if StructName(abc3) != "[]testing.T" {
		t.Error("error")
		return
	}
	fmt.Println("...")
	if ReflectName(t) != "testing.T" {
		t.Error("error")
		return
	}
	if !strings.HasSuffix(ReflectName(TestReflectName), "util.TestReflectName") {
		t.Error("error")
		return
	}
	if !strings.HasSuffix(FuncName(TestReflectName), "util.TestReflectName") {
		t.Error("error")
		return
	}
	if ReflectName(abc) != "[]string" {
		t.Error("error")
		return
	}
	if ReflectName(abc2) != "[]*testing.T" {
		t.Error("error")
		return
	}
	if ReflectName(abc3) != "[]testing.T" {
		t.Error("error")
		return
	}
}

func TestCallStack(t *testing.T) {
	fmt.Println(CallStatck())
}
