package util

import (
	"errors"
	"fmt"
	"reflect"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
)

type Array struct {
	vals []interface{}
}

func CreateArray(slen int) *Array {
	ary := &Array{}
	ary.vals = make([]interface{}, 0, slen)
	return ary
}
func (a *Array) Add(val interface{}) {
	a.vals = append(a.vals, val)
}
func (a *Array) Del(idx int) {
	copy(a.vals[idx:len(a.vals)], a.vals[idx+1:len(a.vals)])
	a.vals = a.vals[0 : len(a.vals)-1]
}
func (a *Array) At(idx int) interface{} {
	return a.vals[idx]
}
func (a *Array) Len() int {
	return len(a.vals)
}
func (a *Array) Ary() []interface{} {
	return a.vals
}

type Pair struct {
	Left  interface{}
	Right interface{}
}

func Err(f string, args ...interface{}) error {
	return errors.New(fmt.Sprintf(f, args...))
}

type WaitGroup struct {
	sync.WaitGroup
	c int32
}

func (w *WaitGroup) Add(i int) {
	w.WaitGroup.Add(i)
	atomic.AddInt32(&w.c, int32(i))
}
func (w *WaitGroup) Done() {
	w.WaitGroup.Done()
	atomic.AddInt32(&w.c, int32(-1))
}
func (w *WaitGroup) Size() int {
	return int(w.c)
}

func ParseInt(s string) (int, error) {
	val, err := strconv.ParseInt(s, 10, 32)
	return int(val), err
}

func ParseInt64(s string) (int64, error) {
	return strconv.ParseInt(s, 10, 64)
}

func ParseInt64s(ss []string) ([]int64, error) {
	is := []int64{}
	for _, s := range ss {
		i, err := ParseInt64(s)
		if err != nil {
			return nil, err
		}
		is = append(is, i)
	}
	return is, nil
}

func ParseInt64s2(s, sep string) ([]int64, error) {
	return ParseInt64s(strings.Split(s, sep))
}

func ParseInts(ss []string) ([]int, error) {
	is := []int{}
	for _, s := range ss {
		i, err := ParseInt(s)
		if err != nil {
			return nil, err
		}
		is = append(is, i)
	}
	return is, nil
}
func ParseInts2(s, sep string) ([]int, error) {
	return ParseInts(strings.Split(s, sep))
}
func CPU() int {
	i := runtime.NumCPU()
	if i < 2 {
		return i
	} else {
		return i - 1
	}
}

func AryS2Map(vals []string) map[string]bool {
	var res = map[string]bool{}
	for _, val := range vals {
		res[val] = true
	}
	return res
}

type SorterLess interface {
	Less(a, b interface{}, desc bool) bool
}

type SorterLessF func(a, b interface{}, desc bool) bool

func (s SorterLessF) Less(a, b interface{}, desc bool) bool {
	return s(a, b, desc)
}

type Sorter struct {
	less SorterLess
	vals reflect.Value
	Desc bool
}

func NewSorter(less SorterLess, vals interface{}) *Sorter {
	var rfVal = reflect.ValueOf(vals)
	if rfVal.Kind() != reflect.Slice {
		panic("the vals type is not slice")
	}
	return &Sorter{
		less: less,
		vals: rfVal,
	}
}

func NewIntSorter(vals interface{}) *Sorter {
	return NewSorter(SorterLessF(IntLess), vals)
}

func NewFloatSorter(vals interface{}) *Sorter {
	return NewSorter(SorterLessF(FloatLess), vals)
}

func NewStringSorter(vals interface{}) *Sorter {
	return NewSorter(SorterLessF(StringLess), vals)
}

func NewMapIntSorter(key string, vals interface{}) *Sorter {
	return NewSorter(MapIntLess(key), vals)
}

func NewMapFloatSorter(key string, vals interface{}) *Sorter {
	return NewSorter(MapFloatLess(key), vals)
}

func NewMapStringSorter(key string, vals interface{}) *Sorter {
	return NewSorter(MapFloatLess(key), vals)
}

