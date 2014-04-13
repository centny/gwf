package igtest

import (
	"fmt"
	"github.com/Centny/Cny4go/util"
	"os"
	"testing"
)

type A struct {
	Iv string `m2s:"IV"`
}

func TestIg(t *testing.T) {
	b := NewBuilder2()
	b.StoreData("F64", float64(111))
	b.StoreData("I64", int64(111))
	av := map[string]interface{}{
		"IV": "kkkkkk",
	}
	b.StoreData("AV", av)
	//
	b.Append3(empty_f)
	b.Append2("F64", int64_f)
	b.Append2("I64", int64_f)
	b.Append2("F64", uint64_f)
	b.Append2("F64", string_f)
	b.Append2("AV,AV", av_f)
	b.Append2("F64,F64", t_f)
	b.Append2("@b,@t,F64,F64", b_f)
	if err := b.Exec(); err != nil {
		t.Error(err.Error())
	}
	//test error
	b.cls_task()
	b.Append2("AV", 111).Exec() //func error
	b.cls_task()
	b.Append2("AV", empty_f).Exec() //argument not equal
	b.cls_task()
	b.Append2("I6", int64_f).Exec()
	b.cls_task()
	b.Append2("I64", int_f).Exec()

	//end
	b.DeleteData("I64")
	b.cls_dset()

}
func empty_f() {
	fmt.Println("empty")
}
func int_f(v int) {
	fmt.Println(v)
}
func int64_f(v int64) {
	fmt.Println(v)
}
func uint64_f(v uint64) {
	fmt.Println(v)
}
func string_f(v string) {
	fmt.Println(v)
}
func av_f(v1 *A, v2 A) {
	fmt.Println(v1, v2)
}
func t_f(v1 int64, v2 string) {
	fmt.Println(v1, v2)
}
func b_f(b *Builder, t Task, v1 int64, v2 string) {
	fmt.Println(b, t, v1, v2)
}

func TestIgE(t *testing.T) {
	inf := "/tmp/kkk.json"
	outf := "/tmp/aaa.json"
	os.Remove(inf)
	os.Remove(outf)
	b := NewBuilder(inf, outf)
	//
	fmt.Println(b.Exec())
	//
	util.FWrite(inf, "kkk")
	fmt.Println(b.Exec())
	os.Remove(inf)
	//
	//
	util.FWrite(inf, "{\"abc\":111}")
	fmt.Println(b.Exec())
	os.Remove(outf)
	//
	//
	b.Out = "/lkkkk"
	fmt.Println(b.Exec())
	//
	b.In = ""
	b.Dset = nil
	fmt.Println(b.Exec())
}

func TestIgE2(t *testing.T) {
	os.Args = append(os.Args, "@ig-i=/tmp/kkk.json")
	os.Args = append(os.Args, "@ig-o=/tmp/aaa.json")
	b := NewBuilder2()
	if b.In != "/tmp/kkk.json" || b.Out != "/tmp/aaa.json" {
		t.Error("not match")
	}
}

func TestArg(t *testing.T) {
	fmt.Println(os.Args)
}
