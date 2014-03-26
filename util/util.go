package util

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
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
func readAllStr(r io.Reader) (string, error) {
	if r == nil {
		return "", nil
	}
	bys, err := ioutil.ReadAll(r)
	if err != nil {
		return "", nil
	}
	return string(bys), nil
}

var HTTPClient http.Client

func HGet(ufmt string, args ...interface{}) (string, error) {
	res, err := HTTPClient.Get(fmt.Sprintf(ufmt, args...))
	if err != nil {
		return "", err
	}
	return readAllStr(res.Body)
}
func HGet2(ufmt string, args ...interface{}) (Map, error) {
	data, err := HGet(ufmt, args...)
	if len(data) < 1 || err != nil {
		return nil, err
	}
	return Json2Map(data)
}
func HPost(url string, fields map[string]string) (string, error) {
	return HPostF(url, fields, "", "")
}
func HPostF(url string, fields map[string]string, fkey string, fp string) (string, error) {
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)
	defer bodyWriter.Close()
	for k, v := range fields {
		bodyWriter.WriteField(k, v)
	}
	if len(fkey) > 0 {
		fileWriter, err := bodyWriter.CreateFormFile(fkey, fp)
		if err != nil {
			return "", err
		}
		fh, err := os.Open(fp)
		if err != nil {
			return "", err
		}
		defer fh.Close()
		_, err = io.Copy(fileWriter, fh)
		if err != nil {
			return "", err
		}
	}
	res, err := HTTPClient.Post(url, bodyWriter.FormDataContentType(), bodyBuf)
	if err != nil {
		return "", err
	}
	return readAllStr(res.Body)
}
func HPost2(url string, fields map[string]string) (Map, error) {
	data, err := HPost(url, fields)
	if len(data) < 1 || err != nil {
		return nil, err
	}
	return Json2Map(data)
}
func HTTPGet(ufmt string, args ...interface{}) string {
	res, _ := HGet(ufmt, args...)
	return res
}

func HTTPGet2(ufmt string, args ...interface{}) Map {
	res, _ := HGet2(ufmt, args...)
	return res
}

func HTTPPost(url string, fields map[string]string) string {
	res, _ := HPost(url, fields)
	return res
}

func HTTPPost2(url string, fields map[string]string) Map {
	res, _ := HPost2(url, fields)
	return res
}

func Map2Query(m Map) string {
	vs := url.Values{}
	for k, v := range m {
		vs.Add(k, v.(string))
	}
	return vs.Encode()
}

func Json2Map(data string) (Map, error) {
	md := Map{}
	d := json.NewDecoder(strings.NewReader(data))
	err := d.Decode(&md)
	if err != nil {
		return nil, err
	}
	return md, nil
}
