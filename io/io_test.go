package io

import (
	"bufio"
	"bytes"
	"fmt"
	"testing"
)

func TestReadLine(t *testing.T) {
	f := func(end bool) {
		bf := bytes.NewBufferString("abc\ndef\nghi\n")
		r := bufio.NewReader(bf)
		for {
			bys, err := ReadLine(r, 10000, end)
			// bys, isp, err := r.ReadLine()
			fmt.Println(string(bys), err)
			if err != nil {
				break
			}
		}
	}
	f(true)
	f(false)

	f2 := func(end bool) {
		bf := bytes.NewBufferString("abc\ndef\nghi\n")
		r := bufio.NewReader(bf)
		var last int64
		for {
			bys, err := ReadLineLast(r, 10000, end, &last)
			// bys, isp, err := r.ReadLine()
			fmt.Println(string(bys), err)
			if err != nil {
				break
			}
		}
	}
	f2(true)
	f2(false)
	bf := bytes.NewBufferString("abc\ndef\nghi\n")
	r := bufio.NewReader(bf)
	_, err := ReadLine(r, 2, true)
	if err == nil {
		t.Error("error")
	}
}

func TestReadFull(t *testing.T) {
	r := bufio.NewReader(&Sw{})
	buf := make([]byte, 3)
	var las int64
	ReadFull(r, buf, &las)
	fmt.Println(string(buf))
	ReadFull(r, buf, &las)
}

type Sw struct {
	i int
}

func (s *Sw) Read(p []byte) (n int, err error) {
	if s.i < 1 {
		s.i = 1
		p[0] = 'A'
		return 1, nil
	} else if s.i < 2 {
		s.i = 2
		p[0] = 'B'
		p[1] = 'C'
		return 2, nil
	} else {
		return 0, fmt.Errorf("erro")
	}
}
