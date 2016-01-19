package netw

import (
	"fmt"
	"testing"
	"time"
)

func TestBufCon(t *testing.T) {
	con := NewBufCon4()
	for i := 0; i < 10; i++ {
		con.Write([]byte("abc->\n"))
	}
	buf := make([]byte, 6)
	for {
		rlen, err := con.Read(buf)
		if err != nil {
			break
		}
		fmt.Print(string(buf[:rlen]))
	}
	fmt.Println(con.RemoteAddr())
	fmt.Println(con.LocalAddr())
	fmt.Println(con.RemoteAddr().Network())
	con.Close()
	con.SetDeadline(time.Now())
	con.SetReadDeadline(time.Now())
	con.SetWriteDeadline(time.Now())
	NewBufCon3("data")
}
