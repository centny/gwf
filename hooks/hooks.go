package hooks

import (
	"errors"
	"fmt"
)

type Hook interface {
	Call(action string, v interface{}, args ...interface{}) (interface{}, error)
}

var hooks_ map[string][]Hook = map[string][]Hook{}
var mk_err map[string]bool = map[string]bool{}

func Call(name, action string, v interface{}, args ...interface{}) (interface{}, error) {
	if _, ok := mk_err[fmt.Sprintf("%v-%v", name, action)]; ok {
		return v, errors.New("Mock Hook error")
	}
	fs, ok := hooks_[name]
	if !ok || len(fs) < 1 {
		return v, nil
	}
	var val interface{} = v
	var err error = nil
	for _, f := range fs {
		val, err = f.Call(action, val, args...)
		if err != nil {
			return val, err
		}
	}
	return val, nil
}
func AddHook(name string, hk Hook) {
	if fs, ok := hooks_[name]; ok {
		fs = append(fs, hk)
		hooks_[name] = fs
	} else {
		hooks_[name] = []Hook{hk}
	}
}
func SetMockErr(name, action string) {
	mk_err[fmt.Sprintf("%v-%v", name, action)] = true
}
func ClsMockErr() {
	mk_err = map[string]bool{}
}

type NameHooks struct {
	Actions map[string][]ActionHook
}

func (n *NameHooks) Call(action string, v interface{}, args ...interface{}) (interface{}, error) {
	fs, ok := n.Actions[action]
	if !ok || len(fs) < 1 {
		return v, nil
	}
	var val interface{} = v
	var err error = nil
	for _, f := range fs {
		val, err = f.Call(val, args...)
		if err != nil {
			return val, err
		}
	}
	return val, nil
}
func (n *NameHooks) AddHook(action string, hk ActionHook) {
	if fs, ok := n.Actions[action]; ok {
		fs = append(fs, hk)
		n.Actions[action] = fs
	} else {
		n.Actions[action] = []ActionHook{hk}
	}
}

//
type ActionHook interface {
	Call(v interface{}, args ...interface{}) (interface{}, error)
}

func NewNameHooks() *NameHooks {
	nh := &NameHooks{}
	nh.Actions = map[string][]ActionHook{}
	return nh
}
func NewNameHooks2(action string, ah ActionHook) *NameHooks {
	nh := &NameHooks{}
	nh.Actions = map[string][]ActionHook{}
	nh.AddHook(action, ah)
	return nh
}
