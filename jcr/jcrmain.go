package jcr

import (
	"fmt"
	"os"
)

func Usage() {
	fmt.Println("Usage:jcr app|start -d <html directory> -o <output directory> -ex <exclude> -in <include> -js <jcr.js location> -f <configure file>")
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
	js := ""
	cf := ""
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
		case "-f":
			cf = os.Args[i+1]
			i++
		default:
			fmt.Println("invalid option:" + os.Args[i])
		}
	}
	switch os.Args[1] {
	case "start":
		if len(cf) < 1 {
			Usage()
			return
		}
		StartSrv(cf)
	default:
		if len(dir) < 1 || len(js) < 1 {
			Usage()
			return
		}
		RunAppend(dir, ex, in, out, js)
	}
}
