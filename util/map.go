package util

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type Validable interface {
	ValidF(f string, args ...interface{}) error
}

//type define to map[string]interface{}
type Map map[string]interface{}

func IntVal(v interface{}) int64 {
	val, err := IntValV(v)
	if err == nil {
		return val
	} else {
		return 0
	}
}

func IntValV(v interface{}) (int64, error) {
	if v == nil {
		return 0, Err("arg value is null")
	}
	k := reflect.TypeOf(v)
	if k.Name() == "Time" {
		t := v.(time.Time)
		return Timestamp(t), nil
	}
	switch k.Kind() {
	case reflect.Int:
		return int64(v.(int)), nil
	case reflect.Int8:
		return int64(v.(int8)), nil
	case reflect.Int16:
		return int64(v.(int16)), nil
	case reflect.Int32:
		return int64(v.(int32)), nil
	case reflect.Int64:
		return v.(int64), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return int64(UintVal(v)), nil
	case reflect.Float32, reflect.Float64:
		return int64(FloatVal(v)), nil
	case reflect.String:
		if fv, err := strconv.ParseInt(v.(string), 10, 64); err == nil {
			return fv, nil
		} else {
			return 0, err
		}
	case reflect.Struct:
		if k.Name() == "Time" {
			return Timestamp(v.(time.Time)), nil
		} else {
			return 0, Err("incompactable kind(%v)", k.Kind())
		}
	default:
		return 0, Err("incompactable kind(%v)", k.Kind())
	}
}

func UintVal(v interface{}) uint64 {
	val, err := UintValV(v)
	if err == nil {
		return val
	} else {
		return 0
	}
}

func UintValV(v interface{}) (uint64, error) {
	if v == nil {
		return 0, Err("arg value is null")
	}
	k := reflect.TypeOf(v)
	switch k.Kind() {
	case reflect.Uint:
		return uint64(v.(uint)), nil
	case reflect.Uint8:
		return uint64(v.(uint8)), nil
	case reflect.Uint16:
		return uint64(v.(uint16)), nil
	case reflect.Uint32:
		return uint64(v.(uint32)), nil
	case reflect.Uint64:
		return v.(uint64), nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return uint64(IntVal(v)), nil
	case reflect.Float32, reflect.Float64:
		return uint64(FloatVal(v)), nil
	case reflect.String:
		if fv, err := strconv.ParseUint(v.(string), 10, 64); err == nil {
			return fv, nil
		} else {
			return 0, err
		}
	default:
		return 0, Err("incompactable kind(%v)", k.Kind().String())
	}
}

func FloatVal(v interface{}) float64 {
	val, err := FloatValV(v)
	if err == nil {
		return val
	} else {
		return 0
	}
}

func FloatValV(v interface{}) (float64, error) {
	if v == nil {
		return 0, Err("arg value is null")
	}
	k := reflect.TypeOf(v)
	switch k.Kind() {
	case reflect.Float32:
		return float64(v.(float32)), nil
	case reflect.Float64:
		return float64(v.(float64)), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return float64(UintVal(v)), nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return float64(IntVal(v)), nil
	case reflect.String:
		if fv, err := strconv.ParseFloat(v.(string), 64); err == nil {
			return fv, nil
		} else {
			return 0, err
		}
	default:
		return 0, Err("incompactable kind(%v)", k.Kind().String())
	}
}
func StrVal(v interface{}) string {
	if v == nil {
		return ""
	}
	switch reflect.TypeOf(v).Kind() {
	case reflect.String:
		return v.(string)
	default:
		return fmt.Sprintf("%v", v)
	}
}

func StrVal2(v interface{}) string {
	if v == nil {
		return ""
	}
	switch reflect.TypeOf(v).Kind() {
	case reflect.String:
		return v.(string)
	case reflect.Slice:
		vals := reflect.ValueOf(v)
		var vs = []string{}
		for i := 0; i < vals.Len(); i++ {
			vs = append(vs, fmt.Sprintf("%v", vals.Index(i).Interface()))
		}
		return strings.Join(vs, ",")
	default:
		return fmt.Sprintf("%v", v)
	}
}

func MapVal(v interface{}) Map {
	if mv, ok := v.(Map); ok {
		return mv
	} else if mv, ok := v.(map[string]interface{}); ok {
		return Map(mv)
	} else {
		return nil
	}
}

func AryVal(v interface{}) []interface{} {
	if vals, ok := v.([]interface{}); ok {
		return vals
	}
	vals := reflect.ValueOf(v)
	if vals.Kind() != reflect.Slice {
		return nil
	}
	var vs = []interface{}{}
	for i := 0; i < vals.Len(); i++ {
		vs = append(vs, vals.Index(i).Interface())
	}
	return vs
}

func AryMapVal(v interface{}) []Map {
	var vals = AryVal(v)
	if vals == nil {
		return nil
	}
	var ms = []Map{}
	for _, val := range vals {
		var mv = MapVal(val)
		if mv == nil {
			return nil
		} else {
			ms = append(ms, mv)
		}
	}
	return ms
}

