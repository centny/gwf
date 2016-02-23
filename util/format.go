package util

import (
	"bufio"
	"bytes"
	"fmt"
	"regexp"
	"strings"
)

func ParseSectionF(pre, sub, data string) Map {
	var pb_pre = regexp.MustCompile(fmt.Sprintf("^\\%v[^\\%v/]*\\%v$", pre, pre, sub))
	var pb_sub = regexp.MustCompile(fmt.Sprintf("^\\%v/[^\\%v]*\\%v$", pre, pre, sub))
	var key_pre, key_sub = "", ""
	var res = Map{}
	var buf = bytes.NewBuffer(nil)
	var bdata = bufio.NewReader(bytes.NewBufferString(data))
	for {
		var line_, err = ReadLine(bdata, 1024, false)
		if err != nil {
			break
		}
		var line = strings.Trim(string(line_), " \t\n")
		if len(line) < 1 {
			continue
		}
		if pb_pre.MatchString(line) {
			buf = bytes.NewBuffer(nil)
			key_pre = strings.Trim(line, fmt.Sprintf(" %v%v \t", pre, sub))
		} else if pb_sub.MatchString(line) {
			key_sub = strings.Trim(line, fmt.Sprintf(" %v%v/ \t", pre, sub))
			var data = strings.Trim(buf.String(), " \t\n")
			if key_pre == key_sub {
				res.SetVal(key_pre, data)
			} else {
				fmt.Printf("invalid section by prefix(%v),subffix(%v),data(%v)", key_pre, key_sub, data)
			}
			key_pre, key_sub = "", ""
			buf = bytes.NewBuffer(nil)
		} else {
			buf.WriteString(line + "\n")
		}
	}
	return res
}
