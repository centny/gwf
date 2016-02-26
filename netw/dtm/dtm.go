//Package dtm provider distributed task control manager.
//
//the DTM_S/DTCM_S is the server side impl.
//
//the DTM_C is the client side impl.
package dtm

import (
	"bytes"
	"fmt"
	"github.com/Centny/gwf/log"
	"github.com/Centny/gwf/netw"
	"github.com/Centny/gwf/netw/impl"
	"github.com/Centny/gwf/netw/rc"
	"github.com/Centny/gwf/pool"
	"github.com/Centny/gwf/routing"
	"github.com/Centny/gwf/util"
	"os"
	"os/exec"
	"strings"
	"sync"
	"sync/atomic"
)

const (
	CMD_M_PROC = 10
	CMD_M_DONE = 20
)

//the DTM handler
type DTM_S_H interface {
	rc.RC_Login_h
	netw.ConHandler
	//process event
	OnProc(d *DTM_S, cid, tid string, rate float64)
	//start event
	OnStart(d *DTM_S, cid, tid, cmds string)
	//stop event
	OnStop(d *DTM_S, cid, tid string)
	//done event
	OnDone(d *DTM_S, args util.Map, cid, tid string, code int, err string, used int64)
	//check and return minial used client id
	MinUsedCid(d *DTM_S, args ...interface{}) string
}

//the default DTM process handler
type DTM_S_Proc struct {
	Rates map[string]map[string]float64 `json:"rates"`  //process rate
	AllC  int                           `json:"all_c"`  //all client count
	TaskC map[string]int                `json:"task_c"` //client count by id
	//
	proc_l sync.RWMutex
	cid    int64
}

//new the default DTM process handler
func NewDTM_S_Proc() *DTM_S_Proc {
	return &DTM_S_Proc{
		Rates:  map[string]map[string]float64{},
		AllC:   0,
		TaskC:  map[string]int{},
		proc_l: sync.RWMutex{},
		cid:    0,
	}
}

//process event
func (d *DTM_S_Proc) OnProc(dtm *DTM_S, cid, tid string, rate float64) {
	if _, ok := d.Rates[cid]; ok {
		d.Rates[cid][tid] = rate
	}
}

//start event
func (d *DTM_S_Proc) OnStart(dtm *DTM_S, cid, tid, cmds string) {
	d.proc_l.Lock()
	defer d.proc_l.Unlock()
	if _, ok := d.Rates[cid]; !ok {
		d.Rates[cid] = map[string]float64{}
	}
	d.Rates[cid][tid] = 0
	d.TaskC[cid] += 1
	d.AllC += 1
}

//stop event
func (d *DTM_S_Proc) OnStop(dtm *DTM_S, cid, tid string) {
}

//done event
func (d *DTM_S_Proc) OnDone(dtm *DTM_S, args util.Map, cid, tid string, code int, err string, used int64) {
	d.proc_l.Lock()
	defer d.proc_l.Unlock()
	if tv, ok := d.Rates[cid]; ok {
		if _, ok := tv[tid]; ok {
			d.TaskC[cid] -= 1
			d.AllC -= 1
		}
		delete(tv, tid)
		d.Rates[cid] = tv
	}
}

//login event
func (d *DTM_S_Proc) OnLogin(rc *impl.RCM_Cmd, token string) (string, error) {
	d.proc_l.Lock()
	defer d.proc_l.Unlock()
	cid := atomic.AddInt64(&d.cid, 1)
	cid_ := fmt.Sprintf("N-%v", cid)
	d.TaskC[cid_] = 0
	return cid_, nil
}

//connection event
func (d *DTM_S_Proc) OnConn(c netw.Con) bool {
	c.SetWait(true)
	return true
}

func (d *DTM_S_Proc) OnClose(c netw.Con) {
	d.proc_l.Lock()
	defer d.proc_l.Unlock()
	var cid = c.Kvs().StrVal("cid")
	if len(cid) > 0 {
		delete(d.TaskC, cid)
		delete(d.Rates, cid)
	}
}

