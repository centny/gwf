package main

import (
	"fmt"
	"github.com/Centny/gwf/netw"
	"github.com/Centny/gwf/netw/handler"
	"github.com/Centny/gwf/pool"
	"github.com/Centny/gwf/util"
	"os"
	"runtime"
	"strings"
	"time"
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
	return true
}
func (r *RC_V_M_s) OnClose(c *netw.Con) {
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
	if os.Args[1] == "-s" {
		run_s()
		return
	}
	if len(os.Args) < 4 {
		fmt.Println("less 4")
		return
	}
	run_c(os.Args[1], os.Args[2], os.Args[3])
}

func run_s() {
	p := pool.NewBytePool(8, 1024)
	vms := handler.NewRC_V_M_S(&RC_V_M_s{})
	vms.AddHFunc("join", join)
	vms.AddHFunc("replace", replace)
	ts := handler.NewChan_Json_S(vms)
	l := netw.NewListener(p, ":7686", ts)
	l.T = 500
	vms.AddHFunc("exit", func(r *handler.RC_V_M_S, rc *handler.RC_V_Cmd, args *util.Map) (interface{}, error) {
		go func() {
			time.Sleep(time.Second)
			ts.Stop()
			l.Close()
		}()
		return &res_v{
			Res: "OK",
		}, nil
	})
	ts.Run(5)
	err := l.Run()
	if err != nil {
		panic(err.Error())
	}
	ts.Wait()
	l.Wait()
}
func join(r *handler.RC_V_M_S, rc *handler.RC_V_Cmd, args *util.Map) (interface{}, error) {
	var arg arg_v
	args.ToS(&arg)
	return &res_v{
		Res: arg.A + arg.B,
	}, nil
}
func replace(r *handler.RC_V_M_S, rc *handler.RC_V_Cmd, args *util.Map) (interface{}, error) {
	var arg arg_v
	args.ToS(&arg)
	return &res_v{
		Res: strings.Replace(arg.A, arg.B, "+++", -1),
	}, nil
}

type arg_v struct {
	A string `m2s:"a" json:"a"`
	B string `m2s:"b" json:"b"`
}
type res_v struct {
	Res string `m2s:"res" json:"res"`
}

func run_c(name, a, b string) {
	fmt.Println("R:", name, a, b)
	p := pool.NewBytePool(8, 1024)
	tc := NewRC()
	c := netw.NewNConPool(p, "127.0.0.1:7686", tc)
	err := c.Dail()
	if err != nil {
		panic(err.Error())
	}
	tc.Start()
	var res res_v
	err = tc.Exec(name, &arg_v{
		A: a,
		B: b,
	}, &res)
	if err != nil {
		panic(err.Error())
	}
	fmt.Println(res.Res)
	fmt.Println("...end...")
}
