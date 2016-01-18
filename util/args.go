package util

import (
	"bytes"
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

func ParseArgs(arg string) []string {
	var args []string
	var tbuf []byte = []byte{}
	add_a := func() {
		if len(tbuf) > 0 {
			args = append(args, string(tbuf))
			tbuf = []byte{}
		}
	}
	//
	var last byte = 0
	buf := bytes.NewBufferString(arg)
	for {
		bys, err := buf.ReadByte()
		if err != nil {
			add_a()
			break
		}
		if last == 0 {
			switch bys {
			case '\t':
				add_a()
			case ' ':
				add_a()
			case '\'':
				add_a()
				last = bys
			case '"':
				add_a()
				last = bys
			default:
				tbuf = append(tbuf, bys)
			}
		} else {
			if last == bys {
				add_a()
				last = 0
			} else {
				tbuf = append(tbuf, bys)
			}
		}

	}
	return args
}
