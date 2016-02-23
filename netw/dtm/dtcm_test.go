package dtm

import (
	"fmt"
	"github.com/Centny/gwf/pool"
	"github.com/Centny/gwf/routing/httptest"
	"github.com/Centny/gwf/util"
	"runtime"
	"sync"
	"testing"
	"time"
)

type Mem_err struct {
	E1, E2, E3, E4 error
	Data           map[string]*Task
}

func MemErrDbc(uri, name string) (*Mem_err, error) {
	return &Mem_err{Data: map[string]*Task{}}, nil
}
func (m *Mem_err) Add(t *Task) error {
	m.Data[t.Id] = t
	return m.E1
}
func (m *Mem_err) Update(t *Task) error {
	m.Data[t.Id] = t
	return m.E2
}
func (m *Mem_err) Del(t *Task) error {
	delete(m.Data, t.Id)
	return m.E3
}
func (m *Mem_err) List() ([]*Task, error) {
	var ts []*Task
	for _, task := range m.Data {
		ts = append(ts, task)
	}
	return ts, m.E4
}

type dtcm_s_h struct {
	cc  map[string]int
	lck sync.RWMutex
	E   error
}

func new_dtcm_s_h() *dtcm_s_h {
	return &dtcm_s_h{
		cc:  map[string]int{},
		lck: sync.RWMutex{},
	}
}
func (d *dtcm_s_h) OnStart(dtcm *DTCM_S, task *Task) {
	d.lck.Lock()
	defer d.lck.Unlock()
	d.cc[task.Id] = 1
}
func (d *dtcm_s_h) OnDone(dtcm *DTCM_S, task *Task) error {
	d.lck.Lock()
	defer d.lck.Unlock()
	delete(d.cc, task.Id)
	return d.E
}