//minial used client id
func (d *DTM_S_Proc) MinUsedCid(dtm *DTM_S, args ...interface{}) string {
	var tcid string = ""
	var min int = 999
	for cid, tc := range d.TaskC {
		if tc < min {
			tcid = cid
			min = tc
		}
	}
	return tcid
}

//total count
func (d *DTM_S_Proc) Total() int {
	return d.AllC
}

//process status
func (d *DTM_S_Proc) SrvHTTP(hs *routing.HTTPSession) routing.HResult {
	hs.JsonRes(d)
	return routing.HRES_RETURN
}

//the distributed task manager server impl
type DTM_S struct {
	*rc.RC_Listener_m
	H        DTM_S_H //the distributed task manager handler
	sequence int64   //sequence
}

//new the distributed task manager server impl
func NewDTM_S(bp *pool.BytePool, addr string, h DTM_S_H, rcm *impl.RCM_S, v2b netw.V2Byte, b2v netw.Byte2V, na impl.NAV_F) *DTM_S {
	sh := &DTM_S{
		H:        h,
		sequence: 0,
	}
	obdh := impl.NewOBDH()
	obdh.AddF(CMD_M_PROC, sh.OnProc)
	obdh.AddF(CMD_M_DONE, sh.OnDone)
	lm := rc.NewRC_Listener_m(bp, addr, netw.NewCCH(h, obdh), rcm, v2b, b2v, na)
	lm.LCH = h
	sh.RC_Listener_m = lm
	return sh
}

//new the distributed task manager server impl by json impl
func NewDTM_S_j(bp *pool.BytePool, addr string, h DTM_S_H) *DTM_S {
	rcm := impl.NewRCM_S_j()
	return NewDTM_S(bp, addr, h, rcm, impl.Json_V2B, impl.Json_B2V, impl.Json_NAV)
}

//process event impl handler
func (d *DTM_S) OnProc(c netw.Cmd) int {
	var args util.Map
	_, err := c.V(&args)
	if err != nil {
		log.E("DTM_S OnProc convert arguments error(%v)", err)
		return -1
	}
	var tid string
	var rate float64
	err = args.ValidF(`
		tid,R|S,L:0;
		rate,R|F,R:-0.99;
		`, &tid, &rate)
	if err != nil {
		log.E("DTM_S OnProc receive bad arguments detail(%v)", err)
		return -1
	}
	d.H.OnProc(d, d.ConCid(c), tid, rate)
	return 0
}

//done event impl handler
func (d *DTM_S) OnDone(c netw.Cmd) int {
	var args util.Map
	_, err := c.V(&args)
	if err != nil {
		log.E("DTM_S OnDone convert arguments error(%v)", err)
		return -1
	}
	var code int
	var tid string
	var err_m string
	var used int64
	err = args.ValidF(`
		code,R|I,R:-999;
		tid,R|S,L:0;
		err,O|S,L:0;
		used,R|I,R:-1;
		`, &code, &tid, &err_m, &used)
	if err != nil {
		log.E("DTM_S OnDone receive bad arguments detail(%v)", err)
		return -1
	}
	d.H.OnDone(d, args, d.ConCid(c), tid, code, err_m, used)
	return 0
}

//connection event
// func (d *DTM_S) OnConn(c netw.Con) bool {
// 	c.SetWait(true)
// 	return true
// }

// //connection event
// func (d *DTM_S) OnClose(c netw.Con) {

// }

//start task by command and special client id.
//return task id
//return error when start task fail.
func (d *DTM_S) StartTask(cid, tid, cmds string) error {
	tc := d.CmdC(cid)
	if tc == nil {
		return util.Err("DTM_S StartTask by cid(%v) error->client not found", cid)
	}
	res, err := tc.Exec_m("start_task", util.Map{
		"tid":  tid,
		"cmds": cmds,
	})
	if err != nil {
		return util.Err("DTM_S StartTask executing by tid(%v),cmds(%v) on client(%v) error->%v", tid, cmds, cid, err)
	}
	if res.IntVal("code") == 0 {
		d.H.OnStart(d, cid, tid, cmds)
		return nil
	} else {
		return util.Err("DTM_S StartTask executing by tid(%v),cmds(%v) on client(%v)  error(%v)->%v",
			tid, cmds, cid, res.IntVal("code"), res.StrVal("err"))
	}
}

