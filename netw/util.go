package netw

func Writev(c Con, val interface{}) (int, error) {
	bys, err := c.V2B()(val)
	if err == nil {
		return c.Writeb(bys)
	} else {
		return 0, err
	}
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
