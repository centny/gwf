package main

import (
	"fmt"
	"github.com/Centny/gwf/routing"
	"github.com/Centny/gwf/util"
	"os"
	"runtime"
	"strings"
)

func main() {
	runtime.GOMAXPROCS(util.CPU())
	var addr, www = ":8090", "./"
	if len(os.Args) > 1 {
		addr = os.Args[1]
	}
	if addr == "-h" {
		fmt.Println("Usage: webts address www")
		return
	}
	if len(os.Args) > 2 {
		www = os.Args[2]
	}
	if !strings.HasSuffix(www, "/") {
		www = www + "/"
	}
	WWW = www
	routing.Shared.Sb = routing.NewSrvSessionBuilder("", "/", "sid", 100000, 8000)
	Hand("", routing.Shared)
	fmt.Printf("listen web server on addr(%v),www(%v)\n", addr, www)
	fmt.Println(routing.ListenAndServe(addr))
}
