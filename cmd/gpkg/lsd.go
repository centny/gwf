package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type Lsd struct {
	root string
	P    string
	M    string
	I    string
	Ms   map[string]bool
}

func NewLsd(m, i, p string) *Lsd {
	return &Lsd{
		M:  m,
		I:  i,
		P:  p,
		Ms: map[string]bool{},
	}
}
func (l *Lsd) WalkFunc(path string, info os.FileInfo, err error) error {
	if err != nil {
		return err
	}
	if info.IsDir() {
		return nil
	}
	if !strings.HasSuffix(path, ".go") {
		return nil
	}
	if strings.HasSuffix(l.root, "/") {
		path = strings.TrimPrefix(path, l.root)
	} else {
		path = strings.TrimPrefix(path, l.root+"/")
	}
	if len(l.P) > 0 && !strings.HasPrefix(path, l.P) {
		return nil
	}
	dir, _ := filepath.Split(path)
	dir = strings.TrimSuffix(dir, "/")
	// fmt.Println(l.I, dir, regexp.MustCompile(l.I).MatchString(dir))
	if len(l.I) > 0 && regexp.MustCompile(l.I).MatchString(dir) {
		return nil
	}
	if len(l.M) > 0 && !regexp.MustCompile(l.M).MatchString(dir) {
		return nil
	}
	if len(dir) > 0 && !l.Ms[dir] {
		l.Ms[dir] = true
	}
	return nil
}

func (l *Lsd) Walk(root string) error {
	l.root = root
	return filepath.Walk(root, l.WalkFunc)
}
func (l *Lsd) Print(a string) {
	for m, _ := range l.Ms {
		fmt.Println(a + m)
	}
}
func (l *Lsd) JoinPrint(sep, a string) {
	var tms []string
	for k, _ := range l.Ms {
		tms = append(tms, a+k)
	}
	fmt.Println(strings.Join(tms, sep))
}
