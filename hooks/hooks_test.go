package hooks

import (
	"github.com/Centny/gwf/util"
	"testing"
)

type Abc struct {
	Err bool
}

func (a *Abc) Call(v interface{}, args ...interface{}) (interface{}, error) {
	if a.Err {
		return nil, util.Err("some error")
	} else {
		return nil, nil
	}
}
func TestHooks(t *testing.T) {
	nh := NewNameHooks2("A1", &Abc{})
	nh1 := NewNameHooks()
	nh1.AddHook("B1", &Abc{})
	nh1.AddHook("B2", &Abc{})
	nh1.AddHook("B2", &Abc{Err: true})
	AddHook("A", nh)
	AddHook("A", nh)
	AddHook("B", nh1)
	Call("A", "A1", nil)
	Call("B", "B1", nil)
	Call("B", "B2", nil)
	Call("B", "B3", nil)
	Call("B", "NOTFOUND", nil)
	Call("NOTFOUND", "NOTFOUND", nil)
	SetMockErr("A", "A1")
	Call("A", "A1", nil)
	ClsMockErr()
}