func AryStrVal(v interface{}) []string {
	var vals = AryVal(v)
	if vals == nil {
		return nil
	}
	var ms = []string{}
	for _, val := range vals {
		ms = append(ms, StrVal(val))
	}
	return ms
}

func AryIntVal(v interface{}) []int {
	as := AryVal(v)
	if as == nil {
		return nil
	}
	is := []int{}
	for _, v := range as {
		iv, err := IntValV(v)
		if err != nil {
			return nil
		}
		is = append(is, int(iv))
	}
	return is
}

func (m Map) UintVal(key string) uint64 {
	if v, ok := m[key]; ok {
		return UintVal(v)
	} else {
		return 0
	}
}
func (m Map) IntVal(key string) int64 {
	if v, ok := m[key]; ok {
		return IntVal(v)
	} else {
		return 0
	}
}
func (m Map) IntValV(key string, d int64) int64 {
	if v, ok := m[key]; ok {
		val, err := IntValV(v)
		if err == nil {
			return val
		} else {
			return d
		}
	} else {
		return d
	}

}
func (m Map) FloatVal(key string) float64 {
	if v, ok := m[key]; ok {
		return FloatVal(v)
	} else {
		return 0
	}
}
func (m Map) StrVal(key string) string {
	if v, ok := m[key]; ok {
		return StrVal(v)
	} else {
		return ""
	}
}
func (m Map) StrValV(key, d string) string {
	if v, ok := m[key]; ok {
		val := StrVal(v)
		if len(val) < 0 {
			return d
		} else {
			return val
		}
	} else {
		return ""
	}
}
func (m Map) StrVal2(key string) string {
	if v, ok := m[key]; ok {
		return StrVal2(v)
	} else {
		return ""
	}
}
func (m Map) MapVal(key string) Map {
	if v, ok := m[key]; ok {
		return MapVal(v)
	} else {
		return nil
	}
}
func (m Map) AryVal(key string) []interface{} {
	if v, ok := m[key]; ok {
		return AryVal(v)
	} else {
		return nil
	}
}
func (m Map) AryMapVal(key string) []Map {
	if v, ok := m[key]; ok {
		return AryMapVal(v)
	} else {
		return nil
	}
}
func (m Map) AryStrVal(key string) []string {
	if v, ok := m[key]; ok {
		return AryStrVal(v)
	} else {
		return nil
	}
}
func (m Map) AryIntVal(key string) []int {
	if v, ok := m[key]; ok {
		return AryIntVal(v)
	} else {
		return nil
	}
}
func (m Map) Val(key string) interface{} {
	if v, ok := m[key]; ok {
		return v
	} else {
		return nil
	}
}
func (m Map) SetVal(key string, val interface{}) {
	if val == nil {
		delete(m, key)
	} else {
		m[key] = val
	}
}

func (m Map) UintValP(path string) uint64 {
	v, _ := m.ValP(path)
	return UintVal(v)
}
func (m Map) IntValP(path string) int64 {
	v, _ := m.ValP(path)
	return IntVal(v)
}
func (m Map) IntValPv(path string, d int64) int64 {
	v, _ := m.ValP(path)
	val, err := IntValV(v)
	if err == nil {
		return val
	} else {
		return d
	}
}
func (m Map) FloatValP(path string) float64 {
	v, _ := m.ValP(path)
	return FloatVal(v)
}
func (m Map) StrValP(path string) string {
	v, _ := m.ValP(path)
	return StrVal(v)
}
func (m Map) StrValPv(path, d string) string {
	v, _ := m.ValP(path)
	val := StrVal(v)
	if len(val) < 1 {
		return d
	} else {
		return val
	}
}
func (m Map) StrValP2(path string) string {
	v, _ := m.ValP(path)
	return StrVal2(v)
}
func (m Map) MapValP(path string) Map {
	v, _ := m.ValP(path)
	return MapVal(v)
}
func (m Map) AryValP(path string) []interface{} {
	v, _ := m.ValP(path)
	return AryVal(v)
}
func (m Map) AryMapValP(path string) []Map {
	v, _ := m.ValP(path)
	return AryMapVal(v)
}
func (m Map) AryStrValP(path string) []string {
	v, _ := m.ValP(path)
	return AryStrVal(v)
}
func (m Map) AryIntValP(path string) []int {
	v, _ := m.ValP(path)
	return AryIntVal(v)
}
func (m Map) ValP(path string) (interface{}, error) {
	path = strings.TrimPrefix(path, "/")
	keys := strings.Split(path, "/")
	return m.valP(keys)
}
func (m Map) valP(keys []string) (interface{}, error) {
	count := len(keys)
	var tv interface{} = m
	for i := 0; i < count; i++ {
		if tv == nil {
			break
		}
		switch reflect.TypeOf(tv).Kind() {
		case reflect.Slice: //if array
			ary, ok := tv.([]interface{}) //check if valid array
			if !ok {
				return nil, errors.New(fmt.Sprintf(
					"invalid array(%v) in path(/%v),expected []interface{}",
					reflect.TypeOf(tv).String(), strings.Join(keys[:i+1], "/"),
				))
			}
			if keys[i] == "@len" { //check if have @len
				return len(ary), nil //return the array length
			}
			idx, err := strconv.Atoi(keys[i]) //get the target index.
			if err != nil {
				return nil, errors.New(fmt.Sprintf(
					"invalid array index(/%v)", strings.Join(keys[:i+1], "/"),
				))
			}
			if idx >= len(ary) || idx < 0 { //check index valid
				return nil, errors.New(fmt.Sprintf(
					"array out of index in path(/%v)", strings.Join(keys[:i+1], "/"),
				))
			}
			tv = ary[idx]
			continue
		case reflect.Map: //if map
			tm := MapVal(tv) //check map covert
			if tm == nil {
				return nil, errors.New(fmt.Sprintf(
					"invalid map in path(/%v)", strings.Join(keys[:i], "/"),
				))
			}
			tv = tm.Val(keys[i])
			continue
		default: //unknow type
			return nil, errors.New(fmt.Sprintf(
				"invalid type(%v) in path(/%v)",
				reflect.TypeOf(tv).Kind(), strings.Join(keys[:i], "/"),
			))
		}
	}
	if tv == nil { //if valud not found
		return nil, errors.New(fmt.Sprintf(
			"value not found in path(/%v)", strings.Join(keys, "/"),
		))
	} else {
		return tv, nil
	}
}

