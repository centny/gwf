package rcmd

import (
	"fmt"

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
}

func NewControl(alias string) *Control {
	return &Control{
		Alias: alias,
	}
}

func (c *Control) Start(rcaddr, token string) (err error) {
	auto := rc.NewAutoLoginH(token)
	auto.Args = util.Map{"alias": c.Alias}
	c.R = rc.NewRC_Runner_m_j(pool.BP, rcaddr, netw.NewCCH(auto, c))
	c.R.Name = c.Alias
	auto.Runner = c.R
	c.R.Start()
	return c.R.Valid()
}

func (c *Control) StartCmd(cids, shell, cmds, logfile string) (res util.Map, err error) {
	res, err = c.R.VExec_m("start", util.Map{
		"cids":    cids,
		"shell":   shell,
		"cmds":    cmds,
		"logfile": logfile,
	})
	return
}

func (c *Control) StopCmd(cids, tid string) (res util.Map, err error) {
	res, err = c.R.VExec_m("stop", util.Map{
		"cids": cids,
		"tid":  tid,
	})
	return
}

func (c *Control) List(cids string) (res util.Map, err error) {
	res, err = c.R.VExec_m("list", util.Map{
		"cids": cids,
	})
	return
}

//OnClose see ConHandler for detail
func (c *Control) OnClose(con netw.Con) {
}

//OnCmd see ConHandler for detail
func (c *Control) OnCmd(con netw.Cmd) int {
	fmt.Println(string(con.Data()))
	return 0
}

// func (c *Control) Wait() {
// 	c.R.Wait()
// }
