package pool

import (
	"fmt"
	"github.com/Centny/gwf/util"
	"math/rand"
	"runtime"
	"sync"
	"testing"
	"time"
)

func TestBytePool(t *testing.T) {
	tbp()
	go func() {
		defer func() {
			fmt.Println(recover())
		}()
		NewBytePool(0, 0)
	}()
}

var bp_wg sync.WaitGroup

func run_bp(bp *BytePool) {
	bp_wg.Add(1)
	for i := 0; i < 1000; i++ {
		// ll := []*list.Element{}
		// for i := 0; i < 10; i++ {
		iv := rand.Intn(102400) + 1
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
func tbp() {
	runtime.GOMAXPROCS(util.CPU())
	bp := NewBytePool(8, 102400)
	bp.T = 1000
	go func() {
		for {
			time.Sleep(time.Second)
			fmt.Println(bp.GC())
		}
	}()
	b_t := util.Now()
	for i := 0; i < 1000; i++ {
		go run_bp(bp)
	}
	time.Sleep(time.Millisecond)
	bp_wg.Wait()
	e_t := util.Now()
	fmt.Println("used time:", e_t-b_t, "size:", bp.Size())
	//
	// bp.Alloc(0)
	tv := bp.Alloc(8)
	fmt.Println(tv)
	bp.Free(tv)
	bp.Free(nil)
	bp.Free([]byte{})
}

func TestTr(t *testing.T) {
	bys1 := make([]byte, 5)
	bys1[0] = 1
	bys1[1] = 2
	bys1[2] = 3
	bys2 := bys1[:3]
	fmt.Println(bys1, bys2)
	fmt.Println(&bys1[0], &bys2[0])
}
