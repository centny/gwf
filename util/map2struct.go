package util

import (
	"fmt"
	"math"
	"os"
	"reflect"
	"strings"
	"time"
)

//type define to map[string]interface{}
type Map map[string]interface{}

func (m Map) UintVal(key string) uint64 {
	if v, ok := m[key]; ok {
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
			return uint64(m.IntVal(key))
		case reflect.Float32, reflect.Float64:
			return uint64(m.FloatVal(key))
		default:
			return math.MaxUint64
		}
	} else {
		return math.MaxUint64
	}
}
func (m Map) IntVal(key string) int64 {
	if v, ok := m[key]; ok {
		switch reflect.TypeOf(v).Kind() {
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
			return int64(m.UintVal(key))
		case reflect.Float32, reflect.Float64:
			return int64(m.FloatVal(key))
		default:
			return math.MaxInt64
		}
	} else {
		return math.MaxInt64
	}
}
func (m Map) FloatVal(key string) float64 {
	if v, ok := m[key]; ok {
		switch reflect.TypeOf(v).Kind() {
		case reflect.Float32:
			return float64(v.(float32))
		case reflect.Float64:
			return float64(v.(float64))
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			return float64(m.UintVal(key))
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return float64(m.IntVal(key))
		default:
			return math.MaxFloat64
		}
	} else {
		return math.MaxFloat64
	}
}
func (m Map) StrVal(key string) string {
	if v, ok := m[key]; ok {
		switch reflect.TypeOf(v).Kind() {
		case reflect.String:
			return v.(string)
		default:
			return fmt.Sprintf("%v", v)
		}
	} else {
		return ""
	}
}

//default date format.
const D_DATEFORMAT string = "2006-01-02 15:04:05"

//map to struct.
func M2S(m Map, dest interface{}) {
	if m == nil || dest == nil || len(m) < 1 {
		return
	}
	//get the reflect type.
	ptype := reflect.TypeOf(reflect.Indirect(reflect.ValueOf(dest)).Interface())
	//get the reflect value.
	pval := reflect.ValueOf(dest).Elem()
	for i := 0; i < ptype.NumField(); i++ {
		f := ptype.Field(i)
		var m2s string = f.Tag.Get("m2s") //get the m2s tag.
		if len(m2s) < 1 {                 //if not m2s tag,using field name.
			m2s = f.Name
		}
		keys := strings.Split(m2s, ",")
		for _, key := range keys {
			v, ok := m[key]
			if !ok || v == nil {
				continue
			}
			vty := reflect.TypeOf(v)
			if f.Type.Kind() == vty.Kind() {
				pval.Field(i).Set(reflect.ValueOf(v))
				continue
			}
			if f.Type.Name() == "Time" {
				switch vty.Name() {
				case "string":
					df := f.Tag.Get("tf")
					if len(df) < 1 {
						df = D_DATEFORMAT
					}
					t, err := time.Parse(df, v.(string))
					if err == nil {
						pval.Field(i).Set(reflect.ValueOf(t))
					} else {
						fmt.Fprintln(os.Stderr, err.Error())
					}
				default:
					iv := m.IntVal(key)
					if iv < math.MaxInt64 {

					}
					pval.Field(i).Set(reflect.ValueOf(Time(iv)))
				}
				continue
			}
			iv := m.IntVal(key)
			uv := m.UintVal(key)
			fv := m.FloatVal(key)
			var val reflect.Value
			if iv < math.MaxInt64 {
				switch f.Type.Kind() {
				case reflect.Int:
					val = reflect.ValueOf(int(iv))
				case reflect.Int8:
					val = reflect.ValueOf(int8(iv))
				case reflect.Int16:
					val = reflect.ValueOf(int16(iv))
				case reflect.Int32:
					val = reflect.ValueOf(int32(iv))
				case reflect.Int64:
					val = reflect.ValueOf(int64(iv))
				}
			}
			if uv < math.MaxUint64 {
				switch f.Type.Kind() {
				case reflect.Uint:
					val = reflect.ValueOf(uint(uv))
				case reflect.Uint8:
					val = reflect.ValueOf(uint8(uv))
				case reflect.Uint16:
					val = reflect.ValueOf(uint16(uv))
				case reflect.Uint32:
					val = reflect.ValueOf(uint32(uv))
				case reflect.Uint64:
					val = reflect.ValueOf(uint64(uv))
				}
			}
			if fv < math.MaxFloat64 {
				switch f.Type.Kind() {
				case reflect.Float32:
					val = reflect.ValueOf(float32(fv))
				case reflect.Float64:
					val = reflect.ValueOf(float64(fv))
				}
			}
			if val.IsValid() {
				pval.Field(i).Set(val)
			}
		}
	}
}

//map array to struct array.
func Ms2Ss(ms interface{}, dest interface{}) {
	if ms == nil || dest == nil {
		return
	}
	if reflect.TypeOf(ms).Kind() != reflect.Slice {
		return
	}
	//get the reflect value.
	pval := reflect.Indirect(reflect.ValueOf(dest))
	rval := reflect.Indirect(reflect.ValueOf(dest))
	//get the reflect type.
	ptype := reflect.TypeOf(rval.Interface()).Elem()
	mss := reflect.ValueOf(ms)
	for i := 0; i < mss.Len(); i++ {
		var mv Map
		msv := mss.Index(i).Interface()
		switch reflect.TypeOf(msv).Name() {
		case "Map":
			mv = msv.(Map)
		default:
			mv = Map(msv.(map[string]interface{}))
		}
		if len(mv) < 1 {
			continue
		}
		pv := reflect.New(ptype)
		rv := reflect.Indirect(pv)
		M2S(mv, pv.Interface())
		pval = reflect.Append(pval, rv)
	}
	//reset the slice address to new.
	rval.Set(pval)
}
