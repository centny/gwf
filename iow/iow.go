package iow

import (
	"bufio"
	"encoding/binary"
	"github.com/Centny/gwf/pool"
	"github.com/Centny/gwf/util"
)

func ReadLdata(r *bufio.Reader, h func(bys []byte) error) error {
	buf := make([]byte, 2)
	var last int64
	for {
		err := util.ReadW(r, buf, &last)
		if err != nil {
			return err
		}
		dlen := binary.BigEndian.Uint16(buf)
		if dlen < 1 {
			return util.Err("the data len is zero")
		}
		tbuf := pool.BP.Alloc(int(dlen))
		err = util.ReadW(r, tbuf, &last)
		if err == nil {
			err = h(tbuf)
		}
		if err != nil {
			return err
		}
	}
}
