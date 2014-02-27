package util

import (
	"fmt"
	"math"
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
	J  string    `m2s:"JV"`
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
	m["JV"] = nil
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

type S2 struct {
	A1 int   `m2s:"B1"`
	A2 int8  `m2s:"B2"`
	A3 int16 `m2s:"B3"`
	A4 int32 `m2s:"B4"`
	A5 int64 `m2s:"B5"`
	A6 int64 `m2s:"C1"`
	A7 int64 `m2s:"C2"`
	//
	B1 uint   `m2s:"A1"`
	B2 uint8  `m2s:"A2"`
	B3 uint16 `m2s:"A3"`
	B4 uint32 `m2s:"A4"`
	B5 uint64 `m2s:"A4"`
	B6 uint64 `m2s:"C1"`
	B7 uint64 `m2s:"C2"`
	//
	C1 float32 `m2s:"A1"`
	C2 float64 `m2s:"A2"`
	C3 float64 `m2s:"B1"`
	C4 float64 `m2s:"B2"`
}

func TestM2S2(t *testing.T) {
	m := make(map[string]interface{})
	m["A1"] = int(1)
	m["A2"] = int8(2)
	m["A3"] = int16(3)
	m["A4"] = int32(4)
	m["A5"] = int64(5)
	//
	m["B1"] = uint(6)
	m["B2"] = uint8(7)
	m["B3"] = uint16(8)
	m["B4"] = uint32(9)
	m["B5"] = uint64(10)
	//
	m["C1"] = float32(11)
	m["C2"] = float64(12)
	//
	var dest S2
	M2S(m, &dest)
	fmt.Println(dest)
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

func TestMap(t *testing.T) {
	m := Map{}
	//
	m["abc"] = "123"
	m["abc2"] = int(1)
	m["float32"] = float32(1)
	m["float64"] = float64(1)
	m["int"] = int(1)
	m["int8"] = int8(1)
	m["int16"] = int16(1)
	m["int32"] = int32(1)
	m["int64"] = int64(1)
	m["uint"] = uint(1)
	m["uint8"] = uint8(1)
	m["uint16"] = uint16(1)
	m["uint32"] = uint32(1)
	m["uint64"] = uint64(1)

	fmt.Println(m.StrVal("abc"))
	fmt.Println(m.StrVal("abc2"))
	fmt.Println(m.StrVal("nf"))
	fmt.Println(m.StrVal("int"))
	//
	fmt.Println(m.IntVal("int"))
	fmt.Println(m.IntVal("int8"))
	fmt.Println(m.IntVal("int16"))
	fmt.Println(m.IntVal("int32"))
	fmt.Println(m.IntVal("int64"))
	fmt.Println(m.IntVal("uint64"))
	fmt.Println(m.IntVal("nf"))
	fmt.Println(m.IntVal("abc"))
	fmt.Println(m.IntVal("float32"))
	fmt.Println(m.IntVal("uint64"))
	//
	fmt.Println(m.UintVal("uint"))
	fmt.Println(m.UintVal("uint8"))
	fmt.Println(m.UintVal("uint16"))
	fmt.Println(m.UintVal("uint32"))
	fmt.Println(m.UintVal("uint64"))
	fmt.Println(m.UintVal("float64"))
	fmt.Println(m.UintVal("nf"))
	fmt.Println(m.UintVal("abc"))
	fmt.Println(m.UintVal("float32"))
	fmt.Println(m.UintVal("int64"))
	//
	fmt.Println(m.FloatVal("float32"))
	fmt.Println(m.FloatVal("float64"))
	fmt.Println(m.FloatVal("int64"))
	fmt.Println(m.FloatVal("nf"))
	fmt.Println(m.FloatVal("abc"))
	fmt.Println(m.FloatVal("int64"))
	fmt.Println(m.FloatVal("uint64"))
	//
}

func TestC(t *testing.T) {
	fv := math.MaxFloat64
	iv := math.MaxInt64
	var uv uint64 = math.MaxUint64
	fmt.Println(uint64(fv))
	fmt.Println(uint64(iv))
	fmt.Println(int64(uv))
	fmt.Println(int64(fv))
	fmt.Println(float64(uv))
	fmt.Println(float64(iv))
	// fmt.Println(int64(fv))
	// fmt.Println(int64(math.MaxFloat64 / 2e8))
	// fmt.Println(int64(math.MaxUint64 / 2))
}
