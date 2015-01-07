package netw

import (
	"github.com/Centny/gwf/pool"
	"net"
)

type NConPool struct {
	Addr string
	*LConPool
}

func NewNConPool(p *pool.BytePool, addr string, h CmdHandler) *NConPool {
	return &NConPool{
		Addr:     addr,
		LConPool: NewLConPool(p, h),
	}
}

func (n *NConPool) Dail() error {
	con, err := net.Dial("tcp", n.Addr)
	if err != nil {
		return err
	}
	n.RunC(con)
	return nil
}

func (n *NConPool) Wait() {
	<-n.LConPool.Wc
}
