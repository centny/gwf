package tutil

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/Centny/gwf/util"
)

func AutoExec2(logf string, precall func(state Perf2, ridx uint64, used int64, rerr error) (uint64, error), call func(uint64) error) (used int64, max, avg uint64, err error) {
	var total, tc, peradd, runmax uint64 = 1000000, 100, 2, 100000
	var pretimeout int64 = 1000
	total, err = strconv.ParseUint(os.Getenv("PERF_TOTOAL"), 10, 64)
	if err != nil {
		total = 1000000
	}
	tc, err = strconv.ParseUint(os.Getenv("PERF_TC"), 10, 64)
	if err != nil {
		tc = 100
	}
	peradd, err = strconv.ParseUint(os.Getenv("PERF_PERADD"), 10, 64)
	if err != nil {
		peradd = 2
	}
	pretimeout, err = strconv.ParseInt(os.Getenv("PERF_TIMEOUT"), 10, 64)
	if err != nil {
		pretimeout = 1000
	}
	runmax, err = strconv.ParseUint(os.Getenv("PERF_MAX"), 10, 64)
	if err != nil {
		runmax = 100000
	}
	fmt.Printf("AutoExec:\n\ttotal:%v,tc:%v,peradd:%v,timeout:%v\n\n", total, tc, peradd, pretimeout)
	perf := NewPerf2()
	perf.ShowState = true
	perf.RunninMax = runmax
	perf.Timeout = pretimeout
	// perf.ExternalStatus = func() string {
	// 	mgo.Plck.Lock()
	// 	defer mgo.Plck.Unlock()
	// 	return fmt.Sprintf(",%v,%v", len(mgo.Pending), mgo.Donc)
	// }
	return perf.AutoExec(total, tc, peradd, logf, precall, call)
}

type Perf2 struct {
	Running        int64
	Max            uint64
	Avg            uint64
	PerUsed        map[float64]int64
	PerUsedMax     int64
	PerUsedMin     int64
	PerUsedAvg     int64
	PerUsedAll     int64
	Done           int64
	Errc           int64
	Used           int64
	lck            *sync.RWMutex
	mrunning       bool
	mwait          *sync.WaitGroup
	stdout, stderr *os.File
	ShowState      bool `json:"-"`
	IncreaseDelay  int64
	ExternalStatus func() string
	RunninMax      uint64
	timeouted      uint64
	fulled         uint64
	Timeout        int64
}

func NewPerf2() *Perf2 {
	return &Perf2{
		lck:            &sync.RWMutex{},
		mwait:          &sync.WaitGroup{},
		IncreaseDelay:  200,
		ExternalStatus: func() string { return "" },
		RunninMax:      8000,
		PerUsed:        map[float64]int64{},
	}
}

type DisVal []string

func (d DisVal) Len() int {
	return len(d)
}

func (d DisVal) Swap(i, j int) {
	d[i], d[j] = d[j], d[i]
}

func (d DisVal) Less(i, j int) bool {
	vali, _ := strconv.ParseFloat(strings.Split(d[i], ":")[0], 10)
	valj, _ := strconv.ParseFloat(strings.Split(d[j], ":")[0], 10)
	return vali < valj
}

func (p *Perf2) String() string {
	p.lck.Lock()
	dis := DisVal{}
	for d, v := range p.PerUsed {
		dis = append(dis, fmt.Sprintf("%v:%v%%", d, int64(100*float64(v)/float64(p.Done))))
	}
	p.lck.Unlock()
	sort.Sort(dis)
	alldis := strings.Join(dis, ", ")
	external := p.ExternalStatus()
	if len(external) > 0 {
		alldis += "\n  " + external
	}
	// alldis += "\n  " + fmt.Sprintf("Max:%v", mgo.UserLockPending2.Max)
	p.NotifyReal()
	return fmt.Sprintf("Used:%v,Done:%v,Errc:%v,TPS:%v,Running:%v,Max:%v,Avg:%v,PerMax:%v,PerMin:%v,PerAvg:%v,Timeout:%v,Fullc:%v\n  %v",
		p.Used, p.Done, p.Errc, p.Done*1000/p.Used, p.Running, p.Max, p.Avg, p.PerUsedMax, p.PerUsedMin, p.PerUsedAvg, p.timeouted, p.fulled, alldis)
}

func (p *Perf2) NotifyReal() {
	realURI := os.Getenv("REAL_URI")
	if len(realURI) > 0 {
		host := os.Getenv("HOSTNAME")
		host = strings.Split(host, ".")[0]
		bys, _ := json.Marshal(map[string]map[string]interface{}{
			host: map[string]interface{}{
				"Used":       p.Used,
				"Done":       p.Done,
				"Errc":       p.Errc,
				"TPS":        p.Done * 1000 / p.Used,
				"Running":    p.Running,
				"Max":        p.Max,
				"Avg":        p.Avg,
				"PerUsedMax": p.PerUsedMax,
				"PerUsedAvg": p.PerUsedAvg,
				"Timeout":    p.timeouted,
				"Fullc":      p.fulled,
			},
		})
		resp, err := http.Post(realURI, "application/json", bytes.NewBuffer(bys))
		if err == nil {
			ioutil.ReadAll(resp.Body)
		}
	}
}

