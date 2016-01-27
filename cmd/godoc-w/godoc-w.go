package main

import (
	"fmt"
	"github.com/Centny/gwf/routing"
	"github.com/Centny/gwf/util"
	"github.com/Centny/gwf/wdoc"
	"net/http"
	"os"
	"runtime"
	"strings"
	"time"
)

var ef func(c int) = os.Exit

func main() {
	var inc string = ""
	var exc string = ""
	var addr string = ""
	var prefix string = ""
	var delay int64 = 100000
	_, options, path := util.Args()
	err := options.ValidF(`
		i,O|S,L:0;
		e,O|S,L:0;
		a,R|S,L:0;
		p,O|S,L:0;
		d,O|I,R:0;
		`, &inc, &exc, &addr, &prefix, &delay)
	if err != nil {
		fmt.Println(err.Error())
		usage()
		ef(1)
		return
	}
	runtime.GOMAXPROCS(util.CPU())
	var wd = "."
	if len(path) > 0 {
		wd = path[0]
	}
	var inc_, exc_ []string
	fmt.Println(inc, exc)
	if len(inc) > 0 {
		inc_ = strings.Split(inc, ",")
	}
	if len(exc) > 0 {
		exc_ = strings.Split(exc, ",")
	}
	pars := wdoc.NewParser()
	pars.Pre = prefix
	go pars.LoopParse(wd, inc_, exc_, time.Duration(delay))
	mux := routing.NewSessionMux2("")
	mux.H("^.*$", pars)
	http.ListenAndServe(addr, mux)
}
func usage() {
	fmt.Println(`Usage:
	godoc-w -inc <include list> -exc <exclude list> -addr <listen addr> -prefix <prefix trim> -delay <check delay> <root path>
			`)
}
