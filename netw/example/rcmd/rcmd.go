package main

import (
	"fmt"
	"github.com/Centny/gwf/netw"
	"github.com/Centny/gwf/netw/handler"
	"github.com/Centny/gwf/pool"
	"github.com/Centny/gwf/util"
	"math/rand"
	"net/http"
	_ "net/http/pprof"
	"os"
	"runtime"
	"runtime/debug"
	"runtime/pprof"
	"sync"
	"time"
	// "sync/atomic"
)

/*--------------------------*/
/*------Simple json RC------*/
//
func NewRC() *handler.RC_V_M_C {
	return handler.NewRC_Json_M_C(func(rc *handler.RC_V_M_C, fname string, args interface{}) (interface{}, error) {
		return map[string]interface{}{
			"name": fname,
			"args": args,
		}, nil
	})
}

type RC_V_M_s struct {
}

func (r *RC_V_M_s) OnConn(c *netw.Con) bool {
	c.SetWait(true)
	return true
}
func (r *RC_V_M_s) OnClose(c *netw.Con) {
	fmt.Println("closing ", c.RemoteAddr().String())
}
func (r *RC_V_M_s) FNAME(rc *handler.RC_V_Cmd) (string, error) {
	return rc.StrVal("name"), nil
}
func (r *RC_V_M_s) FARGS(rc *handler.RC_V_Cmd) (*util.Map, error) {
	mv := rc.MapVal("args")
	return &mv, nil
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU() - 1)
	if len(os.Args) < 2 {
		fmt.Println("less 2")
		return
	}
	switch os.Args[1] {
	case "-s":
		run_s()
	case "-c":
		if len(os.Args) < 4 {
			fmt.Println("use:rcmd -c addr loop|run")
			return
		}
		run_c(os.Args[2], os.Args[3])
	}
}

func run_s() {
	p := pool.NewBytePool(8, 1024)
	p.T = 10000
	go p.GC()
	vms := handler.NewRC_V_M_S(&RC_V_M_s{})
	vms.AddHFunc("plus", plus)
	vms.AddHFunc("sub", sub)
	vms.AddHFunc("panic", panic_c)
	vms.AddHFunc("gc", func(r *handler.RC_H_CMD) (interface{}, error) {
		v1, v2 := p.GC()
		return &res_v{
			Msg: fmt.Sprintf("%v--%v", v1, v2),
		}, nil
	})
	vms.AddHFunc("heap", dump_heap)
	netw.ShowLog = true
	ts := handler.NewChan_Json_S(vms)
	l := netw.NewListener(p, ":7686", ts)
	l.T = 3000
	ts.Run(2)
	err := l.Run()
	if err != nil {
		panic(err.Error())
	}
	// go func() {
	// 	for {
	// 		time.Sleep(5 * time.Second)
	// 		os.Stderr.WriteString("<-------------------------------------------------------------------------------------------------------------------------------->")
	// 		pprof.Lookup("goroutine").WriteTo(os.Stderr, 1)
	// 	}
	// }()
	http.ListenAndServe(":8899", nil)
	ts.Wait()
	l.Wait()
}
func plus(r *handler.RC_H_CMD) (interface{}, error) {
	var arg arg_v
	err := r.ValidF(`
		a,R|I,R:0;
		b,R|I,R:0;
		`, &arg.A, &arg.B)
	if err == nil {
		// fmt.Println(arg)
		return &res_v{
			Res: arg.A + arg.B,
		}, nil
	} else {
		// fmt.Println("err:", err.Error())
		return nil, err
	}
}
func sub(r *handler.RC_H_CMD) (interface{}, error) {
	var arg arg_v
	err := r.ValidF(`
		a,R|I,R:0;
		b,R|I,R:0;
		`, &arg.A, &arg.B)
	if err == nil {
		return &res_v{
			Res: arg.A - arg.B,
		}, nil
	} else {
		return nil, err
	}
}
func panic_c(r *handler.RC_H_CMD) (interface{}, error) {
	pprof.Lookup("goroutine").WriteTo(os.Stderr, 1)
	return &res_v{
		Msg: "SSS->\n" + string(debug.Stack()),
	}, nil
}

var f_cc int = 0

func dump_heap(r *handler.RC_H_CMD) (interface{}, error) {
	fmt.Println("----->")
	f_cc++
	f, err := os.Create(fmt.Sprintf("heap-%v.prof", f_cc))
	if err != nil {
		fmt.Println("----->", err)
		return nil, err
	}
	defer f.Close()
	err = pprof.WriteHeapProfile(f)
	return &res_v{
		Msg: "OK",
	}, err
}

type arg_v struct {
	A int64 `m2s:"a" json:"a"`
	B int64 `m2s:"b" json:"b"`
}
type res_v struct {
	Res int64  `m2s:"res" json:"res"`
	Msg string `m2s:"msg" json:"msg"`
}

func run_c(addr, cmd string) {
	p := pool.NewBytePool(8, 1024)
	tc := NewRC()
	c := netw.NewNConPool(p, addr, tc)
	go p.GC()
	_, err := c.Dail()
	if err != nil {
		panic(err.Error())
	}
	defer c.Close()
	tc.Start()
	defer tc.Stop()
	var errc int64 = 0
	var ecount int64 = 0
	var used int64 = 0
	switch cmd {
	case "run":
		errc, ecount, used = run_go(tc, 80000)
	case "loop":
		for {
			fmt.Println(run_go(tc, 100000))
			time.Sleep(time.Second)
		}
	case "panic":
		var res res_v
		fmt.Println(tc.Exec("panic", nil, &res))
		fmt.Println(res.Msg)
	case "gc":
		var res res_v
		fmt.Println(tc.Exec("gc", nil, &res))
		fmt.Println(res.Msg)
	case "heap":
		var res res_v
		fmt.Println(tc.Exec("heap", nil, &res))
		fmt.Println(res.Msg)
	default:
		return
	}
	// wg.Add(50000 * 2)
	// var adddd_ int64 = 0

	// fmt.Println("----->:", adddd_)
	fmt.Println("used:", used)
	fmt.Println("ecount:", ecount)
	fmt.Println("errc:", errc)
	fmt.Println("...end...")

}

func run_go(tc *handler.RC_V_M_C, count int) (int64, int64, int64) {
	beg := util.Now()
	var ecount int64 = 0
	var errc int64 = 0
	wg := sync.WaitGroup{}
	for i := 0; i < count; i++ {
		go func() {
			wg.Add(1)
			// atomic.AddInt64(&adddd_, 1)
			var res res_v
			arg := &arg_v{
				A: int64(rand.Intn(1000) + 1),
				B: int64(rand.Intn(1000) + 1),
			}
			err := tc.Exec("plus", arg, &res)
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
			var res res_v
			arg := &arg_v{
				A: int64(rand.Intn(1000) + 1),
				B: int64(rand.Intn(1000) + 1),
			}
			err := tc.Exec("sub", arg, &res)
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
