package main

import (
	"fmt"
	"github.com/Centny/gwf/netw"
	"github.com/Centny/gwf/netw/handler"
	"github.com/Centny/gwf/pool"
	"github.com/Centny/gwf/util"
	"math/rand"
	"os"
	"runtime"
	"sync"
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
	fmt.Println("closing ", c.C.RemoteAddr().String())
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
		run_c()
	}
}

func run_s() {
	p := pool.NewBytePool(8, 1024)
	vms := handler.NewRC_V_M_S(&RC_V_M_s{})
	vms.AddHFunc("plus", plus)
	vms.AddHFunc("sub", sub)
	netw.ShowLog = true
	ts := handler.NewChan_Json_S(vms)
	l := netw.NewListener(p, ":7686", ts)
	l.T = 3000
	ts.Run(runtime.NumCPU())
	err := l.Run()
	if err != nil {
		panic(err.Error())
	}
	ts.Wait()
	l.Wait()
}
func plus(r *handler.RC_V_M_S, rc *handler.RC_V_Cmd, args *util.Map) (interface{}, error) {
	var arg arg_v
	err := args.ValidF(`
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
func sub(r *handler.RC_V_M_S, rc *handler.RC_V_Cmd, args *util.Map) (interface{}, error) {
	var arg arg_v
	err := args.ValidF(`
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

type arg_v struct {
	A int64 `m2s:"a" json:"a"`
	B int64 `m2s:"b" json:"b"`
}
type res_v struct {
	Res int64 `m2s:"res" json:"res"`
}

func run_c() {
	p := pool.NewBytePool(8, 1024)
	tc := NewRC()
	c := netw.NewNConPool(p, "127.0.0.1:7686", tc)
	err := c.Dail()
	if err != nil {
		panic(err.Error())
	}
	tc.Start()
	beg := util.Now()
	var errc int64 = 0
	var ecount int64 = 0
	wg := sync.WaitGroup{}
	// wg.Add(50000 * 2)
	// var adddd_ int64 = 0
	for i := 0; i < 80000; i++ {
		go func() {
			wg.Add(1)
			// atomic.AddInt64(&adddd_, 1)
			var res res_v
			arg := &arg_v{
				A: int64(rand.Intn(1000) + 1),
				B: int64(rand.Intn(1000) + 1),
			}
			err = tc.Exec("plus", arg, &res)
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
			err = tc.Exec("sub", arg, &res)
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
	// fmt.Println("----->:", adddd_)
	end := util.Now()
	fmt.Println("beg:", beg)
	fmt.Println("end:", end)
	fmt.Println("used:", end-beg)
	fmt.Println("ecount:", ecount)
	fmt.Println("errc:", errc)
	fmt.Println("...end...")
}
