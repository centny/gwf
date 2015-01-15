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
		log.D_(1, f, args...)
	}
}

//the protocol modes
const H_MOD = "^~^"

//the connection for not data receive
const CON_TIMEOUT int64 = 5000

type NewConF func(cp ConPool, p *pool.BytePool, con net.Conn) Con
type CmdErrF func(c Cmd, code byte, f string, args ...interface{})

//function for covert struct to []byte
type V2Byte func(v interface{}) ([]byte, error)

//function for covert []byte to struct or map.
type Byte2V func(bys []byte, v interface{}) (interface{}, error)

//the base connection event handler
type ConHandler interface {
	//calling when the connection have been connected.
	OnConn(c Con) bool
	//calling when the connection have been closed.
	OnClose(c Con)
}

//the connect data event handler.
type CmdHandler interface {
	//calling when one entire command have been received.
	OnCmd(c Cmd)
}
type CCHandler interface {
	ConHandler
	CmdHandler
}

type CCH struct {
	Con ConHandler
	Cmd CmdHandler
}

func NewCCH(con ConHandler, cmd CmdHandler) *CCH {
	return &CCH{
		Con: con,
		Cmd: cmd,
	}
}
func (cch *CCH) OnConn(c Con) bool {
	return cch.Con.OnConn(c)
}
func (cch *CCH) OnClose(c Con) {
	cch.Con.OnClose(c)
}
func (cch *CCH) OnCmd(c Cmd) {
	cch.Cmd.OnCmd(c)
}

/*

*/
//the command wait handler impl netw.ConHandler.
type CWH struct {
	Wait bool
}

func NewCWH(w bool) *CWH {
	return &CWH{
		Wait: w,
	}
}
func (cwh *CWH) OnConn(c Con) bool {
	if cwh.Wait {
		c.SetWait(cwh.Wait)
	}
	return true
}
func (cwh *CWH) OnClose(c Con) {
}

//the connection struct.
//it will be created when client connected or server received one connection.
type Con interface {
	net.Conn //the base connection
	CP() ConPool
	P() *pool.BytePool //the memory pool
	R() *bufio.Reader  //the buffer reader
	W() *bufio.Writer  //the buffer writer.
	Kvs() util.Map
	Last() int64 //the last update time for data transfer
	SetWait(t bool)
	ReadW(p []byte) error
	// Writeb_(bys ...[]byte) (int, error)
	Writeb(bys ...[]byte) (int, error)
	Writev(val interface{}) (int, error)
	Exec(args interface{}, dest interface{}) (interface{}, error)
	Flush() error
	Waiting() bool
	V2B() V2Byte
	B2V() Byte2V
}
type Con_ struct {
	net.Conn //the base connection
	CP_      ConPool
	P_       *pool.BytePool //the memory pool
	R_       *bufio.Reader  //the buffer reader
	W_       *bufio.Writer  //the buffer writer.
	Kvs_     util.Map
	Last_    int64        //the last update time for data transfer
	Waiting_ int32        //whether in waiting status.
	c_l      sync.RWMutex //connection lock.
	buf      []byte
	V2B_     V2Byte
	B2V_     Byte2V
}

func NewCon(cp ConPool, p *pool.BytePool, con net.Conn) Con {
	return NewCon_(cp, p, con)
}

//new connection.
func NewCon_(cp ConPool, p *pool.BytePool, con net.Conn) *Con_ {
	return &Con_{
		CP_:      cp,
		P_:       p,
		Conn:     con,
		R_:       bufio.NewReader(con),
		W_:       bufio.NewWriter(con),
		Kvs_:     util.Map{},
		Waiting_: 0,
		buf:      make([]byte, 2),
		V2B_: func(v interface{}) ([]byte, error) {
			return nil, util.Err("V2B not implemeted")
		},
		B2V_: func(bys []byte, v interface{}) (interface{}, error) {
			return nil, util.Err("B2V not implemeted")
		},
	}
}
func (c *Con_) CP() ConPool {
	return c.CP_
}
func (c *Con_) P() *pool.BytePool {
	return c.P_
}
func (c *Con_) R() *bufio.Reader {
	return c.R_
}
func (c *Con_) W() *bufio.Writer {
	return c.W_
}
func (c *Con_) Kvs() util.Map {
	return c.Kvs_
}
func (c *Con_) Last() int64 {
	return c.Last_
}

//set the connection waiting status.
//if true,the connection will keep forever.
//if false,the connection will be closed after timeout when not data receive.
func (c *Con_) SetWait(t bool) {
	if t {
		atomic.StoreInt32(&c.Waiting_, 1)
	} else {
		atomic.StoreInt32(&c.Waiting_, 0)
	}
}
func (c *Con_) Waiting() bool {
	return c.Waiting_ > 0
}

