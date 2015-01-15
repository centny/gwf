package main

import (
	"fmt"
	"github.com/Centny/gwf/netw/example/rcmd/cc"
	"github.com/Centny/gwf/netw/example/rcmd/srv"
	"os"
)

func main() {
	if len(os.Args) > 1 && os.Args[1] == "-s" {
		srv.RunSrv()
	} else {
		DoC()
	}
}

func DoC() {
	cc.RunC()
	defer cc.Stop()
	vv, err := cc.List()
	if err != nil {
		panic(err.Error())
	}
	for _, v := range vv {
		fmt.Println(v.V)
	}
}
