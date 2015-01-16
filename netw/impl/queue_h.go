package impl

import (
	"github.com/Centny/gwf/netw"
)

//the queue command handler.
type QueueH struct {
	HS       []netw.CmdHandler    //handlers
	Continue int                  //continue result value.
	OnBreak  func(c netw.Cmd) int //on break callback.
}

//new queue handler.
func NewQueueH() *QueueH {
	qh := &QueueH{
		HS: []netw.CmdHandler{},
	}
	qh.OnBreak = qh.onbreak
	return qh
}
func (q *QueueH) OnCmd(c netw.Cmd) int {
	for _, h := range q.HS {
		res := h.OnCmd(c)
		if res == q.Continue {
			continue
		}
		return q.OnBreak(c)
	}
	return 0
}
func (q *QueueH) onbreak(c netw.Cmd) int {
	c.Done()
	return -1
}

//adding command handler.
func (q *QueueH) AddH(h netw.CmdHandler) {
	q.HS = append(q.HS, h)
}
