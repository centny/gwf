//Package netw provide the base transfer protocol for TCP
//
//it contain the client and server base struct that can be extended by event handler.
//
//Protocol:mod->len->data
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

//whether show debug log or not.
var ShowLog bool = false

func log_d(f string, args ...interface{}) {
	if ShowLog {
		log.D_(4, f, args...)
	}
}

//the protocol modes
const H_MOD = "^~^"

//the connection for not data receive
const CON_TIMEOUT int64 = 5000

//the base connection event handler
type ConHandler interface {
	//calling when the connection have been connected.
	OnConn(c *Con) bool
	//calling when the connection have been closed.
	OnClose(c *Con)
}

//the connect data event handler.
type CmdHandler interface {
	ConHandler
	//calling when one entire command have been received.
	OnCmd(c *Cmd)
}

//the connection struct.
//it will be created when client connected or server received one connection.
type Con struct {
	P       *pool.BytePool //the memory pool
	C       net.Conn       //the base connection
	R       *bufio.Reader  //the buffer reader
	W       *bufio.Writer  //the buffer writer.
	Kvs     util.Map
	Last    int64        //the last update time for data transfer
	waiting int32        //whether in waiting status.
	buf     []byte       //the buffer.
	c_l     sync.RWMutex //connection lock.
}

//new connection.
func NewCon(p *pool.BytePool, con net.Conn) *Con {
	return &Con{
		P:       p,
		C:       con,
		R:       bufio.NewReader(con),
		W:       bufio.NewWriter(con),
		Kvs:     util.Map{},
		waiting: 0,
		buf:     make([]byte, 2),
	}
}

//set the connection waiting status.
//if true,the connection will keep forever.
//if false,the connection will be closed after timeout when not data receive.
func (c *Con) SetWait(t bool) {
	if t {
		atomic.StoreInt32(&c.waiting, 1)
	} else {
		atomic.StoreInt32(&c.waiting, 0)
	}
}

//read the number of the data in p
func (c *Con) ReadW(p []byte) error {
	return util.ReadW(c.R, p, &c.Last)
}

//sending data.
//Data:mod|len|bys...
func (c *Con) Write(bys ...[]byte) error {
	c.c_l.Lock()
	defer c.c_l.Unlock()
	c.W.Write([]byte(H_MOD))
	var tlen uint16 = 0
	for _, b := range bys {
		tlen += uint16(len(b))
	}
	binary.BigEndian.PutUint16(c.buf, tlen)
	c.W.Write(c.buf)
	for _, b := range bys {
		c.W.Write(b)
	}
	return c.W.Flush()
}

//the data commend.
type Cmd struct {
	*Con         //base connection.
	Data  []byte //received data
	data_ []byte
}

//free the memory(Data []byte)
func (c *Cmd) Done() {
	c.P.Free(c.Data)
}

//the connection pool
type LConPool struct {
	T      int64          //the timeout of not data received
	P      *pool.BytePool //the memory pool
	Wg     sync.WaitGroup //wait group.
	H      CmdHandler     //command handler
	Wc     chan int       //the wait chan.
	t_r    bool
	cons   map[net.Conn]*Con
	cons_l sync.RWMutex
}

//new connection pool.
func NewLConPool(p *pool.BytePool, h CmdHandler) *LConPool {
	return &LConPool{
		T:    CON_TIMEOUT,
		P:    p,
		H:    h,
		Wc:   make(chan int),
		cons: map[net.Conn]*Con{},
	}
}

//looping the connection timeout.
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
		log_d("closing %v connection for timeout", len(cons))
		for _, con := range cons {
			con.Close()
		}
		time.Sleep(time.Duration(l.T) * time.Millisecond)
	}
	l.Wc <- 0
}

//close all connection
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

//run one connection by async.
func (l *LConPool) RunC(con net.Conn) {
	go l.RunC_(con)
}

//run on connection by sync.
func (l *LConPool) RunC_(con net.Conn) {
	defer func() {
		con.Close()
		log_d("closing connection(%v)", con.RemoteAddr().String())
	}()
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
			log_d("read head mod from(%v) error:%v", con.RemoteAddr().String(), err.Error())
			break
		}
		if !bytes.HasPrefix(buf, mod) {
			log_d("reading invalid mod(%v) from(%v)", string(buf), con.RemoteAddr().String())
			continue
		}
		dlen := binary.BigEndian.Uint16(buf[mod_l:])
		if dlen < 1 {
			log_d("reading invalid data len for mod(%v) from(%v)", string(buf), con.RemoteAddr().String())
			continue
		}
		dbuf := l.P.Alloc(int(dlen))
		err = c.ReadW(dbuf)
		if err != nil {
			log_d("read data from(%v) error:%v", con.RemoteAddr().String(), err.Error())
			break
		}
		l.H.OnCmd(&Cmd{
			Con:   c,
			Data:  dbuf,
			data_: dbuf,
		})
	}
	l.H.OnClose(c)
}
