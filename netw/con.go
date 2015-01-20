package netw

import (
	"github.com/Centny/gwf/pool"
	"github.com/Centny/gwf/util"
	"net"
)

// type NCon struct {
// 	Addr string
// 	Con
// }

// func NewNCon(con Con, addr string) Con {
// 	return &NCon{
// 		Con:  con,
// 		Addr: addr,
// 	}
// }

//the client connection pool.
type NConPool struct {
	*LConPool //base connection pool.
}

//new client connection pool.
func NewNConPool(p *pool.BytePool, h CCHandler) *NConPool {
	return &NConPool{
		LConPool: NewLConPool(p, h),
	}
}

//dail one connection.
func (n *NConPool) Dail(addr string) (Con, error) {
	con, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}
	cc := n.NewCon(n, n.P, con)
	if !n.H.OnConn(cc) {
		return nil, util.Err("OnConn return false for %v", addr)
	}
	n.RunC(cc)
	return cc, nil
}

func Dail(p *pool.BytePool, addr string, h CCHandler) (*NConPool, Con, error) {
	return DailN(p, addr, h, NewCon)
}
func DailN(p *pool.BytePool, addr string, h CCHandler, ncf NewConF) (*NConPool, Con, error) {
	nc := NewNConPool(p, h)
	nc.NewCon = ncf
	cc, err := nc.Dail(addr)
	return nc, cc, err
}
