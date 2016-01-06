package tutil

import (
	"github.com/Centny/gwf/log"
	"github.com/Centny/gwf/util"
	"os"
	"sync"
	"sync/atomic"
)

func DoPerf(tc int, logf string, call func(int)) (int64, error) {
	return DoPerfV(tc, tc, logf, call)
}

func DoPerfV(total, tc int, logf string, call func(int)) (int64, error) {
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
	var run_call func(int)
	run_call = func(v int) {
		call(v)
		ridx := int(atomic.AddInt32(&tidx_, 1))
		if ridx < total {
			go run_call(ridx)
		}
		ws.Done()
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
	return end - beg, nil
}