func (p *Perf2) AutoExec(total, tc, peradd uint64, logf string, precall func(state Perf2, ridx uint64, used int64, rerr error) (uint64, error), call func(uint64) error) (used int64, max, avg uint64, err error) {
	// percallc := 1
	used, err = p.Exec(total, tc, logf,
		func(state Perf2, ridx uint64, used int64, rerr error) (uint64, error) {
			// time.Sleep(time.Second)
			_, terr := precall(state, ridx, used, rerr)
			if terr == FullError {
				atomic.AddUint64(&p.fulled, 1)
				if uint64(p.Running) < tc {
					return 1, nil
				}
				return 0, nil
			}
			if terr != nil {
				return 0, terr
			}
			if p.Timeout > 0 && used >= p.Timeout {
				atomic.AddUint64(&p.timeouted, 1)
				if uint64(p.Running) < tc {
					return 1, nil
				}
				return 0, nil
			}
			if state.Running > int64(p.RunninMax) {
				return 0, nil
			}
			return peradd, nil
			// beg := utils.Now()
			// terr := precall(idx, state)
			// if terr == nil {
			// 	if state.Running >= p.RunninMax {
			// 		return 0, nil
			// 	}
			// 	pretimeout := utils.Now()-beg < pretimeout
			// 	if pretimeout {
			// 		percallc++
			// 		nadd := peradd*percallc - int(state.Running)
			// 		return nadd, nil
			// 	}
			// 	p.pretimeouted = pretimeout
			// 	if int(state.Running) < tc {
			// 		return 1, nil
			// 	}
			// 	return 0, nil
			// } else if terr == FullError {
			// 	return 0, nil
			// } else {
			// 	return 0, terr
			// }
		}, call)

	max, avg = p.Max, p.Avg
	fmt.Printf(`
		%v
		------Done------
		- total:%v
		- max:%v
		- avg:%v
		- used:%v
		- error:%v
		----------------------
		%v`, p.String(), total, max, avg, used, err, "\n")
	return
}

func (p *Perf2) perdone(perused int64) {
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
	if p.Timeout > 0 {
		multiple := float64(perused) / float64(p.Timeout)
		if multiple <= 0.25 {
			p.PerUsed[0.25]++
		} else if multiple <= 0.5 {
			p.PerUsed[0.5]++
		} else if multiple <= 0.75 {
			p.PerUsed[0.75]++
		} else if multiple <= 1 {
			p.PerUsed[1]++
		} else if multiple <= 1.5 {
			p.PerUsed[1.5]++
		} else if multiple <= 2 {
			p.PerUsed[2]++
		} else {
			floatn := (math.Sqrt(8*multiple+1) - 1) / 2
			n := float64(int64(floatn))
			if n < floatn {
				n++
			}
			p.PerUsed[(n-1)*(n+2)/2+1]++
		}
	}
}
func (p *Perf2) monitor() {
	var allrunning, allc uint64
	p.mrunning = true
	p.mwait.Add(1)
	defer p.mwait.Done()

	beg := util.Now()
	for p.mrunning {
		running := uint64(p.Running)
		allrunning += running
		allc++
		p.Avg = allrunning / allc
		if p.Max < running {
			p.Max = running
		}
		p.Used = util.Now() - beg
		if p.ShowState {
			fmt.Fprintf(p.stdout, "->%v\n", p)
		}
		time.Sleep(time.Second)
	}
}

func (p *Perf2) Exec(total, tc uint64, logf string, increase func(state Perf2, ridx uint64, used int64, rerr error) (uint64, error), call func(uint64) error) (int64, error) {
	p.stdout, p.stderr = os.Stdout, os.Stderr
	if len(logf) > 0 {
		f, err := os.OpenFile(logf, os.O_APPEND|os.O_CREATE|os.O_WRONLY, os.ModePerm)
		if err != nil {
			return 0, err
		}
		os.Stdout = f
		os.Stderr = f
		log.SetOutput(f)
	}
	if tc > total {
		tc = total
	}
	go p.monitor()
	ws := sync.WaitGroup{}
	// ws.Add(total)
	beg := util.Now()
	var tidx_ uint64 = 0
	var run_call func(ridx uint64)
	var run_next func(ridx uint64, used int64, rerr error)
	var err error = nil
	run_call = func(ridx uint64) {
		defer ws.Done()

		perbeg := util.Now()
		terr := call(ridx)
		perused := util.Now() - perbeg
		if terr == nil {
			p.perdone(perused)
		} else {
			atomic.AddInt64(&p.Errc, 1)
		}
		run_next(ridx, perused, terr)
		atomic.AddInt64(&p.Running, -1)

	}
	// increaseLck := sync.RWMutex{}
	run_next = func(ridx uint64, used int64, rerr error) {
		// increaseLck.Lock()
		nc, terr := increase(*p, ridx, used, rerr)
		// increaseLck.Unlock()
		if terr != nil {
			err = terr
			return
		}
		for i := uint64(0); i < nc; i++ {
			ridx := atomic.AddUint64(&tidx_, 1)
			if ridx >= total {
				break
			}
			ws.Add(1)
			atomic.AddInt64(&p.Running, 1)
			go run_call(ridx)
		}
	}
	atomic.AddUint64(&tidx_, tc-1)
	for i := uint64(0); i < tc; i++ {
		ws.Add(1)
		atomic.AddInt64(&p.Running, 1)
		go run_call(i)
	}
	// for err == nil && tidx_ < int32(total) {
	// 	time.Sleep(time.Duration(p.IncreaseDelay) * time.Millisecond)
	// 	run_next(0)
	// }
	ws.Wait()
	end := util.Now()
	if len(logf) > 0 {
		os.Stdout.Close()
		os.Stdout = p.stdout
		os.Stderr = p.stderr
		log.SetOutput(os.Stdout)
	}
	p.mrunning = false
	p.mwait.Wait()
	return end - beg, err
}

func TimeoutSec(state Perf2, ridx uint64, used int64, rerr error) (uint64, error) {
	if rerr != nil {
		return 0, rerr
	}
	if used > 1000 {
		return 0, FullError
	}
	return 1, nil
}
