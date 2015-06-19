package main

import (
	"fmt"
	"os"
	"strings"
)

var ef func(c int) = os.Exit

func main() {
	var j string
	var p string
	var ign, m string
	var a string
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
		case "-i":
			if i < olen-1 {
				ign = os.Args[i+1]
				i++
			}
		case "-m":
			if i < olen-1 {
				m = os.Args[i+1]
				i++
			}
		case "-a":
			if i < olen-1 {
				a = os.Args[i+1]
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
	lsd := NewLsd(m, ign, p)
	paths := strings.Split(path, ",")
	for _, pt := range paths {
		lsd.Walk(pt)
	}
	if len(j) > 0 {
		lsd.JoinPrint(j, a)
	} else {
		lsd.Print(a)
	}
}
func usage() {
	fmt.Println(`Usage:	gpkg [-j <seq>] [-p <prefix>] [-a append prefix] [-i ingore regex] [-m match regex] <base path>`)
}
