package impl

import (
	"encoding/binary"
	"fmt"
	"github.com/Centny/gwf/log"
	"github.com/Centny/gwf/netw"
	"github.com/Centny/gwf/util"
	"math"
	"sync/atomic"
	"time"
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
	return &RC_Err{
		Code: code,
		Data: data,
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
	RCC    uint64
	back_c chan *rc_h_cmd //remote command back chan.
}

func NewRC_C() *RC_C {
	return &RC_C{
		back_c: make(chan *rc_h_cmd, 10000),
	}
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
	r.back_c <- &rc_h_cmd{
		Cmd:   c,
		mark:  ms[0],
		rid:   rid,
		data_: data,
	}
	return 0
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
	req_c chan *chan_c //require chan.
	//
	running bool
	err     error
}

//new on remote command caller.
func NewRC_Con(con netw.Con, bc *RC_C) *RC_Con {
	return &RC_Con{
		Sleep: 100,
		Con:   con,
		bc:    bc,
		req_c: make(chan *chan_c, 10000),
	}
}
func (r *RC_Con) Exec(args interface{}, dest interface{}) (interface{}, error) {
	return r.Exec_(0, false, args, dest)
}
func (r *RC_Con) Exec2(args interface{}) ([]byte, error) {
	return r.ExecV(0, false, args)
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
	if args == nil {
		return nil, util.Err("arg val is nil")
	}
	if r.Con == nil {
		return nil, util.Err("not connected")
	}
	if r.exec_c >= math.MaxUint16 {
		return nil, util.Err("two many exector")
	}
	bys, err := r.V2B()(args)
	if err != nil {
		return nil, err
	}
	tc := &chan_c{
		B:    m,
		BS:   bs,
		C:    make(chan bool),
		Data: bys,
	}
	r.req_c <- tc
	<-tc.C
	close(tc.C)
	if tc.Err == nil {
		defer tc.Back.Done()
		if bs {
			return tc.Back.Data()[1:], nil
		} else {
			return tc.Back.Data(), nil
		}
	} else {
		return nil, tc.Err
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
	log_d("RC_Con starting...")
	go r.Run_()
}

//stop gorutine
func (r *RC_Con) Stop() {
	r.running = false
	log_d("RC_Con stopping...")
}

//run the process of send/receive command(sync).
func (r *RC_Con) Run_() {
	if r.running {
		log.W("RC_Con already running....")
		return
	}
	cm := map[uint16]*chan_c{}
	buf := make([]byte, 3)
	r.running = true
	var trun bool = true
	tk := time.Tick(r.Sleep * time.Millisecond)
	for trun {
		select {
		case cmd := <-r.bc.back_c:
			// log.D("cmd ->%p->%v", r, cmd)
			tid := binary.BigEndian.Uint16(cmd.rid)
			if tc, ok := cm[tid]; ok {
				r.exec_c--
				delete(cm, tid)
				tc.Back = cmd
				if cmd.mark == 0 {
					tc.C <- true
				} else {
					tc.Err = new_rc_err(int(cmd.mark), cmd.Data())
					tc.C <- false
				}
			} else {
				cmd.Done()
				log.W("back chan not found by id(%v) on %v", tid, cm)
			}
		case tc := <-r.req_c:
			con := r.Con
			// log.D("tc ->%p->%v", r, tc)
			r.exec_id++
			binary.BigEndian.PutUint16(buf, r.exec_id)
			buf[2] = 0
			if tc.BS {
				_, tc.Err = con.Writeb(buf, []byte{tc.B}, tc.Data)
			} else {
				_, tc.Err = con.Writeb(buf, tc.Data)
			}
			if tc.Err == nil {
				cm[r.exec_id] = tc
				r.exec_c++
			} else {
				r.err = tc.Err
				tc.C <- false
			}
		case <-tk:
			trun = r.running && r.Con != nil
		}
	}
	log.D("clearing all waiting exec(%v),err(%v)", len(cm), r.err)
	//clear all waiting.
	for _, tc := range cm {
		tc.Err = util.Err("stopped")
		tc.C <- false
	}
	r.running = false
}

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
