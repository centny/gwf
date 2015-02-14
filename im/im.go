package im

import (
	"code.google.com/p/go-uuid/uuid"
	"fmt"
	"github.com/Centny/gwf/im/pb"
	"github.com/Centny/gwf/log"
	"github.com/Centny/gwf/netw"
	"github.com/Centny/gwf/netw/impl"
	"github.com/Centny/gwf/pool"
	"github.com/Centny/gwf/util"
	"github.com/golang/protobuf/proto"
	"net"
	"strings"
	"time"
)

var ShowLog bool = false

func log_d(f string, args ...interface{}) {
	if ShowLog {
		log.D_(1, f, args...)
	}
}

const (
	MK_NIM    = 0
	MK_NRC    = 4
	MK_DIM    = 8
	MK_DRC    = 12
	MK_NODE   = 30
	MK_NODE_M = 31
)

type DbH interface {
	netw.ConHandler
	//
	//
	AddCon(c *Con) error
	DelCon(sid, cid, r string, t byte, ct int) (*Con, error)
	OnCloseCon(c netw.Con, sid, cid string, t byte) error
	//list all connection by target R
	ListCon(rs []string) ([]Con, error)
	//list push task by server id and message id.
	ListPushTask(sid, mid string) (*Msg, []Con, error)
	//
	//
	//find current con user R.
	FUsrR(c netw.Cmd) string
	//list all user R by group R,if gr is nil return all online user R.
	ListUsrR(gr []string) (map[string][]string, error)
	//sift the R to group R and user R.
	Sift(rs []string) ([]string, []string, error)
	//
	//
	AddSrv(srv *Srv) error
	DelSrv(sid string) error
	//find the server by token
	FindSrv(token string) (*Srv, error)
	//list all online server,exclue special server id.
	ListSrv(sid string) ([]Srv, error)
	//
	//
	//user login,return user R.
	OnLogin(c netw.Cmd, args *util.Map) (string, int, error)
	OnLogout(c netw.Cmd, args *util.Map) (string, int, bool, error)
	//
	//
	NewMid() string
	//update the message R status
	Update(mid string, rs map[string]string) error
	//store mesage
	Store(m *Msg) error
	//send unread message
	ListUnread(r string, ct int) ([]Msg, error)
}

//
type Finder interface {
	Find(id string) netw.Con
}
type Sender interface {
	Id() string
	Send(cid string, v interface{}) error
	// SendUnRead(r string) error
}

type MultiFinder struct {
	FS []Finder
}

func NewMultiFinder(fs ...Finder) *MultiFinder {
	return &MultiFinder{
		FS: fs,
	}
}
func (m *MultiFinder) Find(id string) netw.Con {
	for _, f := range m.FS {
		cc := f.Find(id)
		if cc == nil {
			continue
		} else {
			return cc
		}
	}
	return nil
}

type MarkConPoolSender struct {
	Mark []byte
	CP   Finder
	Id_  string
	// EC   uint64
	// lck  sync.RWMutex
}

func NewMarkConPoolSender(mark []byte, cp Finder, sid string) *MarkConPoolSender {
	return &MarkConPoolSender{
		Mark: mark,
		CP:   cp,
		Id_:  sid,
	}
}
func (m *MarkConPoolSender) Id() string {
	return m.Id_
}
func (m *MarkConPoolSender) Send(cid string, v interface{}) error {
	cc := m.CP.Find(cid)
	if cc == nil {
		return util.Err("con not found by id(%v) in pool(%v)", cid, m.Id_)
	} else {
		return m.SendC(cc, v)
	}
}

func (m *MarkConPoolSender) SendC(con netw.Con, v interface{}) error {
	bys, err := con.V2B()(v)
	if err != nil {
		return err
	}
	// mm := v.(*Msg)
	// fmt.Println(fmt.Sprintf("%v", mm), string(bys))
	// m.lck.Lock()
	// defer m.lck.Unlock()
	_, err = con.Writeb(m.Mark, bys)
	// if err == nil || vv < len(bys) {
	// atomic.AddUint64(&m.EC, 1)
	// }
	return err
}

