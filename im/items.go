package im

import (
	"fmt"
	"github.com/Centny/gwf/im/pb"
	// "github.com/Centny/gwf/log"
	"github.com/Centny/gwf/netw"
)

const (
	MS_DONE    = "D"
	MS_PENDING = "P:" //message not send
	// MS_ERR     = "E->:" //message send error.
)

type Msg struct {
	netw.Cmd `bson:"-" json:"-"`
	pb.ImMsg `bson:",inline"`
	//
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
	Id    string `json:"-" bson:"_id"` //the bson id
	Sid   string `json:"-"`            //server id
	Cid   string `json:"-"`            //connection id
	R     string `json:"r"`            //the receive SN
	S     string `json:"s"`            //the IM receiver status.
	T     byte   `json:"t"`            //the connection type in TCP/WS.
	C     int    `json:"c"`            //the connect category
	Token string `json:"token"`        //the login token
}

//online server
type Srv struct {
	Sid     string `bson:"_id" json:"sid"` //server id
	Host    string `json:"host"`           //server addr
	Port    int    `json:"port"`           //server port.
	WsAddr  string `json:"ws_addr"`        //server port.
	PubHost string `json:"pub_host"`       //server public port.
	PubPort int    `json:"pub_port"`       //server public port.
	Token   string `json:"token"`          //server login token
}

func (s *Srv) Addr() string {
	return fmt.Sprintf("%v:%v", s.PubHost, s.PubPort)
}

// type PCM struct {
// 	R string //the user R
// 	C []*Con //mapping connection.
// 	M []*Msg //mapping message.
// }

// func (p *PCM) Send(s Sender, db DbH) int {
// 	var sc int = 0
// 	for _, con := range p.C {
// 		for _, m := range p.M {
// 			m.D = &p.R
// 			err := s.Send(con.Cid, &m.ImMsg)
// 			if err == nil {
// 				sc++
// 				db.Update(m.GetI(), map[string]string{p.R: "D"})
// 				continue
// 			} else {
// 				log.W("sending unread message(%v) error:%v", m.ImMsg, err.Error())
// 				break
// 			}
// 		}
// 	}
// 	return sc
// }
