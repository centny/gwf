package netw

import (
	"encoding/binary"
	"io"
)

func Writev(c Con, val interface{}) (int, error) {
	bys, err := c.V2B()(val)
	if err == nil {
		return c.Writeb(bys)
	} else {
		return 0, err
	}
}
func Writeb(w io.Writer, bys ...[]byte) (int, error) {
	total, err := w.Write([]byte(H_MOD))
	var tlen uint16 = 0
	for _, b := range bys {
		tlen += uint16(len(b))
	}
	buf := make([]byte, 2)
	binary.BigEndian.PutUint16(buf, tlen)
	tv, err := w.Write(buf)
	total += tv
	for _, b := range bys {
		tv, err = w.Write(b)
		total += tv
	}
	return total, err
}

// func Writeb(c Con, bys ...[]byte) (int, error) {

// }
// func Writebb(c Con, bys1 []byte, bys2 ...[]byte) (int, error) {
// 	tbys := [][]byte{bys1}
// 	tbys = append(tbys, bys2...)
// 	return c.Writeb(tbys...)
// }
// func Writebb2(c Cmd, bys1 []byte, bys2 ...[]byte) (int, error) {
// 	tbys := [][]byte{bys1}
// 	tbys = append(tbys, bys2...)
// 	return c.Write(tbys...)
// }

func V(c Cmd, dest interface{}) (interface{}, error) {
	return c.B2V()(c.Data(), dest)
}
