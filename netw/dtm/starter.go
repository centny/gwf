package dtm

import (
	"github.com/Centny/gwf/pool"
	"github.com/Centny/gwf/util"
)

func StartDTCM_S(fcfg *util.Fcfg, dbc DB_C, h DTCM_S_H) (*DTCM_S, error) {
	var bp = pool.NewBytePool(8, fcfg.IntValV("mcache", 1024000))
	dtcm_s, err := NewDTCM_S_j(bp, fcfg, dbc, h)
	if err == nil {
		err = dtcm_s.Run()
	}
	if err == nil {
		dtcm_s.StartChecker(fcfg.Int64ValV("cdelay", 10000))
	}
	return dtcm_s, err
}

func StartDTM_C(fcfg *util.Fcfg) *DTM_C {
	var bp = pool.NewBytePool(8, fcfg.IntValV("mcache", 1024000))
	dtcm_c := NewDTM_C_j(bp, fcfg.Val("srv_addr"))
	dtcm_c.Cfg = fcfg
	dtcm_c.Start()
	return dtcm_c
}
