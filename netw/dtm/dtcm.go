package dtm

import (
	"fmt"
	"github.com/Centny/gwf/log"
	"github.com/Centny/gwf/netw"
	"github.com/Centny/gwf/netw/impl"
	"github.com/Centny/gwf/pool"
	"github.com/Centny/gwf/util"
	"gopkg.in/mgo.v2/bson"
	"regexp"
	"strings"
	"sync"
)

const (
	TKS_PENDING = "PENDING" //pending
	TKS_RUNNING = "RUNNING" //converting
	TKS_COV_ERR = "COV_ERR" //convert error
	TKS_DONE    = "DONE"    //done
)

//not matched error
var NOT_MATCHED = util.Err("not matched command")

//sub process
type Proc struct {
	Cid    string  `bson:"cid",json:"cid"`       //runner id
	Tid    string  `bson:"tid",json:"tid"`       //task id
	Cmds   string  `bson:"cmds",json:"cmds"`     //commands
	Done   float32 `bson:"done",json:"done"`     //the complete rate
	Msg    string  `bson:"msg",json:"msg"`       //the task exit message
	Time   int64   `bson:"time",json:"time"`     //last update time
	Status string  `bson:"status",json:"status"` //the task status
}

//task
type Task struct {
	Id   string           `bson:"_id",json:"id"`    //the id.
	Args []interface{}    `bson:"args",json:"args"` //source file
	Sid  string           `bson:"sid",json:"sid"`   //the server id.
	Proc map[string]*Proc `bson:"proc",json:"proc"` //the proc status
	Info util.Map         `bson:"info",json:"info"` //the task exit message
}

//check if task is done
func (t *Task) IsDone() bool {
	for _, proc := range t.Proc {
		if proc.Status == TKS_PENDING || proc.Status == TKS_RUNNING {
			return false
		}
	}
	return true
}

//command define
type Cmd struct {
	Name string
	Regs []*regexp.Regexp
	Cmds string
}

//check if command is matched by argumnets
func (c *Cmd) Match(args ...interface{}) bool {
	if len(args) < 1 {
		return false
	}
	for _, reg := range c.Regs {
		if reg.MatchString(fmt.Sprintf("%v", args[0])) {
			return true
		}
	}
	return false
}

//parse command by arguments.
func (c *Cmd) ParseCmd(args ...interface{}) string {
	return fmt.Sprintf(c.Cmds, args...)
}

//parse command from configure
func ParseCmds(cfg *util.Fcfg, cmds []string) ([]*Cmd, error) {
	var cmds_ []*Cmd
	for _, cmd := range cmds {
		var regs_s = cfg.Val(cmd + "/regs")
		if len(regs_s) < 1 {
			return nil, util.Err("regs is empty on command(%v)", cmd)
		}
		var regs = strings.Split(regs_s, "&")
		var regs_ []*regexp.Regexp
		for _, reg := range regs {
			var reg_, err = regexp.Compile(reg)
			if err == nil {
				regs_ = append(regs_, reg_)
				continue
			} else {
				return nil, err
			}
		}
		var cmd_s = cfg.Val(cmd + "/cmds")
		if len(cmd_s) < 1 {
			return nil, util.Err("cmds is empty on command(%v)", cmd)
		}
		cmds_ = append(cmds_, &Cmd{
			Name: cmd,
			Regs: regs_,
			Cmds: cmd_s,
		})
	}
	return cmds_, nil
}

//the database handler
type DbH interface {
	//add task to db
	Add(t *Task) error
	//update task to db
	Update(t *Task) error
	//delete task to db
	Del(t *Task) error
}

//the database creator func
type DB_C func(uri, name string) (DbH, error)

type MemH struct {
	Data map[string]*Task
}

func NewMemH() *MemH {
	return &MemH{
		Data: map[string]*Task{},
	}
}
func MemDbc(uri, name string) (DbH, error) {
	return NewMemH(), nil
}
func (m *MemH) Add(t *Task) error {
	m.Data[t.Id] = t
	return nil
}
func (m *MemH) Update(t *Task) error {
	m.Data[t.Id] = t
	return nil
}
func (m *MemH) Del(t *Task) error {
	delete(m.Data, t.Id)
	return nil
}

//the DTCM handler
type DTCM_S_H interface {
	//start event
	OnStart(dtcm *DTCM_S, task *Task)
	//done event
	OnDone(dtcm *DTCM_S, task *Task) error
}

