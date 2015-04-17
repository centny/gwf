package main

import (
	"fmt"
	"github.com/Centny/gwf/util"
	"os"
)

func main() {
	alen := len(os.Args)
	if alen == 3 {
		fmt_s(os.Args[1], os.Args[2])
	} else if alen == 4 && os.Args[1] == "-f" {
		fmt_s(os.Args[2], os.Args[3])
	} else if alen == 6 && os.Args[1] == "-s" {
		sec_s(os.Args[2], os.Args[3], os.Args[4], os.Args[5])
	} else {
		fmt.Println(`
Usage:
	fcfg [-f] <configure file> <format string>
	fcfg -s <configure file> <selected section> <store file> <store section>
			`)
		os.Exit(1)
	}
}

func fmt_s(tf, fs string) {
	cfg := util.NewFcfg3()
	cfg.ShowLog = false
	err := cfg.InitWithFilePath(tf)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
	fmt.Println(cfg.EnvReplaceV(fs, true))
}
func sec_s(tf, sec, sp, ssec string) {
	cfg, err := util.NewFcfg(tf)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
	cfg.ShowLog = true
	err = cfg.Store(sec, sp, ssec)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
