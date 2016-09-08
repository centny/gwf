package util

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
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

func AryS2Map(vals []string) map[string]bool {
	var res = map[string]bool{}
	for _, val := range vals {
		res[val] = true
	}
	return res
}