//the distribute task control manager
type DTCM_S struct {
	*DTM_S
	*DTM_S_Proc
	H    DTCM_S_H   //the handler
	Db   DbH        //the database handler
	Sid  string     //the server id
	Cmds []*Cmd     //command list
	Cfg  *util.Fcfg //the configure
	//
	task_l   sync.RWMutex      //task lock
	tid2task map[string]*Task  //mapping task id on runner to Task.Id
	tid2proc map[string]string //mapping task id on runner to Task.Proc
}

//new DTCM_S by configure, it will be used by DTM_C as client.
/*	#Example
	[loc]
	#the server id
	sid=s1
	#the command list
	cmds=T1,T2
	#listen address
	addr=:2324
	#max command runner
	max=10
	#the db connection
	db_con=xxx
	#the db name
	db_name=xxx
	#tokens for login
	tokens=abc,a1

	#task
	[T1]
	#the regex for mathec task key
	regs=.mkv&.avi
	#the commmand to runner by format string
	cmds=${CMD_1} %v %v_1.mp4

	[T2]
	regs=.mp4&.mkv
	cmds=${CMD_2} %v %v_2.mp4 xx
*/
func NewDTCM_S(bp *pool.BytePool, cfg *util.Fcfg, dbc DB_C, h DTCM_S_H, rcm *impl.RCM_S, v2b netw.V2Byte, b2v netw.Byte2V, na impl.NAV_F) (*DTCM_S, error) {
	var dbh, err = dbc(cfg.Val("db_con"), cfg.Val("db_name"))
	if err != nil {
		return nil, err
	}
	var cmds = cfg.Val("cmds")
	if len(cmds) < 1 {
		return nil, util.Err("loc/cmds is empty")
	}
	var cmds_ []*Cmd
	cmds_, err = ParseCmds(cfg, strings.Split(cmds, ","))
	if err != nil {
		return nil, err
	}
	var sid, addr string = cfg.Val("sid"), cfg.Val("addr")
	var sh = NewDTM_S_Proc()
	var dtcm = &DTCM_S{
		DTM_S_Proc: sh,
		H:          h,
		Db:         dbh,
		Sid:        sid,
		Cmds:       cmds_,
		Cfg:        cfg,
		task_l:     sync.RWMutex{},
		tid2task:   map[string]*Task{},
		tid2proc:   map[string]string{},
	}
	var dtm = NewDTM_S(bp, addr, dtcm, rcm, v2b, b2v, na)
	dtcm.DTM_S = dtm
	var tokens = cfg.Val("tokens")
	if len(tokens) > 0 {
		dtcm.AddToken2(strings.Split(tokens, ","))
	}
	log.D("create DTCM_S by cmds(%v),tokens(%v), parsing %v commands", cmds, tokens, len(cmds_))
	return dtcm, nil
}

//new DTCM_S by json
func NewDTCM_S_j(bp *pool.BytePool, cfg *util.Fcfg, dbc DB_C, h DTCM_S_H) (*DTCM_S, error) {
	rcm := impl.NewRCM_S_j()
	return NewDTCM_S(bp, cfg, dbc, h, rcm, impl.Json_V2B, impl.Json_B2V, impl.Json_NAV)
}

//add task by info and arguments
func (d *DTCM_S) AddTask(info util.Map, args ...interface{}) error {
	var task = &Task{
		Id:   bson.NewObjectId().Hex(),
		Args: args,
		Sid:  d.Sid,
		Proc: map[string]*Proc{},
		Info: info,
	}
	for _, cmd := range d.Cmds {
		if !cmd.Match(args...) {
			continue
		}
		task.Proc[cmd.Name] = &Proc{
			Cmds:   cmd.ParseCmd(args...),
			Time:   util.Now(),
			Status: TKS_PENDING,
		}
	}
	if len(task.Proc) < 1 {
		return NOT_MATCHED
	}
	var err = d.Db.Add(task)
	if err != nil {
		log.E("DTCM_S add task by args(%v),info(%v) error->%v", args, info, err)
		return err
	}
	var current, max, res = d.do_task(task)
	if res == 0 {
		log.D("DTCM_S add task success by %v matched, task(%v) will be running, current(%v)/max(%v)",
			len(task.Proc), task.Id, current, max)
	} else if res == 1 {
		log.D("DTCM_S add task success by %v matched, but runner is busy now on current(%v)/max(%v), task(%v) will be pending",
			len(task.Proc), current, max, task.Id)
	} else {
		log.W("DTCM_S add task having error(code:%v) by %v matched, running status is current(%v)/max(%v), task(%v) will be pending",
			res, len(task.Proc), current, max, task.Id)
	}
	return nil
}

