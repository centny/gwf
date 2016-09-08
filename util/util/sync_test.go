package util

import "testing"

func TestWaitGroup(t *testing.T) {
	var wg = WaitGroup{}
	wg.Add(10)
	if wg.Size() != 10 {
		t.Error("error")
		return
	}
	wg.Done()
	if wg.Size() != 9 {
		t.Error("error")
		return
	}
}
