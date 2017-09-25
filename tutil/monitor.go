package tutil

import (
	"fmt"
	"math"
	"path/filepath"
	"sort"
	"strings"
	"sync"

	"github.com/Centny/gwf/util"
)

type Statable interface {
	State() (interface{}, error)
}

type State struct {
	Name    string `json:"name"`
	Min     int64  `json:"min"`
	Max     int64  `json:"max"`
	Total   int64  `json:"total"`
	Count   int64  `json:"count"`
	ConcMax int64  `json:"con_max"`
	ConcAvg int64  `json:"con_avg"`
	//
	concAll   uint64
	concCount uint64
}

type Monitor struct {
	Used     map[string]*State
	Pending  map[string]int64
	max      map[string]int64
	lck      sync.RWMutex
	sequence uint64
}

func NewMonitor() *Monitor {
	return &Monitor{
		Used:    map[string]*State{},
		Pending: map[string]int64{},
		max:     map[string]int64{},
		lck:     sync.RWMutex{},
	}
}
func (m *Monitor) Start(name string) string {
	m.lck.Lock()
	defer m.lck.Unlock()
	m.sequence += 1
	var id = fmt.Sprintf("%v/%v", name, m.sequence)
	m.Pending[id] = util.Now()
	m.max[name]++
	old, ok := m.Used[name]
	if !ok {
		old = &State{Name: name, Min: math.MaxInt64}
	}
	old.concAll += uint64(m.max[name])
	old.concCount++
	old.ConcAvg = int64(old.concAll / old.concCount)
	if old.ConcMax < m.max[name] {
		old.ConcMax = m.max[name]
	}
	return id
}

func (m *Monitor) Start_(id string) {
	m.lck.Lock()
	defer m.lck.Unlock()
	m.Pending[id] = util.Now()
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
	m.max[name]--
}

func (m *Monitor) State() (interface{}, error) {
	m.lck.RLock()
	defer m.lck.RUnlock()
	var used = []util.Map{}
	for _, u := range m.Used {
		used = append(used, util.Map{
			"name":     u.Name,
			"min":      u.Min,
			"max":      u.Max,
			"total":    u.Total,
			"count":    u.Count,
			"avg":      u.Total / u.Count,
			"conc_max": u.ConcMax,
			"conc_avg": u.ConcAvg,
		})
	}
	sort.Sort(util.NewMapSorterV(used, "/avg", 0, true))
	//
	var pending = map[string]int64{}
	for k, v := range m.Pending {
		pending[k] = v
	}
	//
	return util.Map{
		"used":    used,
		"pending": pending,
	}, nil
}