//do task, it will check running
func (d *DTCM_S) do_task(t *Task) (int, int, int) {
	d.task_l.Lock()
	var max = d.Cfg.IntValV("max", 100)
	var current = d.DTM_S_Proc.Total()
	if max < current+len(t.Proc) {
		d.task_l.Unlock()
		return current, max, 1
	}
	d.start_task(t)
	var err = d.Db.Update(t)
	if err == nil {
		d.H.OnStart(d, t)
		d.task_l.Unlock()
		return current, max, 0
	} else {
		log.E("DTCM_S update task(%v) error->%v", t.Id, err)
		d.task_l.Unlock()
		d.stop_task(t)
		return current, max, 2
	}
}

//start task
func (d *DTCM_S) start_task(t *Task) {
	var err error
	for cmd, proc := range t.Proc {
		proc.Cid, proc.Tid, err = d.StartTask3(proc.Cmds)
		if err == nil {
			d.tid2task[proc.Tid] = t
			d.tid2proc[proc.Tid] = cmd
			proc.Msg = ""
			proc.Status = TKS_RUNNING
			log.D("DTCM_S start task success on cid(%v),tid(%v) by cmds(%v)", proc.Cid, proc.Tid, proc.Cmds)
		} else {
			proc.Cid = ""
			proc.Tid = ""
			proc.Msg = fmt.Sprintf("start task error->%v", err)
			proc.Status = TKS_COV_ERR
			log.E("DTCM_S start task error->%v", err)
		}
	}
}

//stop task
func (d *DTCM_S) stop_task(t *Task) {
	var err error
	for _, proc := range t.Proc {
		if len(proc.Tid) < 1 {
			continue
		}
		err = d.StopTask(proc.Cid, proc.Tid)
		if err != nil {
			proc.Msg = "STOPPED"
			proc.Status = TKS_COV_ERR
			log.E("DTCM_S stop task error->%v", err)
		}
		delete(d.tid2task, proc.Tid)
		delete(d.tid2proc, proc.Tid)
	}
	log.D("DTCM_S stop task(%v)", t.Id)
}

//process event
func (d *DTCM_S) OnProc(dtm *DTM_S, cid, tid string, rate float64) {
	d.DTM_S_Proc.OnProc(dtm, cid, tid, rate)
}

//start event
func (d *DTCM_S) OnStart(dtm *DTM_S, cid, tid, cmds string) {
	d.DTM_S_Proc.OnStart(dtm, cid, tid, cmds)
}

//stop event
func (d *DTCM_S) OnStop(dtm *DTM_S, cid, tid string) {
	d.DTM_S_Proc.OnStop(dtm, cid, tid)
	d.task_l.Lock()
	defer d.task_l.Unlock()
	d.mark_done(cid, tid, "STOPPED", TKS_COV_ERR)
}

//done event
func (d *DTCM_S) OnDone(dtm *DTM_S, cid, tid string, code int, err string, used int64) {
	d.DTM_S_Proc.OnDone(dtm, cid, tid, code, err, used)
	d.task_l.Lock()
	defer d.task_l.Unlock()
	if code == 0 {
		d.mark_done(cid, tid, "", TKS_DONE)
	} else {
		d.mark_done(cid, tid, fmt.Sprintf("done error (code:%v,err:%v)", code, err), TKS_COV_ERR)
	}
}

//mark task done
func (d *DTCM_S) mark_done(cid, tid, msg, status string) {
	var task = d.tid2task[tid]
	if task == nil {
		log.E("DTCM_S stop task error(not found) by tid(%v)", tid)
		return
	}
	log.D("DTCM_S runner(%v/%v) is stopped on task(%v)", tid, cid, task.Id)
	var proc = task.Proc[d.tid2proc[tid]]
	proc.Cid = ""
	proc.Tid = ""
	proc.Time = util.Now()
	proc.Msg = msg
	proc.Status = status
	delete(d.tid2task, tid)
	delete(d.tid2proc, tid)
	var rerr = d.Db.Update(task)
	if rerr == nil {
		d.check_done(task)
	} else {
		log.E("DTCM_S update task error by %v ->%v", util.S2Json(task), rerr)
	}
}

//check done
func (d *DTCM_S) check_done(task *Task) {
	if !task.IsDone() {
		return
	}
	var err = d.H.OnDone(d, task)
	if err != nil {
		log.E("DTCM_S on done error by %v ->%v", util.S2Json(task), err)
		return
	}
	err = d.Db.Del(task)
	if err != nil {
		log.E("DTCM_S delete task error by %v ->%v", util.S2Json(task), err)
	}
}
