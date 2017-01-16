package im

import (
	"fmt"
	"math/rand"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/Centny/gwf/log"
	"github.com/Centny/gwf/netw"
	"github.com/Centny/gwf/util"
)

//the memory implement DbH interface for testing
type MemDbH struct {
	u_cc uint64
	m_cc uint64
	g_cc uint64
	// mr_n_cc uint64

	Cons  map[string]*Con
	con_l sync.RWMutex
	Srvs  map[string]*Srv
	srv_l sync.RWMutex
	Ms    map[string]*Msg
	ms_l  sync.RWMutex
	U2M   map[string][]string
	u2m_l sync.RWMutex
	//
	Usr   map[string]byte
	u_lck sync.RWMutex
	Grp   map[string][]string
	g_lck sync.RWMutex
	//
	Tokens map[string]string
}

func NewMemDbH() *MemDbH {
	return &MemDbH{
		Cons:   map[string]*Con{},
		Srvs:   map[string]*Srv{},
		Ms:     map[string]*Msg{},
		U2M:    map[string][]string{},
		Grp:    map[string][]string{},
		Usr:    map[string]byte{},
		Tokens: map[string]string{},
	}
}
func (m *MemDbH) OnConn(c netw.Con) bool {
	return true
}

//calling when the connection have been closed.
func (m *MemDbH) OnClose(c netw.Con) {
}
func (m *MemDbH) OnCloseCon(c netw.Con, sid, cid string, t byte) error {
	m.con_l.Lock()
	defer m.con_l.Unlock()
	var pre = fmt.Sprintf("%v%v", sid, cid)
	for key, _ := range m.Cons {
		if strings.HasPrefix(key, pre) {
			delete(m.Cons, key)
		}
	}
	return nil
}
func (m *MemDbH) AddCon(c *Con) error {
	if c == nil {
		panic("Con is nil")
	}
	m.con_l.Lock()
	defer m.con_l.Unlock()
	m.Cons[fmt.Sprintf("%v%v%v%v", c.Sid, c.Cid, c.Uid, c.ConType)] = c
	log.D("adding connection %v", c)
	return nil
}
func (m *MemDbH) DelCon(sid, cid, r string, t byte, ct int) (*Con, error) {
	m.con_l.Lock()
	defer m.con_l.Unlock()
	key := fmt.Sprintf("%v%v%v%v", sid, cid, r, t)
	c := m.Cons[key]
	delete(m.Cons, key)
	log.D("delete connection %v", c)
	panic("sss")
	return c, nil
}
func (m *MemDbH) DelConT(sid, cid, token string, t byte) (*Con, error) {
	panic("not impl")
}

//list all connection by target R
func (m *MemDbH) ListCon(rs []string) ([]Con, error) {
	if m == nil {
		panic(nil)
	}
	m.con_l.Lock()
	defer m.con_l.Unlock()
	rsm := map[string]byte{}
	for _, r := range rs {
		rsm[r] = 1
	}
	ccs := []Con{}
	for _, cc := range m.Cons {
		if _, ok := rsm[cc.Uid]; ok {
			ccs = append(ccs, *cc)
		}
	}
	return ccs, nil
}

//
//
func (m *MemDbH) FUsrR(c netw.Cmd) string {
	return c.Kvs().StrVal("R")
}

// func (m *MemDbH) FindUsrR(uid int64) (string, error) {
// 	return fmt.Sprintf("U-%v", uid), nil
// }

//list all user R by group R
func (m *MemDbH) ListUsrR(msg *Msg, gr []string) (map[string][]string, error) {
	m.g_lck.Lock()
	defer m.g_lck.Unlock()
	trs := map[string][]string{}
	for _, g := range gr {
		if rs, ok := m.Grp[g]; ok {
			trs[g] = rs
		}
	}
	return trs, nil
}
func (m *MemDbH) ListR() ([]string, error) {
	m.u_lck.Lock()
	defer m.u_lck.Unlock()
	var usrs []string = []string{}
	for r, _ := range m.Usr {
		usrs = append(usrs, r)
	}
	return usrs, nil
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
	// srv.Token = "abc"
	// fmt.Println(m, srv)
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
	// fmt.Println(m, m.Srvs)
	for _, srv := range m.Srvs {
		if len(sid) > 0 && srv.Sid == sid {
			continue
		}
		srvs = append(srvs, *srv)
	}
	return srvs, nil
}

