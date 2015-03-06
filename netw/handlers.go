package netw

import (
	"github.com/Centny/gwf/log"
)

//the ConHandler by queue.
type QueueConH struct {
	//all handlers.
	CS []ConHandler
}

//QueueConH creator.
func NewQueueConH(cs ...ConHandler) *QueueConH {
	return &QueueConH{
		CS: cs,
	}
}

//see ConHandler
func (q *QueueConH) OnConn(c Con) bool {
	for _, cc := range q.CS {
		if cc.OnConn(c) {
			continue
		}
		return false
	}
	return true
}

//see ConHandler
func (q *QueueConH) OnClose(c Con) {
	for _, cc := range q.CS {
		cc.OnClose(c)
	}
}

//seqarate CCHandler to ConHandler and CmdHandler
type CCH struct {
	Con ConHandler
	Cmd CmdHandler
	// RCC uint64
}

//CCH creator.
func NewCCH(con ConHandler, cmd CmdHandler) *CCH {
	return &CCH{
		Con: con,
		Cmd: cmd,
	}
}

//see CCHandler
func (cch *CCH) OnConn(c Con) bool {
	return cch.Con.OnConn(c)
}

//see CCHandler
func (cch *CCH) OnClose(c Con) {
	cch.Con.OnClose(c)
}

//see CCHandler
func (cch *CCH) OnCmd(c Cmd) int {
	// atomic.AddUint64(&cch.RCC, 1)
	return cch.Cmd.OnCmd(c)
}

//do nothing handler implement CCHandler
type DoNotH struct {
	C       bool //whether allow connect
	W       bool //whether set wait to connect
	ShowLog bool
}

//DoNotH creator.
func NewDoNotH() *DoNotH {
	return &DoNotH{C: true, ShowLog: false, W: true}
}
func (cch *DoNotH) log_d(f string, args ...interface{}) {
	if cch.ShowLog {
		log.D(f, args...)
	}
}

//see ConHandler
func (cch *DoNotH) OnConn(c Con) bool {
	c.SetWait(cch.W)
	return cch.C
}

//see ConHandler
func (cch *DoNotH) OnClose(c Con) {
}

//see CmdHandler
func (cch *DoNotH) OnCmd(c Cmd) int {
	cch.log_d("DoNoH receiving command (%v)", c)
	return 0
}

//the common wait handler impl netw.ConHandler.
//it only exec SetWait to Con.
type CWH struct {
	Wait bool
}

//CWH creator.
func NewCWH(w bool) *CWH {
	return &CWH{
		Wait: w,
	}
}

//see ConHandler
func (cwh *CWH) OnConn(c Con) bool {
	if cwh.Wait {
		c.SetWait(cwh.Wait)
	}
	return true
}

//see ConHandler
func (cwh *CWH) OnClose(c Con) {
}