//read the number of the data in p
func (c *Con_) ReadW(p []byte) error {
	return util.ReadW(c.R_, p, &c.Last_)
}

//sending data.
//Data:mod|len|bys...
func (c *Con_) Writeb(bys ...[]byte) (int, error) {
	c.c_l.Lock()
	defer c.c_l.Unlock()
	c.W_.Write([]byte(H_MOD))
	var tlen uint16 = 0
	for _, b := range bys {
		tlen += uint16(len(b))
	}
	binary.BigEndian.PutUint16(c.buf, tlen)
	total, _ := c.W_.Write(c.buf)
	var tv int = 0
	for _, b := range bys {
		tv, _ = c.W_.Write(b)
		total += tv
	}
	return total, c.Flush()
}
func (c *Con_) Writev(val interface{}) (int, error) {
	return Writev(c, val)
}
func (c *Con_) Exec(args interface{}, dest interface{}) (interface{}, error) {
	return nil, util.Err("connection not implement Exec")
}
func (c *Con_) Flush() error {
	return c.W_.Flush()
}
func (c *Con_) V2B() V2Byte {
	return c.V2B_
}
func (c *Con_) B2V() Byte2V {
	return c.B2V_
}

type Cmd interface {
	//get the connect.
	Con
	//get the command data.
	Data() []byte
	//done the command, the data []byte will free.
	Done()

	V(dest interface{}) (interface{}, error)
	Err(code byte, f string, args ...interface{})
}

//the data commend.
type Cmd_ struct {
	Con          //base connection.
	Data_ []byte //received data
	data_ []byte
}

func (c *Cmd_) Data() []byte {
	return c.Data_
}

//free the memory(Data []byte)
func (c *Cmd_) Done() {
	c.P().Free(c.data_)
}
func (c *Cmd_) V(dest interface{}) (interface{}, error) {
	return V(c, dest)
}
func (c *Cmd_) Err(code byte, f string, args ...interface{}) {
	c.CP().Err()(c, code, f, args...)
}

type ConPool interface {
	LoopTimeout()
	Close()
	RunC(c Con)
	Err() CmdErrF
}

//the connection pool
type LConPool struct {
	T      int64          //the timeout of not data received
	P      *pool.BytePool //the memory pool
	Wg     sync.WaitGroup //wait group.
	H      CCHandler      //command handler
	Wc     chan int       //the wait chan.
	NewCon NewConF
	t_r    bool
	cons   map[net.Conn]Con
	cons_l sync.RWMutex
	Err_   CmdErrF
}

//new connection pool.
func NewLConPool(p *pool.BytePool, h CCHandler) *LConPool {
	return &LConPool{
		T:      CON_TIMEOUT,
		P:      p,
		H:      h,
		Wc:     make(chan int),
		cons:   map[net.Conn]Con{},
		NewCon: NewCon,
		Err_: func(c Cmd, code byte, f string, args ...interface{}) {
			log.D_(2, f, args...)
		},
	}
}

//looping the connection timeout.
func (l *LConPool) LoopTimeout() {
	l.t_r = true
	for l.t_r {
		cons := []net.Conn{}
		tn := util.Now()
		for con, c := range l.cons {
			if c.Waiting() {
				continue
			}
			if (tn - c.Last()) > l.T {
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
	l.cons = map[net.Conn]Con{}
	l.cons_l.Unlock()
}
func (l *LConPool) add_c(c Con) {
	l.cons_l.Lock()
	l.Wg.Add(1)
	l.cons[c] = c
	l.cons_l.Unlock()
}
func (l *LConPool) del_c(c Con) {
	l.cons_l.Lock()
	l.Wg.Done()
	delete(l.cons, c)
	l.cons_l.Unlock()
}

//run one connection by async.
func (l *LConPool) RunC(con Con) {
	// go func(lll *LConPool, conn net.Conn) {
	go l.RunC_(con)
	// }(l, con)
}

//run on connection by sync.
func (l *LConPool) RunC_(con Con) {
	defer func() {
		log_d("closing connection(%v)", con.RemoteAddr().String())
		con.Close()
	}()
	log_d("running connection(%v)", con.RemoteAddr().String())
	l.add_c(con)
	defer l.del_c(con)
	//
	buf := make([]byte, 5)
	mod := []byte(H_MOD)
	mod_l := len(mod)
	//
	for {
		err := con.ReadW(buf)
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
		err = con.ReadW(dbuf)
		if err != nil {
			log_d("read data from(%v) error:%v", con.RemoteAddr().String(), err.Error())
			break
		}
		l.H.OnCmd(&Cmd_{
			Con:   con,
			Data_: dbuf,
			data_: dbuf,
		})
	}
	l.H.OnClose(con)
}
func (l *LConPool) Err() CmdErrF {
	return l.Err_
}
