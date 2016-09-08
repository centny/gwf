package util

import (
	"fmt"
	"strconv"
	"strings"
)

//Str2Ints parse string seperated by comma to int sclie
func Str2Ints(str string) ([]int, error) {
	return Str2IntsSeq(str, ",")
}

//Str2IntsSeq parse string seperated by seq to int sclie
func Str2IntsSeq(str, seq string) ([]int, error) {
	vals := []int{}
	parts := strings.Split(str, seq)
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if len(part) < 1 {
			continue
		}
		v, err := strconv.ParseInt(part, 10, 64)
		if err == nil {
			vals = append(vals, int(v))
		} else {
			return nil, err
		}
	}
	return vals, nil
}

//Vals2Str parsing value to string seperated by comma
func Vals2Str(vals ...interface{}) string {
	str := ""
	for _, v := range vals {
		str = fmt.Sprintf("%v%d,", str, v)
	}
	return strings.Trim(str, ",")
}

//Vals2Strs parsing value to string slice
func Vals2Strs(vals ...interface{}) []string {
	str := []string{}
	for _, v := range vals {
		str = append(str, fmt.Sprintf("%v", v))
	}
	return str
}
