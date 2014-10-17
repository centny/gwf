package main

import (
	"fmt"
	"github.com/Centny/gwf/ini"
	"os"
)

func main() {
	err := ini.Cmds(os.Args[0:])
	if err == nil {
		os.Exit(0)
	} else {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
