package handler

import (
	"encoding/json"
	"fmt"
	"github.com/Centny/gwf/netw"
	"github.com/Centny/gwf/pool"
	"github.com/Centny/gwf/util"
	"runtime"
	"strings"
	"testing"
	"time"
)

/*--------------------------*/
/*------Simple json RC------*/
//
func NewT_RC() *RC_V_M_C {
	return NewRC_Json_M_C(func(rc *RC_V_M_C, fname string, args interface{}) (interface{}, error) {
		if fname == "cnerr" {
			return nil, util.Err("name error")
		}
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
func (r *RC_V_M_s) FNAME(rc *RC_V_Cmd) (string, error) {
	name := rc.StrVal("name")
	if name == "nerr" {
		return "", util.Err("name error")
	} else {
		return name, nil
	}
}
func (r *RC_V_M_s) FARGS(rc *RC_V_Cmd) (*util.Map, error) {
	name := rc.StrVal("name")
	if name == "aerr" {
		return nil, util.Err("args error")
	} else {
		mv := rc.MapVal("args")
		return &mv, nil
	}
}

func run_s() {
	p := pool.NewBytePool(8, 1024)
	vms := NewRC_V_M_S(&RC_V_M_s{})
	vms.AddHFunc("join", join)
	vms.AddHFunc("replace", replace)
	vms.AddFFunc("^no$", no_f)
	vms.AddHFunc("no", no)
	vms.AddHFunc("ferr", ferr)
	ts := NewChan_Json_S(vms)
	l := netw.NewListener(p, ":7686", ts)
	l.T = 500
	vms.AddHFunc("exit", func(r *RC_H_CMD) (interface{}, error) {
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
func no_f(r *RC_H_CMD, vv interface{}) (bool, interface{}, error) {
	return false, &res_v{
		Res: "return",
	}, nil
}
func join(r *RC_H_CMD) (interface{}, error) {
	var arg arg_v
	r.ToS(&arg)
	return &res_v{
		Res: arg.A + arg.B,
	}, nil
}
func replace(r *RC_H_CMD) (interface{}, error) {
	var arg arg_v
	r.ToS(&arg)
	return &res_v{
		Res: strings.Replace(arg.A, arg.B, "+++", -1),
	}, nil
}
func no(r *RC_H_CMD) (interface{}, error) {
	return nil, nil
}
func ferr(r *RC_H_CMD) (interface{}, error) {
	return no, nil
}

type arg_v struct {
	A string `m2s:"a" json:"a"`
	B string `m2s:"b" json:"b"`
}
type res_v struct {
	Res string `m2s:"res" json:"res"`
}

func exec_c(tc *RC_V_M_C, name, a, b string) {
	var res res_v
	err := tc.Exec(name, &arg_v{
		A: a,
		B: b,
	}, &res)
	if err == nil {
		fmt.Println(res.Res)
	} else {
		fmt.Println("err:", err.Error())
	}

}
func exec_c2(tc *RC_V_M_C, name, a, b string) {
	var res util.Map
	err := tc.Exec(name, &arg_v{
		A: a,
		B: b,
	}, &res)
	if err != nil {
		panic(err.Error())
	}
	fmt.Println(res)
}
func run_c() {
	p := pool.NewBytePool(8, 1024)
	tc := NewT_RC()
	c := netw.NewNConPool(p, "127.0.0.1:7686", tc)
	err := c.Dail()
	if err != nil {
		panic(err.Error())
	}
	tc.Start()
	exec_c(tc, "join", "aaabb", "b")
	exec_c2(tc, "join", "aaabb", "b")
	exec_c(tc, "replace", "aaabb", "b")
	exec_c2(tc, "replace", "aaabb", "b")
	exec_c(tc, "no", "a", "b")
	exec_c(tc, "nerr", "a", "b")
	exec_c(tc, "ferr", "a", "b")
	exec_c(tc, "cnerr", "a", "b")
	exec_c(tc, "aerr", "a", "b")
	exec_c(tc, "not_found", "a", "b")
	//
	tc.RC_C.Write([]byte{1})
	tc.RC_C.Exec([]byte{1})
	tc.RC_V_C.Exec(no, nil)
	// exec_c(tc, "exit", "a", "b")
	c.Close()
	exec_c(tc, "no", "a", "b")
	tc.Stop()
	fmt.Println("...end...")
}

func TestVM(t *testing.T) {
	runtime.GOMAXPROCS(runtime.NumCPU() - 1)
	go run_s()
	time.Sleep(100 * time.Millisecond)
	run_c()
	time.Sleep(2 * time.Second)
	NewRC_V_M_C(nil, json.Marshal, json.Unmarshal)
}
