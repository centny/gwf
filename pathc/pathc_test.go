package main

import (
	"fmt"
	"os"
	"testing"
)

func TestMain(t *testing.T) {
	os.Args = []string{"pathc"}
	main()
	os.Args = []string{"pathc", "-w2p", "C:\\fdsf\\df"}
	main()
	os.Args = []string{"pathc", "-p2w", "/C/fdsf/df"}
	main()
}
func TestPathc(t *testing.T) {
	var tp string = ""
	var sp string = ""
	var cp string = ""
	sp = "C:\\fdsf\\df"
	cp = "/C/fdsf/df"
	tp = w2p(sp)
	if tp != cp {
		t.Error(tp)
		return
	}
	tp = p2w(cp)
	if tp != sp {
		t.Error(tp)
		return
	}
	sp = "c:\\fdsf\\df"
	cp = "/c/fdsf/df"
	tp = w2p(sp)
	if tp != cp {
		t.Error(tp)
		return
	}
	tp = p2w(cp)
	if tp != sp {
		t.Error(tp)
		return
	}
	sp = "c:\\fdsf"
	cp = "/c/fdsf"
	tp = w2p(sp)
	if tp != cp {
		t.Error(tp)
		return
	}
	tp = p2w(cp)
	if tp != sp {
		t.Error(tp)
		return
	}
	sp = "c:\\fdsf;D:\\sdfs"
	cp = "/c/fdsf:/D/sdfs"
	tp = w2p(sp)
	if tp != cp {
		t.Error(tp)
		return
	}
	tp = p2w(cp)
	if tp != sp {
		t.Error(tp)
		return
	}
	fmt.Println("all end")
}
