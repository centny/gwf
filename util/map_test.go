package util

import (
	"fmt"
	"math"
	"reflect"
	"testing"
	"time"
)

func TestMap(t *testing.T) {
	m := Map{}
	//
	m["abc"] = "123"
	m["abd"] = "a123"
	m["abe"] = []byte("akkkk")
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
	m["time"] = time.Now()
	m.SetVal("kkkk", 123)
	m.SetVal("kkkk", nil)

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
	fmt.Println(m.IntVal("abd"))
	fmt.Println(m.IntVal("abe"))
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
	fmt.Println(m.UintVal("abd"))
	fmt.Println(m.UintVal("abe"))
	fmt.Println(m.UintVal("float32"))
	fmt.Println(m.UintVal("int64"))
	//
	fmt.Println(m.FloatVal("float32"))
	fmt.Println(m.FloatVal("float64"))
	fmt.Println(m.FloatVal("int64"))
	fmt.Println(m.FloatVal("nf"))
	fmt.Println(m.FloatVal("abc"))
	fmt.Println(m.FloatVal("abd"))
	fmt.Println(m.FloatVal("abe"))
	fmt.Println(m.FloatVal("int64"))
	fmt.Println(m.FloatVal("uint64"))
	//
	fmt.Println(m.IntVal("time"))
	m.SetVal("amap", Map{})
	m.MapVal("amap")
	m.MapVal("int64")
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
	UintVal(nil)
	IntVal(nil)
	FloatVal(nil)
	StrVal(nil)
}
func TestReflect(t *testing.T) {
	var mm map[string]interface{}
	fmt.Println(reflect.TypeOf(mm).Name())
	var m2 Map
	fmt.Println(reflect.TypeOf(m2).Name())
}

func TestGetMapP(t *testing.T) {
	//data
	m1 := map[string]interface{}{
		"s":   "str",
		"i":   int64(16),
		"f":   float64(16),
		"ary": []interface{}{1, 3, 4},
	}
	m2 := map[string]interface{}{
		"a":   "abc",
		"m":   m1,
		"ary": []interface{}{"1", "3", "4"},
	}
	m3 := map[string]interface{}{
		"b":   "abcc",
		"m":   m2,
		"ary": []interface{}{m1, m2},
	}
	m4 := Map{
		"test": 1,
		"ms":   []interface{}{m1, m2, m3},
		"m3":   m3,
		"ary2": []int{1, 3, 4},
		"me":   map[string]string{"a": "b"},
	}
	var v interface{}
	var err error
	v, err = m4.ValP("/path")
	ApE(t, v, err)
	v, err = m4.ValP("/test")
	Ap(t, v, err)
	v, err = m4.ValP("/ms")
	Ap(t, v, err)
	v, err = m4.ValP("/m3")
	Ap(t, v, err)
	//
	v, err = m4.ValP("/m3/b")
	Ap(t, v, err)
	v, err = m4.ValP("/m3/b2")
	ApE(t, v, err)
	v, err = m4.ValP("/m3/ary")
	Ap(t, v, err)
	v, err = m4.ValP("/ms/1")
	Ap(t, v, err)
	v, err = m4.ValP("/ms/100")
	ApE(t, v, err)
	v, err = m4.ValP("/ms/a")
	ApE(t, v, err)
	v, err = m4.ValP("/ary2/100")
	ApE(t, v, err)
	v, err = m4.ValP("/ms/@len")
	Ap(t, v, err)
	v, err = m4.ValP("/ary2/@len")
	ApE(t, v, err)
	v, err = m4.ValP("/test/abc")
	ApE(t, v, err)
	v, err = m4.ValP("/me/a")
	ApE(t, v, err)
	v, err = m4.ValP("/mekkkk/a")
	ApE(t, v, err)
	m4.MapVal("m3")
	m4.MapVal("m4")
	fmt.Println(m4.StrValP("/test"))
	fmt.Println(m4.UintValP("/test"))
	fmt.Println(m4.IntValP("/test"))
	fmt.Println(m4.FloatValP("/test"))
	//

}
func TestASetMapP(t *testing.T) {
	var v interface{}
	var err error
	m := Map{
		"eary":  []string{},
		"ary":   []interface{}{456},
		"emap":  map[string]string{},
		"ntype": "kkkk",
	}
	m.SetValP("/abc", Map{"a": 1})
	v, err = m.ValP("/abc/a")
	Ap(t, v, err)
	err = m.SetValP("/abcd/abc", 123)
	ApE(t, nil, err)
	err = m.SetValP("/eary/1", 123)
	ApE(t, nil, err)
	err = m.SetValP("/ary/0", 123)
	Ap(t, nil, err)
	err = m.SetValP("/ary/5", 123)
	ApE(t, nil, err)
	err = m.SetValP("/ary/a", 123)
	ApE(t, nil, err)
	err = m.SetValP("/emap/a", 123)
	ApE(t, nil, err)
	err = m.SetValP("/ntype/a", 123)
	ApE(t, nil, err)
	err = m.SetValP("", 123)
	ApE(t, nil, err)
	//
	mv := m.MapValP("/abc")
	v, err = mv.ValP("/a")
	Ap(t, v, err)

	//
}

// func TestValp(t *testing.T) {
// 	m := Map{}
// 	err := m.SetValP("/aa/val/a", "11111")
// 	if err != nil {
// 		t.Error(err.Error())
// 	}
// 	fmt.Println(m.StrValP("/aa/val/a"))
// 	m.SetValP("/aa/val/a", nil)
// 	fmt.Println(m.StrValP("/aa/val/a"))
// }
func Ap(t *testing.T, v interface{}, err error) {
	if err != nil {
		t.Error(err.Error())
		return
	}
	fmt.Println(v)
}
func ApE(t *testing.T, v interface{}, err error) {
	fmt.Println(err)
	if err == nil {
		t.Error("not error")
		return
	}
}

func TestArray2(t *testing.T) {
	fmt.Println([]int{1, 3, 5}[:3])
}
