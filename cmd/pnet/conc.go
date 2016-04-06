package main

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/Centny/gwf/util"
	"os"
	"regexp"
	"sort"
)

type con struct {
	Loc string
	Rem string
}

type cons struct {
	L  bool
	CS []con
}

func (c cons) Less(i, j int) bool {
	if c.L {
		return c.CS[i].Loc < c.CS[j].Loc
	} else {
		return c.CS[i].Rem < c.CS[j].Rem
	}
}

func (c cons) Swap(i, j int) {
	c.CS[i], c.CS[j] = c.CS[j], c.CS[i]
}

func (c cons) Len() int {
	return len(c.CS)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage:conc file")
		return
	}
	var file, err = os.OpenFile(os.Args[1], os.O_RDONLY, os.ModePerm)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()
	var buf = bufio.NewReader(file)
	var reg = regexp.MustCompile("[ \t]+")
	var p, q, r, s uint8
	var port uint16
	var cs = cons{}
	cs.L = len(os.Args) > 2 && os.Args[2] == "L"
	for i := 0; true; i++ {
		bys, err := util.ReadLine(buf, 10240, false)
		if err != nil {
			break
		}
		strs := reg.Split(string(bys), 5)
		if len(strs) < 4 {
			continue
		}
		lbuf := bytes.NewBufferString(strs[2])
		lres, err := fmt.Fscanf(lbuf, "%2x%2x%2x%2x:%x", &p, &q, &r, &s, &port)
		if err != nil || lres != 5 {
			continue
		}
		laddr := fmt.Sprintf("%v.%v.%v.%v:%v", s, r, q, p, port)
		rbuf := bytes.NewBufferString(strs[3])
		rres, err := fmt.Fscanf(rbuf, "%2x%2x%2x%2x:%x", &p, &q, &r, &s, &port)
		if err != nil || rres != 5 {
			continue
		}
		raddr := fmt.Sprintf("%v.%v.%v.%v:%v", s, r, q, p, port)
		cs.CS = append(cs.CS, con{
			Loc: laddr,
			Rem: raddr,
		})
	}
	sort.Sort(cs)
	for _, c := range cs.CS {
		fmt.Printf("%20s  %20s\n", c.Loc, c.Rem)
	}
}
