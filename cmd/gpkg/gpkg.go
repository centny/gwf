package main

import (
	"fmt"
	"os"
)

var ef func(c int) = os.Exit

func main() {
	if len(os.Args) < 2 {
		usage()
		ef(1)
		return
	}
	if os.Args[1] == "-j" {
		if len(os.Args) < 4 {
			usage()
			ef(1)
			return
		}
		lsd := NewLsd()
		err := lsd.Walk(os.Args[3])
		if err != nil {
			fmt.Println(err.Error())
			ef(1)
			return
		}
		lsd.JoinPrint(os.Args[2])
	} else {
		lsd := NewLsd()
		err := lsd.Walk(os.Args[1])
		if err != nil {
			fmt.Println(err.Error())
			ef(1)
			return
		}
		lsd.Print()
	}
}
func usage() {
	fmt.Println(`Usage:	gpkg -j <seq> <base path>`)
}
