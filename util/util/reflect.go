package util

import (
	"fmt"
	"reflect"
	"runtime"
	"strings"
)

//SliceExists check if objext exists on slice, if ary is not slice return false
func SliceExists(ary interface{}, obj interface{}) bool {
	switch reflect.TypeOf(ary).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(ary)
		for i := 0; i < s.Len(); i++ {
			if obj == s.Index(i).Interface() {
				return true
			}
		}
		return false
	default:
		return false
	}
}

//IsType check if value is target type
func IsType(v interface{}, t string) bool {
	t = strings.TrimSpace(t)
	if v == nil || len(t) < 1 {
		return false
	}
	return reflect.Indirect(reflect.ValueOf(v)).Type().Name() == t
}

//Join join the slice to string by seqerated
func Join(v interface{}, seq string) string {
	if v == nil {
		return ""
	}
	vtype := reflect.TypeOf(v)
	if vtype.Kind() != reflect.Slice {
		return ""
	}
	vval := reflect.ValueOf(v)
	if vval.Len() < 1 {
		return ""
	}
	val := fmt.Sprintf("%v", vval.Index(0).Interface())
	for i := 1; i < vval.Len(); i++ {
		val += fmt.Sprintf("%v%v", seq, vval.Index(i).Interface())
	}
	return val
}

//CallStatck return the call stack string by current line
func CallStatck() string {
	buf := make([]byte, 102400)
	blen := runtime.Stack(buf, false)
	return string(buf[0:blen])
}

//StructName return the struct type name by reflect
func StructName(v interface{}) string {
	return reflect.Indirect(reflect.ValueOf(v)).Type().String()
}

//FuncName return the func type name by reflect
func FuncName(v interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(v).Pointer()).Name()
}

//ReflectName return the value type name by reflect
func ReflectName(v interface{}) string {
	var val = reflect.Indirect(reflect.ValueOf(v))
	if val.Kind() == reflect.Func {
		return runtime.FuncForPC(val.Pointer()).Name()
	}
	return val.Type().String()
}
