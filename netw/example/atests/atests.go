package main

import (
	"fmt"
	"github.com/Centny/gwf/netw"
	"github.com/Centny/gwf/pool"
	"math/rand"
	"runtime"
	"sync"
	"time"
)

type th_s struct {
	ccs_l sync.RWMutex
	ccs   map[string]*netw.Con
	ccs_r map[*netw.Con]string
}

func (t *th_s) OnConn(c *netw.Con) bool {
	t.ccs_l.Lock()
	defer t.ccs_l.Unlock()
	vkey := fmt.Sprintf("%v", len(t.ccs))
	t.ccs[vkey] = c
	t.ccs_r[c] = vkey
	fmt.Println("Count>>:", len(t.ccs_r))
	return true
}
func (t *th_s) OnCmd(c *netw.Cmd) {
	defer c.Done()
	cc := string(c.Data)
	switch cc {
	case "L":
		c.Con.SetWait(true)
		c.Write([]byte("OK"))
	default:
		l := len(t.ccs)
		if l > 1 {
			tcc := t.ccs[fmt.Sprintf("%v", rand.Intn(l))]
			if tcc != nil {
				tcc.Write(c.Data)
			}
		}
	}
}
func (t *th_s) OnClose(c *netw.Con) {
	t.ccs_l.Lock()
	defer t.ccs_l.Unlock()
	delete(t.ccs, t.ccs_r[c])
	delete(t.ccs_r, c)
}
func (t *th_s) Show() {
	for {
		fmt.Println("Count:", len(t.ccs_r))
		time.Sleep(5 * time.Second)
	}
}
func main() {
	runtime.GOMAXPROCS(runtime.NumCPU() - 1)
	netw.ShowLog = false
	p := pool.NewBytePool(8, 1024)
	ts := &th_s{
		ccs:   map[string]*netw.Con{},
		ccs_r: map[*netw.Con]string{},
	}
	l := netw.NewListener(p, ":7686", ts)
	go ts.Show()
	err := l.Run()
	if err != nil {
		panic(err.Error())
	}
	l.Wait()
}
