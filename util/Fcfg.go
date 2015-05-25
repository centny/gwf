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
type Fcfg struct {
	Map     map[string]interface{}
	ShowLog bool
	sec     string
	Lines   []string
	Seces   []string
	SecLn   map[string]int
}

func NewFcfg(uri string) (*Fcfg, error) {
	uri = strings.Trim(uri, " \t")
	cfg := &Fcfg{
		Map:     Map{},
		ShowLog: true,
		SecLn:   map[string]int{},
	}
	if strings.HasPrefix(uri, "http://") {
		return cfg, cfg.InitWithURL(uri)
	} else {
		return cfg, cfg.InitWithFilePath(uri)
	}
}
func NewFcfg2(data string) (*Fcfg, error) {
	cfg := &Fcfg{
		Map:     Map{},
		ShowLog: true,
		SecLn:   map[string]int{},
	}
	return cfg, cfg.InitWithData(data)
}
func NewFcfg3() *Fcfg {
	return &Fcfg{
		Map:     Map{},
		ShowLog: true,
		SecLn:   map[string]int{},
	}
}
func (f *Fcfg) slog(fs string, args ...interface{}) {
	if f.ShowLog {
		fmt.Println(fmt.Sprintf(fs, args...))
	}
}

//get the value by key.
func (f *Fcfg) Val(key string) string {
	if val, ok := f.Map[key]; ok {
		return val.(string)
	} else if val, ok := f.Map["loc/"+key]; ok {
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
	for k, v := range f.Map {
		sdata = fmt.Sprintf("%v %v=%v\n", sdata, k, v)
	}
	return sdata
}

func (f *Fcfg) Print() {
	fmt.Println(f.Show())
}
func (f *Fcfg) PrintSec(sec string) {
	sdata := ""
	for k, v := range f.Map {
		if strings.HasPrefix(k, sec) {
			sdata = fmt.Sprintf("%v %v=%v\n", sdata, k, v)
		}
	}
	fmt.Println(sdata)
}

//set the value by key and value.
func (f *Fcfg) SetVal(key string, val string) *Fcfg {
	f.Map[key] = val
	return f
}

//delete the value by key.
func (f *Fcfg) Del(key string) *Fcfg {
	delete(f.Map, key)
	return f
}

//check if exist by key.
func (f *Fcfg) Exist(key string) bool {
	if _, ok := f.Map[key]; ok {
		return true
	} else if _, ok := f.Map["loc/"+key]; ok {
		return true
	} else {
		return false
	}
}

//initial the configure by .properties file.
func (f *Fcfg) InitWithFilePath(fp string) error {
	f.slog("loading local configure->%v", fp)
	turl, _ := nurl.Parse(fp)
	qs := turl.Query()
	for k, _ := range qs {
		f.SetVal(k, qs.Get(k))
	}
	fp = turl.Path
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
		f.Lines = append(f.Lines, line)
		line = strings.Trim(line, " ")
		if len(line) < 1 {
			continue
		}
		err = f.exec(base, line)
		if err != nil {
			return err
		}
	}
	return nil
}

func (f *Fcfg) exec(base, line string) error {
	ps := strings.Split(line, "#")
	if len(ps) < 1 || len(ps[0]) < 1 {
		return nil
	}
	line = strings.Trim(ps[0], " \t")
	if regexp.MustCompile("^\\[[^\\]]*\\][\t ]*$").MatchString(line) {
		sec := strings.Trim(line, "\t []")
		f.sec = sec + "/"
		f.Seces = append(f.Seces, sec)
		f.SecLn[sec] = len(f.Lines)
		return nil
	}
	if !strings.HasPrefix(line, "@") {
		ps = strings.SplitN(line, "=", 2)
		if len(ps) < 2 {
			f.slog("not value key found:%v", ps[0])
		} else {
			key := f.sec + f.EnvReplace(strings.Trim(ps[0], " "))
			val := f.EnvReplace(strings.Trim(ps[1], " "))
			f.Map[key] = val
		}
		return nil
	}
	line = strings.TrimPrefix(line, "@")
	ps = strings.SplitN(line, ":", 2)
	if len(ps) < 2 || len(ps[1]) < 1 {
		f.slog("%v", f.EnvReplace(line))
		return nil
	}
	ps[0] = strings.ToLower(strings.Trim(ps[0], " \t"))
	ps[0] = f.EnvReplace(ps[0])
	if ps[0] == "l" {
		ps[1] = strings.Trim(ps[1], " \t")
		if len(ps[1]) < 1 {

		}
		return f.load(base, ps[1])
	}
	if cs := strings.SplitN(ps[0], "==", 2); len(cs) == 2 {
		if cs[0] == cs[1] {
			return f.exec(base, ps[1])
		} else {
			return nil
		}
	}
	if cs := strings.SplitN(ps[0], "!=", 2); len(cs) == 2 {
		if cs[0] != cs[1] {
			return f.exec(base, ps[1])
		} else {
			return nil
		}
	}
	//all other will print line.
	f.slog("%v", f.EnvReplace(line))
	return nil
}
func (f *Fcfg) load(base, line string) error {
	line = f.EnvReplaceV(line, true)
	line = strings.Trim(line, "\t ")
	if len(line) < 1 {
		return nil
	}
	if !(strings.HasPrefix(line, "http://") || filepath.IsAbs(line)) {
		line = base + line
	}
	cfg, err := NewFcfg(line)
	if err == nil {
		f.Merge(cfg)
	}
	return err
}

