package util

import (
	"bufio"
	"bytes"
	"fmt"
	"math"
	nurl "net/url"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

//Author:Centny
//
//the file configure
//
type Fcfg map[string]interface{}

func NewFcfg(uri string) (*Fcfg, error) {
	uri = strings.Trim(uri, " \t")
	cfg := &Fcfg{}
	if strings.HasPrefix(uri, "http://") {
		return cfg, cfg.InitWithURL(uri)
	} else {
		return cfg, cfg.InitWithFilePath(uri)
	}
}
func NewFcfg2(data string) (*Fcfg, error) {
	cfg := &Fcfg{}
	return cfg, cfg.InitWithData(data)
}

//get the value by key.
func (f *Fcfg) Val(key string) string {
	if val, ok := (*f)[key]; ok {
		return val.(string)
	} else {
		return ""
	}
}

//get the int value by key.
func (f *Fcfg) IntVal(key string) int {
	if !f.Exist(key) {
		return math.MaxInt8
	}
	val, err := strconv.Atoi(f.Val(key))
	if err != nil {
		return math.MaxInt8
	}
	return val
}

//get the float value by key.
func (f *Fcfg) FloatVal(key string) float64 {
	if !f.Exist(key) {
		return math.MaxFloat64
	}
	val, err := strconv.ParseFloat(f.Val(key), 64)
	if err != nil {
		return math.MaxFloat64
	}
	return val
}

//
func (f *Fcfg) Show() string {
	sdata := ""
	for k, v := range *f {
		sdata = fmt.Sprintf("%v\t%v=%v\n", sdata, k, v)
	}
	return sdata
}

//set the value by key and value.
func (f *Fcfg) SetVal(key string, val string) *Fcfg {
	(*f)[key] = val
	return f
}

//delete the value by key.
func (f *Fcfg) Del(key string) *Fcfg {
	delete(*f, key)
	return f
}

//check if exist by key.
func (f *Fcfg) Exist(key string) bool {
	_, ok := (*f)[key]
	return ok
}

//initial the configure by .properties file.
func (f *Fcfg) InitWithFilePath(fp string) error {
	if !Fexists(fp) {
		return Err("file(%v) not found", fp)
	}
	fh, err := os.Open(fp)
	if err != nil {
		return err
	}
	defer fh.Close()
	return f.InitWithFile(fh)
}

//initial the configure by .properties format reader.
func (f *Fcfg) InitWithReader(reader *bufio.Reader) error {
	return f.InitWithReader2("", reader)
}
func (f *Fcfg) InitWithReader2(base string, reader *bufio.Reader) error {
	for {
		//read one line
		bys, err := ReadLine(reader, 10000, false)
		if err != nil {
			break
		}
		//
		line := string(bys)
		line = strings.Trim(line, " ")
		if len(line) < 1 {
			continue
		}
		ps := strings.Split(line, "#")
		if len(ps) < 1 || len(ps[0]) < 1 {
			continue
		}
		line = ps[0]
		if strings.HasPrefix(line, "@") {
			line = strings.TrimPrefix(line, "@")
			if !(strings.HasPrefix(line, "http://") || filepath.IsAbs(line)) {
				line = base + line
			}
			cfg, err := NewFcfg(line)
			if err == nil {
				f.Merge(cfg)
				continue
			} else {
				return err
			}
		} else {
			ps = strings.SplitN(line, "=", 2)
			if len(ps) < 2 {
				fmt.Println("not value key found:", ps[0])
				continue
			}
			key := f.EnvReplace(strings.Trim(ps[0], " "))
			val := f.EnvReplace(strings.Trim(ps[1], " "))
			(*f)[key] = val
		}
	}
	return nil
}

//initial the configure by .properties file.
func (f *Fcfg) InitWithFile(tfile *os.File) error {
	reader := bufio.NewReader(tfile)
	dir, _ := filepath.Split(tfile.Name())
	return f.InitWithReader2(dir, reader)
}

//initial the configure by network .properties URL.
func (f *Fcfg) InitWithURL(url string) error {
	fmt.Println("loading remote configure->" + url)
	sres, err := HGet(url)
	if err == nil {
		turl, _ := nurl.Parse(url)
		turl.Path, _ = filepath.Split(turl.Path)
		return f.InitWithReader2(
			fmt.Sprintf("%v://%v%v", turl.Scheme, turl.Host, turl.Path),
			bufio.NewReader(bytes.NewBufferString(sres)))
	} else {
		return err
	}
}
func (f *Fcfg) InitWithData(data string) error {
	return f.InitWithReader(bufio.NewReader(bytes.NewBufferString(data)))
}

//replace tartget patter by ${key} with value in configure map or system environment value.
func (f *Fcfg) EnvReplace(val string) string {
	reg, _ := regexp.Compile("\\$\\{[^\\}]*\\}")
	var rval string = ""
	mhs := reg.FindAll([]byte(val), -1)
	for i := 0; i < len(mhs); i++ {
		bys := mhs[i]
		ptn := string(bys)
		bys = bys[2 : len(bys)-1]
		if len(bys) < 1 {
			continue
		}
		key := string(bys)
		if f.Exist(key) {
			rval = f.Val(key)
		} else {
			rval = os.Getenv(key)
		}
		if len(rval) < 1 {
			continue
		}
		val = strings.Replace(val, ptn, rval, 1)
	}
	return val
}

//merge another configure.
func (f *Fcfg) Merge(t *Fcfg) {
	if t == nil {
		return
	}
	for k, v := range map[string]interface{}(*t) {
		(*f)[k] = v
	}
}
