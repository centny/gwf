package im

import (
	"fmt"
	"github.com/Centny/gwf/im/pb"
	"github.com/Centny/gwf/netw"
)

const (
	MS_DONE    = "D"
	MS_PENDING = "P"    //message not send
	MS_ERR     = "E->:" //message send error.
)

type Msg struct {
	netw.Cmd `bson:"-" json:"-"`
	pb.ImMsg
	Ms map[string]string `json:"-"` //send status for user R.
}

// type DsMsg struct {
// 	M  pb.ImMsg          `json:"m"`
// 	RC map[string]string `json:"rc"`
// }

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
	Sid   string `bson:"_id" json:"sid"` //server id
	Host  string `json:"host"`           //server addr
	Port  int    `json:"port"`           //server port.
	Token string `json:"token"`          //server login token
}

func (s *Srv) Addr() string {
	return fmt.Sprintf("%v:%v", s.Host, s.Port)
}
