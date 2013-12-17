package util

import (
	"errors"
	"os"
	"path/filepath"
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
		_, err = os.Create(path)
		return err
	}
	defer f.Close()
	fi, err := f.Stat()
	if err != nil {
		return err
	}
	if fi.IsDir() {
		return errors.New("can't touch path")
	}
	return nil
}
