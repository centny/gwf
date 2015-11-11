package pool

import (
	"fmt"
	"testing"
)

func TestQueue(t *testing.T) {
	qe := NewQueue(true)
	qe.Push(1, 2, 3)
	if qe.Len() != 3 {
		t.Error("error")
		return
	}
	if qe.Pollv().(int) != 1 {
		t.Error("error")
		return
	}
	if qe.Pollv().(int) != 2 {
		t.Error("error")
		return
	}
	if qe.Len() != 1 {
		t.Error("error")
		return
	}
	qe.Push(4, 5)
	if qe.Len() != 3 {
		t.Error("error")
		return
	}
	if qe.Pollv().(int) != 3 {
		t.Error("error")
		return
	}
	if qe.Pollv().(int) != 4 {
		t.Error("error")
		return
	}
	if qe.Pollv().(int) != 5 {
		t.Error("error")
		return
	}
	if qe.Len() != 0 {
		t.Error("error")
		return
	}
	if qe.Poll() != nil {
		t.Error("error")
		return
	}
	func() {
		defer func() {
			if err := recover(); err == nil {
				t.Error("error")
			}
		}()
		qe.Pollv()
	}()
	fmt.Println("done...")
}
