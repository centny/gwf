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
func Writev2(c Con, bys []byte, val interface{}) (int, error) {
	bys2, err := c.V2B()(val)
	if err == nil {
		return c.Writeb(bys, bys2)
	} else {
		return 0, err
	}
}
func Writeh(w io.Writer, bys ...[]byte) (int, error) {
	total, err := w.Write([]byte(H_MOD))
	if err != nil {
		return 0, err
	}
	var tlen uint16 = 0
	for _, b := range bys {
		tlen += uint16(len(b))
	}
	buf := make([]byte, 2)
	binary.BigEndian.PutUint16(buf, tlen)
	tv, err := w.Write(buf)
	if err != nil {
		return 0, err
	}
	total += tv
	for _, b := range bys {
		tv, err = w.Write(b)
		if err != nil {
			return 0, err
		}
		total += tv
	}
	return total, err
}
func Writel(w io.Writer, bys ...[]byte) (int, error) {
	bys = append(bys, []byte("\n"))
	return Writen(w, bys...)
}
func Writen(w io.Writer, bys ...[]byte) (int, error) {
	total, tv := 0, 0
	var err error
	for _, b := range bys {
		tv, err = w.Write(b)
		if err != nil {
			return 0, err
		}
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
