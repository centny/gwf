package util

import (
	"bufio"
	"bytes"
	"fmt"
	nurl "net/url"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
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
	Base    string
}

func NewFcfg(uri string) (*Fcfg, error) {
	uri = strings.Trim(uri, " \t")
	cfg := &Fcfg{
		Map:     Map{},
		ShowLog: true,
		SecLn:   map[string]int{},
	}
	return cfg, cfg.InitWithUri(uri)
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
func NewFcfg4(src *Fcfg) *Fcfg {
	cfg := &Fcfg{
		Map:     Map{},
		ShowLog: true,
		SecLn:   map[string]int{},
	}
	cfg.InitWithCfg(src)
	return cfg
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

func (f *Fcfg) Val2(key, def string) string {
	val := f.Val(key)
	if len(val) > 0 {
		return val
	} else {
		return def
	}
}

//get the int value by key.
func (f *Fcfg) IntVal(key string) int {
	return f.IntValV(key, 0)
}

func (f *Fcfg) IntValV(key string, d int) int {
	if !f.Exist(key) {
		return d
	}
	val, err := strconv.Atoi(f.Val(key))
	if err != nil {
		return d
	}
	return val
}
func (f *Fcfg) Int64Val(key string) int64 {
	return f.Int64ValV(key, 0)
}
func (f *Fcfg) Int64ValV(key string, d int64) int64 {
	if !f.Exist(key) {
		return d
	}
	val, err := IntValV(f.Val(key))
	if err != nil {
		return d
	}
	return val
}

//get the float value by key.
func (f *Fcfg) FloatVal(key string) float64 {
	return f.FloatValV(key, 0)
}
func (f *Fcfg) FloatValV(key string, d float64) float64 {
	if !f.Exist(key) {
		return d
	}
	val, err := strconv.ParseFloat(f.Val(key), 64)
	if err != nil {
		return d
	}
	return val
}

func (f *Fcfg) FileModeV(key string, d os.FileMode) os.FileMode {
	if !f.Exist(key) {
		return d
	}
	val, err := strconv.ParseUint(f.Val(key), 8, 32)
	if err != nil {
		return d
	}
	return os.FileMode(val)
}

//
func (f *Fcfg) Show() string {
	return f.String()
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
func (f *Fcfg) Range(sec string, cb func(key string, val interface{})) {
	for k, v := range f.Map {
		if strings.HasPrefix(k, sec) {
			cb(strings.TrimPrefix(k, sec+"/"), v)
		}
	}
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
func (f *Fcfg) InitWithUri(uri string) error {
	if strings.HasPrefix(uri, "http://") {
		return f.InitWithURL(uri)
	} else if strings.HasPrefix(uri, "https://") {
		return f.InitWithURL(uri)
	} else {
		return f.InitWithFilePath(uri)
	}
}

func (f *Fcfg) InitWithUri2(uri string, wait bool) error {
	if strings.HasPrefix(uri, "http://") {
		return f.InitWithURL2(uri, wait)
	} else if strings.HasPrefix(uri, "https://") {
		return f.InitWithURL2(uri, wait)
	} else {
		return f.InitWithFilePath2(uri, wait)
	}
}

func (f *Fcfg) InitWithCfg(cfg *Fcfg) {
	for k, v := range cfg.Map {
		f.Map[k] = v
	}
}

//initial the configure by .properties file.
func (f *Fcfg) InitWithFilePath(fp string) error {
	return f.InitWithFilePath2(fp, true)
}
func (f *Fcfg) InitWithFilePath2(fp string, wait bool) error {
	f.slog("loading local configure->%v", fp)
	var fps = strings.Split(fp, "?")
	fp = fps[0]
	if len(fps) > 1 {
		turl, _ := nurl.Parse("/abc?" + fps[1])
		qs := turl.Query()
		for k, _ := range qs {
			f.SetVal(k, qs.Get(k))
		}
	}
	for !Fexists(fp) {
		if wait {
			f.slog("file(%v) not found", fp)
			time.Sleep(3 * time.Second)
			continue
		} else {
			return Err("file(%v) not found", fp)
		}
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
	if len(base) > 0 {
		f.Base = base
	}
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
func (f *Fcfg) InitWithConfReader2(reader *bufio.Reader) error {
	var key string = ""
	var val string = ""
	for {
		//read one line
		bys, err := ReadLine(reader, 10000, false)
		if err != nil {
			if len(key) > 0 {
				f.Map[key] = strings.Trim(val, "\n")
				key, val = "", ""
			}
			break
		}
		line := string(bys)
		if regexp.MustCompile("^\\[[^\\]]*\\][\t ]*$").MatchString(line) {
			sec := strings.Trim(line, "\t []")
			if len(key) > 0 {
				f.Map[key] = strings.Trim(val, "\n")
				key, val = "", ""
			}
			key = sec
		} else {
			val += line + "\n"
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
	if len(line) < 1 {
		return nil
	}
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
	if !(strings.HasPrefix(line, "http://") || strings.HasPrefix(line, "https://") || filepath.IsAbs(line)) {
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
	if len(dir) < 1 {
		dir = "."
	}
	dir, _ = filepath.Abs(dir)
	if !strings.HasSuffix(dir, "/") {
		dir += string(filepath.Separator)
	}
	if strings.HasSuffix(tfile.Name(), ".conf") {
		return f.InitWithConfReader2(reader)
	} else {
		return f.InitWithReader2(dir, reader)
	}
}

//initial the configure by network .properties URL.
func (f *Fcfg) InitWithURL(url string) error {
	return f.InitWithURL2(url, true)
}
func (f *Fcfg) hget(url string) (sres string, err error) {
	code, sres, err := HGet3(url)
	if err == nil && code != 200 {
		err = fmt.Errorf("status code(%v)", code)
	}
	return
}
func (f *Fcfg) InitWithURL2(url string, wait bool) error {
	f.slog("loading remote configure->%v", url)
	var sres string
	var err error
	for {
		sres, err = f.hget(url)
		if err == nil {
			f.slog("loading remote configure(%v) success", url)
			break
		}
		f.slog("loading remote configure(%v):%v", url, err.Error())
		if wait {
			time.Sleep(3 * time.Second)
			continue
		} else {
			break
		}
	}
	if err == nil {
		turl, _ := nurl.Parse(url)
		turl.Path, _ = filepath.Split(turl.Path)
		if strings.HasSuffix(turl.Path, ".conf") {
			return f.InitWithConfReader2(bufio.NewReader(bytes.NewBufferString(sres)))
		} else {
			return f.InitWithReader2(
				fmt.Sprintf("%v://%v%v", turl.Scheme, turl.Host, turl.Path),
				bufio.NewReader(bytes.NewBufferString(sres)))
		}
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
			} else if key == "C_PWD" {
				rval = f.Base
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
	for _, s := range t.Seces {
		if _, ok := f.SecLn[s]; ok {
			continue
		}
		f.Seces = append(f.Seces, s)
	}
}

func (f *Fcfg) Merge2(sec string, t *Fcfg) {
	for k, v := range t.Map {
		f.Map[sec+"/"+k] = v
	}
	if _, ok := f.SecLn[sec]; !ok {
		f.Seces = append(f.Seces, sec)
	}
}

func (f *Fcfg) Strip(sec string) *Fcfg {
	var cfg = NewFcfg3()
	for k, v := range f.Map {
		if !strings.HasPrefix(k, sec) {
			continue
		}
		cfg.Map["loc"+strings.TrimPrefix(k, sec)] = v
	}
	return cfg
}

func (f *Fcfg) String() string {
	buf := bytes.NewBuffer(nil)
	keys, locs := []string{}, []string{}
	for k, _ := range f.Map {
		if strings.HasPrefix(k, "loc/") {
			locs = append(locs, k)
		} else {
			keys = append(keys, k)
		}
	}
	sort.Sort(sort.StringSlice(keys))
	for _, k := range keys {
		buf.WriteString(fmt.Sprintf("%v=%v\n", k, f.Map[k]))
	}
	for _, k := range locs {
		buf.WriteString(fmt.Sprintf("%v=%v\n", k, f.Map[k]))
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
