package dtm

import (
	"fmt"
	"github.com/Centny/gwf/log"
	"github.com/Centny/gwf/routing"
)

type ProcH interface {
	routing.Handler
	Key() string
}

type NormalProc struct {
	*DTM_C
	HKey string
}

func NewNormalProc(dtmc *DTM_C) *NormalProc {
	return &NormalProc{
		HKey:  "^/proc(\\?.*)?",
		DTM_C: dtmc,
	}
}
func (n *NormalProc) Key() string {
	return n.HKey
}
func (n *NormalProc) SrvHTTP(hs *routing.HTTPSession) routing.HResult {
	log.D("DTM_C HandleProc reiceve process %v", hs.R.URL.Query().Encode())
	var tid string
	var rate float64
	err := hs.ValidCheckVal(`
		tid,R|S,L:0;
		`+n.Cfg.Val2("proc_key", "process")+`,R|F,R:-0.001;`, &tid, &rate)
	if err != nil {
		hs.W.Write([]byte(fmt.Sprintf("DTM_C HandleProc receive bad arguments->%v", err.Error())))
		return routing.HRES_RETURN
	}
	err = n.NotifyProc(tid, rate)
	if err != nil {
		log.E("DTM_C HandleProc send process info by tid(%v),rate(%v) err->%v", tid, rate, err)
	}
	hs.W.Write([]byte("OK"))
	return routing.HRES_RETURN
}
