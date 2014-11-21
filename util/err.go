package util

import (
	"errors"
	"fmt"
)

type Error struct {
	Type string `json:"type"`
	In   error  `json:"in"`
	Msg  string `json:"msg"`
}

func (e *Error) Error() string {
	if len(e.Msg) > 0 {
		return e.Msg
	}
	if e.In == nil {
		return ""
	} else {
		return e.In.Error()
	}
}
func (e *Error) String() string {
	return e.Error()
}
func NewErr(typ string, in error) *Error {
	return &Error{
		Type: typ,
		In:   in,
		Msg:  "",
	}
}
func NewErr2(typ string, f string, args ...interface{}) *Error {
	return &Error{
		Type: typ,
		In:   nil,
		Msg:  fmt.Sprintf(f, args...),
	}
}

var NOT_FOUND = errors.New("NOT FOUND")
