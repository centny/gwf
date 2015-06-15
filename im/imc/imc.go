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
	"strconv"
	"strings"
)

func usage() {
	fmt.Println(`
Usage: imr -t <login token> [-s <server address> | -l <server list url>] [-m R|C|T] [-L D|I|W|E|N] [-g groups] [-c message count] [-P push url] [-p push user] [-T timeout]
   -t <login token> target login token.
   -s <server address> special target im server address.
   -l <server list url> the im server list URL,it will rand one server to logn.
   -m R auto reply model,it will auto reply to sender by the sampe content,[default].
   -m C client mode by user input,the input command is "receive content" which split by one empty string.
   -m T running test module
   -L D|I|W|E|N log level or log file path,default is D
   -g <groups> test group split by ,
   -c <8> the message count,default 8
   -P <http://127.0.0.1/api/doPcm?s=%v&r=%v&c=%v&t=%v> the push server format url
   -p <U-1> the push user R
   -T <8000> timeout millisecond, default 8s.
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
	var g, P, p string
	var c int = 8
	var T int64 = 8000
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
		case "-g":
			g = os.Args[i+1]
		case "-P":
			P = os.Args[i+1]
		case "-p":
			p = os.Args[i+1]
		case "-c":
			tc, err := strconv.ParseInt(os.Args[i+1], 10, 64)
			if err != nil {
				fmt.Println(err.Error())
				usage()
				return 1
			}
			c = int(tc)
		case "-T":
			tc, err := strconv.ParseInt(os.Args[i+1], 10, 64)
			if err != nil {
				fmt.Println(err.Error())
				usage()
				return 1
			}
			T = tc
		case "-h":
			usage()
			return 1
		}
	}
	runtime.GOMAXPROCS(util.CPU())
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
	if len(t) < 1 {
		usage()
		return 1
	}
	if m == "T" {
		err := run_do_imc(s, l, t, g, P, p, c, T)
		if err == nil {
			return 0
		} else {
			fmt.Println(err.Error())
			return 1
		}
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
	var fc func()
	if m == "C" {
		imc.OnM = func(i *im.IMC, c netw.Cmd, m *pb.ImMsg) int {
			fmt.Println(string(m.GetC()))
			return 0
		}
		fc = func() {
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
		}
		log.I("imr running on mode C")
	} else if m == "R" {
		imc.OnM = func(i *im.IMC, c netw.Cmd, m *pb.ImMsg) int {
			m.R = []string{m.GetS()}
			c.Writev(m)
			fmt.Println(string(m.GetC()))
			return 0
		}
		fc = nil
		log.I("imr running on mode R")
	} else {
		fmt.Println("unknow model")
		usage()
		return 1
	}
	//
	imc.Start()
	imc.LC.Wait()
	if imc.Logined() {
		imc.StartHB()
	}
	if fc != nil {
		go fc()
	}
	<-imc.WC
	imc = nil
	return 0
}

func run_do_imc(s, l, t, g, P, p string, c int, timeout int64) error {
	log.D("running do imc--->")
	var srv string
	if len(l) > 1 {
		srv = l
	} else if len(s) > 1 {
		srv = s
	} else {
		return util.Err("the server parameter is empty")
	}
	gs := []string{}
	if len(g) > 0 {
		gs = strings.Split(g, ",")
	}
	di := im.NewDoImc(srv, len(l) > 1, strings.Split(t, ","), gs, c, P, p)
	err := di.Do()
	if err != nil {
		return err
	}
	return di.Check2(1000, timeout)
}
