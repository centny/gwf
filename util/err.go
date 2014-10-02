package util

import (
	"fmt"
)

type Error struct {
	Type string `json:"type"`
	In   error  `json:"in"`
	Msg  string `json:"msg"`
}

func (e *Error) Error() string {
	return e.Msg
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
