package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage:pathc <-w2p|-p2w> path")
		return
	}
	switch os.Args[1] {
	case "-p2w":
		fmt.Println(p2w(os.Args[2]))
	default:
		fmt.Println(w2p(os.Args[2]))
	}
}

func w2p(path string) string {
	ns := []string{}
	ws := strings.Split(path, ";")
	for _, w := range ws {
		p := regexp.MustCompile("[a-zA-Z]\\:").ReplaceAllStringFunc(w, func(o string) string {
			return "/" + strings.Replace(o, ":", "", -1)
		})
		p = strings.Replace(p, "\\", "/", -1)
		ns = append(ns, p)
	}
	return strings.Join(ns, ":")
}

func p2w(path string) string {
	ns := []string{}
	ps := strings.Split(path, ":")
	for _, p := range ps {
		w := regexp.MustCompile("\\/[a-zA-Z]\\/").ReplaceAllStringFunc(p, func(o string) string {
			sss := strings.Replace(o, "/", "", -1) + ":\\"
			return sss
		})
		w = strings.Replace(w, "/", "\\", -1)
		ns = append(ns, w)
	}
	return strings.Join(ns, ";")
}