//
//
//user login,return user R.
func (m *MemDbH) OnLogin(r netw.Cmd, args *util.Map) (int, string, string, int, error) {
	m.u_lck.Lock()
	defer m.u_lck.Unlock()
	if args.Exist("token") {
		var ur string
		if tr, ok := m.Tokens[args.StrVal("token")]; ok {
			ur = tr
		} else {
			ur = fmt.Sprintf("U-%v", atomic.AddUint64(&m.u_cc, 1))
		}
		m.Usr[ur] = 1
		log.D("user login by R(%v)", ur)
		r.Kvs().SetVal("R", ur)
		return 0, ur, ur, 1, nil
	} else {
		log.D("user login fail for token not found")
		return 301, "", "", 0, util.Err("login fail:token not found")
	}
}
func (m *MemDbH) OnLogout(r netw.Cmd, args *util.Map) (string, string, int, bool, error) {
	m.u_lck.Lock()
	defer m.u_lck.Unlock()
	rv := r.Kvs().StrVal("R")
	if _, ok := m.Usr[rv]; ok {
		delete(m.Usr, rv)
		log.D("user logout by R(%v)", r)
		return rv, "", 1, true, nil
	} else {
		log.D("user logout fail:R not found")
		return "", "", 0, false, util.Err("login fail:R not found")
	}
}

//
//
//update the message R status
// func (m *MemDbH) Update(mid string, rs map[string]string) error {
// 	m.ms_l.Lock()
// 	defer m.ms_l.Unlock()
// 	if tm, ok := m.Ms[mid]; ok {
// 		for r, s := range rs {
// 			tm.Ms[r] = s
// 		}
// 		m.Ms[mid] = tm
// 		return nil
// 	} else {
// 		return util.Err("message not found by id(%v)", mid)
// 	}
// }
func (m *MemDbH) NewMid() string {
	return fmt.Sprintf("M-%v", atomic.AddUint64(&m.m_cc, 1))
}

//store mesage
func (m *MemDbH) Store(ms *Msg) error {
	m.u2m_l.Lock()
	m.ms_l.Lock()
	defer func() {
		m.ms_l.Unlock()
		m.u2m_l.Unlock()
	}()
	m.Ms[ms.GetI()] = ms
	m.U2M[ms.GetS()] = append(m.U2M[ms.GetS()], ms.GetI())
	// if len(ms.Ms) < 1 {
	// 	panic("message MS is empty")
	// }
	return nil
}
func (m *MemDbH) MarkRecv(r, a string, mids []string) error {
	if len(mids) < 1 {
		return util.Err("the message is empty")
	}
	m.ms_l.Lock()
	defer m.ms_l.Unlock()
	for _, mid := range mids {
		if msg, ok := m.Ms[mid]; ok {
			for _, mss := range msg.Ms[r] {
				if mss.R == a {
					mss.S = MS_DONE
				}
			}
		} else {
			// atomic.AddUint64(&m.mr_n_cc, 1)
			return util.Err("the message not found by id(%v)", mid)
		}
	}
	return nil
}

func (m *MemDbH) RandGrp() (string, int, []string) {
	if len(m.Grp) < 1 {
		return "", 0, nil
	}
	m.g_lck.RLock()
	defer m.g_lck.RUnlock()
	gs := []string{}
	for gr, _ := range m.Grp {
		gs = append(gs, gr)
	}
	g := gs[rand.Intn(len(gs))]
	return g, len(m.Grp[g]), m.Grp[g]
}
func (m *MemDbH) RandUsr(r string) []string {
	ulen := len(m.Usr)
	if ulen < 2 {
		return []string{}
	}
	usrs, _ := m.ListR()
	um := map[string]byte{}
	tlen := rand.Intn(ulen)%16 + 2
	for i := 0; i <= tlen; i++ {
		tr := usrs[rand.Intn(ulen)]
		if tr == r {
			continue
		}
		um[tr] = 1
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
		usrs, _ := m.ListR()
		g := fmt.Sprintf("G-%v", atomic.AddUint64(&m.g_cc, 1))
		us := []string{}
		tlen := rand.Intn(len(m.Usr)) + 1
		mu := map[string]bool{}
		for i := 0; i < tlen; i++ {
			mu[usrs[rand.Intn(len(m.Usr))]] = true
		}
		for u, _ := range mu {
			us = append(us, u)
		}
		m.Grp[g] = us
	}
}
func (m *MemDbH) Show_() (uint64, uint64, uint64, uint64, uint64) {
	mlen := uint64(len(m.Ms))
	var rlen uint64 = 0
	var plen uint64 = 0
	var elen uint64 = 0
	var dlen uint64 = 0
	for _, m := range m.Ms {
		rlen += uint64(len(m.Ms))
		for _, s := range m.Ms {
			for _, mss := range s {
				switch mss.S {
				case MS_DONE:
					dlen++
				case MS_PENDING:
					plen++
				default:
					elen++
				}
			}
		}
	}
	return mlen, rlen, plen, elen, dlen
}

