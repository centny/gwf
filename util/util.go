package util

import (
	"bufio"
	"errors"
	"os"
	"path/filepath"
	"time"
)

func Fexists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	return !os.IsNotExist(err)
}

func FTouch(path string) error {
	f, err := os.Open(path)
	if err != nil {
		p := filepath.Dir(path)
		if !Fexists(p) {
			err := os.MkdirAll(p, os.ModePerm)
			if err != nil {
				return err
			}
		}
		f, err = os.Create(path)
		if f != nil {
			defer f.Close()
		}
		return err
	}
	defer f.Close()
	fi, err := f.Stat()
	if err != nil {
		return err
	}
	if fi.IsDir() {
		return errors.New("can't touch path")
	}
	return nil
}

func ReadLine(r *bufio.Reader, limit int, end bool) ([]byte, error) {
	var isPrefix bool = true
	var bys []byte
	var tmp []byte
	var err error
	for isPrefix {
		tmp, isPrefix, err = r.ReadLine()
		if err != nil {
			return nil, err
		}
		bys = append(bys, tmp...)
	}
	if end {
		bys = append(bys, '\n')
	}
	return bys, nil
}

func Timestamp(t time.Time) int64 {
	return t.UnixNano() / 1e6
}
func Time(timestamp int64) time.Time {
	return time.Unix(0, timestamp*1e6)
}
