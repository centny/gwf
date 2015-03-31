package main

import (
	"fmt"
	"github.com/Centny/gwf/util"
	"os"
)

func main() {
	if len(os.Args) < 5 {
		fmt.Println(`Usage:sr_c <url> <sr file> <app id> <app version>`)
		os.Exit(1)
	}
	res, err := util.HPostF2(os.Args[1], map[string]string{
		"aid":  os.Args[3],
		"ver":  os.Args[4],
		"exec": "A",
	}, "sr_f", os.Args[2])
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	} else if res.IntVal("code") == 0 {
		fmt.Println("OK")
		os.Exit(0)
	} else {
		fmt.Println(res.StrVal("dmsg"))
		os.Exit(1)
	}
}
