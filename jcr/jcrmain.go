package jcr

import (
	"fmt"
	"os"
)

func Usage() {
	fmt.Println(`Usage:jcr app|start options
	-d <html directory> 
	-o <output directory> 
	-ex <exclude>
	-in <include>
	-js <jcr.js location>
	-n <coverage file name>
	-p <listen port>`)
}
func Run() {
	if len(os.Args) < 2 {
		Usage()
		return
	}
	dir := ""
	ex := ""          //exclude
	in := ".*\\.html" //include
	out := "out"
	js := "http://localhost:5457/jcr/jcr.js"
	port := ":5457"
	name := "coverage_jcr"
	alen := len(os.Args) - 1
	for i := 2; i < alen; i++ {
		switch os.Args[i] {
		case "-d":
			dir = os.Args[i+1]
			i++
		case "-ex":
			ex = os.Args[i+1]
			i++
		case "-in":
			in = os.Args[i+1]
			i++
		case "-o":
			out = os.Args[i+1]
			i++
		case "-js":
			js = os.Args[i+1]
			i++
		case "-p":
			port = os.Args[i+1]
			i++
		case "-n":
			name = os.Args[i+1]
			i++
		default:
			fmt.Println("invalid option:" + os.Args[i])
		}
	}
	switch os.Args[1] {
	case "start":
		// if len(name) < 1 || len(out) < 1 || len(port) < 1 {
		// 	Usage()
		// 	return
		// }
		RunSrv(name, out, port)
	default:
		if len(dir) < 1 || len(js) < 1 {
			Usage()
			return
		}
		RunAppend(dir, ex, in, out, js)
	}
}
