package util

import (
	"testing"
)

func TestZip(t *testing.T) {
	err := Zip("t.zip", "./", "./zip.go", "./util.go")
	if err != nil {
		t.Error(err.Error())
		return
	}
	err = Unzip("t.zip", "/tmp")
	if err != nil {
		t.Error(err.Error())
		return
	}
	Unzip("sdfsdf", "sfs")
	Unzip("util.go", "sfs")
	Zip("/sss/t.zip", "./", "./zip.go", "./util.go")
	Zip("t.zip", "./", "sfsdfs/", "/fsdfsko")
}
