package im

import (
	"github.com/Centny/gwf/im/pb"
	"github.com/Centny/gwf/log"
	"github.com/Centny/gwf/netw"
	"github.com/Centny/gwf/netw/impl"
	"github.com/Centny/gwf/pool"
	"github.com/Centny/gwf/util"
	"math/rand"
	"sync/atomic"
	"time"
)

type IMC struct {
	*netw.NConRunner
	C     *impl.RC_Con
	P     *pool.BytePool
	OnM   func(i *IMC, c netw.Cmd, m *pb.ImMsg) int
	Token string
	IC    Con
	MCon  *impl.OBDH_Con
	LC    chan int
	//
	obdh  *impl.OBDH
	tc    *impl.RC_C
	hbing bool
	HBT   time.Duration
	RC    uint64 //receive message count.
}

func NewIMC(srv, token string) *IMC {
	p := pool.NewBytePool(8, 1024000)
	imc := &IMC{
		OnM: func(i *IMC, c netw.Cmd, m *pb.ImMsg) int {
			return 0
		},
		obdh:  impl.NewOBDH(),
		tc:    impl.NewRC_C(),
		P:     p,
		Token: token,
		LC:    make(chan int),
		HBT:   1000 * time.Millisecond,
	}
	imc.obdh.AddH(MK_NRC, imc.tc)
	imc.obdh.AddH(MK_NIM, imc)
	imc.NConRunner = netw.NewNConRunnerN(p, srv, impl.NewChanH2(imc.obdh, 5), IM_NewCon)
	imc.NConRunner.TickData = []byte{}
	imc.ConH = imc
	imc.C = impl.NewRC_Con(nil, imc.tc) //initial con after connected.
	imc.C.Start()
	log_d("creating IMC by %v", srv)
	return imc
}
func NewIMC2(srv *Srv, token string) *IMC {
	return NewIMC(srv.Addr(), token)
}
func NewIMC3(srvs []Srv, token string) *IMC {
	srv := srvs[rand.Intn(len(srvs))]
	return NewIMC2(&srv, token)
}
func NewIMC4(sl, token string) (*IMC, error) {
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
	return NewIMC3(ssl, token), nil
}
func (i *IMC) OnCmd(c netw.Cmd) int {
	var msg pb.ImMsg
	_, err := c.V(&msg)
	if err != nil {
		log.E("convert values(%v) to IM msg error:%v", c.Data(), err.Error())
		return -1
	}
	log.D("receive message->%v", msg)
	atomic.AddUint64(&i.RC, 1)
	return i.OnM(i, c, &msg)
}
func (i *IMC) OnConn(c netw.Con) bool {
	go i.login(c) //must async for exec remove command.
	return true
}
func (i *IMC) login(c netw.Con) {
	i.C.Con = impl.NewOBDH_Con(MK_NRC, c)
	//
	var res util.Map
	_, err := i.C.Execm(MK_NRC_LI, map[string]interface{}{
		"token": i.Token,
	}, &res)
	if err != nil {
		i.StopRunner()
		log.E("IM login by token(%v) err->%v", i.Token, err)
		i.LC <- 1
		return
	}
	if res.IntVal("code") != 0 {
		i.StopRunner()
		log.E("IM login by token(%v) err->%v", i.Token, res)
		i.LC <- 1
		return
	}
	i.MCon = impl.NewOBDH_Con(MK_NIM, c)
	util.Json2S(util.S2Json(res.Val("res")), &i.IC)
	c.SetWait(true)
	log.D("IMC login succes by token(%v)->%v", i.Token, i.IC)
	i.LC <- 0
}
func (i *IMC) HB(data string) (string, error) {
	var res util.Map
	_, err := i.C.Execm(MK_NRC_HB, map[string]interface{}{
		"D": data,
	}, &res)
	return res.StrVal("D"), err
}
func (i *IMC) rhb(delay time.Duration) {
	var times_ time.Duration = 0
	i.hbing = true
	for i.hbing {
		d, err := i.HB("D->")
		if err == nil && d == "D->" {
			times_++
			log.D("HB(%v) success, will retry after %v", d, times_*delay)
			time.Sleep(times_ * delay)
		} else {
			times_ = 0
			log.W("HB(%v) error->%v", d, err)
		}
	}
}
func (i *IMC) StartHB() {
	go i.rhb(i.HBT)
}
func (i *IMC) SMS(s string, t int, c string) (int, error) {
	return i.SMS_V([]string{s}, t, []byte(c))
}
func (i *IMC) SMS_V(rs []string, t int, c []byte) (int, error) {
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
}

func (i *IMC) Close() {
	if i.C != nil {
		i.C.Stop()
	}
	if i.NConRunner != nil {
		i.StopRunner()
	}
}
