package routing

import (
	"fmt"
	"regexp"
	"testing"
)

func TestNoCache(t *testing.T) {
	NewAllNoCacheDir("www")
	nnd := NewNoCacheDir("../test")
	nnd.Add(regexp.MustCompile("^.*\\.html(\\?.*)?$"))
	f, _ := nnd.Open("test.html")
	fi, _ := f.Stat()
	fmt.Println(fi.ModTime())
	nnd.Open("test.go")
	nnd.Open("tt.html")
}
