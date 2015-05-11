package main

import (
	"fmt"
	tlog "github.com/Centny/gwf/log"
	"github.com/Centny/gwf/tools"
	"github.com/Centny/gwf/util"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

var exit = os.Exit

func usage() {
	fmt.Print(`Usage: mexec [-min=0] [-max=8] [-total=8] [-file=save path] [-id] [-log] [-MT=10000] bin args
   -min minimum executer.
   -max maximum executer.
   -total total executer.
   -file save execute reslut to file, default is mexec.json.
   -id append executer id to the end of arguments.
   -log show log.
   -MT keep max time, default is 10000 ms.
`)
}
func main() {
	runtime.GOMAXPROCS(runtime.NumCPU() - 1)
	var err error
	var id bool = false
	var log bool = false
	var file string
	var min, max, total int = 0, 8, 8
	var args []string = []string{}
	var MT int64 = 10000
	var for_brk = false
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
		case strings.HasPrefix(arg, "-log"):
			log = true
		case strings.HasPrefix(arg, "-id"):
			id = true
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
	if len(file) < 1 {
		exk.Save(os.Stdout)
		os.Stdout.WriteString("\n")
		return
	}
	err = exk.SaveP(file)
	if err != nil {
		fmt.Println(err.Error())
		exit(1)
	}
}
