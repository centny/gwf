package util

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"strings"
)

var HTTPClient = HClient{
	Client: http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	},
}

type HClient struct {
	http.Client
}

func (h *HClient) HGet(ufmt string, args ...interface{}) (string, error) {
	_, str, err := h.HGet_H(map[string]string{}, ufmt, args...)
	return str, err
}
func (h *HClient) HGet_H(header map[string]string, ufmt string, args ...interface{}) (int, string, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf(ufmt, args...), nil)
	if err != nil {
		return 0, "", err
	}
	for k, v := range header {
		req.Header.Set(k, v)
	}
	res, err := h.Do(req)
	if err != nil {
		return 0, "", err
	}
	str, err := readAllStr(res.Body)
	return res.StatusCode, str, err
}
func (h *HClient) HGet2(ufmt string, args ...interface{}) (Map, error) {
	data, err := h.HGet(ufmt, args...)
	if len(data) < 1 || err != nil {
		return nil, err
	}
	return Json2Map(data)
}
func (h *HClient) HPost(url string, fields map[string]string) (string, error) {
	return h.HPostF(url, fields, "", "")
}
func (h *HClient) HPostF(url string, fields map[string]string, fkey string, fp string) (string, error) {
	_, str, err := h.HPostF_H(url, fields, map[string]string{}, fkey, fp)
	return str, err
}
func (h *HClient) HPostF_H(url string, fields map[string]string, header map[string]string, fkey string, fp string) (int, string, error) {
	ctype, bodyBuf, err := CreateFormBody(fields, fkey, fp)
	if err != nil {
		return 0, "", err
	}
	req, err := http.NewRequest("POST", url, bodyBuf)
	if err != nil {
		return 0, "", err
	}
	req.Header.Set("Content-Type", ctype)
	for k, v := range header {
		req.Header.Set(k, v)
	}
	res, err := h.Do(req)
	if err != nil {
		return 0, "", err
	}
	str, err := readAllStr(res.Body)
	return res.StatusCode, str, err
}

func (h *HClient) HPostN(url string, ctype string, buf io.Reader) (int, string, error) {
	req, err := http.NewRequest("POST", url, buf)
	if err != nil {
		return 0, "", err
	}
	req.Header.Set("Content-Type", ctype)
	res, err := h.Do(req)
	if err != nil {
		return 0, "", err
	}
	str, err := readAllStr(res.Body)
	return res.StatusCode, str, err
}

func (h *HClient) HPost2(url string, fields map[string]string) (Map, error) {
	data, err := h.HPost(url, fields)
	if len(data) < 1 || err != nil {
		return nil, err
	}
	return Json2Map(data)
}
func (h *HClient) HPostF2(url string, fields map[string]string, fkey string, fp string) (Map, error) {
	data, err := h.HPostF(url, fields, fkey, fp)
	if len(data) < 1 || err != nil {
		return nil, err
	}
	return Json2Map(data)
}
func (h *HClient) HTTPGet(ufmt string, args ...interface{}) string {
	res, _ := h.HGet(ufmt, args...)
	return res
}

func (h *HClient) HTTPGet2(ufmt string, args ...interface{}) Map {
	res, _ := h.HGet2(ufmt, args...)
	return res
}

func (h *HClient) HTTPPost(url string, fields map[string]string) string {
	res, _ := h.HPost(url, fields)
	return res
}

func (h *HClient) HTTPPost2(url string, fields map[string]string) Map {
	res, _ := h.HPost2(url, fields)
	return res
}

func (h *HClient) DLoad(spath string, header map[string]string, ufmt string, args ...interface{}) error {
	req, err := http.NewRequest("GET", fmt.Sprintf(ufmt, args...), nil)
	if err != nil {
		return err
	}
	for k, v := range header {
		req.Header.Set(k, v)
	}
	res, err := h.Do(req)
	if err != nil {
		return err
	}
	if res.StatusCode != 200 {
		return Err("http response code:%v", res.StatusCode)
	}
	f, err := os.OpenFile(spath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.ModePerm)
	if err != nil {
		return err
	}
	defer f.Close()
	buf := bufio.NewWriter(f)
	_, err = io.Copy(buf, res.Body)
	buf.Flush()
	return err
}

