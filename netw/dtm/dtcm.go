package dtm

import (
	"fmt"
	"github.com/Centny/gwf/log"
	"github.com/Centny/gwf/netw"
	"github.com/Centny/gwf/netw/impl"
	"github.com/Centny/gwf/pool"
	"github.com/Centny/gwf/routing"
	"github.com/Centny/gwf/util"
	"math"
	"reflect"
	"regexp"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

const (
	TKS_PENDING  = "PENDING"  //pending
	TKS_RUNNING  = "RUNNING"  //converting
	TKS_COV_ERR  = "COV_ERR"  //convert error
	TKS_DONE     = "DONE"     //done
	TKS_DONE_ERR = "DONE_ERR" //done error
)

//not matched error
// var NOT_MATCHED = util.Err("not matched command")

type NotMatchedErr struct {
	Msg string
}

func NewNotMatchedErr(format string, args ...interface{}) *NotMatchedErr {
	return &NotMatchedErr{
		Msg: fmt.Sprintf(format, args...),
	}
}
func (n *NotMatchedErr) Error() string {
	return n.Msg
}

func IsNotMatchedErr(err error) bool {
	return reflect.Indirect(reflect.ValueOf(err)).Type().Name() == "NotMatchedErr"
}

//sub process
type Proc struct {
	Cid    string      `bson:"cid" json:"cid"`       //runner id
	Tid    string      `bson:"tid" json:"tid"`       //task id
	Cmds   string      `bson:"cmds" json:"cmds"`     //commands
	Done   float32     `bson:"done" json:"done"`     //the complete rate
	Code   int         `bson:"code" json:"code"`     //the task exit code
	Msg    string      `bson:"msg" json:"msg"`       //the task exit message
	Try    int32       `bson:"try" json:"try"`       //the task try time
	Time   int64       `bson:"time" json:"time"`     //last update time
	Status string      `bson:"status" json:"status"` //the task status
	Res    interface{} `bson:"res" json:"res"`       //the task status
}

//task
type Task struct {
	Id     string           `bson:"_id" json:"id"`        //the id.
	Args   []interface{}    `bson:"args" json:"args"`     //source file
	Sid    string           `bson:"sid" json:"sid"`       //the server id.
	Proc   map[string]*Proc `bson:"proc" json:"proc"`     //the proc status
	Info   interface{}      `bson:"info" json:"info"`     //the task exit message
	Mid    string           `bson:"mid" json:"mid"`       //the running machine id
	Runc   int              `bson:"runc" json:"runc"`     //the run times
	Next   int64            `bson:"next" json:"next"`     //last update time
	Time   int64            `bson:"time" json:"time"`     //last update time
	Code   int              `bson:"code" json:"code"`     //the task exit code
	Status string           `bson:"status" json:"status"` //the task try time
}

//check if task is done
func (t *Task) IsDone() bool {
	for _, proc := range t.Proc {
		if proc.Status == TKS_PENDING || proc.Status == TKS_RUNNING || proc.Status == TKS_COV_ERR {
			return false
		}
	}
	return true
}

func (t *Task) CheckDone() (int, bool) {
	var max_try int32 = -1
	for _, proc := range t.Proc {
		if proc.Status == TKS_PENDING || proc.Status == TKS_RUNNING {
			return 0, false
		}
		if proc.Status == TKS_COV_ERR {
			if max_try < proc.Try {
				max_try = proc.Try
			}
		}
		if proc.Code == 413 {
			max_try = math.MaxInt32
		}
	}
	return int(max_try), max_try < 0
}

func (t *Task) IsRunning() bool {
	for _, proc := range t.Proc {
		if proc.Status == TKS_RUNNING {
			return true
		}
	}
	return false
}

func ParseRegs(sec string, cfg *util.Fcfg) ([]*regexp.Regexp, error) {
	var regs_s = cfg.Val(sec + "/regs")
	if len(regs_s) < 1 {
		return nil, util.Err("%v/regs is empty on client(%v)", sec, sec)
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
	return regs_, nil
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
	var cfg = util.NewFcfg3()
	for idx, arg := range args {
		cfg.SetVal(fmt.Sprintf("v%v", idx), fmt.Sprintf("%v", arg))
	}
	return cfg.EnvReplaceV(c.Cmds, false)
}

//parse command from configure
func ParseCmds(cfg *util.Fcfg, cmds []string) ([]*Cmd, error) {
	var cmds_ []*Cmd
	for _, cmd := range cmds {
		regs_, err := ParseRegs(cmd, cfg)
		if err != nil {
			return nil, util.Err("parse regs on command(%v) error->%v", cmd, err)
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

type Client struct {
	Name  string           `json:"name"`
	Regs  []*regexp.Regexp `json:"-"`
	Max   int              `json:"max"`
	Os    string           `json:"os"`
	Token map[string]int   `json:"token"`
}

func (c *Client) Match(args ...interface{}) bool {
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

//parse command from configure
func ParseClients(cfg *util.Fcfg, clients []string) (map[string]*Client, error) {
	var clients_ = map[string]*Client{}
	for _, client := range clients {
		regs_, err := ParseRegs(client, cfg)
		if err != nil {
			return nil, util.Err("parse regs on client(%v) error->%v", client, err)
		}
		var token_s = cfg.Val(client + "/token")
		if len(token_s) < 1 {
			return nil, util.Err("%v/tokens is empty on client(%v)", client, client)
		}
		var token_m = map[string]int{}
		for idx, token := range strings.Split(token_s, ",") {
			token_m[token] = idx + 1
		}
		if _, ok := clients_[client]; ok {
			return nil, util.Err("client by name(%v) is repeat", client)
		}
		clients_[client] = &Client{
			Name:  client,
			Regs:  regs_,
			Token: token_m,
			Os:    cfg.Val("c_os"),
			Max:   cfg.IntValV(client+"/max", 8),
		}
	}
	return clients_, nil
}

func ParseAbsL(cfg *util.Fcfg, abs_l []string) ([]Abs, error) {
	var abs = []Abs{}
	for _, sec := range abs_l {
		var creator, err = CreateAbs(sec, cfg)
		if err != nil {
			return abs, err
		}
		abs = append(abs, creator)
	}
	return abs, nil
}

//the database handler
type DbH interface {
	//add task to db
	Add(t *Task) error
	//update task to db
	Update(t *Task) error
	//delete task to db
	Del(t *Task) error
	//list task from db
	List(running []string, status string, skip, limit int) (int, []*Task, error)
	//find task
	Find(id string) (*Task, error)
	//
	ClearSyncTask() error
}

//the database creator func
type DB_C func(uri, name string) (DbH, error)

type MemH struct {
	Errs map[string]error
	Data map[string]*Task
}

func NewMemH() *MemH {
	return &MemH{
		Errs: map[string]error{},
		Data: map[string]*Task{},
	}
}
func MemDbc(uri, name string) (DbH, error) {
	return NewMemH(), nil
}
func (m *MemH) Add(t *Task) error {
	m.Data[t.Id] = t
	return m.Errs["Add"]
}
func (m *MemH) Update(t *Task) error {
	m.Data[t.Id] = t
	return m.Errs["Update"]
}
func (m *MemH) Del(t *Task) error {
	delete(m.Data, t.Id)
	return m.Errs["Del"]
}
func (m *MemH) List(running []string, status string, skip, limit int) (int, []*Task, error) {
	var ts []*Task
	for _, task := range m.Data {
		if task.Status == status || len(status) < 1 {
			ts = append(ts, task)
		}
	}
	return len(m.Data), ts, m.Errs["List"]
}
func (m *MemH) Find(id string) (*Task, error) {
	return m.Data[id], m.Errs["Find"]
}
func (m *MemH) ClearSyncTask() error {
	return nil
}

type DoNoneH struct {
}

func NewDoNoneH() *DoNoneH {
	return &DoNoneH{}
}
func (d *DoNoneH) OnStart(dtcm *DTCM_S, task *Task) {
	log.D("DoNoneH task(%v) is started", task.Id)
}
func (d *DoNoneH) OnDone(dtcm *DTCM_S, task *Task) error {
	log.D("DoNoneH task(%v) is done with->\n%v\n", task.Id, util.S2Json(task))
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
	H       DTCM_S_H           //the handler
	Db      DbH                //the database handler
	Sid     string             //the server id
	Cmds    []*Cmd             //command list
	LocCmds []*Cmd             //the local command list
	AbsL    []Abs              //the argument builder list
	LocAbsL []Abs              //the argument builder list
	Clients map[string]*Client //client list

	Cfg   *util.Fcfg //the configure
	Local *util.Fcfg
	T2C   map[string]*Client //mapping token to clients
	//
	task_l   sync.RWMutex      //task lock
	tasks    map[string]*Task  //mapping Task.Id to Task
	tid2task map[string]*Task  //mapping task id on runner to Task.Id
	tid2proc map[string]string //mapping task id on runner to Task.Proc
	//
	running bool
	run_c   chan int
	//
	sequence uint64
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
func NewDTCM_S(bp *pool.BytePool, cfg *util.Fcfg, dbc DB_C, h DTCM_S_H, nd impl.ND_F, vna impl.VNA_F, v2b netw.V2Byte, b2v netw.Byte2V, nav impl.NAV_F) (*DTCM_S, error) {
	var dbh, err = dbc(cfg.Val("db_con"), cfg.Val("db_name"))
	if err != nil {
		return nil, err
	}
	dbh.ClearSyncTask()
	//
	var cmds = cfg.Val("cmds")
	if len(cmds) < 1 {
		return nil, util.Err("loc/cmds is empty")
	}
	var cmds_ []*Cmd
	var cmds_s = strings.Split(cmds, ",")
	log.D("DTCM_S parsing cmds by names(%v)", cmds_s)
	cmds_, err = ParseCmds(cfg, cmds_s)
	if err != nil {
		return nil, err
	}
	//
	var loc_cmds_ []*Cmd
	var loc_cmds = cfg.Val("loc_cmds")
	if len(loc_cmds) > 0 {
		var loc_cmds_s = strings.Split(loc_cmds, ",")
		log.D("DTCM_S parsing local cmds by names(%v)", loc_cmds_s)
		loc_cmds_, err = ParseCmds(cfg, loc_cmds_s)
		if err != nil {
			return nil, err
		}
	}
	//
	var clients = cfg.Val("clients")
	if len(clients) < 1 {
		return nil, util.Err("loc/clients is empty")
	}
	var clients_ map[string]*Client
	var clients_s = strings.Split(clients, ",")
	log.D("DTCM_S parsing clients by names(%v)", clients_s)
	clients_, err = ParseClients(cfg, clients_s)
	if err != nil {
		return nil, err
	}
	//
	var abs = []Abs{}
	var abs_l = cfg.Val("abs_l")
	if len(abs_l) > 0 {
		abs, err = ParseAbsL(cfg, strings.Split(abs_l, ","))
		if err != nil {
			return nil, err
		}
	}
	//
	var loc_abs = []Abs{}
	var loc_abs_l = cfg.Val("loc_abs_l")
	if len(loc_abs_l) > 0 {
		loc_abs, err = ParseAbsL(cfg, strings.Split(loc_abs_l, ","))
		if err != nil {
			return nil, err
		}
	}
	//
	var sid, addr string = cfg.Val("sid"), cfg.Val("addr")
	var sh = NewDTM_S_Proc()
	var dtcm = &DTCM_S{
		DTM_S_Proc: sh,
		H:          h,
		Db:         dbh,
		Sid:        sid,
		Cmds:       cmds_,
		LocCmds:    loc_cmds_,
		AbsL:       abs,
		LocAbsL:    loc_abs,
		Clients:    clients_,
		Cfg:        cfg,
		Local:      cfg.Strip("Local"),
		T2C:        map[string]*Client{},
		task_l:     sync.RWMutex{},
		tasks:      map[string]*Task{},
		tid2task:   map[string]*Task{},
		tid2proc:   map[string]string{},
		run_c:      make(chan int, 1),
	}
	var dtm = NewDTM_S(bp, addr, dtcm, nd, vna, v2b, b2v, nav)
	dtcm.DTM_S = dtm
	var tokens = []string{}
	for _, client := range clients_ {
		for token, v := range client.Token {
			if oc, ok := dtcm.T2C[token]; ok {
				return nil, util.Err("token(%v) is repeat on clinet(%v,%v)", token, oc.Name, client.Name)
			}
			dtcm.AddToken3(token, v)
			dtcm.T2C[token] = client
			tokens = append(tokens, token)
		}
	}
	dtcm.Name = "DTCM_S"
	log.D("create DTCM_S by cmds(%v),clients(%v),tokens(%v) parsing %v commands", cmds, clients, tokens, len(cmds_))
	return dtcm, nil
}

//new DTCM_S by json
func NewDTCM_S_j(bp *pool.BytePool, cfg *util.Fcfg, dbc DB_C, h DTCM_S_H) (*DTCM_S, error) {
	// rcm := impl.NewRCM_S_j()
	return NewDTCM_S(bp, cfg, dbc, h, impl.Json_ND, impl.Json_VNA, impl.Json_V2B, impl.Json_B2V, impl.Json_NAV)
}

func (d *DTCM_S) NewTask(id, info interface{}, args ...interface{}) *Task {
	var task = &Task{
		Id:     fmt.Sprintf("%v", id),
		Args:   args,
		Sid:    d.Sid,
		Proc:   map[string]*Proc{},
		Info:   info,
		Time:   util.Now(),
		Status: TKS_RUNNING,
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
	return task
}

func (d *DTCM_S) MatchArgs(info interface{}, args ...interface{}) bool {
	if len(args) < 1 {
		return false
	}
	return d.MatchArgsV(args[0], info, args...)
}
func (d *DTCM_S) MatchArgsV(id, info interface{}, args ...interface{}) bool {
	for _, abs := range d.AbsL {
		if abs.Match(d, id, info, args...) {
			return true
		}
	}
	return false
}

//add task by info and arguments
func (d *DTCM_S) AddTask(info interface{}, args ...interface{}) error {
	if len(args) < 1 {
		return util.Err("at least one argumnet is setted")
	}
	return d.AddTaskV(args[0], info, args...)
}
func (d *DTCM_S) BuildArgs(id, info interface{}, args ...interface{}) (interface{}, interface{}, []interface{}, error) {
	for _, abs := range d.AbsL {
		if abs.Match(d, id, info, args...) {
			return abs.Build(d, id, info, args...)
		}
	}
	return nil, nil, nil, NewNotMatchedErr("DTCM_S not abs matched by id(%v),info(%v),args(%v)", id, util.S2Json(info), util.S2Json(args))
}
func (d *DTCM_S) AddTaskV(id, info interface{}, args ...interface{}) error {
	var err error
	id, info, args, err = d.BuildArgs(id, info, args...)
	if err != nil {
		log.E("%v", err)
		return err
	}
	task, err := d.Db.Find(fmt.Sprintf("%v", id))
	if err != nil {
		err = util.Err("DTCM_S find task by id(%v) error->%v",
			id, err)
		return err
	}
	if task != nil {
		err = util.Err("DTCM_S add task fail by id(%v),args(%v)->the task is already exist->%v",
			id, util.S2Json(args), util.S2Json(task))
		log.E("%v", err)
		return err
	}
	task = d.NewTask(id, info, args...)
	if len(task.Proc) < 1 {
		err = NewNotMatchedErr("DTCM_S not command matched by id(%v),info(%v),args(%v)", id, util.S2Json(info), util.S2Json(args))
		// log.E("%v", err)
		return err
	}
	var current, max, start = d.chk_task(task)
	if start {
		task.Mid = util.MID()
	}
	err = d.Db.Add(task)
	if err != nil {
		log.E("DTCM_S add task by args(%v),info(%v) error->%v", args, info, err)
		return err
	}
	if !start {
		log.D("DTCM_S add task success by %v matched, but runner is busy now on current(%v)/max(%v)/clients(%v), task(%v) will be pending",
			len(task.Proc), current, max, len(d.TaskC), task.Id)
		return nil
	}
	var res = d.do_task(task)
	if res == 0 {
		log.D("DTCM_S add task success by %v matched, task(%v) will be running, current(%v)/max(%v)/clients(%v)",
			len(task.Proc), task.Id, current, max, len(d.TaskC))
	} else {
		log.W("DTCM_S add task having error(code:%v) by %v matched, running status is current(%v)/max(%v)/clients(%v), task(%v) will be pending",
			res, len(task.Proc), current, max, len(d.TaskC), task.Id)
	}
	return nil
}

func (d *DTCM_S) AddTaskH(hs *routing.HTTPSession) routing.HResult {
	var tid string = hs.RVal("tid")
	var args_s string = hs.RVal("args")
	if len(args_s) < 1 {
		return hs.MsgResE3(2, "arg-err", "args argument is empty")
	}
	var args_a = []interface{}{}
	for _, arg := range strings.Split(args_s, ",") {
		args_a = append(args_a, arg)
	}
	if len(tid) < 1 {
		tid = fmt.Sprintf("%v", args_a[0])
	}
	var err = d.AddTaskV(tid, nil, args_a...)
	if err == nil {
		return hs.MsgRes("OK")
	} else {
		err = util.Err("AddTask error->%v", err)
		log.E("%v", err)
		return hs.MsgResErr2(3, "srv-err", err)
	}
}

func (d *DTCM_S) chk_task(t *Task) (int, int, bool) {
	var max = d.Cfg.IntValV("max", 100)
	var current = d.DTM_S_Proc.Total()
	return current, max, max > current+len(t.Proc)
}

// //do task, it will check running
// func (d *DTCM_S) do_task(t *Task) (int, int, int) {
// 	d.task_l.Lock()
// 	defer d.task_l.Unlock()
// 	return d.do_task_(t)
// }
func (d *DTCM_S) do_task(t *Task) int {
	d.start_task(t)
	if !t.IsRunning() {
		t.Mid = ""
		t.Next = util.Now() + int64(math.Pow(2, float64(2*t.Runc-1))*1000)
	}
	var err = d.Db.Update(t)
	if err == nil {
		d.H.OnStart(d, t)
		return 0
	} else {
		log.E("DTCM_S update task(%v) error->%v", t.Id, err)
		go d.stop_task(t)
		return 2
	}
}

func (d *DTCM_S) min_used_cid(t *Task, proc *Proc) (string, string) {
	var tcid string = ""
	var min int = 999
	var nm_c, nm_m = 0, 0
	for cid, tc := range d.TaskC {
		var msgc = d.MsgC(cid)
		if msgc == nil {
			log.E("DTCM_S client not found by id(%v)", cid)
			continue
		}
		var token = msgc.Kvs().StrVal("token")
		var client = d.T2C[token]
		if !client.Match(t.Args...) {
			nm_c += 1
			continue
		}
		if tc >= client.Max {
			nm_m += 1
			continue
		}
		if tc < min {
			tcid = cid
			min = tc
		}
	}
	if len(tcid) > 0 {
		//matched
		return tcid, ""
	}
	if nm_m > 0 {
		//client matched, but busy
		return "", "all client is busy"
	} else {
		//not matched client
		return "", "not matched client"
	}
}

//start task
func (d *DTCM_S) start_task(t *Task) int {
	var err error
	var msg string
	var busy int = 0
	t.Runc += 1
	for cmd, proc := range t.Proc {
		if proc.Status == TKS_DONE {
			continue
		}
		d.task_l.Lock()
		proc.Cid, msg = d.min_used_cid(t, proc)
		if len(proc.Cid) < 1 {
			proc.Status = TKS_PENDING
			proc.Cid = ""
			log.D("DTCM_S select min used client fail with %v by args(%v), process will be pending", msg, util.S2Json(t.Args))
			busy += 1
			d.task_l.Unlock()
			continue
		}
		proc.Tid = fmt.Sprintf("T-%v", atomic.AddUint64(&d.sequence, 1))
		d.tasks[t.Id] = t
		d.tid2task[proc.Tid] = t
		d.tid2proc[proc.Tid] = cmd
		proc.Msg = ""
		proc.Status = TKS_RUNNING
		d.task_l.Unlock()
		log.D("DTCM_S starting runner(%v/%v) on task(%v) by cmds(\n\t%v\n)", proc.Tid, proc.Cid, t.Id, proc.Cmds)
		err = d.StartTask(proc.Cid, proc.Tid, proc.Cmds)
		if err == nil {
			log.D("DTCM_S start runner(%v/%v) success on task(%v) by cmds(\n\t%v\n)", proc.Tid, proc.Cid, t.Id, proc.Cmds)
		} else {
			d.task_l.Lock()
			delete(d.tasks, t.Id)
			delete(d.tid2task, proc.Tid)
			delete(d.tid2proc, proc.Tid)
			proc.Cid = ""
			proc.Tid = ""
			proc.Msg = fmt.Sprintf("start task error->%v", err)
			proc.Status = TKS_COV_ERR
			d.task_l.Unlock()
			log.E("DTCM_S start task error by %v->%v", proc.Cmds, err)
		}
	}
	return busy
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
	delete(d.tasks, t.Id)
	log.D("DTCM_S stop task(%v)", t.Id)
}

func (d *DTCM_S) OnLogin(rc *impl.RCM_Cmd, token string) (string, error) {
	var cid, err = d.DTM_S_Proc.OnLogin(rc, token)
	if err == nil {
		log.D("DTCM_S login by token(%v) as client(%v)", token, util.S2Json(d.T2C[token]))
	}
	return cid, err
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
	d.mark_done(nil, cid, tid, 0, "STOPPED", TKS_COV_ERR)
}

//done event
func (d *DTCM_S) OnDone(dtm *DTM_S, args util.Map, cid, tid string, code int, err string, used int64) {
	d.DTM_S_Proc.OnDone(dtm, args, cid, tid, code, err, used)
	d.task_l.Lock()
	defer d.task_l.Unlock()
	if code == 0 {
		d.mark_done(args, cid, tid, code, "", TKS_DONE)
	} else {
		d.mark_done(args, cid, tid, code, fmt.Sprintf("done error (code:%v,err:%v)", code, err), TKS_COV_ERR)
	}
	go d.do_checker_(1)
	slog("DTCM_S done success with cid(%v),tid(%v),code(%v),err(%v),used(%v)", cid, tid, code, err, used)
}

func (d *DTCM_S) OnClose(c netw.Con) {
	d.task_l.Lock()
	defer d.task_l.Unlock()
	defer d.DTM_S_Proc.OnClose(c)
	var cid = c.Kvs().StrVal("cid")
	var trs = d.Rates[cid]
	if trs == nil || len(trs) < 1 {
		return
	}
	log.W("DTCM_S client(%v) is closed, will clear %v running process", cid, len(trs))
	var tids = map[string][]string{}
	var tasks = map[string]*Task{}
	for tid, _ := range trs {
		var task = d.tid2task[tid]
		if task == nil {
			continue
		}
		tasks[task.Id] = task
		tids[task.Id] = append(tids[task.Id], tid)
	}
	for tid, task := range tasks {
		d.mark_done_v(nil, cid, task, tids[tid], 0, "Runner Closed", TKS_COV_ERR)
	}
}

//mark task done
func (d *DTCM_S) mark_done(res interface{}, cid, tid string, code int, msg, status string) {
	var task = d.tid2task[tid]
	if task == nil {
		log.E("DTCM_S mark done task error(not found) by tid(%v)", tid)
		return
	}
	log.D("DTCM_S runner(%v/%v) is done on task(%v)", tid, cid, task.Id)
	d.mark_done_v(res, cid, task, []string{tid}, code, msg, status)
}
func (d *DTCM_S) mark_done_v(res interface{}, cid string, task *Task, tids []string, code int, msg, status string) {
	for _, tid := range tids {
		var proc = task.Proc[d.tid2proc[tid]]
		if proc == nil {
			continue
		}
		proc.Cid = ""
		proc.Tid = ""
		proc.Time = util.Now()
		proc.Code = code
		proc.Msg = msg
		proc.Status = status
		proc.Res = res
		atomic.AddInt32(&proc.Try, 1)
		delete(d.tid2task, tid)
		delete(d.tid2proc, tid)
	}
	task.Code = code
	if !task.IsRunning() {
		task.Mid = ""
		task.Next = util.Now() + 3 ^ int64(2*task.Runc-1)*1000
	}
	var rerr = d.Db.Update(task)
	if rerr == nil {
		d.check_done(task)
		return
	}
	log.E("DTCM_S update task error by %v ->%v", util.S2Json(task), rerr)
	if !task.IsRunning() {
		delete(d.tasks, task.Id)
		log.W("DTCM_S remove running task(%v) for update task error(%v), it will move to pending pool", util.S2Json(task), rerr)
	}
}

//check done
func (d *DTCM_S) check_done(task *Task) {
	if task.IsRunning() {
		return
	}
	delete(d.tasks, task.Id)
	mtry, done := task.CheckDone()
	if done {
		d.do_done(task)
		return
	}
	cmtry := d.Cfg.IntValV("max_try", 100)
	if mtry < cmtry {
		return
	}
	log.E("DTCM_S convert task(%v) error->too many try times(%v/%v)", util.S2Json(task), mtry, cmtry)
	task.Status = TKS_COV_ERR
	var rerr = d.Db.Update(task)
	if rerr != nil {
		log.E("DTCM_S update task error by %v ->%v", util.S2Json(task), rerr)
	}
}

//do done
func (d *DTCM_S) do_done(task *Task) {
	var err = d.H.OnDone(d, task)
	if err == nil {
		err = d.Db.Del(task)
		if err != nil {
			log.E("DTCM_S delete task error by %v ->%v", util.S2Json(task), err)
		}
		return
	}
	log.E("DTCM_S on done error by %v ->%v", util.S2Json(task), err)
	task.Status = TKS_DONE_ERR
	var rerr = d.Db.Update(task)
	if rerr != nil {
		log.E("DTCM_S update task error by %v ->%v", util.S2Json(task), rerr)
	}
}

//start checker
func (d *DTCM_S) StartChecker(delay int64) {
	go d.loop_checker(delay)
}

//stop checker
func (d *DTCM_S) StopChecker() {
	d.running = false
	<-d.run_c
}

//loop checker
func (d *DTCM_S) loop_checker(delay int64) {
	d.running = true
	for d.running {
		d.do_checker()
		var tdelay = delay
		for tdelay > 0 && d.running {
			time.Sleep(200 * time.Millisecond)
			tdelay -= 200
		}
	}
	log.D("DTCM_S the checker loop is done")
	d.run_c <- 0
}

//do checker
func (d *DTCM_S) do_checker() {
	// var max = d.Cfg.IntValV("max", 100)
	d.do_checker_(1)
}

func (d *DTCM_S) do_checker_(max int) {
	defer func() {
		var err = recover()
		if err != nil {
			log.E("DTCM_S receive panic error->%v", err)
		}
	}()
	if len(d.TaskC) < 1 {
		log.D("DTCM_S do check success, but active client is empty")
		return
	}
	d.task_l.Lock()
	var rids = []string{}
	for rid, _ := range d.tasks {
		rids = append(rids, rid)
	}
	d.task_l.Unlock()
	total, ts, err := d.Db.List(rids, TKS_RUNNING, 0, max)
	if err != nil {
		log.E("DTCM_S do check error->%v", err)
		return
	}
	if len(ts) < 1 {
		// log.D("DTCM_S do check succes and task is empty")
		return
	}
	log.D("DTCM_S do check succes and %v task found, will try run %v task", total, max)
	for _, task := range ts {
		if _, ok := d.tasks[task.Id]; ok {
			continue
		}
		if task.IsDone() {
			d.task_l.Lock()
			d.do_done(task)
			d.task_l.Unlock()
			continue
		}
		var busy = d.start_task(task)
		if busy > 0 {
			log.D("DTCM_S all runner is busy, checker is break")
			break
		}
	}
	log.D("DTCM_S do check done...")
}

func (d *DTCM_S) SrvHTTP(hs *routing.HTTPSession) routing.HResult {
	var total, ts, err = d.Db.List(nil, hs.RVal("status"), 0, d.Cfg.IntValV("max", 100))
	var err_msg = ""
	if err != nil {
		err_msg = err.Error()
		log.E("DTCM_S list fail with error(%v)", err)
	}
	var res = util.Map{
		"proc":      d.DTM_S_Proc,
		"total":     total,
		"running_c": len(d.tasks),
		"err":       err_msg,
	}
	if hs.RVal("tasks") == "1" {
		res["tasks"] = ts
	}
	if hs.RVal("running") == "1" {
		res["running"] = d.tasks
	}
	return hs.JRes(res)
}

func (d *DTCM_S) Hand(pre string, mux *routing.SessionMux) {
	mux.HFunc("^"+pre+"/addTask(\\?.*)?", d.AddTaskH)
	mux.H("^"+pre+"/status(\\?.*)?", d)
}

func (d *DTCM_S) TaskRate(id string) (float64, map[string]float64, error) {
	var task = d.tasks[id]
	if task == nil {
		return 0, nil, util.Err("DTCM_S task is not found by id(%v)", id)
	}
	var total float64 = 0
	var res = map[string]float64{}
	for key, proc := range task.Proc {
		var rate = d.DTM_S.H.Rate(d.DTM_S, proc.Cid, proc.Tid)
		total += rate
		res[key] = rate
	}
	if len(res) > 0 {
		total = total / float64(len(res))
	}
	return total, res, nil
}
