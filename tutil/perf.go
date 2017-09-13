package tutil

import (
	"fmt"
	"os"
	"sync"
	"sync/atomic"

	"github.com/Centny/gwf/log"
	"github.com/Centny/gwf/util"
)

var FullError = fmt.Errorf("runner is almost full")

func DoPerf(tc int, logf string, call func(int)) (int64, error) {
	return DoPerfV(tc, tc, logf, call)
}
func DoPerfV(total, tc int, logf string, call func(int)) (int64, error) {
	return DoPerfV_(total, tc, logf,
		func(idx, running int) (int, error) {
			return 1, nil
		}, func(i int) error {
			call(i)
			return nil
		})
}

func DoAutoPerfV(total, tc, peradd int, logf string, precall func(idx, running int) error, call func(int) error) (int64, error) {
	return DoPerfV_(total, tc, logf,
		func(idx, running int) (int, error) {
			terr := precall(idx, running)
			if terr == nil {
				return peradd, nil
			} else if terr == FullError {
				return 0, nil
			} else {
				return 0, terr
			}
		}, call)
}

func DoPerfV_(total, tc int, logf string, increase func(idx, running int) (int, error), call func(int) error) (int64, error) {
	stdout := os.Stdout
	stderr := os.Stderr
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
	ws := sync.WaitGroup{}
	ws.Add(total)
	beg := util.Now()
	var tidx_ int32 = 0
	var run_call, run_next func(int)
	var err error = nil
	var running int32
	run_call = func(v int) {
		atomic.AddInt32(&running, 1)
		terr := call(v)
		atomic.AddInt32(&running, -1)
		if terr != nil {
			err = terr
		}
		if err == nil {
			run_next(v)
		}
		ws.Done()
	}
	run_next = func(v int) {
		nc, terr := increase(v, int(running))
		if terr != nil {
			err = terr
			return
		}
		for i := 0; i < nc; i++ {
			ridx := int(atomic.AddInt32(&tidx_, 1))
			if ridx >= total {
				break
			}
			go run_call(ridx)
		}
	}
	atomic.AddInt32(&tidx_, int32(tc-1))
	for i := 0; i < tc; i++ {
		go run_call(i)
	}
	ws.Wait()
	end := util.Now()
	if len(logf) > 0 {
		os.Stdout.Close()
		os.Stdout = stdout
		os.Stderr = stderr
		log.SetWriter(os.Stdout)
	}
	return end - beg, err
}
