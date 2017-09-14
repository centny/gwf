package tutil

import (
	"fmt"
	"os"
	"sync"
	"sync/atomic"
	"time"

	"github.com/Centny/gwf/log"
	"github.com/Centny/gwf/util"
)

var FullError = fmt.Errorf("runner is almost full")

func DoPerf(tc int, logf string, call func(int)) (int64, error) {
	return DoPerfV(tc, tc, logf, call)
}
func DoPerfV(total, tc int, logf string, call func(int)) (int64, error) {
	return DoPerfV_(total, tc, logf,
		func(i int) error {
			call(i)
			return nil
		})
}

func DoAutoPerfV(total, tc, peradd int, logf string, pretimeout int64, precall func(idx int, state Perf) error, call func(int) error) (used int64, max, avg int, err error) {
	perf := NewPerf()
	return perf.AutoExec(total, tc, peradd, logf, pretimeout, precall, call)
}

func DoPerfV_(total, tc int, logf string, call func(int) error) (int64, error) {
	return DoAutoPerfV_(total, tc, logf,
		func(idx int, state Perf) (int, error) {
			return 1, nil
		}, call)
}

func DoAutoPerfV_(total, tc int, logf string, increase func(idx int, state Perf) (int, error), call func(int) error) (int64, error) {
	perf := NewPerf()
	return perf.Exec(total, tc, logf, increase, call)
}

type Perf struct {
	Running        int32
	Max            int32
	Avg            int32
	PerUsedMax     int64
	PerUsedMin     int64
	PerUsedAvg     int64
	PerUsedAll     int64
	Done           int64
	lck            *sync.RWMutex
	mrunning       bool
	mwait          *sync.WaitGroup
	stdout, stderr *os.File
	ShowState      bool `json:"-"`
}

func NewPerf() *Perf {
	return &Perf{
		lck:   &sync.RWMutex{},
		mwait: &sync.WaitGroup{},
	}
}

func (p *Perf) AutoExec(total, tc, peradd int, logf string, pretimeout int64, precall func(idx int, state Perf) error, call func(int) error) (used int64, max, avg int, err error) {
	used, err = p.Exec(total, tc, logf,
		func(idx int, state Perf) (int, error) {
			beg := util.Now()
			terr := precall(idx, state)
			if terr == nil {
				if util.Now()-beg < pretimeout {
					return peradd, nil
				}
				if int(state.Running) < tc {
					return 1, nil
				}
				return 0, nil
			} else if terr == FullError {
				if int(state.Running) < tc {
					return 1, nil
				}
				return 0, nil
			} else {
				return 0, terr
			}
		}, call)
	max, avg = int(p.Max), int(p.Avg)
	return
}

func (p *Perf) perdone(perused int64) {
	p.lck.Lock()
	defer p.lck.Unlock()
	p.Done++
	p.PerUsedAll += perused
	p.PerUsedAvg = p.PerUsedAll / p.Done
	if p.PerUsedMax < perused {
		p.PerUsedMax = perused
	}
	if p.PerUsedMin == 0 || p.PerUsedMin > perused {
		p.PerUsedMin = perused
	}
}
func (p *Perf) monitor() {
	var allrunning, allc int32
	p.mrunning = true
	p.mwait.Add(1)
	for p.mrunning {
		running := p.Running
		allrunning += running
		allc++
		p.Avg = allrunning / allc
		if p.Max < running {
			p.Max = running
		}
		if p.ShowState {
			fmt.Printf("State:%v\n", util.S2Json(p))
		}
		time.Sleep(time.Second)
	}
	p.mwait.Done()
}

func (p *Perf) Exec(total, tc int, logf string, increase func(idx int, state Perf) (int, error), call func(int) error) (int64, error) {
	p.stdout, p.stderr = os.Stdout, os.Stderr
	if len(logf) > 0 {
		f, err := os.OpenFile(logf, os.O_APPEND|os.O_CREATE|os.O_WRONLY, os.ModePerm)
		if err != nil {
			return 0, err
		}
		os.Stdout = f
		os.Stderr = f
		log.SetWriter(f)
	}
	if tc > total {
		tc = total
	}
	go p.monitor()
	ws := sync.WaitGroup{}
	// ws.Add(total)
	beg := util.Now()
	var tidx_ int32 = 0
	var run_call, run_next func(int)
	var err error = nil
	run_call = func(v int) {
		perbeg := util.Now()
		terr := call(v)
		atomic.AddInt32(&p.Running, -1)
		perused := util.Now() - perbeg
		if terr == nil {
			p.perdone(perused)
		} else {
			err = terr
		}
		if err == nil {
			run_next(v)
		}
		ws.Done()
	}
	var increaselck = sync.RWMutex{}
	run_next = func(v int) {
		increaselck.Lock()
		defer increaselck.Unlock()
		nc, terr := increase(v, *p)
		if terr != nil {
			err = terr
			return
		}
		for i := 0; i < nc; i++ {
			ridx := int(atomic.AddInt32(&tidx_, 1))
			if ridx >= total {
				break
			}
			ws.Add(1)
			atomic.AddInt32(&p.Running, 1)
			go run_call(ridx)
		}
	}
	atomic.AddInt32(&tidx_, int32(tc-1))
	for i := 0; i < tc; i++ {
		ws.Add(1)
		atomic.AddInt32(&p.Running, 1)
		go run_call(i)
	}
	ws.Wait()
	end := util.Now()
	if len(logf) > 0 {
		os.Stdout.Close()
		os.Stdout = p.stdout
		os.Stderr = p.stderr
		log.SetWriter(os.Stdout)
	}
	p.mrunning = false
	p.mwait.Wait()
	return end - beg, err
}
