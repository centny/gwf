package util

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"reflect"
	"strings"
	"time"
)

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
	pval := reflect.ValueOf(dest)
	if pval.Kind() != reflect.Ptr {
		panic("dest values is not a ptr")
	}
	pval = pval.Elem()
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
						pval.Field(i).Set(reflect.ValueOf(Time(iv)))
					}
				}
				continue
			} else if f.Type.Name() == "string" && vty.Name() == "Time" {
				df := f.Tag.Get("tf")
				if len(df) < 1 {
					df = D_DATEFORMAT
				}
				pval.Field(i).Set(reflect.ValueOf(v.(time.Time).Format(df)))
				continue
			}
			it := f.Tag.Get("it")
			var iv int64
			var uv uint64
			var fv float64
			if it == "Y" && vty.Name() == "string" {
				df := f.Tag.Get("tf")
				if len(df) < 1 {
					df = D_DATEFORMAT
				}
				t, err := time.Parse(df, v.(string))
				if err != nil {
					fmt.Fprintln(os.Stderr, err.Error())
					continue
				}
				ts := Timestamp(t)
				iv = ts
				uv = uint64(ts)
				fv = float64(ts)
			} else {
				iv = m.IntVal(key)
				uv = m.UintVal(key)
				fv = m.FloatVal(key)
			}
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
		fmt.Println("ms or dest is nil")
		return
	}
	if reflect.TypeOf(ms).Kind() != reflect.Slice {
		fmt.Println("not slice")
		return
	}
	//get the reflect value.
	pval := reflect.Indirect(reflect.ValueOf(dest))
	rval := reflect.Indirect(reflect.ValueOf(dest))
	//get the reflect type.
	ptype := reflect.TypeOf(rval.Interface()).Elem()
	isptr := ptype.Kind() == reflect.Ptr
	if isptr {
		ptype = ptype.Elem()
	}
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
		if isptr {
			pval = reflect.Append(pval, rv.Addr())
		} else {
			pval = reflect.Append(pval, rv)
		}
	}
	//reset the slice address to new.
	rval.Set(pval)
}

func Json2S(data string, dest interface{}) error {
	m, err := Json2Map(data)
	if err != nil {
		return err
	}
	M2S(m, dest)
	return nil
}

func Json2Ss(data string, dest interface{}) error {
	ms, err := Json2Ary(data)
	if err != nil {
		return err
	}
	Ms2Ss(ms, dest)
	return nil
}
func S2Json(v interface{}) string {
	bys, _ := json.Marshal(v)
	return string(bys)
}
