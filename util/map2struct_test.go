package util

import (
	"fmt"
	"testing"
	"time"
)

type S1 struct {
	VB string
	A  string    `m2s:"VA"`
	B  string    `m2s:"VB"`
	C  time.Time `m2s:"T" tf:"2006-01-02 15:04:05`
	D  time.Time `m2s:"T"`
	E  time.Time `m2s:"T_L"`
	F  time.Time `m2s:"VA"`
	G  time.Time `m2s:"GT"`
	H  time.Time `m2s:"HT"`
	I  time.Time `m2s:"IT"`
}

func TestM2S(t *testing.T) {
	tt := 1392636100998
	m := make(map[string]interface{})
	m["VA"] = "S1_A"
	m["VB"] = "S2_B"
	m["T"] = "2014-02-17 11:50:05"
	m["T_L"] = tt
	m["GT"] = time.Now()
	m["HT"] = int32(tt)
	m["IT"] = int64(tt)
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
	fmt.Println(Timestamp(dest.E), tt)
	if int64(tt) != Timestamp(dest.E) {
		t.Error("value not corrent ...")
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

func TestTime(t *testing.T) {
	ti := time.Now().UnixNano() / (1e6)
	fmt.Println(ti)
	fmt.Println(time.Unix(ti/1e3, (ti % 1e3)))
}