func TestDtcm(t *testing.T) {
	runtime.GOMAXPROCS(util.CPU())
	bp := pool.NewBytePool(8, 10240000)
	// netw.ShowLog = true
	// impl.ShowLog = true
	var cfg = util.NewFcfg3()
	var err = cfg.InitWithFilePath2("dtcm.properties", false)
	if err != nil {
		t.Error(err.Error())
		return
	}
	var sh = new_dtcm_s_h()
	dtms, err := StartDTCM_S(cfg, MemDbc, sh)
	if err != nil {
		t.Error(err.Error())
		return
	}
	ts := httptest.NewServer2(dtms)
	var cfg_c = util.NewFcfg3()
	err = cfg_c.InitWithFilePath2("dtcm_c.properties", false)
	if err != nil {
		t.Error(err.Error())
		return
	}
	StartDTM_C(cfg_c)
	//
	func() {
		var cfg_c_x = util.NewFcfg3()
		err = cfg_c_x.InitWithFilePath2("dtcm_c.properties", false)
		if err != nil {
			t.Error(err.Error())
			return
		}
		cfg_c_x.SetVal("token", "ax1")
		StartDTM_C(cfg_c_x)
	}()
	time.Sleep(time.Second)
	fmt.Println("---->")
	var rc = 10
	for i := 0; i < rc; i++ {
		err = dtms.AddTask(nil, "abc.mkv", "abc")
		if err != nil {
			t.Error(err.Error())
			return
		}
	}
	for i := 0; i < rc; i++ {
		err = dtms.AddTask(nil, "abc.mp4", "xxx")
		if err != nil {
			t.Error(err.Error())
			return
		}
	}
	for {
		if len(sh.cc) > 0 {
			fmt.Println("waiting->", len(sh.cc), "->", ts.URL)
			time.Sleep(2 * time.Second)
		} else {
			break
		}
	}
	// time.Sleep(5 * time.Second)
	fmt.Println(sh.cc)
	if len(sh.cc) != 0 {
		fmt.Println(sh.cc, "-->")
		t.Error("error")
		return
	}
	err = dtms.AddTask(nil, "xxds", "sd")
	if err == nil {
		t.Error("error")
		return
	}
	err = dtms.AddTask(nil, "xxds.xx", "sd")
	if err != nil {
		t.Error(err.Error())
		return
	}
	fmt.Println("xxxx----->")
	err = dtms.AddTask(nil, "exit", "sd")
	if err != nil {
		t.Error(err.Error())
		return
	}
	time.Sleep(2 * time.Second)
	// ts = httptest.NewServer2(dtms)
	fmt.Println("---->")
	fmt.Println(ts.G(""))
	fmt.Println("---->")
	//test other
	dtms.OnProc(dtms.DTM_S, "cid", "tid", 100)
	dtms.OnStop(dtms.DTM_S, "cid", "ss")
	//
	//test error
	fmt.Println("test error...")
	//
	dtms.stop_task(&Task{
		Args: []interface{}{"abc.mkv", "abc"},
		Proc: map[string]*Proc{
			"xxx": &Proc{
				Tid: "xxx",
			},
		},
	})
	//
	dtms.stop_task(&Task{
		Args: []interface{}{"abc.mkv", "abc"},
		Proc: map[string]*Proc{
			"xx": &Proc{},
		},
	})
	dtms.start_task(&Task{
		Args: []interface{}{"abc.mkv", "abc"},
		Proc: map[string]*Proc{
			"xx": &Proc{
				Status: TKS_DONE,
			},
		},
	})
	var task = &Task{
		Args: []interface{}{"abc.mkv", "abc"},
		Proc: map[string]*Proc{
			"xx": &Proc{
				Status: TKS_PENDING,
			},
		},
	}
	dtms.start_task(task)
	if task.Proc["xx"].Status != TKS_COV_ERR {
		t.Error("error")
		return
	}
	//
	tdb, _ := MemErrDbc("", "")
	dtms.Db = tdb
	fmt.Println("test error...1")
	//
	var tt = dtms.NewTask(nil, "abc.mkv", "abc")
	tdb.Data[tt.Id] = tt
	time.Sleep(2 * time.Second)
	tdb.Data[tt.Id] = tt
	time.Sleep(2 * time.Second)
	tdb.E4 = util.Err("mock error")
	time.Sleep(2 * time.Second)
	//
	tdb.E2 = util.Err("mock error")
	err = dtms.AddTask(nil, "abc.mkv", "fsf")
	if err != nil {
		t.Error("error")
		return
	}
	fmt.Println("test error...2")
	tdb.E1 = util.Err("mock error")
	err = dtms.AddTask(nil, "abc.mkv", "fsf")
	if err == nil {
		t.Error("error")
		return
	}
	//
	fmt.Println("test error...3")
	_, err = NewDTCM_S_j(bp, cfg, func(uri, name string) (DbH, error) {
		return nil, util.Err("error")
	}, &dtcm_s_h{})
	if err == nil {
		t.Error("error")
		return
	}
	cfg.SetVal("T1/regs", "")
	_, err = NewDTCM_S_j(bp, cfg, MemDbc, &dtcm_s_h{})
	if err == nil {
		t.Error("error")
		return
	}
	cfg.SetVal("cmds", "")
	_, err = NewDTCM_S_j(bp, cfg, MemDbc, &dtcm_s_h{})
	if err == nil {
		t.Error("error")
		return
	}
	//
	tdb.E3 = util.Err("error")
	dtms.check_done(&Task{
		Proc: map[string]*Proc{
			"xxx": &Proc{
				Status: TKS_DONE,
			},
		},
	})
	sh.E = util.Err("error")
	dtms.check_done(&Task{
		Proc: map[string]*Proc{
			"xxx": &Proc{
				Status: TKS_DONE,
			},
		},
	})

	//
	var xx = &Cmd{}
	xx.Match()
	//
	dtms.StopChecker()
	//
	fmt.Println("done...")
}

