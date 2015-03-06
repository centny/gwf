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

// func (o *OBDH_Con) Exec(dest interface{}, args interface{}) (interface{}, error) {
// 	return nil, util.Err("connection not implement Exec")
// }

/*


*/
//
type obdh_cmd struct {
	netw.Cmd
	data_ []byte
	mark  byte
}

func (o *obdh_cmd) Data() []byte {
	return o.data_
}
func (o *obdh_cmd) Writeb(bys ...[]byte) (int, error) {
	tbys := [][]byte{[]byte{o.mark}}
	tbys = append(tbys, bys...)
	return o.Cmd.Writeb(tbys...)
}
func (o *obdh_cmd) Writev(val interface{}) (int, error) {
	return netw.Writev(o, val)
}
func (o *obdh_cmd) V(dest interface{}) (interface{}, error) {
	return netw.V(o, dest)
}

//
type OBDH_HF func(c netw.Cmd) int

func (o OBDH_HF) OnCmd(c netw.Cmd) int {
	return o(c)
}

//
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
		log.W("receive empty command data(%v) from %v in(%v)", c.Data(), c.RemoteAddr().String(), c.CP().Id())
		return -1
	}
	log_d("OBDH receive data:%v", string(c.Data()))
	mark, data := util.SplitTwo(c.Data(), 1)
	if hh, ok := o.HS[mark[0]]; ok {
		c.SetErrd(3)
		return hh.OnCmd(&obdh_cmd{
			Cmd:   c,
			data_: data,
			mark:  mark[0],
		})
	} else {
		c.Done()
		log.W("mark(%v,%v) not found in(%v) from %v in(%v)", mark[0], string(c.Data()), o.HS, c.RemoteAddr().String(), c.CP().Id())
		return -1
	}
}

func (o *OBDH) AddH(mark byte, h netw.CmdHandler) {
	o.HS[mark] = h
}
func (o *OBDH) AddF(mark byte, f func(c netw.Cmd) int) {
	o.HS[mark] = OBDH_HF(f)
}
