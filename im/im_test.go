package im

import (
	"fmt"
	"github.com/Centny/gwf/im/pb"
	"github.com/Centny/gwf/netw"
	"github.com/Centny/gwf/netw/impl"
	"github.com/Centny/gwf/pool"
	"github.com/Centny/gwf/util"
	"math/rand"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

type rec_msg struct {
	cc uint64
}

func (r *rec_msg) OnCmd(c netw.Cmd) int {
	defer c.Done()
	atomic.AddUint64(&r.cc, 1)
	if len(c.Data()) == 3 {
		panic("data error")
	}
	return 0
}

var crun bool = false
var s_cc_c uint64 = 0
var m_cc_c uint64 = 0
var hr_cc_c uint64 = 0
var h_cc_c uint64 = 0
var client_c uint64 = 0
var cc_ws sync.WaitGroup
var cc_ws2 sync.WaitGroup

func run_im_nc(p *pool.BytePool, db *MemDbH, rm *rec_msg) {
	srvs, err := db.ListSrv("")
	if err != nil {
		panic(err.Error())
	}
	if len(srvs) < 1 {
		panic("not service")
	}
	srv := srvs[rand.Intn(len(srvs))]
	if len(srv.Token) < 1 {
		panic("token is empty")
	}
	obdh := impl.NewOBDH()
	//
	//
	l, con, err := netw.DailN(p, srv.Addr(), netw.NewCCH(netw.NewDoNotH(), impl.NewChanH2(obdh, 3)), IM_NewCon)
	if err != nil {
		panic(err.Error())
	}
	//
	//
	tc := impl.NewRC_C()
	mcon := impl.NewOBDH_Con(MK_NODE_C, con)
	rcon := impl.NewRC_Con(mcon, tc)
	rcon.Start()
	//
	obdh.AddH(MK_NODE_C, tc)
	obdh.AddH(MK_NIM, rm)
	//
	//
	//
	var res util.Map
	_, err = rcon.Execm(MK_NDC_NLI,
		&NodeV{
			V: util.Map{
				"token": srv.Token,
			},
			B: "abc",
		}, &res)
	if err != nil {
		panic(err.Error())
	}
	if res.MapVal("v").IntVal("code") != 0 {
		fmt.Println(res)
		panic(res.StrVal("res"))
	}
	rcon.Execm(MK_NDC_ULI,
		&NodeV{
			V: util.Map{
				"token": "abc",
			},
			B: "abc",
		}, &res)
	if res.MapVal("v").IntVal("code") != 0 {
		fmt.Println(res)
		panic(res.StrVal("res"))
	}
	// fmt.Println(res.MapVal("v"))
	rcon.Execm(MK_NDC_UUR, &NodeV{
		V: util.Map{
			"token": srv.Token,
			"R":     res.MapVal("v").MapVal("res").StrVal("r"),
		},
		B: "abc",
	}, &res)
	if res.MapVal("v").IntVal("code") != 0 {
		fmt.Println(res)
		panic(res.StrVal("res"))
	}
	// fmt.Println("----->")
	atomic.AddUint64(&m_cc_c, 1) //marking for auto create unread message.
	atomic.AddUint64(&s_cc_c, 1)
	//
	msgc := impl.NewOBDH_Con(MK_NODE_M, con)
	//
	// atomic.AddUint64(&m_cc_c, 1)
	// atomic.AddUint64(&s_cc_c, 1)
	var tt uint32 = 0
	var s string = "U-1"
	msgc.Writev(
		&pb.ImMsg{
			S: &s,
			R: []string{"S-Robot"},
			T: &tt,
			C: []byte{1, 2, 4},
		})
	cc_ws.Done()
	cc_ws2.Wait()
	rcon.Execm(MK_NDC_ULO, map[string]string{}, &res)
	// time.Sleep(1 * time.Second)
	l.Close()
	con.Close()
}
func run_im_c(p *pool.BytePool, db *MemDbH, rm *rec_msg) {
	srvs, err := db.ListSrv("")
	if err != nil {
		panic(err.Error())
	}
	if len(srvs) < 1 {
		panic("not service")
	}
	srv := srvs[rand.Intn(len(srvs))]
	obdh := impl.NewOBDH()
	//
	//
	l, con, err := netw.DailN(p, srv.Addr(), netw.NewCCH(netw.NewDoNotH(), impl.NewChanH2(obdh, 5)), IM_NewCon)
	if err != nil {
		panic(err.Error())
	}
	//
	//MK_NRC
	tc := impl.NewRC_C()
	mcon := impl.NewOBDH_Con(MK_NRC, con)
	rcon := impl.NewRC_Con(mcon, tc)
	rcon.Start()
	//
	obdh.AddH(MK_NRC, tc)
	obdh.AddH(MK_NIM, rm)
	//
	//
	//
	var res util.Map
	_, err = rcon.Execm(MK_NRC_LI, map[string]interface{}{
		"token": "abc",
	}, &res)
	if err != nil {
		panic(err.Error())
	}
	if res.IntVal("code") != 0 {
		fmt.Println(res)
		panic(res.StrVal("res"))
	}
	//
	_, err = rcon.Execm(MK_NRC_UR, map[string]interface{}{}, &res)
	// fmt.Println("----->")
	atomic.AddUint64(&m_cc_c, 1) //marking for auto create unread message.
	atomic.AddUint64(&s_cc_c, 1)
	//
	msgc := impl.NewOBDH_Con(MK_NIM, con)
	//
	// atomic.AddUint64(&m_cc_c, 1)
	// atomic.AddUint64(&s_cc_c, 1)
	var tt uint32 = 0
	msgc.Writev(
		&pb.ImMsg{
			R: []string{"S-Robot"},
			T: &tt,
			C: []byte{1, 2, 4},
		})
	//
	//
	// nodec_m := impl.NewOBDH_Con(MK_NODE_M, con)
	// nodec := impl.NewOBDH_Con(MK_NODE, con)

	for i := 0; i < 1000; i++ {
		rs := []string{}
		uc := 0
		if i%2 == 0 {
			rs = db.RandUsr()
			uc = len(rs)
		} else {
			gs, uc_ := db.RandGrp()
			rs = []string{gs}
			uc = uc_
		}
		if len(rs) < 1 {
			fmt.Println("user not found")
			time.Sleep(500 * time.Millisecond)
			continue
		}
		_, err := msgc.Writev(
			&pb.ImMsg{
				R: rs,
				T: &tt,
				C: []byte{1, 2, 4},
			})
		if err != nil {
			panic(err.Error())
		}
		atomic.AddUint64(&s_cc_c, uint64(uc))
		atomic.AddUint64(&m_cc_c, 1)
		time.Sleep(time.Millisecond)
	}
	cc_ws.Done()
	cc_ws2.Wait()
	// rmcon.ExecRes("LO", nil)
	rcon.Execm(MK_NRC_LO, map[string]string{}, &res)
	if err != nil {
		panic(err.Error())
	}
	if res.IntVal("code") != 0 {
		fmt.Println(res)
		panic(res.StrVal("res"))
	}
	time.Sleep(1 * time.Second)
	l.Close()
	atomic.AddUint64(&hr_cc_c, tc.RCC)
	// atomic.AddUint64(&h_cc_c, tcch.RCC)
}
func show_cc(db *MemDbH, rm *rec_msg, p *pool.BytePool) {
	for {
		time.Sleep(4 * time.Second)
		fmt.Printf("Waiting->M:%v, R:%v==S:%v, HR:%v, H:%v, MS:%v\n", m_cc_c, rm.cc, s_cc_c, hr_cc_c, h_cc_c, p.Size())
	}

}
func wait_rec(db *MemDbH, rm *rec_msg) {
	for {
		time.Sleep(4 * time.Second)
		m, _, _, _, d := db.Show()
		if m < m_cc_c {
			fmt.Printf("Waiting msg(%v),done(%v)\n", m, m_cc_c)
			continue
		}
		if rm.cc < d {
			fmt.Printf("Waiting rec(%v),done(%v)\n", d, rm.cc)
			continue
		} else {
			break
		}
	}
	cc_ws2.Done()

}
func run_c(db *MemDbH, p *pool.BytePool, rm *rec_msg) {
	crun = true
	client_c = 400
	cc_ws.Add(400)
	cc_ws2.Add(1)
	for i := 0; i < 4; i++ {
		for i := 0; i < 50; i++ {
			go run_im_c(p, db, rm)
			go run_im_nc(p, db, rm)
			time.Sleep(time.Millisecond)
		}
		time.Sleep(2 * time.Second)
	}
	go show_cc(db, rm, p)
	cc_ws.Wait()
	wait_rec(db, rm)
	// }
	time.Sleep(3 * time.Second)
}
func run_s(db *MemDbH, p *pool.BytePool) {
	psrv := NewPushSrv(p, ":5598", "Push", netw.NewDoNotH(), db)
	err := psrv.Run()
	if err != nil {
		panic(err.Error())
	}
	ls := []*Listener{}
	go func() {
		for len(psrv.Cons()) < 1 || len(db.Grp) < 1 {
			time.Sleep(500 * time.Millisecond)
		}
		for i := 0; i < 5; i++ {
			for i := 0; i < 5; i++ {
				gr, uc := db.RandGrp()
				psrv.PushN("U-1", gr, "abc", 0)
				atomic.AddUint64(&s_cc_c, uint64(uc))
				atomic.AddUint64(&m_cc_c, 1)
			}
			for i := 0; i < 5; i++ {
				ur := db.RandUsr()
				psrv.PushN("U-1", strings.Join(ur, ","), "abc", 0)
				atomic.AddUint64(&s_cc_c, uint64(len(ur)))
				atomic.AddUint64(&m_cc_c, 1)
			}
		}
		time.Sleep(3 * time.Second)
	}()
	for i := 0; i < 5; i++ {
		l := NewListner2(db, fmt.Sprintf("S-vv-%v", i), p, 9890+i)
		l.T = 30000
		l.PushSrvAddr = "127.0.0.1:5598"
		err = l.Run()
		if err != nil {
			panic(err.Error())
		}
		ls = append(ls, l)
		// go func() {
		// 	for {
		// 		time.Sleep(time.Second)
		// 		fmt.Println("-->", l.NIM.DC, l.NIM.SS.(*MarkConPoolSender).EC)
		// 	}
		// }()
		time.Sleep(3 * time.Second)
	}
	cc_ws2.Wait()

	time.Sleep(2 * time.Second)
	for _, l := range ls {
		l.Close()
	}
}
func TestIm(t *testing.T) {
	runtime.GOMAXPROCS(runtime.NumCPU() - 1)
	// ShowLog = true
	// impl.ShowLog = true
	// netw.ShowLog = true
	db := NewMemDbH()
	rm := &rec_msg{}
	p := pool.NewBytePool(8, 102400)
	go db.GrpBuilder()
	go run_s(db, p)
	time.Sleep(100 * time.Millisecond)
	run_c(db, p, rm)
	fmt.Printf("Done->M:%v, R:%v==S:%v, HR:%v, H:%v\n", m_cc_c, rm.cc, s_cc_c, hr_cc_c, h_cc_c)
	p.T = 10
	time.Sleep(100 * time.Millisecond)
	p.GC()
	fmt.Println("MS:", p.Size())
	m, r, pv, e, d := db.Show()
	if m != m_cc_c || (r-pv-e) < (rm.cc-client_c) || d != (rm.cc-client_c) || s_cc_c < (rm.cc-client_c) || r < s_cc_c {
		t.Error(fmt.Sprintf("%v,%v,%v,%v,%v", m != m_cc_c, (r-pv-e) < (rm.cc-client_c), d != (rm.cc-client_c), s_cc_c < (rm.cc-client_c), r < s_cc_c))
	}
	time.Sleep(4 * time.Second)
}
