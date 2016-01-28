package util

import (
	"os"
	"path/filepath"
	"regexp"
)

func FilterDir(root string, inc []string, exc []string) []string {
	pathes := []string{}
	reg_i := []*regexp.Regexp{}
	reg_e := []*regexp.Regexp{}
	for _, i := range inc {
		reg_i = append(reg_i, regexp.MustCompile(i))
	}
	for _, e := range exc {
		reg_e = append(reg_e, regexp.MustCompile(e))
	}
	m_reg := func(v string, regs []*regexp.Regexp) bool {
		for _, reg := range regs {
			if reg.MatchString(v) {
				return true
			}
		}
		return false
	}
	filepath.Walk(root, func(path string, fi os.FileInfo, err error) error {
		if fi == nil || !fi.IsDir() {
			return nil
		}
		if m_reg(path, reg_i) {
			pathes = append(pathes, path)
			return nil
		} else if m_reg(path, reg_e) {
			return nil
		} else {
			pathes = append(pathes, path)
			return nil
		}
	})
	return pathes
}
