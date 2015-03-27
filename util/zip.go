package util

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func Zip(zf string, base string, fs ...string) error {
	zf_o, err := os.OpenFile(zf, os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		return err
	}
	defer zf_o.Close()
	w := zip.NewWriter(zf_o)
	for _, f := range fs {
		name := strings.TrimPrefix(f, base)
		fw, err := w.Create(name)
		if err != nil {
			return err
		}
		fr, err := os.Open(f)
		if err != nil {
			return err
		}
		_, err = io.Copy(fw, fr)
		if err != nil {
			return err
		}
	}
	return w.Close()
}
func Unzip(zf string, out string) error {
	r, err := zip.OpenReader(zf)
	if err != nil {
		return err
	}
	defer r.Close()
	for _, f := range r.File {
		src, err := f.Open()
		if err != nil {
			return err
		}
		sp := filepath.Join(out, f.Name)
		FTouch(sp)
		dst, err := os.OpenFile(sp, os.O_CREATE|os.O_WRONLY, os.ModePerm)
		if err != nil {
			src.Close()
			return err
		}
		_, err = io.Copy(dst, src)
		if err != nil {
			src.Close()
			dst.Close()
			return err
		}
		src.Close()
		dst.Close()
	}
	return nil
}
