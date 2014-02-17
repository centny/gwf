package util

import (
	"fmt"
	"testing"
)

type S1 struct {
	A string `m2s:"VA"`
	B string `m2s:"VB"`
}

func TestM2S(t *testing.T) {
	m := make(map[string]interface{})
	m["VA"] = "S1_A"
	m["VB"] = "S2_B"
	m1 := make(map[string]interface{})
	m1["VA"] = "S3_A"
	m1["VB"] = "S4_B"
	m3 := make(map[string]interface{})
	//
	mary := make([]Map, 0, 2)
	mary = append(mary, m)
	mary = append(mary, m1)
	mary2 := append(mary, m3)
	//
	//
	var dest S1
	M2S(m, &dest)
	if dest.A != "S1_A" || dest.B != "S2_B" {
		t.Error("value invalid ...")
		return
	}
	var dests []S1
	Ms2Ss(mary, &dests)
	if len(dests) != 2 {
		t.Error("result count is invalid ...")
		return
	}
	fmt.Println(dests)
	var dests2 []S1
	Ms2Ss(mary2, &dests2)
	if len(dests) != 2 {
		t.Error("result count is invalid ...")
		return
	}
	fmt.Println(dests2)
}

func TestM2SErr(t *testing.T) {
	m := make(map[string]interface{})
	m["VA"] = "S1_A"
	m["VB"] = "S2_B"
	m1 := make(map[string]interface{})
	m1["VA"] = "S3_A"
	m1["VB"] = "S4_B"
	m3 := make(map[string]interface{})
	//
	mary := make([]Map, 0, 2)
	mary = append(mary, m)
	mary = append(mary, m1)
	mary2 := make([]Map, 0, 2)
	//
	//
	var dest S1
	M2S(nil, &dest)
	M2S(m, nil)
	M2S(m3, &dest)
	//
	var dests []S1
	Ms2Ss(nil, &dests)
	Ms2Ss(mary, nil)
	Ms2Ss(mary2, &dests)
}