//start task by special task id and commands
func (d *DTM_S) StartTask2(tid, cmds string, args ...interface{}) (string, error) {
	cid := d.H.MinUsedCid(d, args...)
	if len(cid) < 1 {
		return "", util.Err("DTM_S StartTask2 by cmds(%v) error->not logined client found by calling MinUsedCid", cmds)
	}
	return cid, d.StartTask(cid, tid, cmds)
}

//start task by command, it will select the client which the number of task is minimal to run the task.
//return the client id to run task and task id
//return error when start task fail.
func (d *DTM_S) StartTask3(cmds string, args ...interface{}) (cid string, tid string, err error) {
	tid = fmt.Sprintf("T-%v", atomic.AddInt64(&d.sequence, 1))
	cid, err = d.StartTask2(tid, cmds, args...)
	return
}

func (d *DTM_S) StartTask4(cid, cmds string) (string, error) {
	var tid = fmt.Sprintf("T-%v", atomic.AddInt64(&d.sequence, 1))
	return tid, d.StartTask(cid, tid, cmds)
}

//stop task by client id and task id.
//return error when stop task fail.
func (d *DTM_S) StopTask(cid, tid string) error {
	tc := d.CmdC(cid)
	if tc == nil {
		return util.Err("DTM_S StopTask by cid(%v) error->client not found", cid)
	}
	res, err := tc.Exec_m("stop_task", util.Map{
		"tid": tid,
	})
	if err != nil {
		return util.Err("DTM_S StopTask executing by tid(%v) error->%v", tid, err)
	}
	if res.IntVal("code") == 0 {
		d.H.OnStop(d, cid, tid)
		return nil
	} else {
		return util.Err("DTM_S StopTask executing by tid(%v) error(%v)->%v", tid, res.IntVal("code"), res.StrVal("err"))
	}
}

//wait the task done by client id and task id.
//return error when stop task fail.
func (d *DTM_S) WaitTask(cid, tid string) error {
	tc := d.CmdC(cid)
	if tc == nil {
		return util.Err("DTM_S WaitTask by cid(%v) error->client not found", cid)
	}
	res, err := tc.Exec_m("wait_task", util.Map{
		"tid": tid,
	})
	if err != nil {
		return util.Err("DTM_S WaitTask executing by tid(%v) error->%v", tid, err)
	}
	if res.IntVal("code") == 0 {
		return nil
	} else {
		return util.Err("DTM_S WaitTask executing by tid(%v) error(%v)->%v", tid, res.IntVal("code"), res.StrVal("err"))
	}
}

//new the distributed task manager client impl
type DTM_C struct {
	*rc.RC_Runner_m
	Cfg     *util.Fcfg           //configure
	Tasks   map[string]*exec.Cmd //running task
	tasks_l sync.RWMutex
	tasks_c map[string]chan string
}

//new the distributed task manager client impl
//
//it will parse command by DTM_C.Cfg
//
func NewDTM_C(bp *pool.BytePool, addr string, rcm *impl.RCM_S, v2b netw.V2Byte, b2v netw.Byte2V, na impl.NAV_F) *DTM_C {
	ch := &DTM_C{
		Cfg:     util.NewFcfg3(),
		Tasks:   map[string]*exec.Cmd{},
		tasks_l: sync.RWMutex{},
		tasks_c: map[string]chan string{},
	}
	cr := rc.NewRC_Runner_m(bp, addr, ch, rcm, v2b, b2v, na)
	ch.RC_Runner_m = cr
	cr.AddHFunc("start_task", ch.StartTask)
	cr.AddHFunc("wait_task", ch.WaitTask)
	cr.AddHFunc("stop_task", ch.StopTask)
	return ch
}

