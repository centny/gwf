package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/Centny/gwf/log"

	"github.com/Centny/gwf/util"

	"github.com/Centny/gwf/netw/rcmd"
)

const (
	DefaultAddr = ":2984"
)

func main() {
	if len(os.Args) < 2 {
		runControl("127.0.0.1"+DefaultAddr, "Ctrl-local", "local")
		return
	}
	if len(os.Args) == 2 {
		switch os.Args[1] {
		case "-m":
			runMaster(DefaultAddr, `{"Slave-abc":1,"Ctrl-abc":1}`)
		case "-c":
			runControl("127.0.0.1"+DefaultAddr, "Ctrl-local", "local")
		default:
			printUsage(1)
		}
	}
	switch os.Args[1] {
	case "-m":
		runMaster(os.Args[2:]...)
	case "-c":
		runControl(os.Args[2:]...)
	case "-s":
		runSlave(os.Args[2:]...)
	case "-h":
		printUsage(0)
	default:
		printUsage(1)
	}
}

func runControl(args ...string) {
	if len(args) < 2 {
		printUsage(1)
		return
	}
	logpath := "/tmp/rc_control.log"
	logfile, err := os.OpenFile(logpath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		panic(err)
	}
	defer logfile.Close()
	log.SetWriter(logfile)
	rcaddr, token := args[0], args[1]
	alias, _ := os.Hostname()
	if len(args) > 2 {
		alias = args[2]
	}
	if len(alias) < 1 {
		alias = "control"
	}
	SyncHistory()
	fmt.Printf("run control by rcaddr(%v),token(%v),alias(%v)\n", rcaddr, token, alias)
	fmt.Printf("connecting...\n")
	err = rcmd.StartControl(alias, rcaddr, token)
	if err != nil {
		panic(err)
	}
	// fmt.Printf("master is connected\n")
	// rcmd.SharedControl.Wait()
	// stdin := bufio.NewReader(os.Stdin)
	var cids = ""
	for {
		// time.Sleep(time.Second)
		baseline, err := String("> ")
		if err == io.EOF {
			break
		}
		line := strings.TrimSpace(baseline)
		if len(line) < 1 {
			continue
		}
		StoreHistory(baseline)
		if strings.HasPrefix(line, "@ls") {
			res, err := rcmd.SharedControl.List(cids)
			if err == nil {
				fmt.Println(util.S2Json(res))
			} else {
				fmt.Printf("@ls cmd fail with %v\n", err)
			}
			continue
		}
		if strings.HasPrefix(line, "@start") {
			line = strings.TrimPrefix(line, "@start")
			line = strings.TrimSpace(line)
			if len(line) < 1 {
				fmt.Printf("@start <commands> [<args>] [>log]\n")
				continue
			}
			parts := strings.SplitN(line, ">", 2)
			cmds := parts[0]
			logfile := ""
			if len(parts) > 1 {
				logfile = strings.TrimSpace(parts[1])
			}
			res, err := rcmd.SharedControl.StartCmd(cids, "", cmds, logfile)
			if err == nil {
				fmt.Println(util.S2Json(res))
			} else {
				fmt.Printf("@start cmd fail with %v\n", err)
			}
			continue
		}
		if strings.HasPrefix(line, "@eval") {
			line = strings.TrimPrefix(line, "@eval")
			line = strings.TrimSpace(line)
			if len(line) < 1 {
				fmt.Printf("@eval <script file> [<args>] [>log]\n")
				continue
			}
			parts := strings.SplitN(line, ">", 2)
			cmds := strings.TrimSpace(parts[0])
			logfile := ""
			if len(parts) > 1 {
				logfile = strings.TrimSpace(parts[1])
			}
			cmdsParts := strings.SplitN(cmds, " ", 2)
			shellStr, err := ioutil.ReadFile(cmdsParts[0])
			if err != nil {
				fmt.Printf("read shell file(%v) fail with %v\n", cmdsParts[0], err)
				continue
			}
			argsStr := ""
			if len(cmdsParts) > 1 {
				argsStr = cmdsParts[1]
			}
			res, err := rcmd.SharedControl.StartCmd(cids, string(shellStr), argsStr, logfile)
			if err == nil {
				fmt.Println(util.S2Json(res))
			} else {
				fmt.Printf("@eval shell fail with %v\n", err)
			}
			continue
		}
		if strings.HasPrefix(line, "@exec") {
			line = strings.TrimPrefix(line, "@exec")
			line = strings.TrimSpace(line)
			if len(line) < 1 {
				fmt.Printf("@exec <script file> [<args>]\n")
				continue
			}
			cmdsParts := strings.SplitN(line, " ", 2)
			shellStr, err := ioutil.ReadFile(cmdsParts[0])
			if err != nil {
				fmt.Printf("read shell file(%v) fail with %v\n", cmdsParts[0], err)
				continue
			}
			argsStr := ""
			if len(cmdsParts) > 1 {
				argsStr = cmdsParts[1]
			}
			res, err := rcmd.SharedControl.RunCmd(cids, string(shellStr), argsStr)
			if err == nil {
				buf := bytes.NewBuffer(nil)
				for key, val := range res {
					fmt.Fprintf(buf, "%v:\n%v\n", key, val)
				}
				fmt.Printf("%v\n", string(buf.Bytes()))
			} else {
				fmt.Printf("@exec shell fail with %v\n", err)
			}
			continue
		}
		if strings.HasPrefix(line, "@stop") {
			line = strings.TrimPrefix(line, "@stop")
			line = strings.TrimSpace(line)
			if len(line) < 1 {
				fmt.Printf("@stop [<cid>] [<tid>]\n")
				continue
			}
			line = regexp.MustCompile("[ ]+").ReplaceAllString(line, " ")
			parts := strings.SplitN(line, " ", 2)
			var res util.Map
			switch len(parts) {
			case 1:
				res, err = rcmd.SharedControl.StopCmd(cids, parts[0])
			default:
				res, err = rcmd.SharedControl.StopCmd(parts[0], parts[1])
			}
			if err == nil {
				fmt.Println(util.S2Json(res))
			} else {
				fmt.Printf("@stop shell fail with %v\n", err)
			}
			continue
		}
		if strings.HasPrefix(line, "@select") {
			line = strings.TrimPrefix(line, "@select")
			line = strings.TrimSpace(line)
			if len(line) < 1 {
				fmt.Printf("@select [<all or cids>]\n")
				continue
			}
			switch line {
			case "":
				fmt.Printf("current selected:%v\n", cids)
			case "all":
				cids = ""
				fmt.Printf("all client selected\n")
			default:
				cids = line
			}
			continue
		}
		if strings.HasPrefix(line, "@help") {
			printCtrlUsage()
			continue
		}
		if strings.HasPrefix(line, "@exit") {
			break
		}
		if strings.HasPrefix(line, "exit") {
			break
		}
		//default as run command
		if strings.HasPrefix(line, "@run") {
			line = strings.TrimPrefix(line, "@run")
			line = strings.TrimSpace(line)
			if len(line) < 1 {
				fmt.Printf("@run <commands> [<args>]\n")
				continue
			}
		}
		res, err := rcmd.SharedControl.RunCmd(cids, "", line)
		if err == nil {
			buf := bytes.NewBuffer(nil)
			for key, val := range res {
				fmt.Fprintf(buf, "%v:\n%v\n", key, val)
			}
			fmt.Printf("%v\n", string(buf.Bytes()))
		} else {
			fmt.Printf("@run cmd fail with %v\n", err)
		}
	}
	rcmd.StopControl()
}

