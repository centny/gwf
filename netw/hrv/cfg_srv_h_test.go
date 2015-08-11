package hrv

import (
	"fmt"
	"github.com/Centny/gwf/pool"
	"testing"
)

func TestCfgSrvH(t *testing.T) {
	bp := pool.NewBytePool(8, 1024000)
	cfg, err := NewCfgSrvH(bp, "conf.properties")
	if err != nil {
		t.Error(err.Error())
		return
	}
	go cfg.Run()
	lg := NewHrvC_j(bp, "127.0.0.1:8234", "http://localhost")
	lg.Start()
	lg.Login("token", "name", "alias")
	lg.Login("token", "xx", "alias")
	cfg.S.Close()
	lg.Close()
	//
	cfg.OnCmd(nil)
	run_err1(bp)
	run_err2(bp)
	NewCfgSrvH(bp, "sss")
}

func run_err1(bp *pool.BytePool) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()
	cfg, _ := NewCfgSrvH(bp, "conf.properties")
	cfg.SetVal("ADDR", "sfs")
	cfg.Run()
}
func run_err2(bp *pool.BytePool) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()
	cfg, _ := NewCfgSrvH(bp, "conf.properties")
	cfg.SetVal("ADDR", ":9234")
	cfg.SetVal("HTTP", "sfs")
	cfg.Run()
}
