package impl

import (
	"container/list"
	"encoding/binary"
	"fmt"

	"github.com/Centny/gwf/log"
	"github.com/Centny/gwf/netw"
	// "github.com/Centny/gwf/pool"
	"math"
	"sync"
	"sync/atomic"
	"time"

	"github.com/Centny/gwf/util"
)

type Runner interface {
	Start()
	Stop()
}
type RC_Err struct {
	Code int
	Data []byte
}

func new_rc_err(code int, data []byte) *RC_Err {
	var tdata = make([]byte, len(data))
	copy(tdata, data)
	return &RC_Err{
		Code: code,
		Data: tdata,
	}
}
func (r *RC_Err) Error() string {
	return string(r.Data)
}
func (r *RC_Err) String() string {
	return fmt.Sprintf("RC error(%v):%v", r.Code, string(r.Data))
}
func B2V_Copy(bys []byte, v interface{}) (interface{}, error) {
	tbys := make([]byte, len(bys))
	copy(tbys, bys)
	return tbys, nil
}
func V2B_Byte(v interface{}) ([]byte, error) {
	if bys, ok := v.([]byte); ok {
		return bys, nil
	} else {
		return nil, util.Err("only []byte support")
	}
}

type rc_h_cmd struct {
	netw.Cmd
	mark  byte
	rid   []byte
	data_ []byte
}

func (r *rc_h_cmd) Data() []byte {
	return r.data_
}

func (r *rc_h_cmd) Writeb(bys ...[]byte) (int, error) {
	tbys := [][]byte{r.rid, []byte{r.mark}}
	tbys = append(tbys, bys...)
	return r.Cmd.Writeb(tbys...)
}
func (r *rc_h_cmd) Writev(val interface{}) (int, error) {
	return netw.Writev(r, val)
}
func (r *rc_h_cmd) V(dest interface{}) (interface{}, error) {
	return netw.V(r, dest)
}
func (r *rc_h_cmd) Err(code byte, f string, args ...interface{}) {
	r.mark = code
	r.Writeb([]byte(fmt.Sprintf(f, args...)))
	r.Cmd.Err(code, f, args...)
}

type RC_C struct {
	RCC uint64
	// back_c chan *rc_h_cmd //remote command back chan.
	on_cmd func(*rc_h_cmd)
}

func NewRC_C() *RC_C {
	return &RC_C{}
}
func (r *RC_C) OnCmd(c netw.Cmd) int {
	// log_d("RC_C OnCmd")
	atomic.AddUint64(&r.RCC, 1)
	if len(c.Data()) < 3 {
		c.Done()
		c.Err(1, "the cmd []byte(%v) len less 3, expect more", c.Data())
		return -1
	}
	rid, ms, data := util.SplitThree(c.Data(), 2, 3)
	log_d("RC_C receive data:%v", data)
	// r.back_c <-
	r.on_cmd(&rc_h_cmd{
		Cmd:   c,
		mark:  ms[0],
		rid:   rid,
		data_: data,
	})
	return 0
}
func (r *RC_C) Close() {
	// close(r.back_c)
}

//the chan command for each calling.
type chan_c struct {
	B    byte
	BS   bool
	C    chan bool //call back chan
	Data []byte    //the calling data.
	Back netw.Cmd  //the call back command.
	Err  error     //the occur error
}

//the remote command caller.
type RC_Con struct {
	Sleep time.Duration //select sleep time.
	//
	netw.Con //base connection
	//

	exec_c  uint16 //exec count.
	exec_id uint16 //exec id
	//
	bc *RC_C //remote command back chan.
	//
	// req_c chan *chan_c //require chan.
	req_l sync.RWMutex
	//
	// close_c chan int //
	// start_c chan int
	//
	running bool
	err     error
	// wg      sync.WaitGroup
	run_c *list.List
	// run_l   sync.RWMutex
	cids   map[uint16]*chan_c
	cids_l sync.RWMutex
	Cts    map[uint16]int64
}

