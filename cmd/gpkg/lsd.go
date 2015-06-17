package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Lsd struct {
	root string
	Ms   map[string]bool
}

func NewLsd() *Lsd {
	return &Lsd{
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
	dir, _ := filepath.Split(path)
	dir = strings.TrimSuffix(dir, "/")
	if len(dir) > 0 && !l.Ms[dir] {
		l.Ms[dir] = true
	}
	return nil
}

func (l *Lsd) Walk(root string) error {
	l.root = root
	return filepath.Walk(root, l.WalkFunc)
}
func (l *Lsd) Print() {
	for m, _ := range l.Ms {
		fmt.Println(m)
	}
}
func (l *Lsd) JoinPrint(sep string) {
	var tms []string
	for k, _ := range l.Ms {
		tms = append(tms, k)
	}
	fmt.Println(strings.Join(tms, sep))
}
