package rcmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sync"

	"github.com/Centny/gwf/netw"
	"github.com/Centny/gwf/netw/rc"
	"github.com/Centny/gwf/pool"
	"github.com/Centny/gwf/util"

	"github.com/Centny/gwf/log"
	"github.com/Centny/gwf/netw/impl"
)

var SharedSlave *Slave
var BASH = "/bin/bash"
var LOGFILE = "/tmp/r_%v.log"

func StartSlave(alias, rcaddr, token string) (err error) {
	SharedSlave = NewSlave(alias)
	err = SharedSlave.Start(rcaddr, token)
	return
}

type Slave struct {
	Alias      string
	R          *rc.RC_Runner_m
	running    map[string]*Task
	runningLck *sync.RWMutex
	// runningSeq uint64
}

func NewSlave(alias string) *Slave {
	return &Slave{
		Alias:      alias,
		running:    map[string]*Task{},
		runningLck: &sync.RWMutex{},
	}
}

func (s *Slave) Start(rcaddr, token string) (err error) {
	auto := rc.NewAutoLoginH(token)
	auto.Args = util.Map{"alias": s.Alias}
	s.R = rc.NewRC_Runner_m_j(pool.BP, rcaddr, netw.NewCCH(auto, netw.NewDoNotH()))
	s.R.Name = s.Alias
	s.R.AddHFunc("start", s.RcStartCmdH)
	s.R.AddHFunc("stop", s.RcStopCmdH)
	s.R.AddHFunc("list", s.RcListCmdH)
	auto.Runner = s.R
	s.R.Start()
	return
}

func (s *Slave) RcStartCmdH(rc *impl.RCM_Cmd) (res interface{}, err error) {
	var tid, cmds string
	var logfile string
	err = rc.ValidF(`
		tid,R|S,L:0;
		cmds,R|S,L:0;
		logfile,O|S,L:0;
		`, &tid, &cmds, &logfile)
	if err != nil {
		return
	}
	res = ""
	task := NewTask(tid, cmds, logfile, s)
	err = task.Start()
	return task.ID, err
}

func (s *Slave) RcStopCmdH(rc *impl.RCM_Cmd) (res interface{}, err error) {
	var tid string
	err = rc.ValidF(`
		tid,R|S,L:0;
		`, &tid)
	if err != nil {
		return
	}
	s.runningLck.RLock()
	task, ok := s.running[tid]
	s.runningLck.RUnlock()
	if ok {
		err = task.Stop()
	} else {
		err = fmt.Errorf("task is not found by tid(%v)", tid)
	}
	res = "done"
	return
}

func (s *Slave) RcListCmdH(rc *impl.RCM_Cmd) (res interface{}, err error) {
	runningIds := []string{}
	s.runningLck.RLock()
	for id := range s.running {
		runningIds = append(runningIds, id)
	}
	s.runningLck.RUnlock()
	return strings.Join(runningIds, ","), nil
}

func (s *Slave) OnTaskDone(task *Task, err error) {
	message := fmt.Sprintf("%v: task(%v) done with %v", s.R.Name, task.ID, err)
	s.R.MC.Writeb([]byte(message))
	return
}

// func (s *Slave) Wait() {
// 	s.R.Wait()
// }

type Task struct {
	ID      string
	Cmd     *exec.Cmd
	Out     *os.File
	StrCmds string
	LogFile string
	Err     error
	wait    chan int
	slave   *Slave
}

func NewTask(tid, cmds, logfile string, slave *Slave) (task *Task) {
	return &Task{
		ID:      tid,
		StrCmds: cmds,
		LogFile: logfile,
		wait:    make(chan int, 1),
		slave:   slave,
	}
}
func (t *Task) Start() (err error) {
	t.slave.runningLck.Lock()
	// t.slave.runningSeq++
	// t.ID = fmt.Sprintf("#%v", t.slave.runningSeq)
	t.slave.running[t.ID] = t
	t.slave.runningLck.Unlock()
	if len(t.LogFile) < 1 {
		t.LogFile = fmt.Sprintf(LOGFILE, t.ID)
	}
	log.I("creating task by cmds(%v) and logging to file(%v)", t.StrCmds, t.LogFile)
	t.Cmd = exec.Command(BASH, "-c", t.StrCmds)
	t.Out, err = os.OpenFile(t.LogFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, os.ModePerm)
	if err != nil {
		log.E("start task by cmds(%v) fail with open log file(%v) error:%v", t.StrCmds, t.LogFile, err)
		return
	}
	t.Cmd.Stdout, t.Cmd.Stderr = t.Out, t.Out
	err = t.Cmd.Start()
	if err != nil {
		t.Out.Close()
		t.Out = nil
		log.E("start task by cmds(%v) fail with start error:%v", t.StrCmds, err)
		return
	}
	log.I("start task(#%v) by cmds(%v) success and loggin to file(%v)", t.ID, t.StrCmds, t.LogFile)
	go func() {
		t.Err = t.Cmd.Wait()
		t.slave.OnTaskDone(t, t.Err)
		t.slave.runningLck.Lock()
		delete(t.slave.running, t.ID)
		t.slave.runningLck.Unlock()
		t.Out.Close()
		log.I("task(#%v) is done with error(%v)", t.ID, err)
		t.wait <- 1
	}()
	return
}

func (t *Task) Stop() (err error) {
	if t.Cmd != nil && t.Cmd.Process != nil {
		err = t.Cmd.Process.Kill()
		<-t.wait
	}
	return
}
