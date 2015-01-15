package netw

import (
	"github.com/Centny/gwf/pool"
	"github.com/Centny/gwf/util"
	"net"
)

//the client connection pool.
type NConPool struct {
	Addr      string //target address.
	*LConPool        //base connection pool.
}

//new client connection pool.
func NewNConPool(p *pool.BytePool, addr string, h CCHandler) *NConPool {
	return &NConPool{
		Addr:     addr,
		LConPool: NewLConPool(p, h),
	}
}

//dail one connection.
func (n *NConPool) Dail() (Con, error) {
	con, err := net.Dial("tcp", n.Addr)
	if err != nil {
		return nil, err
	}
	cc := n.NewCon(n, n.P, con)
	if !n.H.OnConn(cc) {
		return nil, util.Err("OnConn return false for %v", n.Addr)
	}
	n.RunC(cc)
	return cc, nil
}

func Dail(p *pool.BytePool, addr string, h CCHandler) (*NConPool, Con, error) {
	return DailN(p, addr, h, NewCon)
}
func DailN(p *pool.BytePool, addr string, h CCHandler, ncf NewConF) (*NConPool, Con, error) {
	nc := NewNConPool(p, addr, h)
	nc.NewCon = ncf
	cc, err := nc.Dail()
	return nc, cc, err
}
