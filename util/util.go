package util

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"time"
)

func Fexists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
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
	fi, _ := f.Stat()
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
func AryExist(ary interface{}, obj interface{}) bool {
	switch reflect.TypeOf(ary).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(ary)
		for i := 0; i < s.Len(); i++ {
			if obj == s.Index(i).Interface() {
				return true
			}
		}
		return false
	default:
		return false
	}
}
func readAllStr(r io.Reader) string {
	if r == nil {
		return ""
	}
	bys, err := ioutil.ReadAll(r)
	if err != nil {
		return ""
	}
	return string(bys)
}

var HTTPClient http.Client

func HTTPGet(ufmt string, args ...interface{}) string {
	res, err := HTTPClient.Get(fmt.Sprintf(ufmt, args...))
	if err != nil {
		return ""
	}
	return readAllStr(res.Body)
}

func HTTPGet2(ufmt string, args ...interface{}) map[string]interface{} {
	data := HTTPGet(ufmt, args...)
	if len(data) < 1 {
		return nil
	}
	md := map[string]interface{}{}
	d := json.NewDecoder(strings.NewReader(data))
	err := d.Decode(&md)
	if err != nil {
		return nil
	}
	return md
}
