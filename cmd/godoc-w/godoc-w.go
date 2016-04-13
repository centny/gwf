package main

import (
	"fmt"
	"github.com/Centny/gwf/routing"
	"github.com/Centny/gwf/routing/filter"
	"github.com/Centny/gwf/util"
	"github.com/Centny/gwf/wdoc"
	"io/ioutil"
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
	var www string = "."
	var prefix string = ""
	var delay int64 = 100000
	var out string = ""
	var cmdf string = "pandoc %v -s --highlight-style tango"
	_, options, path := util.Args()
	err := options.ValidF(`
		inc,O|S,L:0;
		exc,O|S,L:0;
		addr,O|S,L:0;
		prefix,O|S,L:0;
		delay,O|I,R:0;
		out,O|S,L:0;
		www,O|S,L:0;
		cmdf,O|S,L:0;
		`, &inc, &exc, &addr, &prefix, &delay, &out, &www, &cmdf)
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
	if len(inc) > 0 {
		inc_ = strings.Split(inc, ",")
	}
	if len(exc) > 0 {
		exc_ = strings.Split(exc, ",")
	}
	pars := wdoc.NewParser(prefix, "/doc", cmdf)
	if len(addr) > 0 {
		go pars.LoopParse(wd, inc_, exc_, time.Duration(delay))
		mux := routing.NewSessionMux2("")
		mux.H("^.*$", filter.NewCORS2("*"))
		mux.H("^/doc.*$", pars)
		mux.Handler("^.*$", http.FileServer(http.Dir(www)))
		fmt.Println(http.ListenAndServe(addr, mux))
	} else if len(out) > 0 {
		err := pars.ParseDir(wd, inc_, exc_)
		if err != nil {
			fmt.Println(err.Error())
			ef(1)
			return
		}
		var res = pars.ToM()
		res.RateV()
		bys, _ := res.Marshal()
		err = ioutil.WriteFile(out, bys, os.ModePerm)
		if err != nil {
			fmt.Println(err.Error())
			ef(1)
			return
		}
	} else {
		fmt.Println("-addr or -out must be setted")
		ef(1)
	}
}
func usage() {
	fmt.Println(`Usage:
	godoc-w -inc <include list> -exc <exclude list> -addr <listen addr> -prefix <prefix trim> -delay <check delay> -out <coverage output file> <root path>
			`)
}
