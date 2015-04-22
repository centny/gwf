package im

import (
	"bytes"
	"github.com/Centny/gwf/log"
	"github.com/Centny/gwf/netw"
)

const (
	WIM_SEQ = "<~>"
)

type WIM_Cmd struct {
	netw.Cmd
	C     string
	data_ []byte
}

func (w *WIM_Cmd) Data() []byte {
	return w.data_
}
func (w *WIM_Cmd) Writeb(bys ...[]byte) (int, error) {
	tbys := [][]byte{[]byte(w.C)}
	tbys = append(tbys, []byte(WIM_SEQ))
	tbys = append(tbys, bys...)
	// log.D("->>>>%v%v%v", string(tbys[0]), string(tbys[1]), string(tbys[2]))
	return w.Cmd.Writeb(tbys...)
}
func (w *WIM_Cmd) Writev(val interface{}) (int, error) {
	return netw.Writev(w, val)
}
func (w *WIM_Cmd) V(dest interface{}) (interface{}, error) {
	return netw.V(w, dest)
}

type WIM_Rh struct {
	*NIM_Rh
}

func (n *WIM_Rh) OnCmd(c netw.Cmd) int {
	tbys := bytes.SplitN(c.Data(), []byte(WIM_SEQ), 2)
	if len(tbys) < 2 {
		log.E("invalid command(%v) from (%v)", string(c.Data()), c.RemoteAddr().String())
		return -1
	}
	tcmd := &WIM_Cmd{
		Cmd:   c,
		C:     string(tbys[0]),
		data_: tbys[1],
	}
	switch tcmd.C {
	case "li":
		return n.LI(tcmd)
	case "lo":
		return n.LO(tcmd)
	case "ur":
		return n.UR(tcmd)
	case "m":
		return n.NIM_Rh.OnCmd(tcmd)
	default:
		log.E("unknow command(%v) from (%v)", tcmd.C, c.RemoteAddr().String())
		return -1
	}
}