func (m *MemDbH) Show() (uint64, uint64, uint64, uint64, uint64) {
	mlen, rlen, plen, elen, dlen := m.Show_()
	fmt.Printf("M:%v, R(%v)-P(%v)-E(%v)=%v, D:%v\n", mlen, rlen, plen, elen, rlen-plen-elen, dlen)
	return mlen, rlen, plen, elen, dlen
}

func (m *MemDbH) ListUnread(r string, ct int, args util.Map) ([]*Msg, error) {
	m.ms_l.Lock()
	defer m.ms_l.Unlock()
	var ms = []*Msg{}
	for _, msg := range m.Ms {
		for tr, mms := range msg.Ms {
			for _, mss_ := range mms {
				if tr == r && mss_.S == MS_PENDING {
					ms = append(ms, msg)
				}
			}
		}
	}
	return ms, nil
}
func (m *MemDbH) DoSync(uid string, ct int, args util.Map, send func(ms []*Msg) error) error {
	ms, err := m.ListUnread(uid, ct, args)
	if err != nil {
		log.E("ListUnread by R(%v),ct(%v) error:%v", uid, ct, err.Error())
		return err
	}
	return send(ms)
}
func (m *MemDbH) ListPushTask(sid, mid string) (*Msg, []Con, error) {
	m.con_l.Lock()
	m.ms_l.Lock()
	defer func() {
		m.ms_l.Unlock()
		m.con_l.Unlock()
	}()
	msg, ok := m.Ms[mid]
	if !ok {
		return nil, nil, util.Err("message not found by id(%v)", mid)
	}
	cons := []Con{}
	for r, v := range msg.Ms {
		pc := 0
		for _, mss := range v {
			if mss.S == MS_PENDING {
				pc++
			}
		}
		if pc < 1 {
			continue
		}
		for _, cc := range m.Cons {
			if cc.Uid == r && cc.Sid == sid {
				cons = append(cons, *cc)
			}
		}
	}
	return msg, cons, nil
}

func (m *MemDbH) AddGrp(grp string, users []string) {
	m.g_lck.Lock()
	defer m.g_lck.Unlock()
	m.Grp[grp] = users
}
func (m *MemDbH) DelGrp(grp string) {
	m.g_lck.Lock()
	defer m.g_lck.Unlock()
	delete(m.Grp, grp)
}
func (m *MemDbH) AddTokens(tokens map[string]string) {
	m.u_lck.Lock()
	defer m.u_lck.Unlock()
	for token, tr := range tokens {
		m.Tokens[token] = tr
	}
}
func (m *MemDbH) DelTokens(tokens []string) {
	m.u_lck.Lock()
	defer m.u_lck.Unlock()
	for _, token := range tokens {
		delete(m.Tokens, token)
	}
}
func (m *MemDbH) ClearMsg(urs []string) {
	m.u_lck.Lock()
	m.u2m_l.Lock()
	m.ms_l.Lock()
	defer func() {
		m.ms_l.Unlock()
		m.u2m_l.Unlock()
		m.u_lck.Unlock()
	}()
	for _, ur := range urs {
		var mids = m.U2M[ur]
		if len(mids) < 1 {
			continue
		}
		for _, mid := range mids {
			delete(m.Ms, mid)
		}
		delete(m.Usr, ur)
		delete(m.U2M, ur)
	}
}

func (m *MemDbH) DoOffline(offline map[string][]string, msg *Msg) error {
	return nil
}

// func (m *MemDbH) ListPcm(sid string) ([]*PCM, error) {
// 	if len(m.Cons) < 1 {
// 		return []*PCM{}, nil
// 	}
// 	cons := []*Con{}
// 	rsss := []string{}
// 	for _, cc := range m.Cons {
// 		rsss = append(rsss, cc.R)
// 		cons = append(cons, cc)
// 	}
// 	msg := &Msg{}
// 	var ss string = "S-Robot"
// 	var tt uint32 = 0
// 	var ii string = m.NewMid()
// 	msg.ImMsg = pb.ImMsg{
// 		I: &ii,
// 		S: &ss,
// 		R: rsss,
// 		T: &tt,
// 		C: []byte("Robot PCM Message"),
// 	}
// 	m.Store(msg)
// 	return []*PCM{
// 		&PCM{
// 			C: cons,
// 			M: []*Msg{msg},
// 		},
// 	}, nil
// }
