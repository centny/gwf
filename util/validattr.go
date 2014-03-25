package util

import (
	"errors"
	"fmt"
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
			options := strings.Split(lrs[1], "-")
			if AryExist(options, ds) {
				return ds, nil
			} else {
				return nil, errors.New(fmt.Sprintf("invalid value(%s) for options(%s)", ds, lrs[1]))
			}
		case "L": //length limit
			l_len, err := strconv.Atoi(lrs[1])
			if err != nil {
				return nil, errors.New(fmt.Sprintf("invalid len limit:%s,err:(%s)", lrs[1], err.Error()))
			}
			if len(ds) < l_len {
				return ds, nil
			} else {
				return nil, errors.New(fmt.Sprintf("value must less %d, but %d", l_len, len(ds)))
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
			rgs := strings.Split(lrs[1], "-")
			if len(rgs) < 2 {
				return nil, errors.New(fmt.Sprintf("invalid range limit:%s", lrs[1]))
			}
			beg, err := strconv.ParseFloat(rgs[0], 64)
			if err != nil {
				return nil, errors.New(fmt.Sprintf("invalid range begin number(%s)", rgs[0]))
			}
			end, err := strconv.ParseFloat(rgs[1], 64)
			if err != nil {
				return nil, errors.New(fmt.Sprintf("invalid range end number option(%s)", rgs[1]))
			}
			if beg < ds && end > ds {
				return ds, nil
			} else {
				return nil, errors.New(fmt.Sprintf("value must match %f<%f<%f", beg, ds, end))
			}
		case "O":
			options := strings.Split(lrs[1], "-")
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
	//define value type function
	validLts := func(ds string) (interface{}, error) {
		switch lts[1] {
		case "S":
			return validStr(ds)
		case "I":
			ids, err := strconv.Atoi(ds)
			if err != nil {
				return nil, errors.New(fmt.Sprintf("invalid value(%s) for type(%s)", lts[1]))
			} else {
				return validNum(float64(ids))
			}
		case "F":
			fds, err := strconv.ParseFloat(ds, 64)
			if err != nil {
				return nil, errors.New(fmt.Sprintf("invalid value(%s) for type(%s)", lts[1]))
			} else {
				return validNum(fds)
			}
		}
		return nil, errors.New(fmt.Sprintf("invalid value type:%s", lts[1]))
	}
	return validLts(data)
}
