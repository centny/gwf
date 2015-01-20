package impl

import (
	"fmt"
	"github.com/Centny/gwf/netw"
	"github.com/Centny/gwf/pool"
	"github.com/Centny/gwf/util"
	"math/rand"
	"runtime"
	"sync"
	"testing"
	"time"
)

func TestExecPer(t *testing.T) {
	runtime.GOMAXPROCS(runtime.NumCPU() - 1)
	vv := make(chan int)
	go execPer_run_s(vv)
	execPer_run_c("127.0.0.1:7686", "run")
	vv <- 1
	time.Sleep(time.Second)
}

type execPer struct {
}

func (r *execPer) OnConn(c netw.Con) bool {
	c.SetWait(true)
	return true
}
func (r *execPer) OnClose(c netw.Con) {
	fmt.Println("closing ", c.RemoteAddr().String())
}

func execPer_run_s(vv chan int) {
	p := pool.NewBytePool(8, 1024)
	p.T = 10000
	netw.ShowLog = true
	go p.GC()
	l, cc, vms := NewChanExecListener_m_j(p, ":7686", &execPer{})
	vms.AddHFunc("plus", execPer_plus)
	vms.AddHFunc("sub", execPer_sub)
	vms.AddHFunc("res", execPer_res)
	l.T = 3000
	cc.Run(2)
	err := l.Run()
	if err != nil {
		panic(err.Error())
	}
	go func() {
		<-vv
		l.Close()
		cc.Stop()
	}()
	cc.Wait()
	l.Wait()

	//
	Json_NewCon(nil, nil, nil)
}
func execPer_plus(r *RCM_Cmd) (interface{}, error) {
	var arg ex_arg_v
	err := r.ValidF(`
		a,R|I,R:0;
		b,R|I,R:0;
		`, &arg.A, &arg.B)
	if err == nil {
		// fmt.Println(arg)
		return &ex_res_v{
			Res: arg.A + arg.B,
		}, nil
	} else {
		// fmt.Println("err:", err.Error())
		return nil, err
	}
}
func execPer_sub(r *RCM_Cmd) (interface{}, error) {
	var arg ex_arg_v
	err := r.ValidF(`
		a,R|I,R:0;
		b,R|I,R:0;
		`, &arg.A, &arg.B)
	if err == nil {
		return &ex_res_v{
			Res: arg.A - arg.B,
		}, nil
	} else {
		return nil, err
	}
}
func execPer_res(r *RCM_Cmd) (interface{}, error) {
	return r.CRes(0, "OKK")
}

type ex_arg_v struct {
	A int64 `m2s:"a" json:"a"`
	B int64 `m2s:"b" json:"b"`
}
type ex_res_v struct {
	Res int64  `m2s:"res" json:"res"`
	Msg string `m2s:"msg" json:"msg"`
}

func execPer_run_c(addr, cmd string) {
	p := pool.NewBytePool(8, 1024)
	go p.GC()
	l, tc, err := ExecDail_m_j(p, addr)
	if err != nil {
		panic(err.Error())
	}
	defer l.Close()
	tc.Start()
	defer tc.Stop()
	errc, ecount, used := execPer_run_go(tc, 80000)
	// wg.Add(50000 * 2)
	// var adddd_ int64 = 0

	// fmt.Println("----->:", adddd_)
	fmt.Println("used:", used)
	fmt.Println("ecount:", ecount)
	fmt.Println("errc:", errc)
	fmt.Println("...end...")

}

func execPer_run_go(tc *RCM_Con, count int) (int64, int64, int64) {
	beg := util.Now()
	var ecount int64 = 0
	var errc int64 = 0
	res, err := tc.ExecRes("res", nil)
	if err != nil {
		panic(err)
	}
	fmt.Println(res.Res, "----->")
	wg := sync.WaitGroup{}
	for i := 0; i < count; i++ {
		go func() {
			wg.Add(1)
			// atomic.AddInt64(&adddd_, 1)
			var res ex_res_v
			arg := &ex_arg_v{
				A: int64(rand.Intn(1000) + 1),
				B: int64(rand.Intn(1000) + 1),
			}
			_, err := tc.Exec("plus", arg, &res)
			if err != nil {
				// fmt.Println(err.Error(), errc)
				errc++
			} else if res.Res != (arg.A + arg.B) {
				fmt.Println(res.Res, arg.A, arg.B)
				ecount++
			}
			wg.Done()
			// atomic.AddInt64(&adddd_, -1)
		}()
		go func() {
			wg.Add(1)
			// atomic.AddInt64(&adddd_, 1)
			var res ex_res_v
			arg := &ex_arg_v{
				A: int64(rand.Intn(1000) + 1),
				B: int64(rand.Intn(1000) + 1),
			}
			_, err := tc.Exec("sub", arg, &res)
			if err != nil {
				// fmt.Println(err.Error(), errc)
				errc++
			} else if res.Res != (arg.A - arg.B) {
				fmt.Println(res.Res, arg.A, arg.B)
				ecount++
			}
			wg.Done()
			// atomic.AddInt64(&adddd_, -1)
		}()
	}
	wg.Wait()
	end := util.Now()
	return errc, ecount, end - beg
}
