package main

import (
	"fmt"
	"os"
)

var ef func(c int) = os.Exit

func main() {
	var j string
	var p string
	var path string
	olen := len(os.Args)
	for i := 1; i < olen; i++ {
		switch os.Args[i] {
		case "-j":
			if i < olen-1 {
				j = os.Args[i+1]
				i++
			}
		case "-p":
			if i < olen-1 {
				p = os.Args[i+1]
				i++
			}
		case "-h":
			usage()
			ef(1)
			return
		default:
			path = os.Args[i]
		}
	}
	if len(path) < 1 {
		usage()
		ef(1)
		return
	}
	lsd := NewLsd(p)
	lsd.Walk(path)
	if len(j) > 0 {
		lsd.JoinPrint(j)
	} else {
		lsd.Print()
	}
}
func usage() {
	fmt.Println(`Usage:	gpkg -j <seq> <base path>`)
}
