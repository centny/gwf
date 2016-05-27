package filter

import (
	"github.com/Centny/gwf/routing"
	"github.com/Centny/gwf/tutil"
	"github.com/Centny/gwf/util"
)

type MonitorH struct {
	States map[string]tutil.Statable
}

func NewMonitorH() *MonitorH {
	return &MonitorH{
		States: map[string]tutil.Statable{},
	}
}

func (m *MonitorH) AddMonitor(key string, s tutil.Statable) *MonitorH {
	m.States[key] = s
	return m
}

func (m *MonitorH) SrvHTTP(hs *routing.HTTPSession) routing.HResult {
	var res = util.Map{}
	for key, s := range m.States {
		val, _ := s.State()
		res[key] = val
	}
	return hs.JRes(res)
}
