package jcr

import (
	"fmt"
	"github.com/Centny/Cny4go/log"
	"github.com/Centny/Cny4go/util"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func RunAppend(dir, ex, in, out, js string) {
	exs := strings.Split(ex, ",")
	ins := strings.Split(in, ",")
	dir, _ = filepath.Abs(dir)
	out, _ = filepath.Abs(out)
	if !util.Fexists(out) {
		err := os.Mkdir(out, os.ModePerm)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
	}
	cover_c := 0
	err := list(dir, out, func(path string) error {
		if match(exs, path) {
			return nil
		}
		if !match(ins, path) {
			return nil
		}
		defer func() {
			cover_c++
		}()
		opath := out + "/" + strings.TrimPrefix(path, dir)
		odir := filepath.Dir(opath)
		os.MkdirAll(odir, os.ModePerm)
		err := util.FCopy(path, opath)
		if err != nil {
			return err
		}
		return util.FAppend(opath, fmt.Sprintf(`
<script type="text/javascript" src="%s" ></script>
			`, js))
	})
	if err != nil {
		log.E("jcr error:%s", err.Error())
	} else {
		log.D("jcr execute success, %d file is covered...", cover_c)
	}
}
func list(dir string, out string, exec func(path string) error) error {
	return filepath.Walk(dir, func(path string, fi os.FileInfo, err error) error {
		return walk_c(dir, out, path, fi, err, exec)
	})
}
func walk_c(dir string, out string, path string, fi os.FileInfo, err error, exec func(path string) error) error {
	if err != nil {
		log.W("list path error(%v)", err.Error())
		return nil
	}
	if fi.IsDir() {
		return nil
	}
	if strings.HasPrefix(path, out) {
		return nil
	}
	return exec(path)
}
func match(regs []string, path string) bool {
	for _, reg := range regs {
		if len(reg) < 1 {
			continue
		}
		if regexp.MustCompile(reg).MatchString(path) {
			return true
		}
	}
	return false
}
