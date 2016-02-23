package util

import (
	"fmt"
	"testing"
)

func TestParseSectionF(t *testing.T) {
	//
	res := ParseSectionF("[", "]", `
[abc]
123
[/abc]
		`)
	fmt.Println(res)
	if res.StrVal("abc") != "123" {
		t.Error("error")
		return
	}
	//
	res = ParseSectionF("(", ")", `
(abc)
123
(/abc)
		`)
	fmt.Println(res)
	if res.StrVal("abc") != "123" {
		t.Error("error")
		return
	}
	//
	res = ParseSectionF("[", "]", `
(abc)
123
(/abcx)
		`)
	if len(res) > 0 {
		t.Error("error")
		return
	}
}
