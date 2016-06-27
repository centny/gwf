package rc

import (
	"fmt"
	"github.com/Centny/gwf/log"
	"github.com/Centny/gwf/netw"
	"github.com/Centny/gwf/netw/impl"
	"github.com/Centny/gwf/pool"
	"github.com/Centny/gwf/routing"
	"github.com/Centny/gwf/tutil"
	"github.com/Centny/gwf/util"
	"net/http"
	"os"
	"runtime"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

////////////////////////////////////
//////// server handler ////////////
////////////////////////////////////
type rc_s_h struct {
	L   *RC_Listener_m
	cid int64
}

func (r *rc_s_h) OnCmd(c netw.Cmd) int {
	var args util.Map
	_, err := c.V(&args)
	if err != nil {
		log.E("rc_s_h V error:%v", err.Error())
		return -1
	}
	var cmd, name, msg string
	err = args.ValidF(`
		c,R|S,L:0;
		n,O|S,L:0;
		m,O|S,L:0;
		`, &cmd, &name, &msg)
	if err != nil {
		log.E("rc_s_h valid args error:%v", err.Error())
		return -1
	}
	switch cmd {
	case "l":
		//login by message connection.(usage 1)
		if len(name) < 1 {
			log.E("login name is empty")
			return -1
		}
		r.L.AddC_c(name, c)
		log.D("S(m)->login success by name(%v)", name)
		return 0
	case "m":
		log.D("S(m)->receive message(%v) from %v", msg, c.RemoteAddr())
		return 0
	default:
		log.W("S(m)->unknow command(%v)", cmd)
		return -1
	}
}
func (r *rc_s_h) OnConn(c netw.Con) bool {
	c.SetWait(true)
	return true
}
func (r *rc_s_h) OnClose(c netw.Con) {
}

//login by remote command connection(usage 2)
func (r *rc_s_h) Login(rc *impl.RCM_Cmd) (interface{}, error) {
	var name string
	err := rc.ValidF(`
		n,O|S,L:0;
		`, &name)
	if err != nil {
		log.E("rc_s_h Login valid args error:%v", err.Error())
		return nil, err
	}
	if len(name) < 1 {
		log.E("login name is empty")
		return nil, util.Err("login name is empty")
	}
	r.L.AddC_rc(name, rc)
	log.D("S(c)->login success by name(%v)", name)
	return util.Map{
		"code": 0,
	}, nil
}

//calling by message command.
func (r *rc_s_h) CallM(target string) error {
	log.D("call1->%v", target)
	if len(target) > 0 {
		//sending command to target client.
		mc := r.L.MsgC(target)
		if mc == nil {
			return util.Err("connection not found by id(%v)", target)
		}
		_, err := mc.Writev(util.Map{
			"c": "m",
			"m": "server message",
		})
		return err
	} else {
		//sending command to all client.
		for cid, mc := range r.L.MsgCs() {
			_, err := mc.Writev(util.Map{
				"c": "m",
				"m": "server message",
			})
			if err != nil {
				log.E("sending message to %v err:%v", cid, err.Error())
			}
		}
	}
	return nil
}

//calling by remote command.
func (r *rc_s_h) CallC(target string) error {
	log.D("call2->%v", target)
	if len(target) > 0 {
		//sending command to target client.
		cc := r.L.CmdC(target)
		if cc == nil {
			return util.Err("connection not found by id(%v)", target)
		}
		var res []string
		_, err := cc.Exec("list", nil, &res)
		log.D("exec list res->%v,err:%v", res, err)
		return err
	} else {
		//sending command to all client.
		for cid, cc := range r.L.CmdCs() {
			var res []string
			_, err := cc.Exec("list", nil, &res)
			log.D("exec list to %v res->%v,err:%v", cid, res, err)
		}
	}
	return nil
}

//handler all remote command
func (r *rc_s_h) Handle(l *RC_Listener_m) {
	r.L = l
	l.AddHFunc("login", r.Login)
}
func (r *rc_s_h) OnLogin(rc *impl.RCM_Cmd, token string) (string, error) {
	if token == "abc3" {
		return "", util.Err("error")
	}
	cid := atomic.AddInt64(&r.cid, 1)
	time.Sleep(10 * time.Millisecond)
	return fmt.Sprintf("NN-%v", cid), nil
}

////////////////////////////////////
//////// client handler ////////////
////////////////////////////////////
type rc_c_h struct {
	R *RC_Runner_m
}

func (r *rc_c_h) OnCmd(c netw.Cmd) int {
	var args util.Map
	_, err := c.V(&args)
	if err != nil {
		log.E("rc_c_h V error:%v", err.Error())
		return -1
	}
	var cmd, msg string
	err = args.ValidF(`
		c,R|S,L:0;
		m,O|S,L:0;
		`, &cmd, &msg)
	if err != nil {
		log.E("rc_c_h valid args error:%v", err.Error())
		return -1
	}
	switch cmd {
	case "m":
		log.D("C(m)->receive message(%v) from %v", msg, c.RemoteAddr())
		return 0
	default:
		log.W("unknow command(%v)", cmd)
		return -1
	}
}
func (r *rc_c_h) OnConn(c netw.Con) bool {
	return true
}
func (r *rc_c_h) OnClose(c netw.Con) {
}

//client command
func (r *rc_c_h) List(rc *impl.RCM_Cmd) (interface{}, error) {
	log.D("C(c)->receive list command")
	time.Sleep(10 * time.Millisecond)
	return []string{"a", "b", "c"}, nil
}

//handler all client command
func (r *rc_c_h) Handle(run *RC_Runner_m) {
	r.R = run
	r.R.AddHFunc("list", r.List)
}

////////////////////////////////////
//////// testing runner ////////////
////////////////////////////////////
func TestB2V(t *testing.T) {
	var bys = []byte{91, 34, 97, 34, 44, 34, 98, 34, 44, 34, 99, 34, 93}
	fmt.Println("xxx")
	fmt.Println(string(bys), "--->")
}

func TestRc(t *testing.T) {
	test_rc(t, false)
}
func test_rc(t *testing.T, multi bool) {
	runtime.GOMAXPROCS(util.CPU())
	impl.ShowLog = true
	netw.ShowLog = true
	netw.ShowLog_C = true
	bp := pool.NewBytePool(8, 102400)
	//
	//
	//initial server.
	sh := &rc_s_h{}
	lm := NewRC_Listener_m_j(bp, ":10801", sh)
	lm.Multi = multi
	lm.LCH = sh
	lm.RCM_S.M = tutil.NewMonitor()
	lm.RCM_S.ShowSlow = 1
	sh.Handle(lm)
	err := lm.Run()
	if err != nil {
		t.Error(err.Error())
		return
	}
	fmt.Println("=== xxxx->001\n\n\n")
	//
	//
	//initial client
	crs := []*RC_Runner_m{}
	//login by message command.
	for i := 0; i < 5; i++ {
		ch := &rc_c_h{}
		cr := NewRC_Runner_m_j(bp, "127.0.0.1:10801", ch)
		fmt.Println(cr.RC_Runner_m, "-->")
		if multi {
			cr.Uuid = util.UUID()
		}
		cr.StartMonitor()
		cr.SetShowSlow(1)
		ch.Handle(cr)
		cr.Start()
		fmt.Printf("xx->%v\n", "001")
		_, err := cr.Writev(util.Map{
			"c": "l",
			"n": fmt.Sprintf("RC-%v", i),
		})
		fmt.Printf("xx->%v\n", "002")
		_, err = cr.Writeb([]byte(util.S2Json(util.Map{
			"c": "m",
			"m": "server message",
		})))
		fmt.Printf("xx->%v\n", "003")
		_, err = cr.Writev2([]byte{}, util.Map{
			"c": "m",
			"m": "server message",
		})
		if err != nil {
			t.Error(err.Error())
			return
		}
		crs = append(crs, cr)
	}
	fmt.Println("=== xxxx->002\n\n\n")
	//login by remote command.
	for i := 5; i < 10; i++ {
		addrs := []string{}
		// for j := 0; j < i-3; j++ {
		addrs = append(addrs, "127.0.0.1:10801")
		// }
		ch := &rc_c_h{}
		cr := NewRC_Runner_m_j(bp, strings.Join(addrs, ","), ch)
		if multi {
			cr.Uuid = util.UUID()
		}
		cr.StartMonitor()
		cr.SetShowSlow(1)
		ch.Handle(cr)
		cr.Start()
		time.Sleep(time.Second)
		name := fmt.Sprintf("RC-%v", i)
		res, err := cr.VExec_m("login", util.Map{
			"n": name,
		})
		if err != nil {
			t.Error(err.Error())
			return
		}
		log.D("login by name(%v)->%v", name, res.IntVal("code"))
		crs = append(crs, cr)
		val, _ := cr.State()
		fmt.Println("CR->", util.S2Json(val))
		//
		err = cr.HB("abc")
		if err != nil {
			t.Error(err.Error())
			return
		}
	}
	fmt.Println("=== xxxx->003\n\n\n")
	// time.Sleep(1000 * time.Second)
	// return
	//
	//
	//calling target.
	for i := 0; i < 10; i++ {
		fmt.Println("xxxx->004-0")
		err = sh.CallM(fmt.Sprintf("RC-%v", i))
		if err != nil {
			t.Error(err.Error())
			return
		}
		fmt.Println("xxxx->004-1")
		err = sh.CallC(fmt.Sprintf("RC-%v", i))
		if err != nil {
			t.Error(err.Error())
			return
		}
		fmt.Println("xxxx->004-2\n")
	}
	fmt.Println("=== xxxx->004\n\n\n")
	//calling all
	err = sh.CallM("")
	if err != nil {
		t.Error(err.Error())
		return
	}
	err = sh.CallC("")
	if err != nil {
		t.Error(err.Error())
		return
	}
	if !lm.Exist("RC-5") {
		t.Error("not exist")
		return
	}
	//
	val, _ := lm.RCM_S.M.State()
	fmt.Println("LM->", util.S2Json(val))
	//
	//test other
	lm.MsgC("not found")
	lm.CmdC("not found")
	lm.RCH.delc(netw.NewCon_(nil, nil, nil))
	lm.RCH.OnCmd(&netw.Cmd_{Con: netw.NewCon_(nil, nil, nil)})
	//
	//
	//close
	time.Sleep(time.Second)
	for _, cr := range crs {
		cr.Close()
		cr.Stop()
	}
	time.Sleep(time.Second)
	lm.Close()
	time.Sleep(time.Second)
}

type rc_login_h struct {
}

func (r *rc_login_h) OnCmd(c netw.Cmd) int {
	return 0
}
func (r *rc_login_h) OnConn(c netw.Con) bool {
	c.SetWait(true)
	return true
}
func (r *rc_login_h) OnClose(c netw.Con) {
}
func (r *rc_login_h) OnLogin(rc *impl.RCM_Cmd, token string) (string, error) {
	return "", util.Err("error")
}

func TestRcLogin(t *testing.T) {
	runtime.GOMAXPROCS(util.CPU())
	// impl.ShowLog = true
	bp := pool.NewBytePool(8, 102400)
	//
	//
	//initial server.
	sh := &rc_login_h{}
	lm := NewRC_Listener_m_j(bp, ":10801", sh)
	lm.AddToken2([]string{"abc"})
	lm.AddToken(map[string]int{
		"abc1": 1,
		"abc2": 2,
	})
	err := lm.Run()
	if err != nil {
		t.Error(err.Error())
		return
	}
	login := func(token string, tw, terr bool) {
		cr := NewRC_Runner_m_j(bp, "127.0.0.1:10801", sh)
		cr.Start()
		err = cr.Login_(token)
		if tw {
			err = cr.Login_(token)
		}
		if terr && err != nil {
			t.Error(err.Error())
		}
		if !terr && err == nil {
			t.Error("error")
		}
		cr.Stop()
		cr.Wait()
	}
	login("abc", false, true)
	login("abc1", false, true)
	login("abc1", true, false)
	login("", false, false)
	login("xxxxx", false, false)
	lm.LCH = sh
	login("abc", false, false)
	lm.Close()
}

func TestErr(t *testing.T) {
	runtime.GOMAXPROCS(util.CPU())
	bp := pool.NewBytePool(8, 102400)
	cr := NewRC_Runner_m_j(bp, "127.0.0.1:980x1", nil)
	go func() {
		fmt.Println("starting....")
		fmt.Println(cr.Writeb(nil))
		fmt.Println(cr.Writev(nil))
		fmt.Println(cr.Writev2(nil, nil))
	}()
	time.Sleep(100 * time.Millisecond)
	fmt.Println(cr.Waitingc())
	cr.Timeout()
	time.Sleep(100 * time.Millisecond)
	fmt.Println(cr.Waitingc())
	cr.Timeout()
	time.Sleep(100 * time.Millisecond)
	fmt.Println(cr.Waitingc())
	cr.Timeout()
	// cr.Timeout()
	// cr.Timeout()
	cr.Start()
	time.Sleep(time.Second)
}

func pref_exec(rc *impl.RCM_Cmd) (interface{}, error) {
	var dc int64
	err := rc.ValidF(`
		dc,R|I,R:0;
		`, &dc)
	if err != nil {
		log.E("pref_exec valid args error:%v", err.Error())
		return nil, err
	}
	return util.Map{
		"code": 0,
		"data": make([]byte, dc),
	}, nil
}

func pref_rc() (int64, int64, error) {
	os.Remove("rc_t.log")
	bp := pool.NewBytePool(8, 102400)
	lm := NewRC_Listener_m_j(bp, ":10812", netw.NewDoNotH())
	lm.AddHFunc("exec", pref_exec)
	err := lm.Run()
	if err != nil {
		return 0, 0, err
	}
	cr := NewRC_Runner_m_j(bp, "127.0.0.1:10812", netw.NewDoNotH())
	cr.Start()
	var fail int64 = 0
	used, _ := tutil.DoPerf(10000, "rc_t.log", func(i int) {
		res, err := cr.VExec_m("exec", util.Map{"dc": i + 1})
		if err != nil {
			fail++
			fmt.Println(err.Error())
			return
		}
		if res.IntVal("code") != 0 {
			panic("not zero")
		}
	})
	lm.Close()
	cr.Stop()
	return used, fail, nil
}

func pref_exec2(hs *routing.HTTPSession) routing.HResult {
	var dc int64
	err := hs.ValidCheckVal(`
		dc,R|I,R:0;
		`, &dc)
	if err != nil {
		log.E("pref_exec2 valid args error:%v", err.Error())
		return hs.MsgResErr2(1, "arg-err", err)
	} else {
		return hs.MsgRes(make([]byte, dc*100))
	}
}

func pref_http() (int64, int64, error) {
	os.Remove("http_t.log")
	mux := routing.NewSessionMux2("/")
	srv := http.Server{
		Addr:    ":10803",
		Handler: mux,
	}
	go srv.ListenAndServe()
	time.Sleep(100 * time.Microsecond)
	mux.HFunc("^.*$", pref_exec2)
	var fail int64 = 0
	used, _ := tutil.DoPerf(10000, "http_t.log", func(i int) {
		res, err := util.HGet2("http://127.0.0.1:10803?dc=%v", i+1)
		if err != nil {
			fail++
			fmt.Println(err.Error())
			return
		}
		if res.IntVal("code") != 0 {
			panic("not zero")
		}
	})
	return used, fail, nil
}

func TestPerformance(t *testing.T) {
	runtime.GOMAXPROCS(util.CPU())
	// impl.ShowLog = true
	used, fail, err := pref_http()
	if err != nil {
		t.Error(err.Error())
		return
	}
	fmt.Printf(`
------------------------------
HTTP->Used:%vms,Fail:%v
------------------------------

			
`, used, fail)
	//
	used, fail, err = pref_rc()
	if err != nil {
		t.Error(err.Error())
		return
	}
	fmt.Printf(`
------------------------------
RC->Used:%vms,Fail:%v
------------------------------

`, used, fail)
}

func TestRelogin(t *testing.T) {
	bp := pool.NewBytePool(8, 102400)
	rcs := NewRC_Listener_m_j(bp, ":2311", netw.NewDoNotH())
	rcs.AddToken2([]string{"abc"})
	err := rcs.Run()
	if err != nil {
		t.Error(err.Error())
		return
	}
	rl := &AutoLoginH{}
	rcc := NewRC_Runner_m_j(bp, "127.0.0.1:2311", netw.NewCCH(rl, netw.NewDoNotH()))
	rl.Runner = rcc
	rl.Token = "abc"
	rcc.Start()
	time.Sleep(2 * time.Second)
	rcc.RC_Con.Close()
	time.Sleep(time.Second)
}
