package main

import (
	"fmt"
	"github.com/Centny/gwf/util"
	"net/url"
	"os"
	"path/filepath"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: getw url <out file>")
		return
	}
	var turl, err = url.Parse(os.Args[1])
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	var fn string
	if len(os.Args) > 2 {
		fn = os.Args[2]
	} else {
		fn, _ = url.QueryUnescape(turl.Path)
		_, fn = filepath.Split(fn)
	}
	if len(fn) < 1 {
		fn = "index.html"
	}
	fmt.Println("Donwload", os.Args[1], "to", fn)
	err = util.DLoad(fn, os.Args[1])
	if err == nil {
		fmt.Println("OK")
	} else {
		fmt.Println(err)
	}
}
