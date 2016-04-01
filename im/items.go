package im

import (
	"fmt"
	"github.com/Centny/gwf/im/pb"
	// "github.com/Centny/gwf/log"
	"github.com/Centny/gwf/netw"
)

const (
	MS_DONE    = "D"
	MS_PENDING = "P" //message not send
	// MS_ERR     = "E->:" //message send error.
)

//message send status.
type MSS struct {
	R string `bson:"r"`
	S string `bson:"s"`
}
type Msg struct {
	netw.Cmd `bson:"-" json:"-"`
	pb.ImMsg `bson:",inline"`
	//
	Ms    map[string][]*MSS `json:"-"` //send status for user R.
	added map[string]bool   `json:"-" bson:"-"`
}

func (m *Msg) ams(r string, mss *MSS) {
	key := fmt.Sprintf("%v-%v", r, mss.R)
	if _, ok := m.added[key]; ok {
		return
	} else {
		m.Ms[r] = append(m.Ms[r], mss)
		m.added[key] = true
	}
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
	Id        string `json:"-" bson:"_id"`                 //the bson id
	Sid       string `json:"-" bson:"sid"`                 //server id
	Cid       string `json:"-" bson:"cid"`                 //connection id
	Uid       string `json:"uid" bson:"uid"`               //the receiver SN
	Status    string `json:"status" bson:"status"`         //the IM receiver status.
	ConType   byte   `json:"con_type" bson:"con_type"`     //the connection type in TCP/WS.
	LoginType int    `json:"login_type" bson:"login_type"` //the connect category
	Token     string `json:"token" bson:"token"`           //the login token
	Time      int64  `json:"time" bson:"time"`             //the login time.
}

//online server
type Srv struct {
	Sid     string `bson:"_id" json:"sid"`           //server id
	Host    string `bson:"hoost" json:"host"`        //server addr
	Port    int    `bson:"port" json:"port"`         //server port.
	WsAddr  string `bson:"ws_addr" json:"ws_addr"`   //server port.
	PubHost string `bson:"pub_host" json:"pub_host"` //server public port.
	PubPort int    `bson:"pub_port" json:"pub_port"` //server public port.
	Token   string `bson:"token" json:"token"`       //server login token
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

type LGR_Arg struct {
	GR []string `json:"gr"`
}