func NewFieldIntSorter(key string, vals interface{}) *Sorter {
	return NewSorter(FieldIntLess(key), vals)
}

func NewFieldFloatSorter(key string, vals interface{}) *Sorter {
	return NewSorter(FieldFloatLess(key), vals)
}

func NewFieldStringSorter(key string, vals interface{}) *Sorter {
	return NewSorter(FieldFloatLess(key), vals)
}

func (s *Sorter) Len() int {
	return s.vals.Len()
}

func (s *Sorter) Less(i, j int) bool {
	return s.less.Less(s.vals.Index(i).Interface(), s.vals.Index(j).Interface(), s.Desc)
}
func (s *Sorter) Swap(i, j int) {
	x, y := s.vals.Index(i).Interface(), s.vals.Index(j).Interface()
	s.vals.Index(i).Set(reflect.ValueOf(y))
	s.vals.Index(j).Set(reflect.ValueOf(x))
}
func (s *Sorter) Sort(desc bool) {
	s.Desc = desc
	sort.Sort(s)
}

func StringLess(a, b interface{}, desc bool) bool {
	if desc {
		return StrVal(a) > StrVal(b)
	} else {
		return StrVal(a) < StrVal(b)
	}
}

func IntLess(a, b interface{}, desc bool) bool {
	if desc {
		return IntVal(a) > IntVal(b)
	} else {
		return IntVal(a) < IntVal(b)
	}
}

func FloatLess(a, b interface{}, desc bool) bool {
	if desc {
		return FloatVal(a) > FloatVal(b)
	} else {
		return FloatVal(a) < FloatVal(b)
	}
}

type MapIntLess string

func (m MapIntLess) Less(a, b interface{}, desc bool) bool {
	var key = string(m)
	if desc {
		return MapVal(a).IntValP(key) > MapVal(b).IntValP(key)
	} else {
		return MapVal(a).IntValP(key) < MapVal(b).IntValP(key)
	}
}

type MapFloatLess string

func (m MapFloatLess) Less(a, b interface{}, desc bool) bool {
	var key = string(m)
	if desc {
		return MapVal(a).FloatValP(key) > MapVal(b).FloatValP(key)
	} else {
		return MapVal(a).FloatValP(key) < MapVal(b).FloatValP(key)
	}
}

type MapStringLess string

func (m MapStringLess) Less(a, b interface{}, desc bool) bool {
	var key = string(m)
	if desc {
		return MapVal(a).StrValP(key) > MapVal(b).StrValP(key)
	} else {
		return MapVal(a).StrValP(key) < MapVal(b).StrValP(key)
	}
}

type FieldIntLess string

func (f FieldIntLess) Less(a, b interface{}, desc bool) bool {
	var key = string(f)
	if desc {
		return IntVal(reflect.ValueOf(a).FieldByName(key).Interface()) > IntVal(reflect.ValueOf(b).FieldByName(key).Interface())
	}
	return IntVal(reflect.ValueOf(a).FieldByName(key).Interface()) < IntVal(reflect.ValueOf(b).FieldByName(key).Interface())
}

type FieldFloatLess string

func (f FieldFloatLess) Less(a, b interface{}, desc bool) bool {
	var key = string(f)
	if desc {
		return FloatVal(reflect.ValueOf(a).FieldByName(key).Interface()) > FloatVal(reflect.ValueOf(b).FieldByName(key).Interface())
	}
	return FloatVal(reflect.ValueOf(a).FieldByName(key).Interface()) < FloatVal(reflect.ValueOf(b).FieldByName(key).Interface())
}

type FieldStringLess string

func (f FieldStringLess) Less(a, b interface{}, desc bool) bool {
	var key = string(f)
	if desc {
		return StrVal(reflect.ValueOf(a).FieldByName(key).Interface()) > StrVal(reflect.ValueOf(b).FieldByName(key).Interface())
	}
	return StrVal(reflect.ValueOf(a).FieldByName(key).Interface()) < StrVal(reflect.ValueOf(b).FieldByName(key).Interface())
}
