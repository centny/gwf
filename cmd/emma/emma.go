package main

import (
	"fmt"
	"github.com/Centny/gwf/tutil"
	"os"
)

func main() {
	if len(os.Args) < 7 {
		fmt.Println("Usage: emma file name class method block line")
		os.Exit(1)
	}
	err := tutil.Emma(os.Args[1], os.Args[2], os.Args[3], os.Args[4], os.Args[5], os.Args[6])
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