//new on remote command caller.
func NewRC_Con(con netw.Con, bc *RC_C) *RC_Con {
	var rc = &RC_Con{
		Sleep: 100,
		Con:   con,
		bc:    bc,
		// req_c:   make(chan *chan_c, 10000),
		// close_c: make(chan int, 2),
		// start_c: make(chan int, 2),
		// wg:      sync.WaitGroup{},
		run_c: list.New(),
		cids:  map[uint16]*chan_c{},
		Cts:   map[uint16]int64{},
	}
	bc.on_cmd = rc.on_cmd
	return rc
}

func (r *RC_Con) Exec(args interface{}, dest interface{}) (interface{}, error) {
	return r.Exec_(0, true, args, dest)
}
func (r *RC_Con) Exec2(args interface{}) ([]byte, error) {
	return r.ExecV(0, true, args)
}

func (r *RC_Con) Execm(m byte, args interface{}, dest interface{}) (interface{}, error) {
	return r.Exec_(m, true, args, dest)
}
func (r *RC_Con) Exec_(m byte, bs bool, args interface{}, dest interface{}) (interface{}, error) {
	var bys, err = r.ExecV(m, bs, args)
	if err == nil {
		return r.B2V()(bys, dest)
	} else {
		return nil, err
	}
}

//execute one command.
func (r *RC_Con) ExecV(m byte, bs bool, args interface{}) ([]byte, error) {
	r.req_l.Lock()
	//defer r.req_l.Unlock()
	if !r.running {
		r.req_l.Unlock()
		return nil, util.Err("RC is not running")
	}
	if args == nil {
		r.req_l.Unlock()
		return nil, util.Err("arg val is nil")
	}
	if r.Con == nil {
		r.req_l.Unlock()
		return nil, util.Err("not connected")
	}
	if r.exec_c >= math.MaxUint16 {
		r.req_l.Unlock()
		return nil, util.Err("two many exector")
	}
	bys, err := r.V2B()(args)
	if err != nil {
		r.req_l.Unlock()
		return nil, err
	}
	tc := &chan_c{
		B:    m,
		BS:   bs,
		C:    make(chan bool, 2),
		Data: bys,
	}
	// r.req_c <- tc
	// r.run_l.Lock()
	var tc_e = r.run_c.PushBack(tc)
	r.send_c(tc)
	// r.run_l.Unlock()
	r.req_l.Unlock()
	<-tc.C
	r.req_l.Lock()
	r.run_c.Remove(tc_e)
	r.req_l.Unlock()
	close(tc.C)
	if tc.Err == nil {
		var buf []byte
		if bs {
			buf = tc.Back.Data()[1:]
		} else {
			buf = tc.Back.Data()
		}
		tmp := make([]byte, len(buf))
		copy(tmp, buf)
		tc.Back.Done()
		return tmp, nil
	} else {
		return nil, tc.Err
	}
}

func (r *RC_Con) on_cmd(cmd *rc_h_cmd) {
	r.cids_l.Lock()
	tid := binary.BigEndian.Uint16(cmd.rid)
	if tc, ok := r.cids[tid]; ok {
		delete(r.cids, tid)
		r.Cts[tid] = util.Now() - r.Cts[tid]
		r.exec_c--
		r.cids_l.Unlock()
		tc.Back = cmd
		if cmd.mark == 0 {
			tc.C <- true
		} else {
			tc.Err = new_rc_err(int(cmd.mark), cmd.Data())
			tc.C <- false
			cmd.Done()
		}
	} else {
		r.cids_l.Unlock()
		cmd.Done()
		log.W("RC_Con(%p) back chan not found by id(%v) on %v", r, tid, r.cids)
	}
}

// func (r *RC_Con) OnConn(c netw.Con) bool {
// 	return true
// }
// func (r *RC_Con) OnClose(c netw.Con) {
// 	r.Stop()
// }

//run the process of send/receive command(async).
func (r *RC_Con) Start() {
	// log_d("RC_Con starting...")
	r.req_l.Lock()
	r.running = true
	r.req_l.Unlock()
	// r.wg.Add(1)
	// go r.Run_()
	// <-r.start_c
}

//stop gorutine
func (r *RC_Con) Stop() {
	r.req_l.Lock()
	r.running = false
	// r.close_c <- 1
	for ele := r.run_c.Front(); ele != nil; ele = ele.Next() {
		tc := ele.Value.(*chan_c)
		tc.Err = util.Err("stopped")
		tc.C <- false
	}
	r.req_l.Unlock()
	log_d("RC_Con stopping...")
	// r.wg.Wait()
}

