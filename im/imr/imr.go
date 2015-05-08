package main

import (
	"bufio"
	"fmt"
	"github.com/Centny/gwf/im"
	"github.com/Centny/gwf/im/pb"
	"github.com/Centny/gwf/log"
	"github.com/Centny/gwf/netw"
	"github.com/Centny/gwf/util"
	"os"
	"runtime"
	"strings"
)

func usage() {
	fmt.Println(`
Usage: imr -t <login token> [-s <server address> | -l <server list url>] [-m R|C] [-L D|I|W|E|N]
   -t <login token> target login token.
   -s <server address> special target im server address.
   -l <server list url> the im server list URL,it will rand one server to logn.
   -m R auto reply model,it will auto reply to sender by the sampe content,[default].
   -m C client mode by user input,the input command is "receive content" which split by one empty string.
   -L D|I|W|E|N log level or log file path,default is D,
   `)
}

var imc *im.IMC
var ef func(code int) = os.Exit

func main() {
	ef(run())
}
func run() int {
	var l, s, t string
	var m string = "R"
	var L string = "D"
	olen := len(os.Args)
	for i, v := range os.Args {
		if i > olen-2 {
			break
		}
		switch v {
		case "-l":
			l = os.Args[i+1]
		case "-s":
			s = os.Args[i+1]
		case "-t":
			t = os.Args[i+1]
		case "-m":
			m = os.Args[i+1]
		case "-L":
			L = os.Args[i+1]
		case "-h":
			usage()
			return 1
		}
	}
	if len(t) < 1 {
		usage()
		return 1
	}
	if len(s) > 0 {
		imc = im.NewIMC(s, t)
	} else if len(l) > 0 {
		imc_, err := im.NewIMC4(l, t)
		if err != nil {
			fmt.Println(err.Error())
			return 1
		}
		imc = imc_
	} else {
		usage()
		return 1
	}
	runtime.GOMAXPROCS(runtime.NumCPU() - 1)
	switch L {
	case "N":
		log.SetLevel(log.NONE)
	case "E":
		log.SetLevel(log.ERROR)
	case "W":
		log.SetLevel(log.WARNING)
	case "I":
		log.SetLevel(log.INFO)
	case "D":
		log.SetLevel(log.DEBUG)
	default:
		log.SetLevel(log.DEBUG)
		lf, err := os.OpenFile(L, os.O_CREATE|os.O_WRONLY|os.O_APPEND, os.ModePerm)
		if err != nil {
			fmt.Println(err.Error())
			return 1
		}
		defer lf.Close()
		log.SetWriter(lf)
	}
	if m == "C" {
		imc.OnM = func(i *im.IMC, c netw.Cmd, m *pb.ImMsg) int {
			fmt.Println(string(m.GetC()))
			return 0
		}
		go func() {
			<-imc.LC
			buf := bufio.NewReader(os.Stdin)
			last := ""
			for {
				bys, err := util.ReadLine(buf, 10240, false)
				if err != nil {
					break
				}
				cmds := strings.SplitN(string(bys), " ", 2)
				if len(cmds) > 1 {
					imc.SMS(cmds[0], 0, cmds[1])
					last = cmds[0]
					continue
				}
				if len(last) < 1 {
					fmt.Println("command not split by empty string")
				} else {
					imc.SMS(last, 0, cmds[0])
				}
			}
		}()
		log.I("imr running on mode C")
	} else if m == "R" {
		imc.OnM = func(i *im.IMC, c netw.Cmd, m *pb.ImMsg) int {
			m.R = []string{m.GetS()}
			c.Writev(m)
			fmt.Println(string(m.GetC()))
			return 0
		}
		log.I("imr running on mode R")
	} else {
		fmt.Println("unknow model")
		usage()
		return 1
	}
	//
	imc.StartRunner()
	<-imc.WC
	return 0
}