//initial the configure by .properties file.
func (f *Fcfg) InitWithFile(tfile *os.File) error {
	reader := bufio.NewReader(tfile)
	dir, _ := filepath.Split(tfile.Name())
	return f.InitWithReader2(dir, reader)
}

//initial the configure by network .properties URL.
func (f *Fcfg) InitWithURL(url string) error {
	f.slog("loading remote configure->%v", url)
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
func (f *Fcfg) EnvReplace(val string) string {
	return f.EnvReplaceV(val, false)
}

//replace tartget patter by ${key} with value in configure map or system environment value.
func (f *Fcfg) EnvReplaceV(val string, empty bool) string {
	reg := regexp.MustCompile("\\$\\{[^\\}]*\\}")
	var rval string
	val = reg.ReplaceAllStringFunc(val, func(m string) string {
		keys := strings.Split(strings.Trim(m, "${}\t "), ",")
		for _, key := range keys {
			if f.Exist(key) {
				rval = f.Val(key)
			} else {
				rval = os.Getenv(key)
			}
			if len(rval) > 0 {
				break
			}
		}
		if len(rval) > 0 {
			return rval
		}
		if empty {
			return ""
		} else {
			return m
		}
	})
	// var rval string = ""
	// mhs := reg.FindAll([]byte(val), -1)
	// for i := 0; i < len(mhs); i++ {
	// 	bys := mhs[i]
	// 	ptn := string(bys)
	// 	bys = bys[2 : len(bys)-1]
	// 	if len(bys) < 1 {
	// 		continue
	// 	}
	// 	key := string(bys)
	// 	if f.Exist(key) {
	// 		rval = f.Val(key)
	// 	} else {
	// 		rval = os.Getenv(key)
	// 		if len(rval) < 1 {
	// 			continue
	// 		}
	// 	}
	// 	val = strings.Replace(val, ptn, rval, 1)
	// }
	return val
}

//merge another configure.
func (f *Fcfg) Merge(t *Fcfg) {
	if t == nil {
		return
	}
	for k, v := range t.Map {
		f.Map[k] = v
	}
}

func (f *Fcfg) String() string {
	buf := bytes.NewBuffer(nil)
	for k, v := range f.Map {
		buf.WriteString(fmt.Sprintf("%v=%v\n", k, v))
	}
	return buf.String()
}
func (f *Fcfg) Store(sec, fp, tsec string) error {
	var seci int = -1
	for idx, s := range f.Seces {
		if s == sec {
			seci = idx
		}
	}
	if seci < 0 {
		return Err("section not found by %v", sec)
	}
	var beg, end int = f.SecLn[sec], len(f.Lines)
	if seci < len(f.Seces)-1 {
		end = f.SecLn[f.Seces[seci+1]] - 1
	}
	tf, err := os.OpenFile(fp, os.O_CREATE|os.O_APPEND|os.O_WRONLY, os.ModePerm)
	if err != nil {
		return err
	}
	defer tf.Close()
	buf := bufio.NewWriter(tf)
	buf.WriteString("[" + tsec + "]\n")
	for i := beg; i < end; i++ {
		buf.WriteString(f.Lines[i])
		buf.WriteString("\n")
	}
	return buf.Flush()
}
