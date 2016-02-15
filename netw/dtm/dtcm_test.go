package dtm

import (
	"fmt"
	"github.com/Centny/gwf/pool"
	"github.com/Centny/gwf/util"
	"runtime"
	"sync/atomic"
	"testing"
	"time"
)

type Mem_err struct {
	E1, E2, E3 error
	Data       map[string]*Task
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

type dtcm_s_h struct {
	cc int32
	E  error
}

func (d *dtcm_s_h) OnStart(dtcm *DTCM_S, task *Task) {
	atomic.AddInt32(&d.cc, 1)
}
func (d *dtcm_s_h) OnDone(dtcm *DTCM_S, task *Task) error {
	atomic.AddInt32(&d.cc, -1)
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
	var sh = &dtcm_s_h{}
	dtms, err := NewDTCM_S_j(bp, cfg, MemDbc, sh)
	if err != nil {
		t.Error(err.Error())
		return
	}
	err = dtms.Run()
	if err != nil {
		t.Error(err.Error())
		return
	}
	fmt.Println(dtms.Sid)
	dtmc := NewDTM_C_j(bp, "127.0.0.1:2324")
	err = dtmc.Cfg.InitWithFilePath2("dtcm_c.properties", false)
	if err != nil {
		t.Error(err.Error())
		return
	}
	dtmc.Start()
	time.Sleep(time.Second)
	fmt.Println("---->")
	for i := 0; i < 10; i++ {
		err = dtms.AddTask(nil, "abc.mkv", "abc")
		if err != nil {
			t.Error(err.Error())
			return
		}
	}
	for i := 0; i < 10; i++ {
		err = dtms.AddTask(nil, "abc.mp4", "xxx")
		if err != nil {
			t.Error(err.Error())
			return
		}
	}
	time.Sleep(2 * time.Second)
	fmt.Println(sh.cc)
	if sh.cc != 0 {
		t.Error("error")
		return
	}
	err = dtms.AddTask(nil, "xxds", "sd")
	if err == nil {
		t.Error("error")
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
	//test other
	dtms.OnProc(dtms.DTM_S, "cid", "tid", 100)
	dtms.OnStop(dtms.DTM_S, "cid", "ss")
	//
	//test error
	fmt.Println("test error...")
	//
	dtms.stop_task(&Task{
		Proc: map[string]*Proc{
			"xxx": &Proc{
				Tid: "xxx",
			},
		},
	})
	//
	dtms.stop_task(&Task{
		Proc: map[string]*Proc{
			"xx": &Proc{},
		},
	})
	//
	tdb, _ := MemErrDbc("", "")
	dtms.Db = tdb
	fmt.Println("test error...1")
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