func (r *RC_Con) send_c(tc *chan_c) {
	con := r.Con
	// log.D("tc ->%p->%v", r, tc)
	buf := make([]byte, 3)
	r.cids_l.Lock()
	r.exec_id++
	exec_id := r.exec_id
	binary.BigEndian.PutUint16(buf, exec_id)
	buf[2] = 0
	r.cids[exec_id] = tc
	r.Cts[exec_id] = util.Now()
	r.cids_l.Unlock()
	log_d("RC_Con(%p) sending exec id(%v)", r, exec_id)
	if tc.BS {
		_, tc.Err = con.Writeb(buf, []byte{tc.B}, tc.Data)
	} else {
		_, tc.Err = con.Writeb(buf, tc.Data)
	}
	if tc.Err == nil {
		r.cids_l.Lock()
		r.exec_c++
		r.cids_l.Unlock()
	} else {
		r.cids_l.Lock()
		delete(r.cids, exec_id)
		r.Cts[exec_id] = util.Now() - r.Cts[exec_id]
		r.cids_l.Unlock()
		r.err = tc.Err
		tc.C <- false
	}
}

//run the process of send/receive command(sync).
// func (r *RC_Con) Run_() {
// 	defer r.wg.Done()
// 	if r.running {
// 		log.W("RC_Con already running....")
// 		return
// 	}
// 	cm := map[uint16]*chan_c{}
// 	r.running = true
// 	r.start_c <- 1
// 	// var trun bool = true
// 	// tk := pool.NewTick(r.Sleep * time.Millisecond)
// 	for r.running {
// 		select {
// 		case cmd := <-r.bc.back_c:
// 			if cmd == nil {
// 				break
// 			}
// 			// log.D("cmd ->%p->%v", r, cmd)
// 		case tc := <-r.req_c:
// 			if tc == nil {
// 				break
// 			}
// 		case <-r.close_c:
// 			break
// 			// 	trun = r.running && r.Con != nil
// 		}
// 	}
// 	// fmt.Println("xx0")
// 	// pool.PutTick(r.Sleep*time.Millisecond, tk)
// 	//clear all waiting.

// 	r.running = false
// 	close(r.req_c)
// 	close(r.start_c)
// 	close(r.close_c)
// 	log.D("clearing all waiting exec(%v),err(%v)", len(cm), r.err)
// }

type RC_C_H struct {
}

func NewRC_C_H() *RC_C_H {
	return &RC_C_H{}
}
func (r *RC_C_H) OnClose(c netw.Con) {
	if r, ok := c.(Runner); ok {
		r.Stop()
	}
}
func (r *RC_C_H) OnConn(c netw.Con) bool {
	if r, ok := c.(Runner); ok {
		r.Start()
	}
	return true
}

// func NewExecConPool(p *pool.BytePool, h netw.CmdHandler, tc *RC_C, v2b netw.V2Byte, b2v netw.Byte2V) *netw.NConPool {
// 	cch := netw.NewCCH(NewRC_C_H(), h)
// 	np := netw.NewNConPool(p, cch)
// 	np.NewCon = func(cp netw.ConPool, p *pool.BytePool, con net.Conn) netw.Con {
// 		cc := netw.NewCon_(cp, p, con)
// 		cc.V2B_, cc.B2V_ = v2b, b2v
// 		rcc := NewRC_Con(cc, tc)
// 		// cch.Con = rcc
// 		return rcc
// 	}
// 	return np
// }

/*


 */
//the remote command server handler.
type RC_S struct {
	H netw.CmdHandler
}

//new remote command server handler.
func NewRC_S(h netw.CmdHandler) *RC_S {
	return &RC_S{
		H: h,
	}
}
func (r *RC_S) OnCmd(c netw.Cmd) int {
	if len(c.Data()) < 3 {
		c.Done()
		c.Err(1, "the cmd []byte(%v) len less 3, expect more", c.Data())
		return -1
	}
	rid, ms, data := util.SplitThree(c.Data(), 2, 3)
	return r.H.OnCmd(&rc_h_cmd{
		Cmd:   c,
		mark:  ms[0],
		rid:   rid,
		data_: data,
	})
}
