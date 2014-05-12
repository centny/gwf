package util

import (
	"errors"
	"fmt"
	"math"
	"reflect"
	"strconv"
	"strings"
	"time"
)

//type define to map[string]interface{}
type Map map[string]interface{}

func IntVal(v interface{}) int64 {
	if v == nil {
		return math.MaxInt64
	}
	k := reflect.TypeOf(v)
	if k.Name() == "Time" {
		t := v.(time.Time)
		return Timestamp(t)
	}
	switch k.Kind() {
	case reflect.Int:
		return int64(v.(int))
	case reflect.Int8:
		return int64(v.(int8))
	case reflect.Int16:
		return int64(v.(int16))
	case reflect.Int32:
		return int64(v.(int32))
	case reflect.Int64:
		return v.(int64)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return int64(UintVal(v))
	case reflect.Float32, reflect.Float64:
		return int64(FloatVal(v))
	case reflect.String:
		if fv, err := strconv.ParseInt(v.(string), 10, 64); err == nil {
			return fv
		} else {
			return math.MaxInt64
		}
	default:
		return math.MaxInt64
	}
}
func UintVal(v interface{}) uint64 {
	if v == nil {
		return math.MaxUint64
	}
	switch reflect.TypeOf(v).Kind() {
	case reflect.Uint:
		return uint64(v.(uint))
	case reflect.Uint8:
		return uint64(v.(uint8))
	case reflect.Uint16:
		return uint64(v.(uint16))
	case reflect.Uint32:
		return uint64(v.(uint32))
	case reflect.Uint64:
		return v.(uint64)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return uint64(IntVal(v))
	case reflect.Float32, reflect.Float64:
		return uint64(FloatVal(v))
	case reflect.String:
		if fv, err := strconv.ParseUint(v.(string), 10, 64); err == nil {
			return fv
		} else {
			return math.MaxInt64
		}
	default:
		return math.MaxUint64
	}
}
func FloatVal(v interface{}) float64 {
	if v == nil {
		return math.MaxFloat64
	}
	switch reflect.TypeOf(v).Kind() {
	case reflect.Float32:
		return float64(v.(float32))
	case reflect.Float64:
		return float64(v.(float64))
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return float64(UintVal(v))
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return float64(IntVal(v))
	case reflect.String:
		if fv, err := strconv.ParseFloat(v.(string), 64); err == nil {
			return fv
		} else {
			return math.MaxFloat64
		}
	default:
		return math.MaxFloat64
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

func MapVal(v interface{}) Map {
	if mv, ok := v.(Map); ok {
		return mv
	} else if mv, ok := v.(map[string]interface{}); ok {
		return Map(mv)
	} else {
		return nil
	}
}

func (m Map) UintVal(key string) uint64 {
	if v, ok := m[key]; ok {
		return UintVal(v)
	} else {
		return math.MaxUint64
	}
}
func (m Map) IntVal(key string) int64 {
	if v, ok := m[key]; ok {
		return IntVal(v)
	} else {
		return math.MaxInt64
	}
}
func (m Map) FloatVal(key string) float64 {
	if v, ok := m[key]; ok {
		return FloatVal(v)
	} else {
		return math.MaxFloat64
	}
}
func (m Map) StrVal(key string) string {
	if v, ok := m[key]; ok {
		return StrVal(v)
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
func (m Map) FloatValP(path string) float64 {
	v, _ := m.ValP(path)
	return FloatVal(v)
}
func (m Map) StrValP(path string) string {
	v, _ := m.ValP(path)
	return StrVal(v)
}
func (m Map) MapValP(path string) Map {
	v, _ := m.ValP(path)
	return MapVal(v)
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
