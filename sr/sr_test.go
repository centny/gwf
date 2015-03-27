package filter

import (
	"fmt"
	"github.com/Centny/gwf/routing"
	"github.com/Centny/gwf/routing/httptest"
	"github.com/Centny/gwf/util"
	"runtime"
	"sync/atomic"
	"testing"
	"time"
)

type srh struct {
	c int64
}

func (d *srh) Path(hs *routing.HTTPSession, sr *SR) (string, error) {
	if d.c < 1 {
		d.c++
		return "", util.Err("normal err")
	} else {
		return fmt.Sprintf("%v/%v-%v", sr.R, util.Now(), atomic.AddInt64(&d.c, 1)), nil
	}
}

func (d *srh) OnSrF(hs *routing.HTTPSession, sp, sf string) error {
	return util.Err("normal err")
}

type srh_q_h struct {
	b bool
}

func (sr *srh_q_h) Args(s *SRH_Q, hs *routing.HTTPSession, sp, sf string) (util.Map, error) {
	return hs.AllRVal(), nil
}
func (sr *srh_q_h) Proc(s *SRH_Q, i *SRH_Q_I) error {
	if sr.b {
		return util.Err("normal error")
	}
	for _, ev := range i.Evs {
		fmt.Println(ev)
	}
	return nil
}
func TestSr(t *testing.T) {
	runtime.GOMAXPROCS(runtime.NumCPU() - 1)
	sr := NewSR("/tmp")
	ts := httptest.NewServer2(sr)
	ts.G("")
	for i := 0; i < 10; i++ {
		ts.PostF2("", "sr_f", "sr_f.zip", nil)
	}
	ts.PostF2("", "sr_f", "sr.go", nil)
	sr.H = &srh{}
	ts.PostF2("", "sr_f", "sr_f.zip", nil)
	ts.PostF2("", "sr_f", "sr_f.zip", nil)

	//
	sqh := &srh_q_h{}
	sr2, srh_q := NewSR3("/tmp", sqh)
	ts2 := httptest.NewServer2(sr2)
	ts2.PostF2("", "sr_f", "sr_f.zip", nil)
	srh_q.Run(5)
	for i := 0; i < 10; i++ {
		ts2.PostF2("", "sr_f", "sr_f.zip", nil)
	}
	ts2.PostF2("", "sr_f", "sr.go", nil)
	ts2.PostF2("", "sr_f", "sr.zip", nil)
	util.FWrite2("er.dat", []byte{0, 0, 'a', 'b', 'c'})
	util.Zip("er.zip", ".", "./er.dat")
	ts2.PostF2("", "sr_f", "er.zip", nil)
	time.Sleep(500 * time.Millisecond)
	sqh.b = true
	ts2.PostF2("", "sr_f", "sr_f.zip", nil)
	time.Sleep(500 * time.Millisecond)
	srh_q.Stop()
	time.Sleep(time.Second)
}
