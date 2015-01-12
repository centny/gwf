package handler

import (
	"encoding/json"
	"fmt"
	"github.com/Centny/gwf/log"
	"github.com/Centny/gwf/netw"
	"github.com/Centny/gwf/util"
	"reflect"
)

//function for covert struct to []byte
type V2Byte func(v interface{}) ([]byte, error)

//function for covert []byte to struct or map.
type Byte2V func(bys []byte, dest interface{}) error

//
type RC_V_C struct {
	*RC_C
	V2B V2Byte
	B2V Byte2V
}

func NewRC_V_C(v2b V2Byte, b2v Byte2V) *RC_V_C {
	return &RC_V_C{
		V2B:  v2b,
		B2V:  b2v,
		RC_C: NewRC_C(),
	}
}
func NewRC_Json_C() *RC_V_C {
	return NewRC_V_C(json.Marshal, func(bys []byte, dest interface{}) error {
		var err error = nil
		if reflect.Indirect(reflect.ValueOf(dest)).Kind() == reflect.Map {
			err = json.Unmarshal(bys, dest)
		} else {
			mv := util.Map{}
			err = json.Unmarshal(bys, &mv)
			mv.ToS(dest)
		}
		return err
	})
}
func (r *RC_V_C) Exec(v interface{}, dest interface{}) error {
	bys, err := r.V2B(v)
	if err != nil {
		return err
	}
	tc, err := r.RC_C.Exec(bys)
	if err != nil {
		return err
	}
	defer tc.Done()
	return r.B2V(tc.Data, dest)
}

//the remote command server call back command struct.
type RC_V_Cmd struct {
	*RC_Cmd
	*util.Map
}

//the extended command handler.
type RC_V_H interface {
	netw.ConHandler
	//calling when one entire command have been received.
	OnCmd(rc *RC_V_Cmd) (interface{}, error)
}

//the remote command server handler.
type RC_V_S struct {
	H   RC_V_H
	V2B V2Byte
	B2V Byte2V
}

//new remote command server handler.
func NewRC_V_S(h RC_V_H, v2b V2Byte, b2v Byte2V) *RC_V_S {
	return &RC_V_S{
		H:   h,
		V2B: v2b,
		B2V: b2v,
	}
}
func NewRC_Json_S(h RC_V_H) *RC_V_S {
	return NewRC_V_S(h, json.Marshal, json.Unmarshal)
}
func (r *RC_V_S) OnConn(c *netw.Con) bool {
	return r.H.OnConn(c)
}
func (r *RC_V_S) OnClose(c *netw.Con) {
	r.H.OnClose(c)
}
func (r *RC_V_S) w_v(c *RC_Cmd, v interface{}) {
	bys, err := r.V2B(v)
	if err == nil {
		c.Write(bys)
	} else {
		log.E("RC_V_S V2B error:%v", err.Error())
		r.w_err(c, err)
	}
}
func (r *RC_V_S) w_err(c *RC_Cmd, err error) {
	c.Write([]byte(fmt.Sprintf(`{"err":"%v"}`, err.Error())))
}
func (r *RC_V_S) OnCmd(c *RC_Cmd) {
	var mv util.Map
	err := r.B2V(c.Data, &mv)
	if err != nil {
		r.w_err(c, err)
		return
	}
	tv, err := r.H.OnCmd(&RC_V_Cmd{
		RC_Cmd: c,
		Map:    &mv,
	})
	if err == nil {
		r.w_v(c, tv)
	} else {
		r.w_err(c, err)
	}
}
