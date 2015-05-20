package im

import (
	"bufio"
	"bytes"
	"code.google.com/p/go.net/websocket"
	"fmt"
	"github.com/Centny/gwf/im/pb"
	"github.com/Centny/gwf/log"
	"github.com/Centny/gwf/netw"
	"github.com/Centny/gwf/netw/impl"
	"github.com/Centny/gwf/pool"
	"github.com/Centny/gwf/util"
	"math/rand"
	"net/http"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

type rec_msg struct {
	cc  *uint64
	con *impl.OBDH_Con
}

func (r *rec_msg) OnCmd(c netw.Cmd) int {
	defer c.Done()
	add_r_cc()
	var msg pb.ImMsg
	_, err := c.V(&msg)
	if err != nil {
		panic(err)
	}
	if len(msg.GetA()) < 1 {
		panic("----xxxxx--A is empty-->")
	}
	_, err = r.con.Writev(map[string]interface{}{
		"i": msg.GetI(),
		"a": msg.GetA(),
	})
	if err != nil {
		panic(err)
	}
	return 0
}

var crun bool = false
var r_cc_c uint64 = 0
var s_cc_c uint64 = 0 //user count ->s
var m_cc_c uint64 = 0 //message count ->s
// var hr_cc_c uint64 = 0 //command count ->r
// var h_cc_c uint64 = 0
var client_c uint64 = 0
var cc_ws sync.WaitGroup
var cc_ws2 sync.WaitGroup

func add_r_cc() {
	atomic.AddUint64(&r_cc_c, 1)
}

// func run_im_nc(p *pool.BytePool, db *MemDbH, rm *rec_msg) {
// 	srvs, err := db.ListSrv("")
// 	if err != nil {
// 		panic(err.Error())
// 	}
// 	if len(srvs) < 1 {
// 		panic("not service")
// 	}
// 	srv := srvs[rand.Intn(len(srvs))]
// 	if len(srv.Token) < 1 {
// 		panic("token is empty")
// 	}
// 	obdh := impl.NewOBDH()
// 	//
// 	//
// 	l, con, err := netw.DailN(p, srv.Addr(), netw.NewCCH(netw.NewDoNotH(), impl.NewChanH2(obdh, 3)), IM_NewCon)
// 	if err != nil {
// 		panic(err.Error())
// 	}
// 	//
// 	//
// 	rcon := impl.NewOBDH_Con(MK_NODE_C, con)
// 	nli_c := impl.NewOBDH_Con(MK_NDC_NLI, rcon)
// 	uli_c := impl.NewOBDH_Con(MK_NDC_ULI, rcon)
// 	ulo_c := impl.NewOBDH_Con(MK_NDC_ULO, rcon)
// 	uur_c := impl.NewOBDH_Con(MK_NDC_UUR, rcon)
// 	c_cb := func(c netw.Cmd) int {
// 		var na NodeV
// 		_, err := c.B2V()(c.Data(), &na)
// 		if err != nil {
// 			panic(err.Error())
// 		}
// 		if na.V.IntVal("code") != 0 {
// 			fmt.Println(na.V)
// 			panic(na.V.StrVal("res"))
// 		}
// 		return 0
// 	}
// 	c_li := func(c netw.Cmd) int {
// 		var na NodeV
// 		_, err := c.B2V()(c.Data(), &na)
// 		if err != nil {
// 			panic(err.Error())
// 		}
// 		if na.V.IntVal("code") != 0 {
// 			fmt.Println(na.V)
// 			panic(na.V.StrVal("res"))
// 		}
// 		uur_c.Writev(
// 			&NodeV{
// 				V: util.Map{
// 					"token": srv.Token,
// 					"R":     na.V.MapVal("res").StrVal("r"),
// 				},
// 				B: "abc",
// 			})
// 		return 0
// 	}
// 	cmdh := impl.NewOBDH()
// 	cmdh.AddF(MK_NDC_NLI, c_cb)
// 	cmdh.AddF(MK_NDC_ULI, c_li)
// 	cmdh.AddF(MK_NDC_ULO, c_cb)
// 	cmdh.AddF(MK_NDC_UUR, c_cb)
// 	//
// 	obdh.AddH(MK_NODE_C, cmdh)
// 	obdh.AddH(MK_NIM, rm)
// 	//
// 	//
// 	//
// 	nli_c.Writev(
// 		&NodeV{
// 			V: util.Map{
// 				"token": srv.Token,
// 			},
// 			B: "abc",
// 		})
// 	uli_c.Writev(
// 		&NodeV{
// 			V: util.Map{
// 				"token": "abc",
// 			},
// 			B: "abc",
// 		})
// 	// fmt.Println(res.MapVal("v"))
// 	// fmt.Println("----->")
// 	atomic.AddUint64(&m_cc_c, 1) //marking for auto create unread message.
// 	atomic.AddUint64(&s_cc_c, 1)
// 	//
// 	msgc := impl.NewOBDH_Con(MK_NODE_M, con)
// 	//
// 	// atomic.AddUint64(&m_cc_c, 1)
// 	// atomic.AddUint64(&s_cc_c, 1)
// 	var tt uint32 = 0
// 	var s string = "U-1"
// 	msgc.Writev(
// 		&pb.ImMsg{
// 			S: &s,
// 			R: []string{"S-Robot"},
// 			T: &tt,
// 			C: []byte{1, 2, 4},
// 		})
// 	cc_ws.Done()
// 	cc_ws2.Wait()
// 	ulo_c.Writev(
// 		&NodeV{
// 			V: util.Map{},
// 			B: "abc",
// 		})
// 	// time.Sleep(1 * time.Second)
// 	l.Close()
// 	con.Close()
// }
func run_im_w(p *pool.BytePool, db *MemDbH) {
	srvs, err := db.ListSrv("")
	if err != nil {
		panic(err.Error())
	}
	if len(srvs) < 1 {
		panic("not service")
	}
	srv := srvs[rand.Intn(len(srvs))]
	wsc, err := websocket.Dial("ws://127.0.0.1"+srv.WsAddr, "", "http://127.0.0.1"+srv.WsAddr)
	if err != nil {
		panic(err.Error())
	}
	li_c := make(chan int, 10)
	var lr string
	go func() {
		li_c <- 0
		r := bufio.NewReader(wsc)
		for {
			bys, err := util.ReadLine(r, 102400, false)
			if err != nil {
				break
			}
			tbys := bytes.SplitN(bys, []byte(WIM_SEQ), 2)
			switch string(tbys[0]) {
			case "m":
				mv, _ := util.Json2Map(string(tbys[1]))
				if len(mv.StrVal("i")) < 1 {
					panic("i is empty")
				}
				wsc.Write([]byte("mr" + WIM_SEQ + util.S2Json(map[string]interface{}{
					"i": mv.StrVal("i"),
					"a": mv.StrVal("a"),
				}) + "\n"))
				add_r_cc()
				// fmt.Println("m-->", string(tbys[1]))
			case "li":
				mv, _ := util.Json2Map(string(tbys[1]))
				lr = mv.StrValP("res/r")
				li_c <- 1
				// fmt.Println("li-->", string(tbys[1]))
			case "ur":
				//do nothing.
				li_c <- 1
				// fmt.Println("ur-->", string(tbys[1]))
			case "lo":
				//do nothing.
				// fmt.Println("lo-->", string(tbys[1]))
			case "mr":
				//do nothing
			default:
				panic("unknow->" + string(bys))
			}
		}
	}()
	<-li_c
	wsc.Write([]byte("li" + WIM_SEQ + util.S2Json(map[string]interface{}{
		"token": "abc",
	}) + "\n"))
	//
	<-li_c
	wsc.Write([]byte("ur" + WIM_SEQ + "{}\n"))
	atomic.AddUint64(&m_cc_c, 1) //marking for auto create unread message.
	atomic.AddUint64(&s_cc_c, 1)
	var tt uint32 = 0
	// wsc.Write([]byte("m" + WIM_SEQ + util.S2Json(&pb.ImMsg{
	// 	R: []string{"S-Robot"},
	// 	T: &tt,
	// 	C: []byte{1, 2, 4},
	// }) + "\n"))
	<-li_c
	if len(lr) < 1 {
		panic("lr is empty")
	}
	time.Sleep(200 * time.Millisecond)
	for i := 0; i < 100; i++ {
		rs := []string{}
		uc := 0
		if i%2 == 0 {
			rs = db.RandUsr(lr)
			uc = len(rs)
		} else {
			gs, _, urs := db.RandGrp()
			rs = []string{gs}
			uc = 0
			for _, ur := range urs {
				if ur == lr {
					continue
				}
				uc++
			}
		}
		for _, r := range rs {
			if len(strings.Trim(r, "\t ")) < 1 {
				rs = []string{}
				break
			}
		}
		if len(rs) < 1 {
			fmt.Println("user not found")
			time.Sleep(500 * time.Millisecond)
			i--
			continue
		}

		mm := &pb.ImMsg{
			R: rs,
			T: &tt,
			C: []byte{1, 2, 4},
		}
		_, err = wsc.Write([]byte("m" + WIM_SEQ + util.S2Json(mm) + "\n"))
		if err != nil {
			panic(err.Error())
		}
		atomic.AddUint64(&s_cc_c, uint64(uc))
		atomic.AddUint64(&m_cc_c, 1)
		time.Sleep(time.Millisecond)
	}
	cc_ws.Done()
	cc_ws2.Wait()
	wsc.Write([]byte("lo" + WIM_SEQ + "{}\n"))
	time.Sleep(1 * time.Second)
	wsc.Close()
	// atomic.AddUint64(&hr_cc_c, tc.RCC)
	// atomic.AddUint64(&h_cc_c, tcch.RCC)
	// fmt.Print("run_im_w end...")
}
func run_im_c(p *pool.BytePool, db *MemDbH) {
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
	fmt.Println("connect to server by local addres:", con.LocalAddr())
	//
	//MK_NRC
	tc := impl.NewRC_C()
	mcon := impl.NewOBDH_Con(MK_NRC, con)
	rcon := impl.NewRC_Con(mcon, tc)
	rcon.Start()
	//
	rm := &rec_msg{
		con: impl.NewOBDH_Con(MK_NMR, con),
	}
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
	var lr string = res.MapVal("res").StrVal("r")
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
	// msgc.Writev(
	// 	&pb.ImMsg{
	// 		R: []string{"S-Robot"},
	// 		T: &tt,
	// 		C: []byte{1, 2, 4},
	// 	})
	//
	//
	// nodec_m := impl.NewOBDH_Con(MK_NODE_M, con)
	// nodec := impl.NewOBDH_Con(MK_NODE, con)
	time.Sleep(200 * time.Millisecond)
	for i := 0; i < 100; i++ {
		rs := []string{}
		uc := 0
		if i%2 == 0 {
			rs = db.RandUsr(lr)
			uc = len(rs)
		} else {
			gs, _, urs := db.RandGrp()
			rs = []string{gs}
			uc = 0
			for _, ur := range urs {
				if ur == lr {
					continue
				}
				uc++
			}
		}
		for _, r := range rs {
			if len(strings.Trim(r, "\t ")) < 1 {
				rs = []string{}
				break
			}
		}
		if len(rs) < 1 {
			fmt.Println("user not found")
			time.Sleep(500 * time.Millisecond)
			i--
			continue
		}
		mm := &pb.ImMsg{
			R: rs,
			T: &tt,
			C: []byte{1, 2, 4},
		}
		_, err = msgc.Writev(mm)
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
	// atomic.AddUint64(&hr_cc_c, tc.RCC)
	// atomic.AddUint64(&h_cc_c, tcch.RCC)
	// fmt.Print("h_cc_c end...")
}
func show_cc(db *MemDbH, p *pool.BytePool) {
	for {
		time.Sleep(4 * time.Second)
		mlen, rlen, plen, elen, dlen := db.Show_()
		fmt.Printf("Waiting->M(s):%v, R(r):%v==S(s):%v, MemS:%v, mlen(%v), rlen(%v), plen(%v), elen(%v), dlen(%v)\n",
			m_cc_c, r_cc_c, s_cc_c, p.Size(), mlen, rlen, plen, elen, dlen)
	}

}
func wait_rec(db *MemDbH) {
	for {
		time.Sleep(4 * time.Second)
		m, _, _, _, d := db.Show()
		if m < m_cc_c {
			fmt.Printf("Waiting msg(r:%v),done(s:%v)\n", m, m_cc_c)
			continue
		}
		if r_cc_c != (d + db.mr_n_cc) {
			fmt.Printf("Waiting rec(%v),done(%v)\n", r_cc_c, d+db.mr_n_cc)
			continue
		} else {
			break
		}
	}
	cc_ws2.Done()

}
func run_c(db *MemDbH, p *pool.BytePool) {
	crun = true
	xl, yl := 5, 6
	client_c = uint64(xl * yl * 2)
	cc_ws.Add(xl * yl * 2)
	cc_ws2.Add(1)
	for i := 0; i < xl; i++ {
		for j := 0; j < yl; j++ {
			go run_im_c(p, db)
			go run_im_w(p, db)
			// go run_im_nc(p, db, rm)
			time.Sleep(time.Millisecond)
		}
		time.Sleep(2 * time.Second)
	}
	go show_cc(db, p)
	cc_ws.Wait()
	wait_rec(db)
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
				gr, uc, _ := db.RandGrp()
				psrv.PushN("U-1", gr, "abc", 0)
				atomic.AddUint64(&s_cc_c, uint64(uc))
				atomic.AddUint64(&m_cc_c, 1)
			}
			for i := 0; i < 5; i++ {
				ur := db.RandUsr("")
				psrv.PushN("U-1", strings.Join(ur, ","), "abc", 0)
				atomic.AddUint64(&s_cc_c, uint64(len(ur)))
				atomic.AddUint64(&m_cc_c, 1)
			}
		}
		time.Sleep(3 * time.Second)
	}()
	for i := 0; i < 5; i++ {
		l := NewListner3(db, fmt.Sprintf("S-vv-%v", i), p, 9890+i, 1000000)
		l.WsAddr = fmt.Sprintf(":%v", 9870+i)
		l.PushSrvAddr = "127.0.0.1:5598"
		rc := make(chan int)
		go func() {
			rc <- 1
			hs := &http.Server{
				Handler: l.WIM_L.WsS(),
				Addr:    l.WsAddr,
			}
			hs.ListenAndServe()
		}()
		<-rc
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
		log.W("----->NewListner2->%v", 9890+i)
		time.Sleep(time.Duration(i+1) * time.Second)
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
	p := pool.NewBytePool(8, 102400)
	go db.GrpBuilder()
	go run_s(db, p)
	time.Sleep(2000 * time.Millisecond)
	run_c(db, p)
	fmt.Printf("Done->M:%v, R:%v==S:%v\n", m_cc_c, r_cc_c, s_cc_c)
	p.T = 10
	time.Sleep(100 * time.Millisecond)
	p.GC()
	fmt.Println("MS:", p.Size())
	m, r, _, _, d := db.Show()
	if m != m_cc_c || (d+db.mr_n_cc) != r_cc_c || s_cc_c < r_cc_c || r < s_cc_c {
		t.Error(fmt.Sprintf("%v,%v,%v,%v", m != m_cc_c, (d+db.mr_n_cc) != r_cc_c, s_cc_c < r_cc_c, r < s_cc_c))
	}
	time.Sleep(time.Second)
}
