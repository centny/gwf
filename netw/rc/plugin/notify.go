package plugin

import (
	"sync"

	"strings"

	"github.com/Centny/gwf/log"
	"github.com/Centny/gwf/netw"
	"github.com/Centny/gwf/netw/impl"
	"github.com/Centny/gwf/netw/rc"
	"github.com/Centny/gwf/util"
)

// SharedNotify is the shared notify server instance
var SharedNotify = NewNotifySrv(nil)

// NotifyPostMessage is using SharedNotify and ignore notify when SharedNotify is not running
func NotifyPostMessage(m *Message) error {
	if SharedNotify.running {
		return SharedNotify.PostMessage(m)
	}
	return nil
}

// NotifyMessageMark is the message mark to transfter on RC socket channel
var NotifyMessageMark byte = 165

// Message is the notify message
type Message struct {
	ID     string   `bson:"_id" json:"id"`                            //the message id
	Oid    string   `bson:"oid" json:"oid"`                           //the owner id
	Owner  string   `bson:"owner" json:"owner"`                       //the owner type
	Type   string   `bson:"type" json:"type"`                         //the message type
	Attrs  util.Map `bson:"attrs" json:"attrs"`                       //external attributes
	Marked []string `bson:"marked,omitempty" json:"marked,omitempty"` //the key of already mark done
	Count  int      `bson:"count" json:"count"`                       //the done count
	Time   int64    `bson:"time" json:"time"`                         //the create time
}

// NotifyDb is the notify server database interface.
type NotifyDb interface {
	//adding message
	AddMessage(m *Message) error
	//remove message, controling by remove count
	RemoveMessage(id string) error
	//done message
	DoneMessage(mid, key string) (*Message, error)
	//return the message remove count by type
	RemoveCount(mtype string) (int, error)
	//list message by message fields.
	ListMessage(m *Message) ([]*Message, error)
}

// NotifySrv is the RC handler for notify server.
type NotifySrv struct {
	L       *rc.RC_Listener_m
	Db      NotifyDb
	clients map[string]map[string]int
	lck     sync.RWMutex
	msgChan chan *Message
	cidChan chan util.Map
	running bool
}

// NewNotifySrv the createor.
func NewNotifySrv(db NotifyDb) *NotifySrv {
	return &NotifySrv{
		Db:      db,
		clients: map[string]map[string]int{},
		lck:     sync.RWMutex{},
	}
}

// MonitorH is marking the type of message for monitoring.
func (n *NotifySrv) MonitorH(rc *impl.RCM_Cmd) (res interface{}, err error) {
	var mtype string
	var mask = 0
	err = rc.ValidF(`
        type,R|S,L:0;
        mask,O|I,R:0;
    `, &mtype, &mask)
	if err != nil {
		return "", err
	}
	var cid = rc.Id()
	n.lck.Lock()
	for _, t := range strings.Split(mtype, ",") {
		var clients = n.clients[t]
		if clients == nil {
			clients = map[string]int{}
		}
		clients[cid] = mask
		n.clients[t] = clients
		n.cidChan <- util.Map{
			"cid":  cid,
			"type": t,
		}
	}
	n.L.AddC_rc(cid, rc)
	n.lck.Unlock()
	slog("NotifySrv monitor type(%v) by mask(%v) on client(%v) success", mtype, mask, cid)
	return "OK", nil
}

// MarkH is marking message done
func (n *NotifySrv) MarkH(rc *impl.RCM_Cmd) (res interface{}, err error) {
	var mid, key string
	err = rc.ValidF(`
        mid,R|S,L:0;
        key,R|S,L:0;
    `, &mid, &key)
	if err != nil {
		return "", err
	}
	msg, err := n.Db.DoneMessage(mid, key)
	if err != nil {
		log.E("NotifySrv done message by key(%v),id(%v) fail with %v", key, mid, err)
		return "", err
	}
	count, err := n.Db.RemoveCount(msg.Type)
	if err != nil {
		log.E("NotifySrv get remove count by type(%v) fail with %v", msg.Type, err)
		return "", err
	}
	if msg.Count < count {
		slog("NotifySrv mark done(%v/%v) with id(%v),key(%v) success", msg.Count, count, mid, key)
		return "OK", nil
	}
	err = n.Db.RemoveMessage(msg.ID)
	if err != nil {
		log.E("NotifySrv remove message by id(%v) fail with %v", mid, err)
	} else {
		slog("NotifySrv remove message by id(%v) success", mid)
	}
	return "OK", nil
}

// Hand is register the handler to listener.
func (n *NotifySrv) Hand(l *rc.RC_Listener_m) {
	n.L = l
	l.AddHFunc("notify/monitor", n.MonitorH)
	l.AddHFunc("notify/mark", n.MarkH)
}

// PostMessage send notify message
func (n *NotifySrv) PostMessage(m *Message) error {
	var err = n.Db.AddMessage(m)
	if err != nil {
		log.E("NotifySrv add message(%v) fail with %v", util.S2Json(m), err)
		return err
	}
	n.msgChan <- m
	slog("NotifySrv post message(%v) success", util.S2Json(m))
	return nil
}

// Start is starting the async post task
func (n *NotifySrv) Start() {
	if n.running {
		return
	}
	n.msgChan = make(chan *Message, 1000)
	n.cidChan = make(chan util.Map, 1000)
	go n.loopChan()
}

// Stop is stopping the async post task
func (n *NotifySrv) Stop() {
	close(n.msgChan)
	close(n.cidChan)
}

