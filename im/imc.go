package im

import (
	"github.com/Centny/gwf/im/pb"
	"github.com/Centny/gwf/log"
	"github.com/Centny/gwf/netw"
	"github.com/Centny/gwf/netw/impl"
	"github.com/Centny/gwf/pool"
	"github.com/Centny/gwf/tutil"
	"github.com/Centny/gwf/util"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
)

type IMC struct {
	*netw.NConRunner
	C     *impl.RC_Con
	MRC   *impl.OBDH_Con
	P     *pool.BytePool
	OnM   func(i *IMC, c netw.Cmd, m *pb.ImMsg) int
	Token string
	IC    Con
	MCon  *impl.OBDH_Con
	LC    sync.WaitGroup
	//
	obdh    *impl.OBDH
	tc      *impl.RC_C
	hbing   bool
	logined bool
	HBT     time.Duration
	RC      uint64 //receive message count.
	HbLog   bool
	hb_wg   sync.WaitGroup
	chanh   *impl.ChanH
}

func NewIMC(p *pool.BytePool, srv, token string) *IMC {
	imc := &IMC{
		OnM: func(i *IMC, c netw.Cmd, m *pb.ImMsg) int {
			go i.MR(m.GetA(), m.GetI())
			return 0
		},
		obdh:  impl.NewOBDH(),
		tc:    impl.NewRC_C(),
		P:     p,
		Token: token,
		HBT:   1000 * time.Millisecond,
		hb_wg: sync.WaitGroup{},
	}
	imc.chanh = impl.NewChanH2(imc.obdh, 5)
	imc.obdh.AddH(MK_NRC, imc.tc)
	imc.obdh.AddH(MK_NIM, imc)
	imc.NConRunner = netw.NewNConRunnerN(p, srv, imc.chanh, IM_NewCon)
	imc.NConRunner.TickData = []byte{}
	imc.ConH = imc
	imc.C = impl.NewRC_Con(nil, imc.tc) //initial con after connected.
	log_d("creating IMC by %v", srv)
	return imc
}
func NewIMC2(p *pool.BytePool, srv *Srv, token string) *IMC {
	return NewIMC(p, srv.Addr(), token)
}
func NewIMC3(p *pool.BytePool, srvs []Srv, token string) *IMC {
	srv := srvs[rand.Intn(len(srvs))]
	return NewIMC2(p, &srv, token)
}
func NewIMC4(p *pool.BytePool, sl, token string) (*IMC, error) {
	ssm, err := util.HGet2(sl)
	if err != nil {
		return nil, err
	}
	if ssm.IntVal("code") != 0 {
		return nil, util.Err("listSrv(%v) err->%v", sl, ssm)
	}
	var ssl []Srv
	util.Json2Ss(util.S2Json(ssm.Val("data")), &ssl)
	if len(ssl) < 1 {
		return nil, util.Err("im server not found on listSrv(%v) by %v", sl, ssm)
	}
	return NewIMC3(p, ssl, token), nil
}
func NewIMC5(p *pool.BytePool, srv string, ls bool, token string) (*IMC, error) {
	if ls {
		return NewIMC4(p, srv, token)
	} else {
		return NewIMC(p, srv, token), nil
	}
}
func (i *IMC) hblog(f string, args ...interface{}) {
	if i.HbLog {
		log.D_(1, f, args...)
	}
}
func (i *IMC) Start() {
	i.LC.Add(1)
	i.C.Start()
	i.NConRunner.StartRunner()
}
func (i *IMC) OnCmd(c netw.Cmd) int {
	defer c.Done()
	var msg pb.ImMsg
	_, err := c.V(&msg)
	if err != nil {
		log.E("convert values(%v) to IM msg error:%v", c.Data(), err.Error())
		return -1
	}
	// log.D("receive message->%v", msg)
	atomic.AddUint64(&i.RC, 1)
	// go func() {
	// 	if err := i.MR(msg.GetI()); err != nil {
	// 		log.E("mark msg(%v) recv err:%v", msg, err.Error())
	// 	}
	// }()
	return i.OnM(i, c, &msg)
}
func (i *IMC) OnConn(c netw.Con) bool {
	go i.login(c) //must async for exec remove command.
	return true
}
func (i *IMC) login(c netw.Con) {
	defer i.LC.Done()
	log.D("doing login by token(%v)", i.Token)
	i.C.Con = impl.NewOBDH_Con(MK_NRC, c)
	i.MRC = impl.NewOBDH_Con(MK_NMR, c)
	//
	var res util.Map
	_, err := i.C.Execm(MK_NRC_LI, map[string]interface{}{
		"token": i.Token,
	}, &res)
	if err != nil {
		log.E("IM login by token(%v) err->%v", i.Token, err)
		i.logined = false
		i.StopRunner()
		return
	}
	if res.IntVal("code") != 0 {
		log.E("IM login by token(%v) err->%v", i.Token, res)
		i.logined = false
		i.StopRunner()
		return
	}
	i.MCon = impl.NewOBDH_Con(MK_NIM, c)
	util.Json2S(util.S2Json(res.Val("res")), &i.IC)
	c.SetWait(true)
	log.D("IMC login succes by token(%v)->%v", i.Token, i.IC)
	i.logined = true
	// i.UR()
}
func (i *IMC) HB(data string) (string, error) {
	var res util.Map
	_, err := i.C.Execm(MK_NRC_HB, map[string]interface{}{
		"D": data,
	}, &res)
	return res.StrVal("D"), err
}
func (i *IMC) UR() error {
	var res util.Map
	_, err := i.C.Execm(MK_NRC_UR, map[string]interface{}{}, &res)
	if err != nil {
		return err
	}
	if res.IntVal("code") == 0 {
		return nil
	} else {
		return util.Err("%v", res.StrVal("err"))
	}
}
func (i *IMC) MR(a, mid string) error {
	if i.MRC == nil {
		panic("not start connect")
	}
	_, err := i.MRC.Writev(map[string]interface{}{
		"i": mid,
		"a": a,
	})
	return err
}
func (i *IMC) GR(gr []string) (map[string][]string, error) {
	if i.MRC == nil {
		panic("not start connect")
	}
	if len(gr) < 1 {
		return nil, nil
	}
	log_d("sending GR by (%v)", gr)
	var res util.Map
	_, err := i.C.Execm(MK_NRC_GR, map[string]interface{}{
		"gr": gr,
	}, &res)
	if err != nil {
		return nil, err
	}
	if res.IntVal("code") != 0 {
		return nil, util.Err("%v", res.StrVal("err"))
	}
	var ur map[string][]string
	err = util.Json2S(util.S2Json(res.Val("res")), &ur)
	return ur, err
}
func (i *IMC) rhb(delay time.Duration) {
	var times_ time.Duration = 0
	i.hbing = true
	log.D("running HB by delay(%v)...", delay)
	for i.hbing {
		time.Sleep(times_ * delay)
		d, err := i.HB("D->")
		if err == nil && d == "D->" {
			times_++
			i.hblog("HB(%v) success, will retry after %v", d, times_*delay)
		} else {
			times_ = 0
			log.W("HB(D->) error->%v", err)
		}
	}
	log.D("HB is stopped...")
	i.hb_wg.Done()
}
func (i *IMC) StartHB() {
	log.D("IMC starting HB by delay(%v)", i.HBT)
	i.hb_wg.Add(1)
	go i.rhb(i.HBT)
}
func (i *IMC) SMS(s string, t int, c string) (int, error) {
	return i.SMS_V([]string{s}, t, []byte(c))
}
func (i *IMC) SMS_V(rs []string, t int, c []byte) (int, error) {
	if i.MCon == nil {
		panic("not start connect")
	}
	var tt uint32 = uint32(t)
	mm := &pb.ImMsg{
		R: rs,
		T: &tt,
		C: c,
	}
	return i.MCon.Writev(mm)
}
func (i *IMC) OnClose(c netw.Con) {
	log.D("IMC OnClose...")
	i.LC.Add(1)
}
func (i *IMC) Logined() bool {
	return i.logined
}
func (i *IMC) Close() {
	log.D("IMC closing...")
	i.hbing = false
	if i.NConRunner != nil {
		i.StopRunner()
	}
	i.chanh.Stop()
	i.chanh.Wait()
	if i.C != nil {
		i.C.Stop()
	}
	i.tc.Close()
	i.hb_wg.Wait()
}

func (i *IMC) StartMonitor() {
	i.chanh.M = tutil.NewMonitor()
}
func (i *IMC) State() (interface{}, error) {
	if i.chanh.M == nil {
		return nil, nil
	} else {
		return i.chanh.M.State()
	}
}
