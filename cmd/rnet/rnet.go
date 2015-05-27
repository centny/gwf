package main

import (
	"encoding/base64"
	"fmt"
	"github.com/Centny/gwf/log"
	"github.com/Centny/gwf/netw"
	"github.com/Centny/gwf/netw/impl"
	"github.com/Centny/gwf/pool"
	"github.com/Centny/gwf/util"
	"os"
	"runtime"
	"strings"
	"time"
)

type Cmd struct {
	Data []byte
}

func (cmd *Cmd) run_c(c netw.Con) {
	var idx int = 0
	for {
		_, err := c.Writeb(cmd.Data, []byte(fmt.Sprintf("%v", idx)))
		if err != nil {
			log.E("writing data(%v) err(%v) to %v", idx, err.Error(), c.RemoteAddr())
			break
		}
		idx++
		time.Sleep(time.Second)
	}
	fmt.Println("run_c end....")
}
func (cmd *Cmd) OnConn(c netw.Con) bool {
	c.SetWait(true)
	go cmd.run_c(c)
	return true
}
func (cmd *Cmd) OnClose(c netw.Con) {
}
func (cmd *Cmd) OnCmd(c netw.Cmd) int {
	log.I("receive data(%v) from %v", string(c.Data()), c.RemoteAddr())
	return 0
}
func RNet(addr string, d []byte) error {
	netw.ShowLog = true
	impl.ShowLog = true
	p := pool.NewBytePool(8, 1024)
	l := netw.NewListener(p, addr, "N", &Cmd{
		Data: d,
	})
	l.T = 100
	err := l.Run()
	if err != nil {
		return err
	}
	defer l.Close()
	l.Wait()
	return nil
}

//
//
var exit = os.Exit

func usage() {
	fmt.Print(`Usage: rnet [-addr=:12345] [-t=text message] [-b=base64 message]
  -addr listen address.
  -t the plain text to sending to client, default is "S->\n".
  -b the base64 data to sending to client, it will auto decode to []byte.
 `)
}
func main() {
	runtime.GOMAXPROCS(util.CPU())
	var err error
	var addr string = ":12345"
	var data []byte = []byte("S->")
	for _, arg := range os.Args {
		switch {
		case strings.HasPrefix(arg, "-t="):
			data = []byte(strings.TrimPrefix(arg, "-t="))
		case strings.HasPrefix(arg, "-b="):
			data, err = base64.StdEncoding.DecodeString(strings.TrimPrefix(arg, "-b="))
			if err != nil {
				fmt.Println(err.Error())
				usage()
				exit(1)
			}
		case strings.HasPrefix(arg, "-addr="):
			addr = strings.TrimPrefix(arg, "-addr=")
		}
	}
	if len(addr) < 1 || len(data) < 1 {
		usage()
		exit(1)
		return
	}
	err = RNet(addr, data)
	if err != nil {
		fmt.Println(err.Error())
		exit(1)
	}
}
