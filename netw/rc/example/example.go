package main

import (
	"fmt"
	"github.com/Centny/gwf/netw"
	"github.com/Centny/gwf/routing"
	"github.com/Centny/gwf/util"
	"os"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: example <-s|-c] addr")
		return
	}
	netw.ShowLog = true
	switch os.Args[1] {
	case "-s":
		fmt.Println(StartRCSrv(os.Args[2]))
		routing.HFunc("^/test.*$", func(hs *routing.HTTPSession) routing.HResult {
			CallTestClient(hs.RVal("value"))
			fmt.Println("done...")
			return routing.HRES_RETURN
		})
		routing.ListenAndServe(":2422")
	case "-c":
		StartRCClient(os.Args[2])
		var res, err = CallTestSrv(fmt.Sprintf("client %v", util.Now()))
		if err == nil {
			fmt.Println(res.StrVal("data"))
		} else {
			fmt.Println("error->", err)
		}
		Runner.Wait()
	}
}
