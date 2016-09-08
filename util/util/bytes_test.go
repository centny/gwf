package util

import "testing"

func TestSplitByte(t *testing.T) {
	var bys = []byte("1234567890")
	var bysa, bysb = SplitTwo(bys, 5)
	if string(bysa) != "12345" || string(bysb) != "67890" {
		t.Error("error")
		return
	}
	var bys1, bys2, bys3 = SplitThree(bys, 3, 6)
	if string(bys1) != "123" || string(bys2) != "456" || string(bys3) != "7890" {
		t.Error("error")
		return
	}
}
