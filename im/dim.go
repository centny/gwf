package im

import (
	"fmt"
	"github.com/Centny/gwf/im/pb"
	"github.com/Centny/gwf/log"
	"github.com/Centny/gwf/netw"
	"github.com/Centny/gwf/netw/impl"
	"github.com/Centny/gwf/pool"
	"github.com/Centny/gwf/util"
	"sync"
)

//
type DIM_Rh struct {
	Db   DbH
	SS   Sender
	DS   map[string]netw.Con
	ds_l sync.RWMutex
}

func (d *DIM_Rh) OnConn(c netw.Con) bool {
	return true
}
func (d *DIM_Rh) OnClose(c netw.Con) {
	d.ds_l.Lock()
	defer d.ds_l.Unlock()
	delete(d.DS, c.Id())
}
func (d *DIM_Rh) Find(id string) netw.Con {
	return d.DS[id]
}
func (d *DIM_Rh) OnCmd(c netw.Cmd) int {
	defer c.Done()
	var dm pb.DsMsg
	_, err := c.V(&dm)
	if err != nil {
		log.E("convert value(%v) to DsMsg error:%v", string(c.Data()), err.Error())
		return -1
	}
	if len(dm.Rc) < 1 {
		log.E("receive invalid DsMsg(%v)", &dm)
		return -1
	}
	log_d("DIM_Rh recieve message:%v", &dm)
	ms := map[string]string{}
	for _, con := range dm.Rc {
		dm.M.D = con.R
		err = d.SS.Send(con.GetC(), dm.M)
		if err == nil {
			ms[con.GetR()] = MS_DONE
		} else {
			log.E("sending message(%v) to R(%v) in S(%v) err:%v", dm.M, con.GetR(), d.SS.Id(), err.Error())
			ms[con.GetR()] = MS_ERR + err.Error()
		}
	}
	err = d.Db.Update(dm.M.GetI(), ms)
	if err == nil {
		return 0
	} else {
		log.E("update message(%v) in Distribute server err:%v", &dm.M, err.Error())
		return -1
	}
}
func (d *DIM_Rh) Exec(r *impl.RCM_Cmd) (interface{}, error) {
	switch r.Name {
	case "LI":
		return d.LI(r)
	}
	return nil, util.Err("action not found by name(%v)", r.Name)
}
func (d *DIM_Rh) LI(r *impl.RCM_Cmd) (interface{}, error) {
	var token string
	var sid string
	err := r.ValidF(`
		sid,R|S,L:0,server id is empty;
		token,R|S,L:0,token is empty;
		`, &sid, &token)
	if err != nil {
		return r.CRes(1, err.Error())
	}
	srv, err := d.Db.FindSrv(token)
	if err != nil {
		return r.CRes(1, err.Error())
	}
	if srv.Sid != d.SS.Id() {
		errs := fmt.Sprintf("login fail,invalid token(%v) for current server(%v,%v)", token, d.SS.Id(), srv.Token)
		log.W("SLI login(%v)", errs)
		return r.CRes(1, errs)
	}
	d.ds_l.Lock()
	defer d.ds_l.Unlock()
	d.DS[sid] = r.BaseCon()
	r.SetWait(true)
	return r.CRes(0, "OK")
}

type DimPool struct {
	*impl.RC_C_H
	P      *pool.BytePool
	LS     map[string]*netw.NConPool
	DS     map[string]netw.Con
	MC     map[string]*impl.RCM_Con
	V2B    netw.V2Byte
	B2V    netw.Byte2V
	NA     impl.NAV_F
	NewCon netw.NewConF
	Db     DbH
	Sid    string
	DIM    *DIM_Rh
}

func NewDimPool(db DbH, sid string, p *pool.BytePool,
	v2b netw.V2Byte, b2v netw.Byte2V, na impl.NAV_F,
	nc netw.NewConF, dim *DIM_Rh) *DimPool {
	return &DimPool{
		RC_C_H: impl.NewRC_C_H(),
		P:      p,
		LS:     map[string]*netw.NConPool{},
		DS:     map[string]netw.Con{},
		MC:     map[string]*impl.RCM_Con{},
		V2B:    v2b,
		B2V:    b2v,
		NA:     na,
		Db:     db,
		Sid:    sid,
		NewCon: nc,
		DIM:    dim,
	}
}
func (d *DimPool) Dail() error {
	srvs, err := d.Db.ListSrv(d.Sid)
	if err != nil {
		return err
	}
	log.D("DimPool distribution server found(%v)", srvs)
	for _, srv := range srvs {
		if _, ok := d.DS[srv.Sid]; ok {
			continue
		}
		err = d.dail_(&srv)
		if err != nil {
			return err
		}
	}
	return nil
}
func (d *DimPool) dail_(srv *Srv) error {
	obdh := impl.NewOBDH()
	tc := impl.NewRC_C()
	obdh.AddH(MK_DRC, tc)
	obdh.AddH(MK_DIM, d.DIM)
	l, con, err := netw.DailN(d.P, srv.Addr(), netw.NewCCH(d, obdh), d.NewCon)
	if err != nil {
		return err
	}
	mcon := impl.NewOBDH_Con(MK_DRC, con)
	rcon := impl.NewRC_Con(mcon, tc)
	rmcon := impl.NewRCM_Con(rcon, d.NA)
	rmcon.Start()
	res, err := rmcon.ExecRes("LI", util.Map{
		"token": srv.Token,
		"sid":   d.Sid,
	})
	if err != nil {
		l.Close()
		return err
	}
	if res.Code != 0 {
		l.Close()
		return util.Err("Login result:%v", res.Res)
	}
	rmcon.SetId(srv.Sid)
	d.LS[srv.Sid] = l
	d.DS[srv.Sid] = con
	d.MC[srv.Sid] = rmcon
	log.D("dail distribution server(%v) to pool(%v)", srv, d.DS)
	return nil
}

func (d *DimPool) OnClose(c netw.Con) {
	delete(d.LS, c.Id())
	delete(d.DS, c.Id())
	d.RC_C_H.OnClose(c)
}
func (d *DimPool) OnConn(c netw.Con) bool {
	return true
}
func (d *DimPool) Find(id string) netw.Con {
	return d.DS[id]
}
func (d *DimPool) Close() {
	for _, l := range d.LS {
		l.Close()
	}
}
