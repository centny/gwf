package util

import (
	"fmt"
	"os"
	"reflect"
	"time"
)

//type define to map[string]interface{}
type Map map[string]interface{}

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
		var key string = f.Tag.Get("m2s") //get the m2s tag.
		if len(key) < 1 {                 //if not m2s tag,using field name.
			key = f.Name
		}
		if v, ok := m[key]; ok {
			if f.Type.Name() == "Time" {
				switch reflect.TypeOf(v).Name() {
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
				case "Time":
					pval.Field(i).Set(reflect.ValueOf(v))
				case "int64":
					pval.Field(i).Set(reflect.ValueOf(Time(v.(int64))))
				case "int":
					pval.Field(i).Set(reflect.ValueOf(Time(int64(v.(int)))))
				case "int32":
					pval.Field(i).Set(reflect.ValueOf(Time(int64(v.(int32)))))
				}

			} else {
				if v != nil {
					pval.Field(i).Set(reflect.ValueOf(v))
				}
			}
		}
	}
}

//map array to struct array.
func Ms2Ss(ms []Map, dest interface{}) {
	if ms == nil || dest == nil || len(ms) < 1 {
		return
	}
	//get the reflect value.
	pval := reflect.Indirect(reflect.ValueOf(dest))
	rval := reflect.Indirect(reflect.ValueOf(dest))
	//get the reflect type.
	ptype := reflect.TypeOf(rval.Interface()).Elem()
	for _, mv := range ms {
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
