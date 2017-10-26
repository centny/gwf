package pshell

import (
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/Centny/gwf/netw"
	"github.com/Centny/gwf/netw/rc"
	"github.com/Centny/gwf/pool"
	"github.com/Centny/gwf/util"
)

var SharedControl *Control

func StartControl(alias, rcaddr, token string) (err error) {
	SharedControl = NewControl(alias)
	err = SharedControl.Start(rcaddr, token)
	return
}

func StopControl() {
	SharedControl.R.Close()
	SharedControl = nil
}

type Control struct {
	Alias string
	R     *rc.RC_Runner_m
	Outer ShellOuter
}

func NewControl(alias string) *Control {
	return &Control{
		Alias: alias,
	}
}

func (c *Control) Start(rcaddr, token string) (err error) {
	auto := rc.NewAutoLoginH(token)
	auto.Args = util.Map{"alias": c.Alias}
	c.R = rc.NewRC_Runner_m_j(pool.BP, rcaddr, netw.NewCCH(netw.NewQueueConH(auto, c), c))
	c.R.Name = c.Alias
	auto.Runner = c.R
	c.R.Start()
	return c.R.Valid()
}

func (c *Control) List() (res util.Map, err error) {
	res, err = c.R.VExec_m("list", util.Map{})
	return
}

func (c *Control) Exec(cids string, shell, cmds string) (res util.Map, err error) {
	res, err = c.R.VExec_m("exec", util.Map{
		"cids":  cids,
		"shell": shell,
		"cmds":  cmds,
	})
	return
}

func (c *Control) AddSession(name, addr, username, password string) (err error) {
	_, err = c.R.VExec_s("add_session", util.Map{
		"name":     name,
		"addr":     addr,
		"username": username,
		"password": password,
	})
	return
}

//OnConn see ConHandler for detail
func (c *Control) OnConn(con netw.Con) bool {
	//fmt.Println("master is connected")
	return true
}

//OnClose see ConHandler for detail
func (c *Control) OnClose(con netw.Con) {
	//fmt.Println("master is disconnected")
}

//OnCmd see ConHandler for detail
func (c *Control) OnCmd(con netw.Cmd) int {
	var data = util.Map{}
	con.V(&data)
	var n = data.StrVal("n")
	var m = data.StrVal("m")
	lines := strings.Split(m, "\n")
	length := len(lines)
	for i := 0; i < length; i++ {
		line := lines[i]
		if i == length-1 && len(line) < 1 {
			continue
		}
		c.writeOut(n, line)
	}
	return 0
}

func (c *Control) writeOut(name, message string) {
	if c.Outer == nil {
		fmt.Printf("%v:%v\n", name, message)
	} else {
		c.Outer.OnData(name, message)
	}
}

// func (c *Control) Wait() {
// 	c.R.Wait()
// }

type ShellOuter interface {
	OnData(name, message string)
}

type FileShellOuter struct {
	Logf  string
	Multi bool
	lck   sync.RWMutex
	out   map[string]*os.File
}

func NewFileShellOuter(logf string) *FileShellOuter {
	return &FileShellOuter{
		Logf:  logf,
		Multi: strings.Contains(logf, "%v"),
		out:   map[string]*os.File{},
	}
}

func (f *FileShellOuter) OnData(name, message string) {
	f.lck.Lock()
	var logf string
	if f.Multi {
		logf = fmt.Sprintf(f.Logf, name)
	} else {
		logf = f.Logf
	}
	outf, ok := f.out[logf]
	if !ok {
		var err error
		outf, err = os.OpenFile(logf, os.O_APPEND|os.O_CREATE|os.O_WRONLY, os.ModePerm)
		if err != nil {
			f.lck.Unlock()
			return
		}
	}
	f.lck.Unlock()
	outf.WriteString(message + "\n")
}
