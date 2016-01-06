package util

import (
	"errors"
	"fmt"
	"math"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

//valid the data to specified value by limit
//data:target value string
//valLT:target value value type limit,example R|F(required float value)
//	R:required,O:option
//	S:string value,I:integet value,F:float value
//valRT:taget value value range limit,
//example O:1-2-3-4(option value in 1-2-3-4),P:^.*\@.*$(match @)
//	O:option,R:range,P:pattern
//	seperate value by -
//limit_r:if return error when require value not found.
func ValidAttrT(data string, valLT string, valLR string, limit_r bool) (interface{}, error) {
	valLR = strings.Replace(valLR, "%N", ",", -1)
	valLR = strings.Replace(valLR, "%%", "%", -1)
	lrs := strings.SplitN(valLR, ":", 2) //valid value range.
	lts := strings.SplitN(valLT, "|", 2) //valid required type
	if len(lrs) < 2 {
		return nil, errors.New(fmt.Sprintf("invalid range limit:%s", valLR))
	}
	if len(lts) < 2 {
		return nil, errors.New(fmt.Sprintf("invalid type limit:%s", valLT))
	}
	if len(data) < 1 { //chekc the value required.
		if lts[0] == "R" && limit_r {
			return nil, errors.New("data is empty")
		} else {
			return nil, nil
		}
	}
	//define the valid string function.
	validStr := func(ds string) (interface{}, error) {
		//check range limit.
		switch lrs[0] {
		case "O": //option limit.
			options := strings.Split(lrs[1], "~")
			if AryExist(options, ds) {
				return ds, nil
			} else {
				return nil, errors.New(fmt.Sprintf("invalid value(%s) for options(%s)", ds, lrs[1]))
			}
		case "L": //length limit
			slen := int64(len(ds))
			rgs := strings.Split(lrs[1], "~")
			var beg, end int64 = 0, 0
			var err error = nil
			if len(rgs) > 0 && len(rgs[0]) > 0 {
				beg, err = strconv.ParseInt(rgs[0], 10, 64)
				if err != nil {
					return nil, errors.New(fmt.Sprintf("invalid range begin number(%s)", rgs[0]))
				}
			} else {
				beg = 0
			}
			if len(rgs) > 1 && len(rgs[1]) > 0 {
				end, err = strconv.ParseInt(rgs[1], 10, 64)
				if err != nil {
					return nil, errors.New(fmt.Sprintf("invalid range end number option(%s)", rgs[1]))
				}
			} else {
				end = math.MaxInt64
			}
			if beg < slen && end > slen {
				return ds, nil
			} else {
				return nil, errors.New(fmt.Sprintf("string length must match %d<len<%d, but %d", beg, end, slen))
			}
		case "P": //regex pattern limit
			mched, err := regexp.MatchString(lrs[1], ds)
			if err != nil {
				return nil, err
			}
			if mched {
				return ds, nil
			} else {
				return nil, errors.New(fmt.Sprintf("value(%s) not match regex(%s)", ds, lrs[1]))
			}
		}
		//unknow range limit type.
		return nil, errors.New(fmt.Sprintf("invalid range limit %s for string", lrs[0]))
	}
	//define valid number function.
	validNum := func(ds float64) (interface{}, error) {
		//check range limit.
		switch lrs[0] {
		case "R":
			var beg, end float64 = 0, 0
			var err error = nil
			rgs := strings.Split(lrs[1], "~")
			if len(rgs) > 0 && len(rgs[0]) > 0 {
				beg, err = strconv.ParseFloat(rgs[0], 64)
				if err != nil {
					return nil, errors.New(fmt.Sprintf("invalid range begin number(%s)", rgs[0]))
				}
			} else {
				beg = 0
			}
			if len(rgs) > 1 && len(rgs[1]) > 0 {
				end, err = strconv.ParseFloat(rgs[1], 64)
				if err != nil {
					return nil, errors.New(fmt.Sprintf("invalid range end number option(%s)", rgs[1]))
				}
			} else {
				end = math.MaxFloat64
			}
			if beg < ds && end > ds {
				return ds, nil
			} else {
				return nil, errors.New(fmt.Sprintf("value must match %f<val<%f, but %v", beg, end, ds))
			}
		case "O":
			options := strings.Split(lrs[1], "~")
			var oary []float64
			for _, o := range options { //covert to float array.
				v, err := strconv.ParseFloat(o, 64)
				if err != nil {
					return nil, errors.New(fmt.Sprintf("invalid number option(%s)", lrs[1]))
				}
				oary = append(oary, v)
			}
			if AryExist(oary, ds) {
				return ds, nil
			} else {
				return nil, errors.New(fmt.Sprintf("invalid value(%f) for options(%s)", ds, lrs[1]))
			}
		}
		//unknow range limit type.
		return nil, errors.New(fmt.Sprintf("invalid range limit %s for float", lrs[0]))
	}
	//define valid number function.
	validInt := func(ds int64) (interface{}, error) {
		//check range limit.
		switch lrs[0] {
		case "R":
			var beg, end int64 = 0, 0
			var err error = nil
			rgs := strings.Split(lrs[1], "~")
			if len(rgs) > 0 && len(rgs[0]) > 0 {
				beg, err = strconv.ParseInt(rgs[0], 10, 64)
				if err != nil {
					return nil, errors.New(fmt.Sprintf("invalid range begin number(%s)", rgs[0]))
				}
			} else {
				beg = 0
			}
			if len(rgs) > 1 && len(rgs[1]) > 0 {
				end, err = strconv.ParseInt(rgs[1], 10, 64)
				if err != nil {
					return nil, errors.New(fmt.Sprintf("invalid range end number option(%s)", rgs[1]))
				}
			} else {
				end = math.MaxInt64
			}
			if beg < ds && end > ds {
				return ds, nil
			} else {
				return nil, errors.New(fmt.Sprintf("value must match %v<val<%v, but %v", beg, end, ds))
			}
		case "O":
			options := strings.Split(lrs[1], "~")
			var oary []int64
			for _, o := range options { //covert to float array.
				v, err := strconv.ParseInt(o, 10, 64)
				if err != nil {
					return nil, errors.New(fmt.Sprintf("invalid number option(%s)", lrs[1]))
				}
				oary = append(oary, v)
			}
			if AryExist(oary, ds) {
				return ds, nil
			} else {
				return nil, errors.New(fmt.Sprintf("invalid value(%v) for options(%s)", ds, lrs[1]))
			}
		}
		//unknow range limit type.
		return nil, errors.New(fmt.Sprintf("invalid range limit %s for float", lrs[0]))
	}
	//define value type function
	validLts := func(ds string) (interface{}, error) {
		switch lts[1] {
		case "S":
			return validStr(ds)
		case "I":
			ids, err := strconv.ParseFloat(ds, 64)
			if err != nil {
				return nil, errors.New(fmt.Sprintf("invalid value(%s) for type(%s):%v", ds, lts[1], err))
			} else {
				return validInt(int64(ids))
			}
		case "F":
			fds, err := strconv.ParseFloat(ds, 64)
			if err != nil {
				return nil, errors.New(fmt.Sprintf("invalid value(%s) for type(%s):%v", ds, lts[1], err))
			} else {
				return validNum(fds)
			}
		}
		return nil, errors.New(fmt.Sprintf("invalid value type:%s", lts[1]))
	}
	return validLts(data)
}

type AttrFunc func(key string) string

func ValidAttrF(f string, cf AttrFunc, limit_r bool, args ...interface{}) error {
	f = regexp.MustCompile("\\/\\/.*").ReplaceAllString(f, "")
	f = strings.Replace(f, "\n", "", -1)
	f = strings.Trim(f, " \t;")
	if len(f) < 1 {
		return errors.New("format not found")
	}
	trimfs := strings.Split(f, ";")
	if len(trimfs) != len(args) {
		return errors.New("args count is not equal format count")
	}
	for idx, fs := range trimfs {
		fs = strings.Trim(fs, " \t")
		fstr := strings.SplitN(fs, ",", 4)
		if len(fstr) < 3 {
			return errors.New(fmt.Sprintf("format error:%s", fs))
		}
		rval, err := ValidAttrT(cf(fstr[0]), fstr[1], fstr[2], limit_r)
		if err != nil {
			rerr := errors.New(fmt.Sprintf("limit(%s),%s", fs, err.Error()))
			if len(fstr) > 3 {
				// log.D("%s", rerr.Error())
				return errors.New(fstr[3])
			} else {
				return rerr
			}
		}
		if rval == nil {
			continue
		}
		pval := reflect.Indirect(reflect.ValueOf(args[idx]))
		err = ValidSetVal(pval, rval)
		if err != nil {
			return err
		}
	}
	return nil
}

func ValidSetVal(dst reflect.Value, src interface{}) error {
	sk := reflect.TypeOf(src)
	if sk.Kind() == dst.Kind() {
		dst.Set(reflect.ValueOf(src))
		return nil
	}
	var tiv int64
	var tfv float64
	var terr error
	switch dst.Kind() {
	case reflect.Int:
		tiv, terr = IntValV(src)
		if terr == nil {
			iv := int(tiv)
			dst.Set(reflect.ValueOf(iv))
		}
	case reflect.Int16:
		tiv, terr = IntValV(src)
		if terr == nil {
			iv := int16(tiv)
			dst.Set(reflect.ValueOf(iv))
		}
	case reflect.Int32:
		tiv, terr = IntValV(src)
		if terr == nil {
			iv := int32(tiv)
			dst.Set(reflect.ValueOf(iv))
		}
	case reflect.Int64:
		tiv, terr = IntValV(src)
		if terr == nil {
			iv := int64(tiv)
			dst.Set(reflect.ValueOf(iv))
		}
	case reflect.Uint:
		tiv, terr = IntValV(src)
		if terr == nil {
			iv := uint(tiv)
			dst.Set(reflect.ValueOf(iv))
		}
	case reflect.Uint16:
		tiv, terr = IntValV(src)
		if terr == nil {
			iv := uint16(tiv)
			dst.Set(reflect.ValueOf(iv))
		}
	case reflect.Uint32:
		tiv, terr = IntValV(src)
		if terr == nil {
			iv := uint32(tiv)
			dst.Set(reflect.ValueOf(iv))
		}
	case reflect.Uint64:
		tiv, terr = IntValV(src)
		if terr == nil {
			iv := uint64(tiv)
			dst.Set(reflect.ValueOf(iv))
		}
	case reflect.Float32:
		tfv, terr = FloatValV(src)
		if terr == nil {
			fv := float32(tfv)
			dst.Set(reflect.ValueOf(fv))
		}
	case reflect.Float64:
		tfv, terr = FloatValV(src)
		if terr == nil {
			fv := float64(tfv)
			dst.Set(reflect.ValueOf(fv))
		}
	case reflect.String:
		tsv := StrVal(src)
		dst.Set(reflect.ValueOf(tsv))
	}
	if terr == nil {
		return nil
	} else {
		return Err("parse kind(%v) value to kind(%v) value->%v", sk.Kind(), dst.Kind(), terr)
	}
}
