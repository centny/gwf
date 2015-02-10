package impl

import (
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

type RC_V_M_s struct {
}

func (r *RC_V_M_s) OnConn(c netw.Con) bool {
	return true
}
func (r *RC_V_M_s) OnClose(c netw.Con) {
}

func run_s() {
	p := pool.NewBytePool(8, 1024)
	l, vms := NewExecListener_m_j(p, ":7686", &RC_V_M_s{})
	vms.AddHFunc("join", join)
	vms.AddHFunc("replace", replace)
	vms.AddFFunc("^no$", no_f)
	vms.AddFFunc("^to$", to_f)
	vms.AddHFunc("no", no)
	vms.AddHFunc("ferr", ferr)
	vms.AddHFunc("terr", terr)
	vms.AddHFunc("exit", func(r *RCM_Cmd) (interface{}, error) {
		go func() {
			time.Sleep(time.Second)
			l.Close()
		}()
		return &res_v{
			Res: "OK",
		}, nil
	})
	l.T = 500
	err := l.Run()
	if err != nil {
		panic(err.Error())
	}
	l.Wait()
	NewExecListener_m(p, "ssss", nil, nil, nil)
}
func no_f(r *RCM_Cmd) (bool, interface{}, error) {
	return false, &res_v{
		Res: "return",
	}, nil
}
func to_f(r *RCM_Cmd) (bool, interface{}, error) {
	return false, nil, util.Err("filter error")
}
func join(r *RCM_Cmd) (interface{}, error) {
	var arg arg_v
	r.ToS(&arg)
	return &res_v{
		Res: arg.A + arg.B,
	}, nil
}
func replace(r *RCM_Cmd) (interface{}, error) {
	var arg arg_v
	r.ToS(&arg)
	return &res_v{
		Res: strings.Replace(arg.A, arg.B, "+++", -1),
	}, nil
}
func no(r *RCM_Cmd) (interface{}, error) {
	return nil, nil
}
func terr(r *RCM_Cmd) (interface{}, error) {
	return nil, util.Err("rrrrr error")
}
func ferr(r *RCM_Cmd) (interface{}, error) {
	return no, nil
}

type arg_v struct {
	A string `m2s:"a" json:"a"`
	B string `m2s:"b" json:"b"`
}
type res_v struct {
	Res string `m2s:"res" json:"res"`
}

func exec_c(tc *RCM_Con, name, a, b string) {
	var res res_v
	_, err := tc.Exec(name, &arg_v{
		A: a,
		B: b,
	}, &res)
	if err == nil {
		fmt.Println(res.Res)
	} else {
		fmt.Println("err:", err.Error())
	}

}
func exec_c2(tc *RCM_Con, name, a, b string) {
	var res util.Map
	_, err := tc.Exec(name, &arg_v{
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
	np, tc, err := ExecDail_m_j(p, "127.0.0.1:7686", netw.NewCWH(true))
	if err != nil {
		panic(err.Error())
	}
	tc.Start()
	exec_c(tc, "join", "aaabb", "b")
	exec_c2(tc, "join", "aaabb", "b")
	exec_c(tc, "replace", "aaabb", "b")
	exec_c2(tc, "replace", "aaabb", "b")
	exec_c(tc, "no", "a", "b")
	exec_c(tc, "to", "a", "b")
	exec_c(tc, "nerr", "a", "b")
	exec_c(tc, "ferr", "a", "b")
	exec_c(tc, "terr", "a", "b")
	exec_c(tc, "cnerr", "a", "b")
	exec_c(tc, "aerr", "a", "b")
	exec_c(tc, "not_found", "a", "b")
	tc.NAV = func(rc *RCM_Con, name string, args interface{}) (interface{}, error) {
		return nil, util.Err("errrr")
	}
	exec_c(tc, "no", "a", "b")
	tc.NAV = Json_NAV
	// //
	exec_c(tc, "exit", "a", "b")
	fmt.Println(tc.RC_Con.Exec([]byte{1}, nil))
	fmt.Println(tc.RC_Con.Exec(map[string]interface{}{}, nil))
	np.Close()
	exec_c(tc, "no", "a", "b")
	tc.Stop()
	fmt.Println("...end...")
	ExecDail_m(p, "addr", nil, nil, nil, nil)
}
func run_c2() {
	p := pool.NewBytePool(8, 1024)
	rc := NewRC_Runner_m_j("127.0.0.1:7686", p)
	go func() {
		time.Sleep(500 * time.Millisecond)
		rc.Start()
	}()
	fmt.Println("waiting connect--->")
	for i := 0; i < 10; i++ {
		go func(v int) {
			rc.Valid()
			fmt.Println("vvv->", v)
		}(i)
	}
	time.Sleep(time.Second)
}
func TestExecM(t *testing.T) {
	runtime.GOMAXPROCS(runtime.NumCPU() - 1)
	ShowLog = true
	go run_s()
	time.Sleep(100 * time.Millisecond)
	run_c()
	time.Sleep(2 * time.Second)
}
