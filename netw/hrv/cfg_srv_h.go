package hrv

import (
	"github.com/Centny/gwf/log"
	"github.com/Centny/gwf/netw"
	"github.com/Centny/gwf/pool"
	"github.com/Centny/gwf/routing"
	"github.com/Centny/gwf/util"
	"net/http"
	"strings"
)

type CfgSrvH struct {
	*util.Fcfg
	Bp *pool.BytePool
	S  *HrvS
}

//new
func NewCfgSrvH(bp *pool.BytePool, uri string) (*CfgSrvH, error) {
	cfg, err := util.NewFcfg(uri)
	if err == nil {
		return &CfgSrvH{
			Bp:   bp,
			Fcfg: cfg,
		}, nil
	} else {
		return nil, err
	}

}

//init the HRV server.
func (c *CfgSrvH) init() {
	c.S = NewHrvS_j(c.Bp, c.Val("ADDR"))
	c.S.H = c
	c.S.Pre = c.Val("PRE")
	c.S.ShowLog = c.IntVal("LOG") == 1
	for _, h := range strings.Split(c.Val("HEADERS"), ",") {
		c.S.Headers[h] = true
	}
	c.S.SetWww(c.Val("WWW"))
	for k, _ := range c.Map {
		v := c.Val(k)
		switch {
		case strings.HasPrefix(k, "A-"):
			c.S.Args[strings.TrimPrefix(k, "A-")] = v
		case strings.HasPrefix(k, "H-"):
			c.S.Head[strings.TrimPrefix(k, "H-")] = v
		case strings.HasPrefix(k, "P-"):
			c.S.AddPattern(v)
		}
	}
}

//run server
func (c *CfgSrvH) Run() {
	c.init()
	err := c.S.Run()
	if err != nil {
		panic(err.Error())
	}
	mux := routing.NewSessionMux2(c.Val("PRE"))
	mux.HFunc(".*", c.S.Doh)
	http.Handle(c.Val("PRE")+"/", mux)
	log.I("list web on port(%v)", c.Val("HTTP"))
	http.ListenAndServe(c.Val("HTTP"), nil)
}

//
func (c *CfgSrvH) OnLogin(token, name, alias string) error {
	val := c.Val("T-" + name)
	if token == val {
		return nil
	} else {
		return util.Err("invalid")
	}
}

//
func (c *CfgSrvH) OnConn(con netw.Con) bool {
	return true
}
func (c *CfgSrvH) OnClose(con netw.Con) {
}
func (c *CfgSrvH) OnCmd(cmd netw.Cmd) int {
	return -1
}
