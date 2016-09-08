package io

import (
	"bufio"
	"fmt"
	"github.com/Centny/gwf/util/util"
)

// var DEFAULT_MODE os.FileMode = os.ModePerm

//read one line from reader and limit data,
func ReadLine(reader *bufio.Reader, limit int, end bool) ([]byte, error) {
	return ReadLineLast(reader, limit, end, nil)
}

//read one line from reader and limit data, the last is record the last read data time
func ReadLineLast(reader *bufio.Reader, limit int, end bool, last *int64) ([]byte, error) {
	var isPrefix bool = true
	var bys []byte
	var tmp []byte
	var err error
	for isPrefix {
		tmp, isPrefix, err = reader.ReadLine()
		if err != nil {
			return nil, err
		}
		bys = append(bys, tmp...)
		if len(bys) > limit {
			return nil, fmt.Errorf("too long by limt(%v)", limit)
		}
		if last != nil {
			*last = util.Now()
		}
	}
	if end {
		bys = append(bys, '\n')
	}
	return bys, nil
}

func ReadFull(reader *bufio.Reader, buf []byte, last *int64) error {
	length, all := len(buf), 0
	tbuf := buf
	for {
		readed, err := reader.Read(tbuf)
		if err != nil {
			return err
		}
		if last != nil {
			*last = util.Now()
		}
		all += readed
		if all < length {
			tbuf = buf[all:]
			continue
		} else {
			break
		}
	}
	return nil
}
