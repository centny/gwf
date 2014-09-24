package main

import (
	"github.com/Centny/gwf/jcr"
	"os"
)

func main() {
	err := jcr.Run()
	if err != nil {
		os.Exit(1)
	}
}
