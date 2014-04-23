package routing

import (
	"net/http"
	"os"
	"regexp"
	"time"
)

type Dir struct {
	http.Dir
	Inc []*regexp.Regexp
}
type File struct {
	http.File
}
type FileInfo struct {
	os.FileInfo
}

func (d *Dir) Add(m *regexp.Regexp) {
	d.Inc = append(d.Inc, m)
}

func (d *Dir) Open(name string) (http.File, error) {
	rf, err := d.Dir.Open(name)
	if err != nil {
		return rf, err
	}
	for _, inc := range d.Inc {
		if inc.MatchString(name) {
			return &File{File: rf}, nil
		}
	}
	return rf, nil
}

func (f *File) Stat() (os.FileInfo, error) {
	d, err := f.File.Stat()
	return &FileInfo{FileInfo: d}, err
}

func (f *FileInfo) ModTime() time.Time {
	return time.Now()
}

func NewNoCacheDir(path string) *Dir {
	return &Dir{
		Dir: http.Dir(path),
		Inc: []*regexp.Regexp{},
	}
}

func NewAllNoCacheDir(path string) *Dir {
	return &Dir{
		Dir: http.Dir(path),
		Inc: []*regexp.Regexp{regexp.MustCompile("^.*$")},
	}
}
