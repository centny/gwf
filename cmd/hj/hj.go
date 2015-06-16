package main

import (
	"fmt"
	"github.com/Centny/gwf/util"
	"os"
	"reflect"
)

var ef func(code int) = os.Exit

func usage() {
	fmt.Println(`Usage: hj <value path> <http url format> <args>`)
}
func hget(path string, uf string, args ...interface{}) (interface{}, error) {
	res, err := util.HGet2(uf, args...)
	if err != nil {
		return nil, err
	}
	return res.ValP(path)
}

func main() {
	if len(os.Args) < 3 {
		usage()
		ef(1)
		return
	}
	args := []interface{}{}
	for i := 3; i < len(os.Args); i++ {
		args = append(args, os.Args[i])
	}
	val, err := hget(os.Args[1], os.Args[2], args...)
	if err == nil {
		switch reflect.TypeOf(val).Kind().String() {
		case "map":
			fmt.Println(util.S2Json(val))
		case "slice":
			fmt.Println(util.S2Json(val))
		default:
			fmt.Println(val)
		}
		ef(0)
	} else {
		fmt.Println(err.Error())
		ef(1)
	}
}
