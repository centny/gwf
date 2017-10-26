package main

import (
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

	"github.com/Centny/gwf/netw/pshell"
)

const (
	DefaultAddr = ":2734"
)

func main() {
	if len(os.Args) < 2 {
		runControl("127.0.0.1"+DefaultAddr, "Ctrl-local", "local")
		return
	}
	if len(os.Args) == 2 {
		switch os.Args[1] {
		case "-s":
			runServer()
		case "-c":
			runControl("127.0.0.1"+DefaultAddr, "Ctrl-local", "local")
		default:
			printUsage(1)
		}
	}
	switch os.Args[1] {
	case "-s":
		runServer(os.Args[2:]...)
	case "-c":
		runControl(os.Args[2:]...)
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
	logpath := "/tmp/pshell_control.log"
	logfile, err := os.OpenFile(logpath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		panic(err)
	}
	defer logfile.Close()
	log.SetWriter(logfile)
	//
	shellOutF := "/tmp/pshell_out_%v.log"
	shellOut := pshell.NewFileShellOuter(shellOutF)
	//
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
	err = pshell.StartControl(alias, rcaddr, token)
	if err != nil {
		panic(err)
	}
	pshell.SharedControl.Outer = shellOut
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
			res, err := pshell.SharedControl.List()
			if err == nil {
				fmt.Println(util.S2Json(res))
			} else {
				fmt.Printf("@ls cmd fail with %v\n", err)
			}
			continue
		}
		if strings.HasPrefix(line, "@add") {
			line = strings.TrimPrefix(line, "@add")
			line = strings.TrimSpace(line)
			cmds := regexp.MustCompile("[\\ ]+").Split(line, -1)
			if len(cmds) < 4 {
				fmt.Printf("@add <name> <address> <username> <password>\n")
				continue
			}
			err := pshell.SharedControl.AddSession(cmds[0], cmds[1], cmds[2], cmds[3])
			if err == nil {
				fmt.Println("ok")
			} else {
				fmt.Printf("@add shell fail with %v\n", err)
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
			cmds := strings.SplitN(line, " ", 2)
			shellStr, err := ioutil.ReadFile(cmds[0])
			if err != nil {
				fmt.Printf("read shell file(%v) fail with %v\n", cmds[0], err)
				continue
			}
			args := ""
			if len(cmds) > 1 {
				args = cmds[1]
			}
			res, err := pshell.SharedControl.Exec(cids, string(shellStr), args)
			if err == nil {
				fmt.Println(util.S2Json(res))
			} else {
				fmt.Printf("@eval shell fail with %v\n", err)
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
		res, err := pshell.SharedControl.Exec(cids, "", line)
		if err == nil {
			fmt.Println(util.S2Json(res))
		} else {
			fmt.Printf("@run cmd fail with %v\n", err)
		}
	}
	pshell.StopControl()
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

func runServer(args ...string) {
	var conf = "conf/pshell.properties"
	if len(args) > 0 {
		conf = args[0]
	}
	Conf.InitWithFilePath2(conf, true)
	var ts = ReadTokens()
	var hosts = ReadHosts()
	err := pshell.StartServer(Listen(), ts, hosts...)
	if err != nil {
		panic(err)
	}
	pshell.SharedServer.Wait()
}

func printUsage(exit int) {
	_, name := filepath.Split(os.Args[0])
	fmt.Printf(`Usage:
	%v -s <configure file>]	run as server
	%v -c <server addr> <token> [<alias>]	run as control
	%v	run as local control%v`,
		name, name, name, "\n")
	os.Exit(exit)
}

func printCtrlUsage() {
	fmt.Println(`
	@eval <script file> [<args>] [>log] => eval local script to remote host async.
	@select [<all or cids>] => active clients.
	exit => exit the current shell.
	@help	=> show this.`)
}
