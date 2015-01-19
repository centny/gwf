package im

import (
	"fmt"
	"github.com/Centny/gwf/netw"
	"github.com/Centny/gwf/netw/impl"
	"github.com/Centny/gwf/util"
	"math/rand"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

//the memory implement DbH interface for testing
type MemDbH struct {
	u_cc uint64
	m_cc uint64
	g_cc uint64

	Cons  map[string]*Con
	con_l sync.RWMutex
	Srvs  map[string]*Srv
	srv_l sync.RWMutex
	Ms    map[string]*Msg
	ms_l  sync.RWMutex
	//
	Usr   []string
	u_lck sync.RWMutex
	Grp   map[string][]string
}

func NewMemDbH() *MemDbH {
	return &MemDbH{
		Cons: map[string]*Con{},
		Srvs: map[string]*Srv{},
		Ms:   map[string]*Msg{},
		Grp:  map[string][]string{},
		Usr:  []string{},
	}
}
func (m *MemDbH) OnConn(c netw.Con) bool {
	return true
}

//calling when the connection have been closed.
func (m *MemDbH) OnClose(c netw.Con) {
}
func (m *MemDbH) AddCon(c *Con) error {
	if c == nil {
		panic("Con is nil")
	}
	m.con_l.Lock()
	defer m.con_l.Unlock()
	m.Cons[c.Sid+c.Cid] = c
	return nil
}
func (m *MemDbH) DelCon(sid, cid string) error {
	m.con_l.Lock()
	defer m.con_l.Unlock()
	delete(m.Cons, sid+cid)
	return nil
}

//list all connection by target R
func (m *MemDbH) ListCon(rs []string) ([]Con, error) {
	if m == nil {
		panic(nil)
	}
	rsm := map[string]byte{}
	for _, r := range rs {
		rsm[r] = 1
	}
	ccs := []Con{}
	for _, cc := range m.Cons {
		if _, ok := rsm[cc.R]; ok {
			ccs = append(ccs, *cc)
		}
	}
	return ccs, nil
}

//
//
//list all user R by group R
func (m *MemDbH) ListUsrR(gr []string) ([]string, error) {
	trs := []string{}
	for _, g := range gr {
		if rs, ok := m.Grp[g]; ok {
			trs = append(trs, rs...)
		}
	}
	return trs, nil
}

//sift the R to group R and user R.
func (m *MemDbH) Sift(rs []string) ([]string, []string, error) {
	ur, gr := []string{}, []string{}
	for _, r := range rs {
		if strings.HasPrefix(r, "G-") {
			gr = append(gr, r)
		} else {
			ur = append(ur, r)
		}
	}
	return gr, ur, nil
}

//
//
func (m *MemDbH) AddSrv(srv *Srv) error {
	m.srv_l.Lock()
	defer m.srv_l.Unlock()
	m.Srvs[srv.Sid] = srv
	return nil
}
func (m *MemDbH) DelSrv(sid string) error {
	m.srv_l.Lock()
	defer m.srv_l.Unlock()
	delete(m.Srvs, sid)
	return nil
}

//find the server by token
func (m *MemDbH) FindSrv(token string) (*Srv, error) {
	for _, srv := range m.Srvs {
		if srv.Token == token {
			return srv, nil
		}
	}
	return nil, util.Err("server not found by token(%v)", token)
}

//list all online server,exclue special server id.
func (m *MemDbH) ListSrv(sid string) ([]Srv, error) {
	srvs := []Srv{}
	for _, srv := range m.Srvs {
		if srv.Sid == sid {
			continue
		}
		srvs = append(srvs, *srv)
	}
	return srvs, nil
}

//
//
//user login,return user R.
func (m *MemDbH) OnUsrLogin(r *impl.RCM_Cmd) (string, error) {
	m.u_lck.Lock()
	defer m.u_lck.Unlock()
	ur := fmt.Sprintf("U-%v", atomic.AddUint64(&m.u_cc, 1))
	m.Usr = append(m.Usr, ur)
	log_d("user login by R(%v)", ur)
	return ur, nil
}

//
//
//update the message R status
func (m *MemDbH) Update(ms *Msg, rs map[string]string) error {
	m.ms_l.Lock()
	defer m.ms_l.Unlock()
	if tm, ok := m.Ms[ms.Id]; ok {
		for r, s := range rs {
			tm.Ms[r] = s
		}
		m.Ms[tm.Id] = tm
		return nil
	} else {
		return util.Err("message not found by id(%v)", ms.Id)
	}
}

//store mesage
func (m *MemDbH) Store(ms *Msg) error {
	m.ms_l.Lock()
	defer m.ms_l.Unlock()
	ms.Id = fmt.Sprintf("M-%v", atomic.AddUint64(&m.m_cc, 1))
	m.Ms[ms.Id] = ms
	return nil
}

func (m *MemDbH) RandGrp() (string, int) {
	if len(m.Grp) < 1 {
		return "", 0
	}
	gs := []string{}
	for gr, _ := range m.Grp {
		gs = append(gs, gr)
	}
	g := gs[rand.Intn(len(gs))]
	return g, len(m.Grp[g])
}
func (m *MemDbH) RandUsr() []string {
	ulen := len(m.Usr)
	if ulen < 1 {
		return []string{}
	}
	um := map[string]byte{}
	tlen := rand.Intn(ulen)%16 + 1
	for i := 0; i <= tlen; i++ {
		um[m.Usr[rand.Intn(ulen)]] = 1
	}
	tur := []string{}
	for u, _ := range um {
		tur = append(tur, u)
	}
	return tur
}
func (m *MemDbH) GrpBuilder() {
	for {
		time.Sleep(time.Second)
		if len(m.Usr) < 1 {
			continue
		}
		g := fmt.Sprintf("G-%v", atomic.AddUint64(&m.g_cc, 1))
		us := []string{}
		tlen := rand.Intn(len(m.Usr)) + 1
		mu := map[string]bool{}
		for i := 0; i < tlen; i++ {
			mu[m.Usr[rand.Intn(len(m.Usr))]] = true
		}
		for u, _ := range mu {
			us = append(us, u)
		}
		m.Grp[g] = us
	}
}
func (m *MemDbH) Show() (uint64, uint64, uint64, uint64, uint64) {
	mlen := uint64(len(m.Ms))
	var rlen uint64 = 0
	var plen uint64 = 0
	var elen uint64 = 0
	var dlen uint64 = 0
	for _, m := range m.Ms {
		rlen += uint64(len(m.Ms))
		for _, s := range m.Ms {
			if strings.HasPrefix(s, "E-") {
				elen++
			} else if s == MS_PENDING {
				plen++
			} else {
				dlen++
			}
		}
	}
	fmt.Printf("M:%v, R(%v)-P(%v)-E(%v)=%v, D:%v\n", mlen, rlen, plen, elen, rlen-plen-elen, dlen)
	return mlen, rlen, plen, elen, dlen
}