func (m Map) SetValP(path string, val interface{}) error {
	if len(path) < 1 {
		return errors.New("path is empty")
	}
	path = strings.TrimPrefix(path, "/")
	keys := strings.Split(path, "/")
	//
	i := len(keys) - 1
	pv, err := m.valP(keys[:i])
	if err != nil {
		return err
	}
	switch reflect.TypeOf(pv).Kind() {
	case reflect.Slice:
		ary, ok := pv.([]interface{}) //check if valid array
		if !ok {
			return errors.New(fmt.Sprintf(
				"invalid array(%v) in path(/%v),expected []interface{}",
				reflect.TypeOf(pv).String(), strings.Join(keys[:i+1], "/"),
			))
		}
		idx, err := strconv.Atoi(keys[i]) //get the target index.
		if err != nil {
			return errors.New(fmt.Sprintf(
				"invalid array index(/%v)", strings.Join(keys[:i+1], "/"),
			))
		}
		if idx >= len(ary) || idx < 0 { //check index valid
			return errors.New(fmt.Sprintf(
				"array out of index in path(/%v)", strings.Join(keys[:i+1], "/"),
			))
		}
		ary[idx] = val
	case reflect.Map:
		tm := MapVal(pv) //check map covert
		if tm == nil {
			return errors.New(fmt.Sprintf(
				"invalid map in path(/%v)", strings.Join(keys[:i], "/"),
			))
		}
		tm.SetVal(keys[i], val)
	default: //unknow type
		return errors.New(fmt.Sprintf(
			"not map type(%v) in path(/%v)",
			reflect.TypeOf(pv).Kind(), strings.Join(keys[:i], "/"),
		))
	}
	return nil
}

func (m Map) ValidF(f string, args ...interface{}) error {
	return ValidAttrF(f, m.StrValP2, true, args...)
}
func (m Map) ToS(dest interface{}) {
	M2S(m, dest)
}
func (m Map) Exist(key string) bool {
	_, ok := m[key]
	return ok
}
func NewMap(f string) (Map, error) {
	bys, err := ioutil.ReadFile(f)
	if err != nil {
		return nil, err
	}
	var kvs Map = Map{}
	err = json.Unmarshal(bys, &kvs)
	return kvs, err
}

func NewMaps(f string) ([]Map, error) {
	bys, err := ioutil.ReadFile(f)
	if err != nil {
		return nil, err
	}
	var kvs []Map = []Map{}
	err = json.Unmarshal(bys, &kvs)
	return kvs, err
}

type MapSorter struct {
	Maps []Map
	Key  string
	Type int //0 is int,1 is float,2 is string
}

func (m *MapSorter) Len() int {
	return len(m.Maps)
}

func (m *MapSorter) Less(i, j int) bool {
	switch m.Type {
	case 0:
		return m.Maps[i].IntValP(m.Key) < m.Maps[j].IntValP(m.Key)
	case 1:
		return m.Maps[i].FloatValP(m.Key) < m.Maps[j].FloatValP(m.Key)
	default:
		return m.Maps[i].StrValP(m.Key) < m.Maps[j].StrValP(m.Key)
	}
}
func (m *MapSorter) Swap(i, j int) {
	m.Maps[i], m.Maps[j] = m.Maps[j], m.Maps[i]
}