type Listener struct {
	*netw.Listener
	Obdh    *impl.OBDH
	NIM     *NIM_Rh
	DIP     *DimPool
	DIM     *DIM_Rh
	Db      DbH
	P       *pool.BytePool
	Port    int
	Sid     string
	PubHost string
	PubPort int
	Err_    netw.CmdErrF
	//
	PushSrvAddr   string
	PushSrvTick   time.Duration
	PushConRunner *netw.NConRunner
}

func NewListner(db DbH, sid string, p *pool.BytePool, port int, v2b netw.V2Byte, b2v netw.Byte2V, nd impl.ND_F, nav impl.NAV_F, vna impl.VNA_F) *Listener {
	//
	//
	obdh := impl.NewOBDH()
	//
	nim := &NIM_Rh{Db: db, PushChan: make(chan string, 10000)}
	nim_ob := impl.NewOBDH()
	nim.H(nim_ob)
	obdh.AddH(MK_NIM, nim)
	obdh.AddH(MK_NRC, impl.NewRC_S(nim_ob))
	//
	dim := &DIM_Rh{Db: db, DS: map[string]netw.Con{}}
	dim_m := impl.NewRCM_S(nd, vna)
	dim_m.AddHH("LI", dim)
	obdh.AddH(MK_DIM, dim)
	obdh.AddH(MK_DRC, impl.NewRC_S(dim_m))
	//
	ndh := impl.NewOBDH()
	nrh := &NodeRh{NIM: nim}
	nch := &NodeCmds{Db: db, DS: map[string]netw.Con{}}
	nch.H(ndh)
	obdh.AddH(MK_NODE, ndh)
	obdh.AddH(MK_NODE_M, nrh)
	//

	//
	var rl netw.ConPool
	ncf := func(cp netw.ConPool, p *pool.BytePool, con net.Conn) netw.Con {
		cc := netw.NewCon_(rl, p, con)
		cc.V2B_ = v2b
		cc.B2V_ = b2v
		return cc
	}
	dip := NewDimPool(db, sid, p, v2b, b2v, nav, ncf, dim)
	cch := netw.NewCCH(netw.NewQueueConH(dim, nim), obdh)
	l := netw.NewListenerN(p, fmt.Sprintf(":%v", port), sid, cch, ncf)
	// l.LConPool.SetId(sid)
	// l.SetId(sid)
	nim.SS = NewMarkConPoolSender([]byte{MK_NIM}, l, sid)
	dim.SS = nim.SS
	nch.SS = nim.SS
	nim.DS = NewMarkConPoolSender([]byte{MK_DIM}, NewMultiFinder(dim, dip), sid)
	var tl = &Listener{
		Listener:    l,
		Obdh:        obdh,
		NIM:         nim,
		DIP:         dip,
		DIM:         dim,
		Db:          db,
		P:           p,
		Port:        port,
		Sid:         sid,
		PubHost:     "127.0.0.1",
		PubPort:     port,
		PushSrvTick: 30000,
	}
	rl = tl
	return tl
}
func NewListner2(db DbH, sid string, p *pool.BytePool, port int) *Listener {
	return NewListner(db, sid, p, port,
		IM_V2B, IM_B2V, impl.Json_ND, impl.Json_NAV, impl.Json_VNA)
}
func (l *Listener) Run() error {
	err := l.DIP.Dail()
	if err != nil {
		return err
	}
	err = l.Listener.Run()
	if err != nil {
		l.DIP.Close()
		return err
	}
	err = l.Db.AddSrv(&Srv{
		Sid:   l.Sid,
		Host:  l.PubHost,
		Port:  l.PubPort,
		Token: uuid.New(),
	})
	if err != nil {
		l.DIP.Close()
		l.Listener.Close()
		return err
	}
	if len(l.PushSrvAddr) > 0 {
		l.ConPushSrv(l.PushSrvAddr)
		l.NIM.StartPushTask()
	}
	return nil
}

