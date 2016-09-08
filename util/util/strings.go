package util

import (
	"strings"
)

//Split split string to slice, if val is empty string return nil
func Split(val, sep string) []string {
	if len(val) < 1 {
		return nil
	}
	return strings.Split(val, sep)
}

//TrimStrsRepeat trim all string by cutset and check repeat
func TrimStrsRepeat(vals []string, cutset string, repeat bool) []string {
	var exist = map[string]bool{}
	var res = []string{}
	for _, val := range vals {
		val = strings.Trim(val, cutset)
		if len(val) < 1 {
			continue
		}
		if exist[val] && repeat {
			continue
		}
		res = append(res, val)
		exist[val] = true
	}
	return res
}

//TrimStrs trim all string by cutset
func TrimStrs(vals []string, cutset string) []string {
	return TrimStrsRepeat(vals, cutset, false)
}

//Trim trim string by space.
func Trim(s string) string {
	return strings.TrimSpace(s)
}
