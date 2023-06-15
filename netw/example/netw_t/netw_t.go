package main

// import (
// 	"fmt"
// 	"github.com/Centny/gwf/netw"
// 	"github.com/Centny/gwf/netw/handler"
// 	"github.com/Centny/gwf/pool"
// 	"os"
// 	"runtime"
// 	"sync/atomic"
// 	"time"
// )

// type th_s struct {
// 	ccc int64
// }

// func (t *th_s) OnConn(c *netw.Con) bool {
// 	return true
// }
// func (t *th_s) OnCmd(c *netw.Cmd) {
// 	defer c.Done()
// 	c.Write([]byte("S----->"))
// 	atomic.AddInt64(&t.ccc, 1)
// }
// func (t *th_s) OnClose(c *netw.Con) {
// }
// func (t *th_s) show() {
// 	for {
// 		time.Sleep(2 * time.Second)
// 		fmt.Println("CCC:", t.ccc)
// 	}
// }

// ////////

// type th_c struct {
// }

// func (t *th_c) OnConn(c *netw.Con) bool {
// 	c.Write([]byte("L"))
// 	return true
// }
// func (t *th_c) OnCmd(c *netw.Cmd) {
// 	defer c.Done()
// 	// fmt.Println(string(c.Data))
// }

// func (t *th_c) OnClose(c *netw.Con) {

// }

// func run_c() {
// 	tc := &th_c{}
// 	p := pool.NewBytePool(8, 1024)
// 	nc, cc, err := netw.Dail(p, "127.0.0.1:7686", tc)
// 	if err != nil {
// 		fmt.Println("connect fail")
// 		return
// 	}
// 	for i := 0; i < 1000000; i++ {
// 		cc.Write([]byte("CCC->>"))
// 	}
// 	time.Sleep(time.Second)
// 	nc.Close()
// 	fmt.Println("all end ...")
// }
// func run_s() {
// 	netw.ShowLog = false
// 	ts := &th_s{}
// 	ts_h := handler.NewChanH(ts)
// 	ts_h.Run(util.CPU())
// 	p := pool.NewBytePool(8, 1024)
// 	l := netw.NewListener(p, ":7686", ts_h)
// 	err := l.Run()
// 	if err != nil {
// 		panic(err.Error())
// 	}
// 	go ts.show()
// 	l.Wait()
// }

// //////
// func main() {
// 	runtime.GOMAXPROCS(util.CPU())
// 	if len(os.Args) < 2 {
// 		fmt.Println("less 2")
// 		return
// 	}
// 	switch os.Args[1] {
// 	case "-s":
// 		run_s()
// 	case "-c":
// 		run_c()
// 	}
// }
