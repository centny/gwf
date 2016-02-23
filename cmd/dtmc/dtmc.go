package main

import (
	"fmt"
	"github.com/Centny/gwf/netw/dtm"
	"github.com/Centny/gwf/util"
	"os"
)

func main() {
	var cfg = "conf/dtmc.properties"
	if len(os.Args) > 1 {
		cfg = os.Args[1]
	}
	var fcfg = util.NewFcfg3()
	fcfg.InitWithFilePath2(cfg, true)
	var dtmc = dtm.StartDTM_C(fcfg)
	fmt.Println(dtmc.RunProcH())
}