//new the distributed task manager client impl by json
func NewDTM_C_j(bp *pool.BytePool, addr string) *DTM_C {
	rcm := impl.NewRCM_S_j()
	return NewDTM_C(bp, addr, rcm, impl.Json_V2B, impl.Json_B2V, impl.Json_NAV)
}

//command impl
func (d *DTM_C) OnCmd(c netw.Cmd) int {
	return 0
}

//connection event
func (d *DTM_C) OnConn(c netw.Con) bool {
	c.SetWait(true)
	var token = d.Cfg.Val("token")
	if len(token) > 0 {
		go d.Login_(token)
	}
	return true
}

//connection event
func (d *DTM_C) OnClose(c netw.Con) {
}

//start task impl func
func (d *DTM_C) StartTask(rc *impl.RCM_Cmd) (interface{}, error) {
	var tid string
	var cmds string
	err := rc.ValidF(`
		tid,R|S,L:0;
		cmds,R|S,L:0;
		`, &tid, &cmds)
	if err != nil {
		return util.Map{"code": -1, "err": fmt.Sprintf("DTM_C start task calling by bad arguments->%v", err)}, nil
	}
	err = d.run_cmd(tid, cmds)
	if err == nil {
		return util.Map{"code": 0, "tid": tid}, nil
	} else {
		return util.Map{"code": -2, "err": fmt.Sprintf("DTM_C StartTask running command(%v) by tid(%v) error->%v", cmds, tid, err)}, nil
	}
}

//stop task impl func
func (d *DTM_C) StopTask(rc *impl.RCM_Cmd) (interface{}, error) {
	var tid string
	err := rc.ValidF(`
		tid,R|S,L:0;
		`, &tid)
	if err != nil {
		return util.Map{"code": -1, "err": fmt.Sprintf("DTM_C stop task calling by bad arguments->%v", err)}, nil
	}
	runner, ok := d.Tasks[tid]
	if !ok {
		return util.Map{"code": -2, "err": fmt.Sprintf("DTM_C stop task by id(%v) fail(task is not found)", tid)}, nil
	}
	err = runner.Process.Kill()
	if err == nil {
		return util.Map{"code": 0}, nil
	} else {
		return util.Map{"code": -3, "err": fmt.Sprintf("DTM_C kill task by id(%v) error(%v)", tid, err)}, nil
	}
}

//wait task impl func
func (d *DTM_C) WaitTask(rc *impl.RCM_Cmd) (interface{}, error) {
	var tid string
	err := rc.ValidF(`
		tid,R|S,L:0;
		`, &tid)
	if err != nil {
		return util.Map{"code": -1, "err": fmt.Sprintf("DTM_C wait task calling by bad arguments->%v", err)}, nil
	}
	tc, ok := d.tasks_c[tid]
	if !ok {
		return util.Map{"code": -2, "err": fmt.Sprintf("DTM_C wait task by id(%v) fail(task is not found)", tid)}, nil
	}
	msg := <-tc
	if len(msg) < 1 {
		return util.Map{"code": 0}, nil
	} else {
		return util.Map{"code": 3, "err": msg}, nil
	}
}

//run the process http handler
func (d *DTM_C) RunProcH() error {
	addr := d.Cfg.Val("proc_addr")
	if len(addr) < 1 {
		log.I("DTM_C RunProcH listen address configure(proc_addr) is not found, http proccess receiver will not start")
		return nil
	}
	routing.HFunc("^/proc(\\?.*)?$", d.HandleProc)
	routing.Shared.Print()
	log.I("DTM_C RunProcH listen the process handle on addr(%v)", addr)
	return routing.ListenAndServe(addr)
}

