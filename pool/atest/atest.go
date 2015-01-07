package main

import (
	// "container/list"
	"fmt"
	"github.com/Centny/gwf/pool"
	"github.com/Centny/gwf/util"
	"math/rand"
	"runtime"
	"sync"
	"time"
)

func main() {
	TestBytePool()
}

var bp_wg sync.WaitGroup

func run_bp(bp *pool.BytePool) {
	bp_wg.Add(1)
	for i := 0; i < 10000; i++ {
		// ll := []*list.Element{}
		// for i := 0; i < 10; i++ {
		iv := rand.Intn(102400) + 1
		// tv := bp.Alloc(iv)
		// if tv == nil {
		// 	panic("nil")
		// }
		// by := tv.Value.([]byte)
		by := bp.Alloc(iv)
		by[0] = 1
		// ll = append(ll, tv)
		// }
		// for _, tv := range ll {
		bp.Free(by)
		// }
	}
	bp_wg.Done()
}
func TestBytePool() {
	runtime.GOMAXPROCS(runtime.NumCPU() - 1)
	bp := pool.NewBytePool(8, 102400)
	bp.T = 10000
	go func() {
		for {
			time.Sleep(10 * time.Second)
			fmt.Println(bp.GC())
		}
	}()
	b_t := util.Now()
	for i := 0; i < 10000; i++ {
		go run_bp(bp)
	}
	time.Sleep(time.Millisecond)
	bp_wg.Wait()
	e_t := util.Now()
	fmt.Println("used time:", e_t-b_t)
}
