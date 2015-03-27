package iow

import (
	"bufio"
	"fmt"
	"io"
	"testing"
)

type Reader struct {
	c int
}

func (r *Reader) Read(bys []byte) (n int, err error) {
	if r.c > 4 {
		return 0, io.EOF
	}
	bys[0] = 0
	bys[1] = 3
	bys[2] = 'A'
	bys[3] = 'B'
	bys[4] = 'C'
	r.c++
	return 5, nil
}

type Reader2 struct {
	c int
}

func (r *Reader2) Read(bys []byte) (n int, err error) {
	if r.c > 1 {
		bys[0] = 0
		bys[1] = 0
		bys[2] = 'A'
		bys[3] = 'B'
		bys[4] = 'C'
	} else {
		bys[0] = 0
		bys[1] = 3
		bys[2] = 'A'
		bys[3] = 'B'
		bys[4] = 'C'
	}
	r.c++
	return 5, nil
}

type Reader3 struct {
	c int
}

func (r *Reader3) Read(bys []byte) (n int, err error) {
	if r.c > 0 {
		return 0, io.EOF
	} else {
		bys[0] = 0
		bys[1] = 3
		r.c++
		return 2, nil
	}
}
func TestRl(t *testing.T) {
	err := ReadLdata(bufio.NewReader(&Reader{}), func(bys []byte) error {
		fmt.Println(string(bys), len(bys))
		return nil
	})
	fmt.Println(err.Error())
	err = ReadLdata(bufio.NewReader(&Reader2{}), func(bys []byte) error {
		fmt.Println(string(bys), len(bys))
		return nil
	})
	fmt.Println(err.Error())
	err = ReadLdata(bufio.NewReader(&Reader3{}), func(bys []byte) error {
		fmt.Println(string(bys), len(bys))
		return nil
	})
	fmt.Println(err.Error())
	err = ReadLdata(bufio.NewReader(&Reader{}), func(bys []byte) error {
		return io.EOF
	})
	fmt.Println(err.Error())
}
