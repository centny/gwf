package main

import (
	"fmt"
	tlog "github.com/Centny/gwf/log"
	"github.com/Centny/gwf/tools"
	"github.com/Centny/gwf/util"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
)

var exit = os.Exit

func usage() {
	fmt.Print(`Usage: mexec [-min=0] [-max=8] [-total=8] [-file=save path] [-emma=save path] [-id] [-log] [-MT=10000] [-fail=0] bin args
   -min minimum executer.
   -max maximum executer.
   -total total executer.
   -file save execute reslut to file, default is stdout.
   -emma save execute reslut to emma xml format.
   -id append executer id to the end of arguments.
   -log show log.
   -fail failed when errro rate is greater target value.
   -MT keep max time, default is 10000 ms.
`)
}
func main() {
	runtime.GOMAXPROCS(util.CPU())
	var err error
	var id bool = false
	var log bool = false
	var file string
	var emma string
	var min, max, total int = 0, 8, 8
	var args []string = []string{}
	var MT int64 = 10000
	var for_brk = false
	var fail float64 = 0
	for idx, arg := range os.Args {
		if idx == 0 {
			continue
		}
		switch {
		case strings.HasPrefix(arg, "-min="):
			min, err = util.ParseInt(strings.TrimPrefix(arg, "-min="))
			if err != nil {
				fmt.Println(err.Error())
				usage()
				exit(1)
			}
		case strings.HasPrefix(arg, "-max="):
			max, err = util.ParseInt(strings.TrimPrefix(arg, "-max="))
			if err != nil {
				fmt.Println(err.Error())
				usage()
				exit(1)
			}
		case strings.HasPrefix(arg, "-total="):
			total, err = util.ParseInt(strings.TrimPrefix(arg, "-total="))
			if err != nil {
				fmt.Println(err.Error())
				usage()
				exit(1)
			}
		case strings.HasPrefix(arg, "-MT="):
			MT, err = util.ParseInt64(strings.TrimPrefix(arg, "-MT="))
			if err != nil {
				fmt.Println(err.Error())
				usage()
				exit(1)
			}
		case strings.HasPrefix(arg, "-file="):
			file = strings.TrimPrefix(arg, "-file=")
		case strings.HasPrefix(arg, "-emma="):
			emma = strings.TrimPrefix(arg, "-emma=")
		case strings.HasPrefix(arg, "-log"):
			log = true
		case strings.HasPrefix(arg, "-id"):
			id = true
		case strings.HasPrefix(arg, "-fail="):
			fail, err = strconv.ParseFloat(strings.TrimPrefix(arg, "-fail="), 64)
			if err != nil {
				fmt.Println(err.Error())
				usage()
				exit(1)
			}
		default:
			if len(os.Args)-1 >= idx {
				args = os.Args[idx:]
			}
			for_brk = true
			break
		}
		if for_brk {
			break
		}
	}
	if max < 1 || total < 1 || max > total || len(args) < 1 {
		usage()
		exit(1)
		return
	}
	//
	if log {
		tlog.D("run exec by min(%v),max(%v),total(%v),bin(%v),args(%v)", min, max, total, args[0], args[1:])
	}
	exk := tools.NewExeK(min, max, total, args[0], args[1:]...)
	exk.CmdF = func(exe *tools.Exec, exk *tools.ExeK, idx string) *exec.Cmd {
		args := exe.Args
		if id {
			args = append(args, fmt.Sprintf("%v", idx))
		}
		return exec.Command(exe.Bin, args...)
	}
	exk.MT = MT
	exk.ShowLog = log
	exk.Start()
	exk.Wait()
	if len(file) < 1 && len(emma) < 1 {
		exk.Save(os.Stdout)
		os.Stdout.WriteString("\n")
	} else {
		err = exk.SaveP(file, emma)
		if err != nil {
			fmt.Println(err.Error())
			exit(1)
		}
	}
	_, suc, total := exk.Data()
	if float64(total-suc)/float64(total) > fail {
		fmt.Println(fmt.Sprintf("---->mexec is failed(%v) by error(%v)/total(%v)", fail, total-suc, total))
		exit(1)
	}
}
