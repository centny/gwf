package main

import (
	"fmt"
	"github.com/Centny/gwf/netw"
	"github.com/Centny/gwf/pool"
	"github.com/Centny/gwf/routing"
	"net/http"
	"runtime"
	"sync"
	"time"
)

type th_c struct {
	i     int
	name  string
	iii   int
	count int
}

func (t *th_c) OnConn(c *netw.Con) bool {
	c.Write([]byte("L"))
	// fmt.Println("---->")
	return true
}
func (t *th_c) OnCmd(c *netw.Cmd) {
	defer c.Done()
	cc := string(c.Data)
	switch cc {
	case "OK":
		go func(con *netw.Con) {
			for {
				t.iii++
				con.Write([]byte(fmt.Sprintf("%v-%v", t.name, t.iii)))
				time.Sleep(100 * time.Millisecond)
			}
		}(c.Con)
		// fmt.Println("ok--->")
	default:
		t.count++
	}
}

func (t *th_c) OnClose(c *netw.Con) {

}

var wg sync.WaitGroup
var mc_l sync.RWMutex
var mc map[*th_c]*netw.NConPool = map[*th_c]*netw.NConPool{}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU() - 1)
	go run_c(fmt.Sprintf("BC-%v", 1))
	go run_api()
	time.Sleep(time.Millisecond)
	wg.Wait()
}
func run_api() {
	mux := routing.NewSessionMux2("")
	mux.HFunc("^/add_c.*", add_c)
	mux.HFunc("^/list_c.*", list_c)
	http.ListenAndServe(":9987", mux)
}
func list_c(hs *routing.HTTPSession) routing.HResult {
	var msgc_r int64 = 0
	var msgc_s int64 = 0
	var all_tc []map[string]interface{} = []map[string]interface{}{}
	for tc, _ := range mc {
		all_tc = append(all_tc, map[string]interface{}{
			"name":   tc.name,
			"msgc_r": tc.count,
			"msgc_s": tc.iii,
		})
		msgc_s += int64(tc.iii)
		msgc_r += int64(tc.count)
	}
	return hs.MsgRes(map[string]interface{}{
		"count":  len(mc),
		"msgc_r": msgc_r,
		"msgc_s": msgc_s,
		"ls":     all_tc,
	})
}
func add_c(hs *routing.HTTPSession) routing.HResult {
	var ic int64 = 0
	err := hs.ValidCheckVal(`
		ic,R|I,R:0;
		`, &ic)
	if err != nil {
		return hs.MsgResE(1, err.Error())
	}
	var i int64
	for i = 0; i < ic; i++ {
		go run_c(fmt.Sprintf("C-%v", i))
	}
	return hs.MsgRes("OK")
}
func run_c(name string) {
	wg.Add(1)
	tc := &th_c{name: name}
	p := pool.NewBytePool(8, 1024)
	c := netw.NewNConPool(p, "127.0.0.1:7686", tc)
	err := c.Dail()
	if err != nil {
		panic(err.Error())
	}
	if tc == nil {
		panic("nil")
	}
	if c == nil {
		panic("ccccc")
	}
	mc_l.Lock()
	mc[tc] = c
	mc_l.Unlock()
	c.Wait()
	wg.Done()
}
