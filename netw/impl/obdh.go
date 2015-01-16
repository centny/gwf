package impl

import (
	"github.com/Centny/gwf/log"
	"github.com/Centny/gwf/netw"
	"github.com/Centny/gwf/util"
)

type OBDH_Con struct {
	Mark byte
	netw.Con
}

func NewOBDH_Con(mark byte, con netw.Con) *OBDH_Con {
	return &OBDH_Con{
		Mark: mark,
		Con:  con,
	}
}
func (o *OBDH_Con) Writeb(bys ...[]byte) (int, error) {
	tbys := [][]byte{[]byte{o.Mark}}
	tbys = append(tbys, bys...)
	return o.Con.Writeb(tbys...)
}
func (o *OBDH_Con) Writev(val interface{}) (int, error) {
	return netw.Writev(o, val)
}
func (o *OBDH_Con) Exec(dest interface{}, args interface{}) (interface{}, error) {
	return nil, util.Err("connection not implement Exec")
}

/*


*/
//
type obdh_cmd struct {
	netw.Cmd
	data_ []byte
}

func (o *obdh_cmd) Data() []byte {
	return o.data_
}
func (o *obdh_cmd) V(dest interface{}) (interface{}, error) {
	return netw.V(o, dest)
}

type OBDH struct {
	HS map[byte]netw.CmdHandler
}

func NewOBDH() *OBDH {
	return &OBDH{
		HS: map[byte]netw.CmdHandler{},
	}
}
func (o *OBDH) OnCmd(c netw.Cmd) int {
	if len(c.Data()) < 2 {
		c.Done()
		log.W("receive empty command data from %v", c.RemoteAddr().String())
		return -1
	}
	mark, data := util.SplitTwo(c.Data(), 1)
	if hh, ok := o.HS[mark[0]]; ok {
		return hh.OnCmd(&obdh_cmd{
			Cmd:   c,
			data_: data,
		})
	} else {
		c.Done()
		log.W("mark not found(%v) from %v", mark, c.RemoteAddr().String())
		return -1
	}
}

func (o *OBDH) AddH(mark byte, h netw.CmdHandler) {
	o.HS[mark] = h
}