func HGet(ufmt string, args ...interface{}) (string, error) {
	return HTTPClient.HGet(ufmt, args...)
}
func HGet2(ufmt string, args ...interface{}) (Map, error) {
	return HTTPClient.HGet2(ufmt, args...)
}
func HPost(url string, fields map[string]string) (string, error) {
	return HTTPClient.HPost(url, fields)
}
func HPostF(url string, fields map[string]string, fkey string, fp string) (string, error) {
	return HTTPClient.HPostF(url, fields, fkey, fp)
}
func HPostN(url string, ctype string, buf io.Reader) (int, string, error) {
	return HTTPClient.HPostN(url, ctype, buf)
}
func HPost2(url string, fields map[string]string) (Map, error) {
	return HTTPClient.HPost2(url, fields)
}
func HPostF2(url string, fields map[string]string, fkey string, fp string) (Map, error) {
	return HTTPClient.HPostF2(url, fields, fkey, fp)
}
func HTTPGet(ufmt string, args ...interface{}) string {
	return HTTPClient.HTTPGet(ufmt, args...)
}

func HTTPGet2(ufmt string, args ...interface{}) Map {
	return HTTPClient.HTTPGet2(ufmt, args...)
}

func HTTPPost(url string, fields map[string]string) string {
	return HTTPClient.HTTPPost(url, fields)
}

func HTTPPost2(url string, fields map[string]string) Map {
	return HTTPClient.HTTPPost2(url, fields)
}

func DLoad(spath string, ufmt string, args ...interface{}) error {
	return HTTPClient.DLoad(spath, map[string]string{}, ufmt, args...)
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
func Map2Query(m Map) string {
	vs := url.Values{}
	for k, _ := range m {
		vs.Add(k, m.StrVal(k))
	}
	return vs.Encode()
}

func Json2Map(data string) (Map, error) {
	md := Map{}
	d := json.NewDecoder(strings.NewReader(data))
	err := d.Decode(&md)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("invalid json data(%s)", err.Error()))
	}
	return md, nil
}

func Json2Ary(data string) ([]interface{}, error) {
	var ary []interface{}
	d := json.NewDecoder(strings.NewReader(data))
	err := d.Decode(&ary)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("invalid json data(%s)", err.Error()))
	}
	return ary, nil
}
func CreateFileForm(bodyWriter *multipart.Writer, fkey, fp string) error {
	fileWriter, err := bodyWriter.CreateFormFile(fkey, fp)
	if err != nil {
		return err
	}
	fh, err := os.Open(fp)
	if err != nil {
		return err
	}
	defer fh.Close()
	_, err = io.Copy(fileWriter, fh)
	if err != nil {
		return err
	}
	return nil
}
func CreateFormBody(fields map[string]string, fkey string, fp string) (string, *bytes.Buffer, error) {
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)
	for k, v := range fields {
		bodyWriter.WriteField(k, v)
	}
	if len(fkey) > 0 {
		err := CreateFileForm(bodyWriter, fkey, fp)
		if err != nil {
			return "", nil, err
		}
	}
	ctype := bodyWriter.FormDataContentType()
	bodyWriter.Close()
	return ctype, bodyBuf, nil
}

type fs_size interface {
	Size() int64
}

type fs_stat interface {
	Stat() (os.FileInfo, error)
}
type fs_name interface {
	Name() string
}

func FormFSzie(src interface{}) int64 {
	var fsize int64 = 0
	if statInterface, ok := src.(fs_stat); ok {
		fileInfo, _ := statInterface.Stat()
		fsize = fileInfo.Size()
	}
	if sizeInterface, ok := src.(fs_size); ok {
		fsize = sizeInterface.Size()
	}
	return fsize
}

// func FormFName(src interface{}) string {
// 	if nameInterface, ok := src.(fs_name); ok {
// 		return nameInterface.Name()
// 	} else {
// 		return ""
// 	}
// }

func DoWeb(addr, dir string) error {
	return http.ListenAndServe(addr, http.FileServer(http.Dir(dir)))
}
