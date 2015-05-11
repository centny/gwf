package util

import (
	"errors"
	"fmt"
	"strconv"
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