func (n *NotifySrv) loopChan() {
	log.D("NotifySrv start message loop...")
	n.running = true
	for n.running {
		select {
		case msg := <-n.msgChan:
			if msg == nil {
				n.running = false
				break
			}
			n.notifyMessage(msg)
		case cid := <-n.cidChan:
			if cid == nil {
				n.running = false
				break
			}
			n.notifyClient(cid.StrVal("cid"), cid.StrVal("type"))
		}
	}
	log.D("NotifySrv loop is stopped...")
}

func (n *NotifySrv) notifyMessage(m *Message) {
	var cids []string
	n.lck.Lock()
	for cid := range n.clients[m.Type] {
		cids = append(cids, cid)
	}
	n.lck.Unlock()
	for _, cid := range cids {
		var msgc = n.L.MsgC(cid)
		if msgc == nil {
			continue
		}
		msgc.Writev2([]byte{NotifyMessageMark}, m)
	}
}

func (n *NotifySrv) notifyClient(cid, mtype string) {
	var ms, err = n.Db.ListMessage(&Message{
		Type: mtype,
	})
	if err != nil {
		log.E("NotifySrv list message by type(%v) fail with ", mtype, err)
		return
	}
	var msgc = n.L.MsgC(cid)
	if msgc == nil {
		log.E("NotifySrv find message socket client by id(%v) fail with not found ", cid)
		return
	}
	for _, m := range ms {
		msgc.Writev2([]byte{NotifyMessageMark}, m)
	}
}

// NotifyHandler is the notify handler for notify client.
type NotifyHandler interface {
	//on message received
	OnMessage(n *NotifyClient, m *Message) error
}

// NotifyHandlerF is handler func
type NotifyHandlerF func(n *NotifyClient, m *Message) error

// OnMessage see NotifyHandler
func (f NotifyHandlerF) OnMessage(n *NotifyClient, m *Message) error {
	return f(n, m)
}

// NotifyClient is the RC handler for notify client
type NotifyClient struct {
	R *rc.RC_Runner_m
	H NotifyHandler
}

// NewNotifyClient is creator.
func NewNotifyClient(h NotifyHandler) *NotifyClient {
	return &NotifyClient{
		H: h,
	}
}

// SetRunner initial runner.
func (n *NotifyClient) SetRunner(r *rc.RC_Runner_m, obdh *impl.OBDH) {
	n.R = r
	obdh.AddH(NotifyMessageMark, n)
}

// Monitor start monitor by message type and mark
func (n *NotifyClient) Monitor(mtype string, mask int) error {
	var _, err = n.R.VExec_s("notify/monitor", util.Map{
		"type": mtype,
		"mask": mask,
	})
	return err
}

// Mark is marking message is done with key
func (n *NotifyClient) Mark(mid, key string) error {
	var _, err = n.R.VExec_s("notify/mark", util.Map{
		"mid": mid,
		"key": key,
	})
	return err
}

// OnCmd handler message.
func (n *NotifyClient) OnCmd(c netw.Cmd) int {
	var msg = &Message{}
	var _, err = c.V(msg)
	if err != nil {
		log.E("NotifyClient parsing receive message fail with %v", err)
		return -1
	}
	err = n.H.OnMessage(n, msg)
	if err != nil {
		log.E("NotifyClient executing handler fail with %v", err)
		return -1
	}
	slog("NotifyClient executing handler for message(%v) success", util.S2Json(msg))
	return 0
}

// NotifyMemDb is impl to NotifyDb on memory
type NotifyMemDb struct {
	MS    map[string]*Message
	Count map[string]int
	Lck   sync.RWMutex
}

// NewNotifyMemDb is NotifyMemDb creator
func NewNotifyMemDb() *NotifyMemDb {
	return &NotifyMemDb{
		MS:    map[string]*Message{},
		Count: map[string]int{},
		Lck:   sync.RWMutex{},
	}
}

// AddMessage @see NotifyDb
func (n *NotifyMemDb) AddMessage(m *Message) error {
	n.Lck.Lock()
	defer n.Lck.Unlock()
	m.ID = util.UUID()
	n.MS[m.ID] = m
	return nil
}

// RemoveMessage @see NotifyDb
func (n *NotifyMemDb) RemoveMessage(id string) error {
	n.Lck.Lock()
	defer n.Lck.Unlock()
	delete(n.MS, id)
	return nil
}

// DoneMessage @see NotifyDb
func (n *NotifyMemDb) DoneMessage(mid, key string) (msg *Message, err error) {
	n.Lck.Lock()
	defer n.Lck.Unlock()
	msg = n.MS[mid]
	if msg == nil {
		return nil, util.NOT_FOUND
	}
	msg.Count++
	for _, marked := range msg.Marked {
		if marked == key {
			return msg, nil
		}
	}
	msg.Marked = append(msg.Marked, key)
	return msg, nil
}

// RemoveCount @see NotifyDb
func (n *NotifyMemDb) RemoveCount(mtype string) (count int, err error) {
	if val, ok := n.Count[mtype]; ok {
		return val, nil
	}
	return 1, nil
}

// ListMessage @see NotifyDb
func (n *NotifyMemDb) ListMessage(m *Message) (ms []*Message, err error) {
	n.Lck.Lock()
	defer n.Lck.Unlock()
	for _, m := range n.MS {
		if m.Type == m.Type {
			ms = append(ms, m)
		}
	}
	return
}
