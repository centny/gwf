package routing

import (
	"github.com/Centny/gwf/log"
	"net/http"
	"os"
	"regexp"
	"time"
)

type Dir struct {
	http.Dir
	Inc     []*regexp.Regexp
	ShowLog bool
}
type File struct {
	http.File
}
type FileInfo struct {
	os.FileInfo
}

func (d *Dir) log(f string, args ...interface{}) {
	if d.ShowLog {
		log.D(f, args...)
	}
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
			d.log("not cahce for path:%v", name)
			return &File{File: rf}, nil
		}
	}
	d.log("using normal file system for path:%v", name)
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
