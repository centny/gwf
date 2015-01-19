package im

import (
	"github.com/Centny/gwf/netw"
)

const (
	MS_DONE    = "D"
	MS_PENDING = "P"    //message not send
	MS_ERR     = "E->:" //message send error.
)

type Msg struct {
	netw.Cmd `json:"-"`
	Id       string            `json:"id" "_id"`
	S        string            `json:"s"`
	R        []string          `json:"r"`
	T        byte              `json:"t"`
	C        []byte            `json:"c"`
	Ms       map[string]string `json:"-"`
}

type DsMsg struct {
	M  Msg               `json:"m"`
	RC map[string]string `json:"rc"`
}

//connection
type Con struct {
	Sid string `json:"sid"` //server id
	Cid string `json:"cid"` //connection id
	R   string `json:"r"`   //the receive SN
	S   string `json:"s"`   //the IM receiver status.
}

//online server
type Srv struct {
	Sid   string `json:"sid"`   //server id
	Addr  string `json:"addr"`  //server addr
	Token string `json:"token"` //server login token
}