func TestParseCmds(t *testing.T) {
	var cfg, _ = util.NewFcfg2(`
		[T1]
		regs=.(s
		[T2]
		regs=.*
		`)
	var err error
	_, err = ParseCmds(cfg, []string{"xx"})
	if err == nil {
		t.Error("error")
		return
	}
	_, err = ParseCmds(cfg, []string{"T1"})
	if err == nil {
		t.Error("error")
		return
	}
	_, err = ParseCmds(cfg, []string{"T2"})
	if err == nil {
		t.Error("error")
		return
	}
}

func TestParseClients(t *testing.T) {
	var cfg, _ = util.NewFcfg2(`
[C0]
#max command runner
max=10
token=ax1,ax2

[C1]
#max command runner
max=10
token=a1,a2,abc
regs=.m[p4&.mkv

[C2]
#max command runner
max=10
regs=.flac&.wav

[C3]
#max command runner
max=10
token=a1,a6
regs=.flacx&.wav

[C4]
#max command runner
max=10
token=a1,a6
regs=.flacx&.wav
		`)
	var err error
	_, err = ParseClients(cfg, []string{"C0"})
	if err == nil {
		t.Error("error")
		return
	}
	_, err = ParseClients(cfg, []string{"C1"})
	if err == nil {
		t.Error("error")
		return
	}
	_, err = ParseClients(cfg, []string{"C2"})
	if err == nil {
		t.Error("error")
		return
	}
	_, err = ParseClients(cfg, []string{"C3", "C3"})
	if err == nil {
		t.Error("error")
		return
	}
	fmt.Println(err)
	//
	var cc = &Client{}
	cc.Match()
}

func TestNewDTCM_S_Err(t *testing.T) {
	var err error
	var sh = new_dtcm_s_h()
	var cfg *util.Fcfg

	cfg, _ = util.NewFcfg2(`
[loc]
#the server id
sid=s1
#the command list
cmds=T1,T2,T3,T4,T5
#clients
clients=C0,C1,C2,C3
#listen address
addr=:2324
#the db connection
db_con=xxx
#the db name		
db_name=xxx
#
max=8
#check delay
cdelay=500
mcache=1024000

#task
[T1]
#the regex for mathec task key
regs=.mkv&.avi
#the commmand to runner by format string
cmds=${CMD_1} ${v0} ${v1}_1.mp4

[T2]
regs=.mp4&.mkv
cmds=${CMD_2} ${v0} ${v1}_2.mp4 xx

[T3]
regs=.flac&.wav
cmds=${CMD_3} ${v0} ${v1}_3.mp3

[T4]
regs=^.*\.xx$
cmds=${CMD_4} ${v0} ${v1}_3.mp3

[T5]
regs=^exit$
cmds=${CMD_2} ${v0} ${v1}_3.mp3

[C0]
#max command runner
max=10
token=ax1,ax2
regs=.flacx&.wavx

[C1]
#max command runner
max=10
token=a1,a2,abc
regs=.mp4&.mkv&.flac&.wav&.avi&^exit$&^.*\.xx$

[C2]
#max command runner
max=10
token=a3,a4
regs=.flac&.wav

[C3]
#max command runner
max=10
token=a1,a3
regs=.flacx&.wavx

		`)
	//
	_, err = StartDTCM_S(cfg, MemDbc, sh)
	if err == nil {
		t.Error("error")
		return
	}
	fmt.Println(err)
	//
	cfg.SetVal("C2/regs", "")
	_, err = StartDTCM_S(cfg, MemDbc, sh)
	if err == nil {
		t.Error("error")
		return
	}
	fmt.Println(err)
	//
	cfg.SetVal("clients", "")
	_, err = StartDTCM_S(cfg, MemDbc, sh)
	if err == nil {
		t.Error("error")
		return
	}
	fmt.Println(err)
	//
	var cc = &Client{}
	cc.Match()
}

func TestDoNone(t *testing.T) {
	dnh := NewDoNoneH()
	dnh.OnStart(nil, &Task{})
	dnh.OnDone(nil, &Task{})
}
