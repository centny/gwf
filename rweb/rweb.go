package main

import (
	"fmt"
	"github.com/Centny/gwf/util"
	"os"
)

func main() {
	addr := ":80"
	if len(os.Args) > 1 {
		addr = os.Args[1]
	}
	fmt.Println("running on", addr)
	fmt.Println(util.DoWeb(addr, "."))
}
