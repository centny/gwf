package main

import (
	"github.com/Centny/Cny4go/jcr"
	"os"
)

func main() {
	err := jcr.Run()
	if err != nil {
		os.Exit(1)
	}
}