//process http handler impl
func (d *DTM_C) HandleProc(hs *routing.HTTPSession) routing.HResult {
	log.D("DTM_C HandleProc reiceve process %v", hs.R.URL.Query().Encode())
	var tid string
	var rate float64
	err := hs.ValidCheckVal(`
		tid,R|S,L:0;
		`+d.Cfg.Val2("proc_key", "process")+`,R|F,R:-0.001;`, &tid, &rate)
	if err != nil {
		hs.W.Write([]byte(fmt.Sprintf("DTM_C HandleProc receive bad arguments->%v", err.Error())))
		return routing.HRES_RETURN
	}
	err = d.NotifyProc(tid, rate)
	if err != nil {
		log.E("DTM_C HandleProc send process info by tid(%v),rate(%v) err->%v", tid, rate, err)
	}
	hs.W.Write([]byte("OK"))
	return routing.HRES_RETURN
}

//notify process to server
func (d *DTM_C) NotifyProc(tid string, rate float64) error {
	_, err := d.Writev2([]byte{CMD_M_PROC}, util.Map{
		"tid":  tid,
		"rate": rate,
	})
	return err
}

//add task by id and runner
func (d *DTM_C) add_task(tid string, runner *exec.Cmd) chan string {
	d.tasks_l.Lock()
	d.Tasks[tid] = runner
	c := make(chan string)
	d.tasks_c[tid] = c
	d.tasks_l.Unlock()
	return c
}

//delete task by id
func (d *DTM_C) del_task(tid string) {
	d.tasks_l.Lock()
	delete(d.Tasks, tid)
	c := d.tasks_c[tid]
	delete(d.tasks_c, tid)
	close(c)
	d.tasks_l.Unlock()
}

//run command by id and commmand string
func (d *DTM_C) run_cmd(tid, cmds string) error {
	log.I("DTM_C run_cmd running command(\n\t%v\n) by tid(%v)", cmds, tid)
	cfg := util.NewFcfg4(d.Cfg)
	cfg.SetVal("proc_tid", tid)
	cmds = cfg.EnvReplaceV(cmds, false)
	log.D("DTM_C calling command(\n\t%v\n)", cmds)
	beg := util.Now()
	var runner = exec.Command(d.Cfg.Val2("bash_c", "bash"), "-c", cmds)
	runner.Dir = cfg.Val2("proc_ws", ".")
	var env = cfg.Val2("proc_env", "")
	if len(env) > 0 {
		runner.Env = append(os.Environ(), strings.Split(env, ",")...)
	}
	buf := &bytes.Buffer{}
	runner.Stdout = buf
	runner.Stderr = buf
	err := runner.Start()
	if err != nil {
		err = util.Err("DTM_C run_cmd start error->%v", err)
		log.E("%v", err)
		return err
	}
	task_c := d.add_task(tid, runner)
	go func() {
		args := util.Map{"tid": tid}
		err = runner.Wait()
		used := util.Now() - beg
		res := buf.String()
		if err == nil {
			log.D("DTM_C run_cmd by running command(\n\t%v\n) success,used(%vms)->\n%v", cmds, used, res)
			args["code"] = d.cmd_do_res(args, cmds, res)
		} else {
			log.E("DTM_C run_cmd by running command(\n\t%v\n) error(%v)->\n%v", cmds, err, res)
			args["code"] = -1
			args["err"] = err.Error()
		}
		args["used"] = used
		d.Writev2([]byte{CMD_M_DONE}, args)
		task_c <- args.StrVal("err")
		d.del_task(tid)
	}()
	return nil
}

func (d *DTM_C) cmd_do_res(args util.Map, cmds, res string) int {
	var res_a = strings.SplitN(res, "----------------result----------------", 2)
	if len(res_a) < 2 {
		return 0
	}
	var mres = util.ParseSectionF("[", "]", res_a[1])
	var jval = mres.StrVal("json")
	if len(jval) < 1 {
		args["data"] = mres
		return 0
	}
	var jval_m, err = util.Json2Map(jval)
	if err == nil {
		args["data"] = jval_m
		return 0
	} else {
		log.E("DTM_C parse json result on command(\n\t%v\n) by data(%v) error->%v", cmds, jval, err)
		args["data"] = mres
		args["err"] = err.Error()
		return -2
	}
}
