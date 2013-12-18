package util

import (
	"fmt"
	"testing"
)

func TestEnvReplace(t *testing.T) {
	f := &Fcfg{}
	f.SetVal("a", "b111111")
	fmt.Println(f.EnvReplace("sss${a} ${abc} ${da} ${HOME}"))
}

func TestInit(t *testing.T) {
	f := &Fcfg{}
	err := f.InitWithFilePath("src/org.cst.tsim/common/fcfg_data.properties")
	if err != nil {
		fmt.Println("error:", err)
	}
	for key, val := range *f {
		fmt.Println(key, ":", val)
	}
}
func TestValType(t *testing.T) {
	f := &Fcfg{}
	err := f.InitWithFilePath("src/org.cst.tsim/common/fcfg_data.properties")
	if err != nil {
		fmt.Println("error:", err)
	}
	fmt.Println(f.FloatVal("floata"))
	fmt.Println(f.FloatVal("floatb"))
	fmt.Println(f.FloatVal("inta"))
}
