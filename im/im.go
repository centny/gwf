package im

import (
	"code.google.com/p/go-uuid/uuid"
	"fmt"
	"github.com/Centny/gwf/log"
	"github.com/Centny/gwf/netw"
	"github.com/Centny/gwf/netw/impl"
	"github.com/Centny/gwf/pool"
	"github.com/Centny/gwf/util"
	"net"
)

var ShowLog bool = false

func log_d(f string, args ...interface{}) {
	if ShowLog {
		log.D(f, args...)
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
	DelCon(sid, cid, r string, t byte) error
	//list all connection by target R
	ListCon(rs []string) ([]Con, error)
	//
	//
	//list all user R by group R
	ListUsrR(gr []string) ([]string, error)
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
	OnUsrLogin(c netw.Cmd, r *util.Map) (string, error)
	OnUsrLogout(r string, args *util.Map) error
	//
	//
	//update the message R status
	Update(m *Msg, rs map[string]string) error
	//store mesage
	Store(m *Msg) error
}

//
type Finder interface {
	Find(id string) netw.Con
}
type Sender interface {
	Id() string
	Send(cid string, v interface{}) error
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
		return util.Err("con not found by id(%v)", cid)
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
}

func NewListner(db DbH, sid string, p *pool.BytePool, port int, v2b netw.V2Byte, b2v netw.Byte2V, nd impl.ND_F, nav impl.NAV_F, vna impl.VNA_F) *Listener {
	//
	//
	obdh := impl.NewOBDH()
	//
	nim := &NIM_Rh{Db: db}
	nim_m := impl.NewRCM_S(nd, vna)
	nim_m.AddHH("LI", nim)
	nim_m.AddHH("LO", nim)
	obdh.AddH(MK_NIM, nim)
	obdh.AddH(MK_NRC, impl.NewRC_S(nim_m))
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
	l := netw.NewListenerN(p, fmt.Sprintf(":%v", port), cch, ncf)
	nim.SS = NewMarkConPoolSender([]byte{MK_NIM}, l, sid)
	dim.SS = nim.SS
	nch.SS = nim.SS
	nim.DS = NewMarkConPoolSender([]byte{MK_DIM}, NewMultiFinder(dim, dip), sid)
	var tl = &Listener{
		Listener: l,
		NIM:      nim,
		DIP:      dip,
		DIM:      dim,
		Db:       db,
		P:        p,
		Port:     port,
		Sid:      sid,
		PubHost:  "127.0.0.1",
		PubPort:  port,
	}
	rl = tl
	return tl
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
	return nil
}

func (l *Listener) Close() {
	l.Listener.Close()
	l.DIP.Close()
	err := l.Db.DelSrv(l.Sid)
	if err != nil {
		log.E("delete server by sid(%v) err:%v", l.Sid, err.Error())
	}
}
