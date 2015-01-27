package im

import (
	"fmt"
	"github.com/Centny/gwf/netw"
)

const (
	MS_DONE    = "D"
	MS_PENDING = "P"    //message not send
	MS_ERR     = "E->:" //message send error.
)

type Msg struct {
	netw.Cmd `-`
	Id       string            `_id`
	S        string            `json:"s"` //the sender R.
	R        []string          `json:"r"` //logic R
	D        string            `json:"d"` //target user R.
	T        byte              `json:"t"` //type.
	C        []byte            `json:"c"` //the content.
	Ms       map[string]string `json:"-"` //send status for user R.
}

type DsMsg struct {
	M  Msg               `json:"m"`
	RC map[string]string `json:"rc"`
}

const (
	CT_TCP = 0
	CT_WS  = 10
)

//connection
type Con struct {
	Sid string `json:"-"`   //server id
	Cid string `json:"cid"` //connection id
	R   string `json:"r"`   //the receive SN
	S   string `json:"s"`   //the IM receiver status.
	T   byte   `json:'t'`   //the connection type.
}

//online server
type Srv struct {
	Id    string `_id`
	Sid   string `json:"sid"`   //server id
	Host  string `json:"host"`  //server addr
	Port  int    `json:"port"`  //server port.
	Token string `json:"token"` //server login token
}

func (s *Srv) Addr() string {
	return fmt.Sprintf("%v:%v", s.Host, s.Port)
}
