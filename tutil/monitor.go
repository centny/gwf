package tutil

import (
	"fmt"
	"github.com/Centny/gwf/util"
	"math"
	"path/filepath"
	"sort"
	"strings"
	"sync"
)

type Statable interface {
	State() (interface{}, error)
}

type State struct {
	Name  string `json:"name"`
	Min   int64  `json:"min"`
	Max   int64  `json:"max"`
	Total int64  `json:"total"`
	Count int64  `json:"count"`
}

type Monitor struct {
	Used     map[string]*State
	Pending  map[string]int64
	lck      sync.RWMutex
	sequence uint64
}

func NewMonitor() *Monitor {
	return &Monitor{
		Used:    map[string]*State{},
		Pending: map[string]int64{},
		lck:     sync.RWMutex{},
	}
}
func (m *Monitor) Start(name string) string {
	m.lck.Lock()
	defer m.lck.Unlock()
	m.sequence += 1
	var id = fmt.Sprintf("%v/%v", name, m.sequence)
	m.Pending[id] = util.Now()
	return id
}

func (m *Monitor) Done(id string) {
	m.lck.Lock()
	defer m.lck.Unlock()
	beg, ok := m.Pending[id]
	if !ok {
		return
	}
	delete(m.Pending, id)
	name := filepath.Dir(id)
	name = strings.TrimSuffix(name, "/")
	old, ok := m.Used[name]
	if !ok {
		old = &State{Name: name, Min: math.MaxInt64}
	}
	used := util.Now() - beg
	old.Total += used
	old.Count += 1
	if old.Max < used {
		old.Max = used
	}
	if old.Min > used {
		old.Min = used
	}
	m.Used[name] = old
}

func (m *Monitor) State() (interface{}, error) {
	m.lck.RLock()
	defer m.lck.RUnlock()
	var used = []util.Map{}
	for _, u := range m.Used {
		used = append(used, util.Map{
			"name":  u.Name,
			"min":   u.Min,
			"max":   u.Max,
			"total": u.Total,
			"count": u.Count,
			"avg":   u.Total / u.Count,
		})
	}
	sort.Sort(util.NewMapSorterV(used, "/avg", 0, true))
	return util.Map{
		"used":    used,
		"pending": m.Pending,
	}, nil
}
