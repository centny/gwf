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
	switch os.Args[1] {
	case "-p":
		if alen > 3 {
			cfg_p(os.Args[2], os.Args[3])
		} else if alen == 3 {
			cfg_p(os.Args[2], "")
		} else {
			usage()
			os.Exit(1)
		}
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
	fcfg -p <configure file> [section name]
	fcfg [-f] <configure file> <format string>
	fcfg -s <configure file> <selected section> <store file> <store section>
			`)
}
func cfg_p(tf, sec string) {
	cfg := util.NewFcfg3()
	cfg.ShowLog = false
	err := cfg.InitWithFilePath(tf)
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
