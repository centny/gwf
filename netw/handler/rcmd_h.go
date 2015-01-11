package handler

import (
	"encoding/binary"
	"github.com/Centny/gwf/log"
	"github.com/Centny/gwf/netw"
	"github.com/Centny/gwf/util"
	"math"
	"time"
)

//the chan command for each calling.
type chan_c struct {
	C    chan bool //call back chan
	Data []byte    //the calling data.
	Back *netw.Cmd //the call back command.
	Err  error     //the occur error
}

//the remote command caller.
type RC_C struct {
	Sleep time.Duration //select sleep time.
	//
	*netw.Con //base connection
	c_        chan int
	//
	exec_c  uint16 //exec count.
	exec_id uint16 //exec id
	//
	back_c chan *netw.Cmd //remote command back chan.
	//
	req_c chan *chan_c //require chan.
	//
	running bool
}

//new on remote command caller.
func NewRC_C() *RC_C {
	return &RC_C{
		Sleep:  100,
		c_:     make(chan int),
		back_c: make(chan *netw.Cmd),
		req_c:  make(chan *chan_c),
	}
}

//execute one command.
func (r *RC_C) Exec(data []byte) (*netw.Cmd, error) {
	if len(data) < 1 {
		return nil, util.Err("arg data is empty")
	}
	if r.Con == nil {
		return nil, util.Err("not connected")
	}
	if r.exec_c >= math.MaxUint16 {
		return nil, util.Err("two many exector")
	}
	tc := &chan_c{
		C:    make(chan bool),
		Data: data,
	}
	r.req_c <- tc
	<-tc.C
	return tc.Back, tc.Err
}

func (r *RC_C) OnConn(c *netw.Con) bool {
	r.Con = c
	r.c_ <- 0
	return true
}

//wait connect
func (r *RC_C) WaitConn() int {
	return <-r.c_
}
func (r *RC_C) OnCmd(c *netw.Cmd) {
	r.back_c <- c
}
func (r *RC_C) OnClose(c *netw.Con) {
	r.c_ <- 1
	r.Stop()
}

//start all,it will wait the connect completed.
func (r *RC_C) Start() int {
	tv := r.WaitConn()
	if tv == 0 {
		r.Run()
	}
	return tv
}

//stop gorutine
func (r *RC_C) Stop() {
	r.running = false
}

//run the process of send/receive command(async).
func (r *RC_C) Run() {
	go r.Run_()
}

//run the process of send/receive command(sync).
func (r *RC_C) Run_() {
	cm := map[uint16]*chan_c{}
	buf := make([]byte, 2)
	r.running = true
	for r.running && r.Con != nil {
		con := r.Con
		select {
		case tc := <-r.req_c:
			binary.BigEndian.PutUint16(buf, r.exec_id)
			tc.Err = con.Write(buf, tc.Data)
			if tc.Err == nil {
				cm[r.exec_id] = tc
				r.exec_c++
				r.exec_id++
			} else {
				tc.C <- false
			}
		case cmd := <-r.back_c:
			if len(cmd.Data) < 2 {
				log.W("response data is less 2,%v", cmd.Data)
				break
			}
			tbuf, data := cmd.Data[:2], cmd.Data[2:]
			tid := binary.BigEndian.Uint16(tbuf)
			if tc, ok := cm[tid]; ok {
				cmd.Data = data
				r.exec_c--
				delete(cm, tid)
				tc.Back = cmd
				tc.C <- true
			} else {
				log.W("back chan not found by id:%v", tid)
			}
		case <-time.Tick(r.Sleep * time.Millisecond):
		}
	}
	r.running = false
}

//the remote command server call back command struct.
type RC_Cmd struct {
	*netw.Cmd
	rid []byte
}

//rewrite the base function
func (r *RC_Cmd) Write(bys []byte) error {
	return r.Cmd.Write(r.rid, bys)
}

//the extended command handler.
type RC_H interface {
	//calling when the connection have been connected.
	OnConn(c *netw.Con) bool
	//calling when one entire command have been received.
	OnCmd(rc *RC_Cmd)
	//calling when the connection have been closed.
	OnClose(c *netw.Con)
}

//the remote command server handler.
type RC_S struct {
	H RC_H
}

//new remote command server handler.
func NewRC_S(h RC_H) *RC_S {
	return &RC_S{
		H: h,
	}
}
func (r *RC_S) OnConn(c *netw.Con) bool {
	return r.H.OnConn(c)
}
func (r *RC_S) OnCmd(c *netw.Cmd) {
	rid, data := c.Data[:2], c.Data[2:]
	c.Data = data
	r.H.OnCmd(&RC_Cmd{
		Cmd: c,
		rid: rid,
	})
}
func (r *RC_S) OnClose(c *netw.Con) {
	r.OnClose(c)
}
