package main

import (
	"fmt"
	"github.com/Centny/gwf/util"
	"os"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "Usage:fcfg <configure file> <format string>")
		os.Exit(1)
	}
	cfg, err := util.NewFcfg(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
	fmt.Println(cfg.EnvReplace(os.Args[2]))
}
