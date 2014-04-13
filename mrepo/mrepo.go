package main

import (
	"bufio"
	"fmt"
	"github.com/Centny/Cny4go/log"
	"github.com/Centny/Cny4go/util"
	"io"
	"os"
	"strconv"
	"strings"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage:mrepo <out file> <in file> <in file> ...")
		return
	}
	for _, f := range os.Args[2:] {
		AppendF(f)
	}
	StoreCache(os.Args[1])
}

type Cover struct {
	A int
	B int
}

var cache map[string]Cover = map[string]Cover{}

func StoreCache(fname string) {
	f, err := os.OpenFile(fname, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		log.E("open file(%s) error:%s", fname, err.Error())
		return
	}
	defer f.Close()
	buf := bufio.NewWriter(f)
	for k, v := range cache {
		buf.WriteString(fmt.Sprintf("%v %v %v\n", k, v.A, v.B))
	}
	buf.Flush()
}

//append file
func AppendF(fname string) {
	log.D("append file:%s", fname)
	f, err := os.Open(fname)
	if err != nil {
		log.W("open file(%v) error:%s", fname, err.Error())
		return
	}
	defer f.Close()
	buf := bufio.NewReader(f)
	for {
		l, err := util.ReadLine(buf, 102400, false)
		if err == io.EOF {
			break
		} else if err != nil {
			log.W("read file(%s) error:%s", fname, err.Error())
		} else {
			AppendR(string(l))
		}
	}
}

//append one row
func AppendR(row string) {
	if strings.HasPrefix(row, "mode") {
		return
	}
	rvals := strings.Split(row, " ")
	if len(rvals) < 3 {
		log.W("invalid row:%s", row)
		return
	}
	a, err := strconv.Atoi(rvals[1])
	if err != nil {
		log.W("covert row(%v) error:%s ", row, err.Error())
		return
	}
	b, err := strconv.Atoi(rvals[2])
	if err != nil {
		log.W("covert row(%v) error:%s ", row, err.Error())
		return
	}
	if c, ok := cache[rvals[0]]; ok {
		if c.A < 1 {
			c.A = a
		}
		if c.B < 1 {
			c.B = b
		}
		cache[rvals[0]] = c
	} else {
		cache[rvals[0]] = Cover{A: a, B: b}
	}
}
