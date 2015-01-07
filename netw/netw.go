package netw

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"github.com/Centny/gwf/log"
	"github.com/Centny/gwf/pool"
	"github.com/Centny/gwf/util"
	"net"
	"sync"
	"sync/atomic"
	"time"
)

var ShowLog bool = false

func log_d(f string, args ...interface{}) {
	if ShowLog {
		log.D_(4, f, args...)
	}
}

const H_MOD = "^~^"
const CON_TIMEOUT int64 = 5000

type CmdHandler interface {
	OnConn(c *Con) bool
	OnCmd(c *Cmd)
	OnClose(c *Con)
}

type Con struct {
	P       *pool.BytePool
	C       net.Conn
	R       *bufio.Reader
	W       *bufio.Writer
	Last    int64
	waiting int32
	buf     []byte
	c_l     sync.RWMutex
}

func NewCon(p *pool.BytePool, con net.Conn) *Con {
	return &Con{
		P:       p,
		C:       con,
		R:       bufio.NewReader(con),
		W:       bufio.NewWriter(con),
		waiting: 0,
		buf:     make([]byte, 2),
	}
}
func (c *Con) SetWait(t bool) {
	if t {
		atomic.StoreInt32(&c.waiting, 1)
	} else {
		atomic.StoreInt32(&c.waiting, 0)
	}
}
func (c *Con) ReadW(p []byte) error {
	return util.ReadW(c.R, p, &c.Last)
}
func (c *Con) Write(bys []byte) error {
	c.c_l.Lock()
	defer c.c_l.Unlock()
	c.W.Write([]byte(H_MOD))
	binary.BigEndian.PutUint16(c.buf, uint16(len(bys)))
	c.W.Write(c.buf)
	c.W.Write(bys)
	return c.W.Flush()
}

type Cmd struct {
	*Con
	Data []byte
}

func (c *Cmd) Done() {
	c.P.Free(c.Data)
}

type LConPool struct {
	T      int64
	P      *pool.BytePool
	Wg     sync.WaitGroup
	H      CmdHandler
	Wc     chan int
	t_r    bool
	cons   map[net.Conn]*Con
	cons_l sync.RWMutex
}

func NewLConPool(p *pool.BytePool, h CmdHandler) *LConPool {
	return &LConPool{
		T:    CON_TIMEOUT,
		P:    p,
		H:    h,
		Wc:   make(chan int),
		cons: map[net.Conn]*Con{},
	}
}
func (l *LConPool) LoopTimeout() {
	l.t_r = true
	for l.t_r {
		cons := []net.Conn{}
		tn := util.Now()
		for con, c := range l.cons {
			if c.waiting > 0 {
				continue
			}
			if (tn - c.Last) > l.T {
				cons = append(cons, con)
			}
		}
		for _, con := range cons {
			con.Close()
		}
		time.Sleep(time.Duration(l.T) * time.Millisecond)
	}
	l.Wc <- 0
}
func (l *LConPool) Close() {
	l.t_r = false
	l.cons_l.Lock()
	for c, _ := range l.cons {
		c.Close()
	}
	l.cons = map[net.Conn]*Con{}
	l.cons_l.Unlock()
}
func (l *LConPool) add_c(c *Con) {
	l.cons_l.Lock()
	l.Wg.Add(1)
	l.cons[c.C] = c
	l.cons_l.Unlock()
}
func (l *LConPool) del_c(c *Con) {
	l.cons_l.Lock()
	l.Wg.Done()
	delete(l.cons, c.C)
	l.cons_l.Unlock()
}
func (l *LConPool) RunC(con net.Conn) {
	go l.RunC_(con)
}
func (l *LConPool) RunC_(con net.Conn) {
	defer con.Close()
	c := NewCon(l.P, con)
	if !l.H.OnConn(c) {
		return
	}
	l.add_c(c)
	defer l.del_c(c)
	//
	buf := make([]byte, 5)
	mod := []byte(H_MOD)
	mod_l := len(mod)
	//
	for {
		err := c.ReadW(buf)
		if err != nil {
			log.W("read head mod from(%v) error:%v", con.RemoteAddr().String(), err.Error())
			break
		}
		if !bytes.HasPrefix(buf, mod) {
			log.W("reading invalid mod(%v) from(%v)", string(buf), con.RemoteAddr().String())
			continue
		}
		dlen := binary.BigEndian.Uint16(buf[mod_l:])
		if dlen < 1 {
			log.W("reading invalid data len for mod(%v) from(%v)", string(buf), con.RemoteAddr().String())
			continue
		}
		dbuf := l.P.Alloc(int(dlen))
		err = c.ReadW(dbuf)
		if err != nil {
			log_d("read data from(%v) error:%v", con.RemoteAddr().String(), err.Error())
			break
		}
		l.H.OnCmd(&Cmd{
			Con:  c,
			Data: dbuf,
		})
	}
	l.H.OnClose(c)
}
