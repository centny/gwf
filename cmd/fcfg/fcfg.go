package main

import (
	"fmt"
	"github.com/Centny/gwf/util"
	"os"
)

func main() {
	alen := len(os.Args)
	if alen < 2 {
		usage()
		os.Exit(1)
	}
	var cfg_p_ = func(log bool) {
		if alen == 4 {
			cfg_p(os.Args[2], os.Args[3], log)
		} else if alen == 3 {
			cfg_p(os.Args[2], "", log)
		} else {
			usage()
			os.Exit(1)
		}
	}
	switch os.Args[1] {
	case "-p":
		cfg_p_(false)
	case "-c":
		cfg_p_(true)
	case "-s":
		if alen > 5 {
			sec_s(os.Args[2], os.Args[3], os.Args[4], os.Args[5])
		} else {
			usage()
			os.Exit(1)
		}
	case "-f":
		if alen > 3 {
			fmt_s(os.Args[2], os.Args[3])
		} else {
			usage()
			os.Exit(1)
		}
	default:
		if alen == 3 {
			fmt_s(os.Args[1], os.Args[2])
		} else {
			usage()
			os.Exit(1)
		}
	}
}
func usage() {
	fmt.Println(`
Usage:
	fcfg -p <configure file> [section name] [true:show log]   ->print configure
	fcfg -c <configure file> [section name]	                  ->check configure to print and show log
	fcfg [-f] <configure file> <format string>                ->format configure to string
	fcfg -s <configure file> <selected section> <store file> <store section>  ->merge configure.
			`)
}
func cfg_p(tf, sec string, log bool) {
	cfg := util.NewFcfg3()
	cfg.ShowLog = log
	err := cfg.InitWithUri(tf)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
	if len(sec) > 0 {
		cfg.PrintSec(sec)
	} else {
		cfg.Print()
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