var HISTORY = ""

func SetHistory() {
	usr, err := user.Current()
	if err != nil {
		panic(err)
	}
	HISTORY = filepath.Join(usr.HomeDir, ".rcmd_history")
}

func StoreHistory(line string) {
	AddHistory(line)
	err := SaveHistory(HISTORY)
	if err != nil {
		log.E("save history to %v fail with %v", HISTORY, err)
	}
}

func SyncHistory() {
	SetHistory()
	err := LoadHistory(HISTORY)
	if err != nil {
		log.E("load history to %v fail with %v", HISTORY, err)
	}
}

func runMaster(args ...string) {
	if len(args) < 2 {
		printUsage(1)
		return
	}
	rcaddr, tokens := args[0], args[1]
	fmt.Printf("run master by rcaddr(%v),tokens(%v)\n", rcaddr, tokens)
	var ts = map[string]int{}
	util.Json2S(tokens, &ts)
	ts["Ctrl-local"] = 1
	err := rcmd.StartMaster(rcaddr, ts)
	if err != nil {
		panic(err)
	}
	rcmd.SharedMaster.Wait()
}

func runSlave(args ...string) {
	if len(args) < 2 {
		printUsage(1)
		return
	}
	rcaddr, token := args[0], args[1]
	alias, _ := os.Hostname()
	if len(args) > 2 {
		alias = args[2]
	}
	if len(alias) < 1 {
		panic("the slave alias is empty")
	}
	fmt.Printf("run slave by rcaddr(%v),token(%v),alias(%v)\n", rcaddr, token, alias)
	err := rcmd.StartSlave(alias, rcaddr, token)
	if err != nil {
		panic(err)
	}
	// rcmd.SharedSlave.Wait()
	wait := make(chan int)
	<-wait
}

func printUsage(exit int) {
	_, name := filepath.Split(os.Args[0])
	fmt.Printf(`Usage:
	%v -m [<listen addr> <token config>]	run as master
	%v -s <master addr> <token> [<alias>]	run as slave
	%v -c <master addr> <token> [<alias>]	run as control
	%v	run as local control%v`,
		name, name, name, name, "\n")
	os.Exit(exit)
}

func printCtrlUsage() {
	fmt.Println(`
	@ls		=>list the running task.
	@start <commands> [<args>] [>log] => start command.
	@eval <script file> [<args>] [>log] => eval local script to remote host async.
	@exec <script file> [<args>] => eval local script to remote host sync.
	@stop [<cid>] [<tid>] => stop command on special client.
	@stop <tid> => stop command on all client.
	@select [<all or cids>] => active clients.
	@exit => exit the current shell.
	@help	=> show this.`)
}
