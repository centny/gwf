package routing

import (
	"net/http"
	"os"
	"time"
)

type Dir struct {
	http.Dir
}
type File struct {
	http.File
}

func (d Dir) Open(name string) (http.File, error) {
	rf, err := d.Dir.Open(name)
	if err != nil {
		return rf, err
	}
	return &File{
		File: rf,
	}, nil
}

func (f *File) Stat() (os.FileInfo, error) {
	d, err := f.File.Stat()
	if err != nil {
		return d, err
	}
	return &FileInfo{
		FileInfo: d,
	}, nil
}

type FileInfo struct {
	os.FileInfo
}

func (f *FileInfo) ModTime() time.Time {
	return time.Now()
}
