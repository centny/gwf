package util

import (
	"os"
	"strings"
)

func Args() (string, Map, []string) {
	return ArgsV(os.Args)
}

func ArgsV(args []string) (string, Map, []string) {
	alen := len(args)
	args1 := Map{}
	args2 := []string{}
	var arg_k, arg_v string
	for i := 1; i < alen; {
		arg_k = args[i]
		if i < alen-1 {
			arg_v = args[i+1]
		} else {
			arg_v = ""
		}
		if strings.HasPrefix(arg_k, "-") {
			if strings.HasPrefix(arg_v, "-") {
				args1[strings.TrimPrefix(arg_k, "-")] = ""
				i += 1
			} else {
				args1[strings.TrimPrefix(arg_k, "-")] = arg_v
				i += 2
			}
		} else {
			args2 = append(args2, arg_k)
			i += 1
		}
	}
	return args[0], args1, args2
}
