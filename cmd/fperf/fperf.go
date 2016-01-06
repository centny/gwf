package main

import (
	"fmt"
	"github.com/Centny/gwf/log"
	"github.com/Centny/gwf/tutil"
	"github.com/Centny/gwf/util"
	"os"
)

var ef func(c int) = os.Exit

func main() {
	var wd string
	var mode string = "RW"
	var pref string = "test_"
	var bs int64 = 1024
	var count int = 1
	var total, max int = 100, util.CPU()
	var beg, end int = 0, 8
	var clean bool = true

	_, options, path := util.Args()
	err := options.ValidF(`
		M,O|S,O:RW~R~W;
		p,O|S,L:0;
		B,O|I,R:8~1024000;
		c,O|I,R:0;
		t,O|I,R:0;
		m,O|I,R:0;
		b,O|I,R:-1;
		e,O|I,R:0;
		`, &mode, &pref, &bs, &count, &total, &max, &beg, &end)
	if err != nil {
		fmt.Println(err.Error())
		usage()
		ef(1)
		return
	}
	if options.Exist("h") {
		usage()
		ef(0)
		return
	}
	clean = !options.Exist("n")
	if len(path) < 1 {
		wd = "."
	} else {
		wd = path[0]
	}
	tp := tutil.NewFPerf(wd)
	tp.Clear = clean
	var used, bys int64
	switch mode {
	case "RW":
		log.D(`test file system performance by RW mode with options
	pref:%v,
	bs:%v,
	count:%v,
	total:%v,
	max:%v,`, pref, bs, count, total, max)
		used, err = tp.Perf4MultiRw(pref, "", total, max, bs, count)
		bys = int64(total) * bs * int64(count)
	case "R":
		log.D(`test file system performance by R mode with options
	pref:%v,
	max:%v,
	beg:%v,
	end:%v,`, pref, max, beg, end)
		used, bys, err = tp.Perf4MultiR(pref, "", beg, end, max)
	case "W":
		log.D(`test file system performance by W mode with options
	pref:%v,
	bs:%v,
	count:%v,
	total:%v,
	max:%v,`, pref, bs, count, total, max)
		used, err = tp.Perf4MultiW(pref, "", total, max, bs, count)
		bys = int64(total) * bs * int64(count)
	}
	if err != nil {
		fmt.Println(err.Error())
		ef(1)
		return
	}
	if used < 1 {
		used = 1
	}
	fmt.Printf("Used:%vms, Speed:%v\n", used, util.BysSize(bys/used*1000))
}
func usage() {
	fmt.Println(`Usage:fperf [options] <test path>
	-M test mode in RW(read write)/R(only read)/W(only write)
	-p file perfi for read/write
	-B the buffer size between 8 and 1024000
	-c the data(buffer size) count for writing to file
	-t the total file count
	-m the max thread to run test at the same time
	-b the begin index to read file
	-e the end index to write file
	-n not clear the tmp file
	`)
}
