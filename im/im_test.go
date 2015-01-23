package im

import (
	"fmt"
	"github.com/Centny/gwf/netw"
	"github.com/Centny/gwf/netw/impl"
	"github.com/Centny/gwf/pool"
	"math/rand"
	"runtime"
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
var cc_ws sync.WaitGroup
var cc_ws2 sync.WaitGroup

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
	tc := impl.NewRC_C()
	obdh.AddH(MK_NRC, tc)
	obdh.AddH(MK_NIM, rm)
	ch := impl.NewChanH(obdh)
	tcch := netw.NewCCH(netw.NewDoNoH(), ch)
	ch.Run(5)
	l, con, err := netw.DailN(p, srv.Addr(), tcch, impl.Json_NewCon)
	if err != nil {
		panic(err.Error())
	}
	mcon := impl.NewOBDH_Con(MK_NRC, con)
	rcon := impl.NewRC_Con(mcon, tc)
	rmcon := impl.NewRCM_Con(rcon, impl.Json_NAV)
	rmcon.Start()
	res, err := rmcon.ExecRes("LI", map[string]interface{}{
		"token": "abc",
	})
	if res.Code != 0 {
		panic(res.Res)
	}
	msgc := impl.NewOBDH_Con(MK_NIM, con)
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
		_, err := msgc.Writev(&Msg{
			R: rs,
			T: 0,
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
	rmcon.ExecRes("LO", nil)
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
			continue
		}
		if rm.cc < d {
			continue
		} else {
			break
		}
	}
	cc_ws2.Done()

}
func run_c(db *MemDbH, p *pool.BytePool, rm *rec_msg) {
	crun = true
	cc_ws.Add(200)
	cc_ws2.Add(1)
	for i := 0; i < 4; i++ {
		for i := 0; i < 50; i++ {
			go run_im_c(p, db, rm)
			time.Sleep(time.Millisecond)
		}
		time.Sleep(3 * time.Second)
	}
	go show_cc(db, rm, p)
	cc_ws.Wait()
	wait_rec(db, rm)
	// }
	time.Sleep(3 * time.Second)
}
func run_s(db *MemDbH, p *pool.BytePool) {
	ls := []*Listener{}
	for i := 0; i < 5; i++ {
		l := NewListner(db, fmt.Sprintf("S-vv-%v", i), p, 9890+i,
			impl.Json_V2B, impl.Json_B2V, impl.Json_ND, impl.Json_NAV, impl.Json_VNA)
		err := l.Run()
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
		time.Sleep(2 * time.Second)
	}
	cc_ws2.Wait()
	// time.Sleep(2 * time.Second)
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
	p := pool.NewBytePool(8, 1024)
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
	if m != m_cc_c || (r-pv-e) < rm.cc || d != rm.cc || s_cc_c < rm.cc || r < s_cc_c {
		t.Error(fmt.Sprintf("%v,%v,%v,%v,%v", m != m_cc_c, (r-pv-e) < rm.cc, d != rm.cc, s_cc_c < rm.cc, r < s_cc_c))
	}
	time.Sleep(4 * time.Second)
}

// func TestMap(t *testing.T) {
// 	vv := map[string]map[string]string{}
// 	vv["a"]["b"] = "c"
// }
