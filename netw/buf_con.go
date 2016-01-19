package netw

import (
	"bytes"
	"fmt"
	"net"
	"time"
)

type BufAddr struct {
	Con *BufCon
}

func (b *BufAddr) Network() string {
	return fmt.Sprintf("BufCon(%p)", b.Con)
}
func (b *BufAddr) String() string {
	return fmt.Sprintf("BufCon(%p)", b.Con)
}

type BufCon struct {
	*bytes.Buffer
}

func NewBufCon(buf *bytes.Buffer) *BufCon {
	return &BufCon{Buffer: buf}
}
func NewBufCon2(buf []byte) *BufCon {
	return NewBufCon(bytes.NewBuffer(buf))
}
func NewBufCon3(data string) *BufCon {
	return NewBufCon(bytes.NewBufferString(data))
}
func NewBufCon4() *BufCon {
	return NewBufCon2(nil)
}

func (b *BufCon) Close() error {
	return nil
}

func (b *BufCon) LocalAddr() net.Addr {
	return &BufAddr{Con: b}
}

func (b *BufCon) RemoteAddr() net.Addr {
	return &BufAddr{Con: b}
}

func (b *BufCon) SetDeadline(t time.Time) error {
	return nil
}

func (b *BufCon) SetReadDeadline(t time.Time) error {
	return nil
}

func (b *BufCon) SetWriteDeadline(t time.Time) error {
	return nil
}
