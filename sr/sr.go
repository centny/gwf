package filter

import (
	"bufio"
	"fmt"
	"github.com/Centny/gwf/iow"
	"github.com/Centny/gwf/log"
	"github.com/Centny/gwf/routing"
	"github.com/Centny/gwf/sr/pb"
	"github.com/Centny/gwf/util"
	"github.com/golang/protobuf/proto"
	"io"
	"os"
	"path/filepath"
	"sync/atomic"
	"time"
)

type SRH interface {
	Path(hs *routing.HTTPSession, sr *SR) (string, error)
	OnSrF(hs *routing.HTTPSession, aid, ver, sp, sf string) error
	OnSrL(hs *routing.HTTPSession, aid, ver string, last, pn, ps int64) (interface{}, int64, error)
}
type SR struct {
	H SRH
	R string //root store path.
}

func NewSR(r string) *SR {
	return &SR{
		R: r,
		H: &SRH_N{
			c: 0,
		},
	}
}
func NewSR2(r string, h SRH) *SR {
	return &SR{
		R: r,
		H: h,
	}
}
func NewSR3(r string, h SRH_Q_H) (*SR, *SRH_Q) {
	sq := NewSRH_Q(r, h)
	return NewSR2(r, sq), sq
}
func (s *SR) SrvHTTP(hs *routing.HTTPSession) routing.HResult {
	var action string = "A"
	var aid, ver string
	err := hs.ValidCheckVal(`
		aid,R|S,L:0;
		ver,R|S,L:0;
		action,O|S,O:A~L;
		`, &aid, &ver, &action)
	if err != nil {
		return hs.MsgResErr2(1, "arg-err", err)
	}
	switch action {
	case "A":
		return s.AddSr(hs, aid, ver)
	default:
		return s.ListSr(hs, aid, ver)
	}
}
func (s *SR) AddSr(hs *routing.HTTPSession, aid, ver string) routing.HResult {
	sp, err := s.H.Path(hs, s)
	if err != nil {
		return hs.MsgResErr2(1, "arg-err", err)
	}
	sf := fmt.Sprintf("%v/%v/sr_f.zip", s.R, sp)
	_, err = hs.RecF("sr_f", sf)
	if err != nil {
		return hs.MsgResErr2(1, "srv-err", err)
	}
	err = s.H.OnSrF(hs, aid, ver, sp, sf)
	if err == nil {
		return hs.MsgRes("OK")
	} else {
		return hs.MsgResErr2(1, "srv-err", err)
	}
}
func (s *SR) ListSr(hs *routing.HTTPSession, aid, ver string) routing.HResult {
	var last, pn, ps int64 = 0, 0, 20
	err := hs.ValidCheckVal(`
		last,O|I,R:0;
		pn,O|I,R:0;
		ps,O|I,R:0;
		`, &last, &pn, &ps)
	if err != nil {
		return hs.MsgResErr2(1, "arg-err", err)
	}
	data, total, err := s.H.OnSrL(hs, aid, ver, last, pn, ps)
	if err == nil {
		return hs.MsgResP(data, pn, ps, total)
	} else {
		return hs.MsgResErr2(1, "srv-err", err)
	}
}

type SRH_N struct {
	c int64
}

func (s *SRH_N) Path(hs *routing.HTTPSession, sr *SR) (string, error) {
	return fmt.Sprintf("%v-%v", util.Now(), atomic.AddInt64(&s.c, 1)), nil
}

func (s *SRH_N) OnSrF(hs *routing.HTTPSession, aid, ver, sp, sf string) error {
	return nil
}
func (s *SRH_N) OnSrL(hs *routing.HTTPSession, aid, ver string, last, pn, ps int64) (interface{}, int64, error) {
	return []interface{}{}, 0, nil
}

type SRH_Q_I struct {
	Sp  string
	Aid string
	Ver string
	Kvs util.Map
	Evs []*pb.Evn
}
type SRH_Q_H interface {
	Args(s *SRH_Q, hs *routing.HTTPSession, aid, ver, sp, sf string) (util.Map, error)
	Proc(s *SRH_Q, i *SRH_Q_I) error
	ListSr(s *SRH_Q, hs *routing.HTTPSession, aid, ver string, last, pn, ps int64) (interface{}, int64, error)
}
type SRH_Q struct {
	SRH_N
	R       string
	H       SRH_Q_H
	Q       chan *SRH_Q_I
	Running bool
}

func NewSRH_Q(r string, h SRH_Q_H) *SRH_Q {
	return &SRH_Q{
		R: r,
		H: h,
		Q: make(chan *SRH_Q_I, 3000),
	}
}
func (s *SRH_Q) OnSrF(hs *routing.HTTPSession, aid, ver, sp, sf string) error {
	if !s.Running {
		log.W("SRH_Q OnSrF err:Proc is not running")
		return util.Err("SRH_Q not running")
	}
	kvs, err := s.H.Args(s, hs, aid, ver, sp, sf)
	if err == nil {
		s.Q <- &SRH_Q_I{
			Sp:  sp,
			Aid: aid,
			Ver: ver,
			Kvs: kvs,
		}
	}
	return err
}
func (s *SRH_Q) OnSrL(hs *routing.HTTPSession, aid, ver string, last, pn, ps int64) (interface{}, int64, error) {
	return s.H.ListSr(s, hs, aid, ver, last, pn, ps)
}
func (s *SRH_Q) Proc() {
	tick := time.Tick(500)
	for s.Running {
		select {
		case i := <-s.Q:
			s.doproc(i)
		case <-tick:
		}
	}
	log.D("SRH_Q Proc done...")
}
func (s *SRH_Q) doproc(i *SRH_Q_I) {
	sr_p := filepath.Join(s.R, i.Sp)
	sr_f := filepath.Join(s.R, i.Sp, "sr_f.zip")
	err := util.Unzip(sr_f, sr_p)
	if err != nil {
		log.E("unzip %v err:%v", sr_f, err.Error())
		return
	}
	sr_er := filepath.Join(sr_p, "er.dat")
	er_f, err := os.Open(sr_er)
	if err != nil {
		log.E("open er.data file %v err:%v", sr_er, err.Error())
		return
	}
	err = iow.ReadLdata(bufio.NewReader(er_f), func(bys []byte) error {
		var evn pb.Evn
		err = proto.Unmarshal(bys, &evn)
		if err == nil {
			i.Evs = append(i.Evs, &evn)
		}
		return err
	})
	er_f.Close()
	if err != nil && err != io.EOF {
		log.E("Unmarshal er.data file %v err:%v", sr_er, err.Error())
		return
	}
	err = s.H.Proc(s, i)
	if err != nil {
		log.E("Proc SRH_Q_I %v err:%v", i, err.Error())
	}
}
func (s *SRH_Q) Run(c int) {
	s.Running = true
	for i := 0; i < c; i++ {
		go s.Proc()
	}
	log.I("SRH_Q Run %v Proc", c)
}
func (s *SRH_Q) Stop() {
	s.Running = false
	log.I("SRH_Q Stopping Proc")
}
