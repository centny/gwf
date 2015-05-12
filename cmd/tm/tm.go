package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"runtime"
	"time"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU() - 1)
	l, err := net.Listen("tcp", ":12345")
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Fatal(err)
		}
		conn.(*net.TCPConn).SetWriteBuffer(3)
		conn.(*net.TCPConn).SetLinger(-1)
		conn.(*net.TCPConn).SetNoDelay(true)
		go func(c net.Conn) {
			defer c.Close()
			fmt.Println(io.Copy(os.Stdout, c))
		}(conn)
		go func(c net.Conn) {
			for {
				fmt.Println(c.Write([]byte("S->")))
				time.Sleep(time.Second)
			}
		}(conn)
	}
}
