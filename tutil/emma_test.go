package tutil

import (
	"os"
	"testing"
)

func TestEmma(t *testing.T) {
	os.Remove("t.xml")
	err := Emma("t.xml", "abc", "1/10", "1/10", "1/10", "29/100")
	if err != nil {
		t.Error(err.Error())
		return
	}
	err = Emma("t.xml", "abc", "1/10", "1/10", "1/10", "29/100")
	if err != nil {
		t.Error(err.Error())
		return
	}
	err = Emma("t.xml", "abc", "1/10", "1/10", "1/10", "0/0")
	if err != nil {
		t.Error(err.Error())
		return
	}
	Emma("t.xml", "abc", "xxs/ss", "1/10", "1/10", "29/100")
	Emma("t.xml", "abc", "110", "1/10", "1/10", "29/100")
	Emma("t.xml", "abc", "1/10", "110", "1/10", "29/100")
	Emma("t.xml", "abc", "1/10", "1/10", "110", "29/100")
	Emma("t.xml", "abc", "1/10", "1/10", "1/10", "29100")
	cov_val_E("s)")
	cov_val_E("(s)")
	Emma("tutil.go", "abc", "1/10", "1/10", "1/10", "0/0")
}
