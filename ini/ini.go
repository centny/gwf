package ini

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/Centny/gwf/util"
	"io"
	"os"
	"strings"
)

type Ini struct {
	Kvs  map[string]map[string]map[string]int
	Path string
}

func (a *Ini) Add(section, key string, tval ...string) {
	if _, ok := a.Kvs[section]; !ok {
		a.Kvs[section] = map[string]map[string]int{}
	}
	if _, ok := a.Kvs[section][key]; !ok {
		a.Kvs[section][key] = map[string]int{}
	}
	for _, val := range tval {
		val = strings.Trim(val, "\t \n")
		a.Kvs[section][key][val] = 1
	}
}

func (a *Ini) Del(section, key string, tval ...string) error {
	if _, ok := a.Kvs[section]; ok {
		if _, ok := a.Kvs[section][key]; ok {
			for _, val := range tval {
				delete(a.Kvs[section][key], val)
			}
			return nil
		} else {
			return errors.New(fmt.Sprintf("values not found by key %s", key))
		}
	} else {
		return errors.New(fmt.Sprintf("values not found by section %s", section))
	}
}
func (a *Ini) Parse_(r *bufio.Reader) error {
	var sec string
	for {
		bys, err := util.ReadLine(r, 102400, false)
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return err
			}
		}
		line := string(bys)
		line = strings.Trim(line, "\t \n")
		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			sec = strings.Trim(line, "[ \t]")
			continue
		}
		if len(sec) < 1 {
			fmt.Println("ignore line(not section):", line)
			continue
		}
		kvs := strings.SplitN(line, "=", 2)
		if len(kvs) < 2 {
			fmt.Println("ignore line(not key value):", line)
			continue
		}
		a.Add(sec, strings.Trim(kvs[0], "\n\t "), strings.Split(kvs[1], ",")...)
	}
	return nil
}
func (a *Ini) Parse(fp string) error {
	f, err := os.OpenFile(fp, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return err
	}
	defer f.Close()
	return a.Parse_(bufio.NewReader(f))
}
func (a *Ini) Store_(w *bufio.Writer) error {
	for sec, kvs := range a.Kvs {
		w.WriteString(fmt.Sprintf("[%s]\n", sec))
		for key, vals := range kvs {
			vals_ := []string{}
			for val_, _ := range vals {
				vals_ = append(vals_, val_)
			}
			line := fmt.Sprintf(" %s=%s\n", key, strings.Join(vals_, ", "))
			w.WriteString(line)
		}
	}
	return w.Flush()
}
func (a *Ini) Store(fp string) error {
	f, err := os.OpenFile(fp, os.O_WRONLY|os.O_CREATE, os.ModePerm)
	if err != nil {
		return err
	}
	defer f.Close()
	return a.Store_(bufio.NewWriter(f))
}

func NewIni() *Ini {
	return &Ini{
		Kvs: map[string]map[string]map[string]int{},
	}
}

func ParseIni(fp string) (*Ini, error) {
	az := NewIni()
	return az, az.Parse(fp)
}

func Cmds(args []string) error {
	if len(args) > 3 {
		switch args[1] {
		case "add":
			ps := NewIni()
			ps.Parse(args[0])
			ps.Add(args[2], args[3], args[3:]...)
			return ps.Store(args[0])
		case "del":
			ps, err := ParseIni(args[0])
			if err != nil {
				return err
			}
			ps.Del(args[2], args[3], args[3:]...)
			return ps.Store(os.Args[0])
		}
	}
	return errors.New("Usage:inic <path> add|del section val1 val2 val3 ...")
}