func (l *Listener) Close() {
	if l.PushConRunner != nil {
		l.PushConRunner.StopRunner()
		l.PushConRunner.Close()
		l.NIM.Push("") //for stop
	}
	l.Listener.Close()
	l.DIP.Close()
	err := l.Db.DelSrv(l.Sid)
	if err != nil {
		log.E("delete server by sid(%v) err:%v", l.Sid, err.Error())
	}
}

func (l *Listener) ConPushSrv(addr string) {
	l.PushConRunner = netw.NewNConRunnerN(l.P, addr, l, impl.Json_NewCon)
	l.PushConRunner.Tick = l.PushSrvTick
	l.PushConRunner.StartRunner()
	l.PushConRunner.StartTick()
}

func (l *Listener) OnCmd(c netw.Cmd) int {
	var args util.Map
	_, err := c.V(&args)
	if err != nil {
		log.E("convert push args err:%v", err.Error())
		return -1
	}
	mid := args.StrVal("MID")
	if len(mid) < 1 {
		log.E("receive invalid push by:%v", args)
		return -1
	}
	l.NIM.Push(mid)
	log.D("receive on push notification by mid(%v)", mid)
	return 0
}
func IM_V2B(v interface{}) ([]byte, error) {
	switch v.(type) {
	case *pb.ImMsg:
		bys, err := proto.Marshal(v.(*pb.ImMsg))
		if err == nil {
			return bys, nil
		} else {
			log.D("IM_V2B(proto) by v(%v) err:%v", v, err.Error())
			return bys, err
		}
	case *pb.DsMsg:
		bys, err := proto.Marshal(v.(*pb.DsMsg))
		if err == nil {
			return bys, nil
		} else {
			log.D("IM_V2B(proto) by v(%v) err:%v", v, err.Error())
			return bys, err
		}
	default:
		bys, err := impl.Json_V2B(v)
		if err == nil {
			return bys, nil
		} else {
			log.D("IM_V2B(json) by v(%v) err:%v", v, err.Error())
			return bys, err
		}
	}
}
func IM_B2V(bys []byte, v interface{}) (interface{}, error) {
	switch v.(type) {
	case *pb.ImMsg:
		err := proto.Unmarshal(bys, v.(*pb.ImMsg))
		if err == nil {
			return v, nil
		} else {
			log.D("IM_B2V(proto) by []byte(%v) err:%v", bys, err.Error())
			return v, err
		}
	case *pb.DsMsg:
		err := proto.Unmarshal(bys, v.(*pb.DsMsg))
		if err == nil {
			return v, nil
		} else {
			log.D("IM_B2V(proto) by []byte(%v) err:%v", bys, err.Error())
			return v, err
		}
	default:
		_, err := impl.Json_B2V(bys, v)
		if err == nil {
			return v, nil
		} else {
			log.D("IM_B2V(json) by []byte(%v) err:%v", bys, err.Error())
			return v, err
		}
	}
}

func IM_NewCon(cp netw.ConPool, p *pool.BytePool, con net.Conn) netw.Con {
	cc := netw.NewCon_(cp, p, con)
	cc.V2B_ = IM_V2B
	cc.B2V_ = IM_B2V
	return cc
}

func SendUnread(ss Sender, db DbH, r netw.Cmd, rv string, ct int) {
	ms, err := db.ListUnread(rv, ct)
	if err != nil {
		log.E("ListUnread by R(%v),ct(%v) error:%v", rv, ct, err.Error())
		return
	}
	for _, m := range ms {
		m.D = &rv
		av := m.Ms[rv]
		avs := strings.Split(av, MS_SEQ)
		avl := len(avs)
		if avl < 2 {
			log.W("invalid unread message:%v for R(%v)", m, rv)
			continue
		}
		for i := 1; i < avl; i++ {
			m.A = &avs[i]
			err = ss.Send(r.Id(), &m.ImMsg)
			if err != nil {
				log.W("sending unread message(%v) error:%v", &m.ImMsg, err.Error())
				return
			}
		}
		db.Update(m.GetI(), map[string]string{rv: "D"})
	}
	log_d("SendUnread %v messages is sended to %v", len(ms), rv)
}
