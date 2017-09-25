package rcmd

import (
	"fmt"
	"strings"
	"sync"

	"github.com/Centny/gwf/log"
	"github.com/Centny/gwf/netw"
	"github.com/Centny/gwf/netw/impl"
	"github.com/Centny/gwf/netw/rc"
	"github.com/Centny/gwf/pool"
	"github.com/Centny/gwf/util"
)

var SharedMaster *Master

func StartMaster(rcaddr string, ts map[string]int) (err error) {
	SharedMaster = NewMaster()
	err = SharedMaster.Run(rcaddr, ts)
	return
}

type Master struct {
	L          *rc.RC_Listener_m
	running    map[string]int
	runningLck *sync.RWMutex
	runningSeq uint64
}

func NewMaster() *Master {
	return &Master{
		running:    map[string]int{},
		runningLck: &sync.RWMutex{},
	}
}

func (m *Master) Run(rcaddr string, ts map[string]int) (err error) {
	m.L = rc.NewRC_Listener_m_j(pool.BP, rcaddr, m)
	m.L.Name = "Master"
	m.L.LCH = m
	m.L.AddHFunc("start", m.RcStartCmdH)
	m.L.AddHFunc("stop", m.RcStopCmdH)
	m.L.AddHFunc("list", m.RcListCmdH)
	m.L.AddToken(ts)
	err = m.L.Run()
	return
}

func (m *Master) OnLogin(rc *impl.RCM_Cmd, token string) (cid string, err error) {
	parts := strings.SplitN(token, "-", 2)
	if len(parts) < 2 {
		err = fmt.Errorf("the token must having two part split by -, but %v using", token)
		return
	}
	baseCid, err := m.L.RCH.OnLogin(rc, token)
	if err != nil {
		return
	}
	cid = fmt.Sprintf("%v-%v", parts[0], baseCid)
	return
}

func (m *Master) isSlave(cid string) bool {
	return strings.HasPrefix(cid, "Slave-")
}

func (m *Master) isControl(cid string) bool {
	return strings.HasPrefix(cid, "Ctrl-")
}

func (m *Master) RcStartCmdH(rc *impl.RCM_Cmd) (res interface{}, err error) {
	var cmds string
	var logfile string
	err = rc.ValidF(`
		cmds,R|S,L:0;
		logfile,R|S,L:0;
		`, &cmds, &logfile)
	if err != nil {
		return
	}
	m.runningLck.Lock()
	m.runningSeq++
	tid := fmt.Sprintf("#%v", m.runningSeq)
	m.running[tid] = 1
	m.runningLck.Unlock()
	started := util.Map{}
	cmdCs := m.L.CmdCs()
	log.D("master try start cmd(%v) on %v connections", cmds, len(cmdCs))
	for cid, rcm := range cmdCs {
		if !m.isSlave(cid) {
			continue
		}
		alias := rcm.Kvs().StrValV("alias", cid)
		log.D("starting remote by cmds(%v),logfile(%v) to %v", cmds, logfile, alias)
		tid, execErr := rcm.Exec_s("start", util.Map{
			"cmds":    cmds,
			"logfile": logfile,
			"tid":     tid,
		})
		if execErr == nil {
			started[alias] = tid
			log.D("%v: remote command(%v) start success by id(%v) and logging to file(%v) ",
				alias, cmds, tid, logfile)
		} else {
			started[alias] = execErr.Error()
			log.W("%v: remote command(%v) start fail with %v", alias, cmds, execErr)
		}
	}
	return started, nil
}

func (m *Master) RcStopCmdH(rc *impl.RCM_Cmd) (res interface{}, err error) {
	var cid, tid string
	err = rc.ValidF(`
		cid,O|S,L:0;
		tid,R|S,L:0;
		`, &cid, &tid)
	if err != nil {
		return
	}
	cmdCs := map[string]*impl.RCM_Con{}
	if len(cid) > 0 {
		allCmdCs := m.L.CmdCs()
		for realCid, rcm := range allCmdCs {
			alias := rcm.Kvs().StrValV("alias", realCid)
			if alias != cid {
				continue
			}
			if !m.isSlave(realCid) {
				err = fmt.Errorf("can not send command to clien(%v), it is not slave", cid)
				return
			}
			cmdCs[realCid] = rcm
		}
		if len(cmdCs) < 1 {
			err = fmt.Errorf("remote client not found by id(%v)", cid)
			return
		}
	} else {
		cmdCs = m.L.CmdCs()
	}
	result := util.Map{}
	for realCid, rcm := range cmdCs {
		if !m.isSlave(realCid) {
			continue
		}
		alias := rcm.Kvs().StrValV("alias", realCid)
		log.D("stopping remote by tid(%v) to %v", tid, alias)
		_, execErr := rcm.Exec_s("stop", util.Map{
			"tid": tid,
		})
		if execErr == nil {
			result[alias] = "ok"
			log.D("%v: stop remote command by tid(%v) success", alias, tid)
		} else {
			result[alias] = execErr.Error()
			log.W("%v: stop remote command by tid(%v) fail with %v", alias, tid, execErr)
		}
	}
	return result, nil
}

func (m *Master) RcListCmdH(rc *impl.RCM_Cmd) (res interface{}, err error) {
	cmdCs := m.L.CmdCs()
	result := util.Map{}
	for cid, rcm := range cmdCs {
		if !m.isSlave(cid) {
			continue
		}
		alias := rcm.Kvs().StrValV("alias", cid)
		running, execErr := rcm.Exec_s("list", util.Map{})
		if execErr == nil {
			result[alias] = running
		} else {
			result[alias] = execErr.Error()
		}
	}
	return result, nil
}

//OnConn see ConHandler for detail
func (m *Master) OnConn(c netw.Con) bool {
	c.SetWait(true)
	return true
}

//OnClose see ConHandler for detail
func (m *Master) OnClose(c netw.Con) {
}

//OnCmd see ConHandler for detail
func (m *Master) OnCmd(c netw.Cmd) int {
	msgCs := m.L.MsgCs()
	data := c.Data()
	for cid, msg := range msgCs {
		if !m.isControl(cid) {
			continue
		}
		_, err := msg.Writeb(data)
		if err == nil {
			log.D("send messget(%v) to control(%v,%v) success", string(data), cid, msg.Kvs().StrVal("alias"))
		} else {
			log.W("send messget(%v) to control(%v,%v) fail with %v", string(data), cid, msg.Kvs().StrVal("alias"), err)
		}
	}
	return 0
}
