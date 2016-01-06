package tutil

import (
	"fmt"
	"github.com/Centny/gwf/log"
	"io"
	"os"
	"path/filepath"
	"sync/atomic"
)

type FPerf struct {
	Path  string
	Clear bool
}

func NewFPerf(path string) *FPerf {
	return &FPerf{
		Path:  path,
		Clear: true,
	}
}

func (f *FPerf) Write(name string, bs int64, count int) error {
	buf := make([]byte, bs)
	tf, err := os.OpenFile(filepath.Join(f.Path, name), os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		return err
	}
	defer tf.Close()
	for i := 0; i < count; i++ {
		_, err = tf.Write(buf)
		if err != nil {
			return err
		}
	}
	return nil
}

func (f *FPerf) Read(name string) (int64, error) {
	buf := make([]byte, 10240)
	tf, err := os.Open(filepath.Join(f.Path, name))
	if err != nil {
		return 0, err
	}
	defer tf.Close()
	var readed int64 = 0
	var tr int = 0
	for {
		if tr, err = tf.Read(buf); err != nil {
			break
		} else {
			readed += int64(tr)
		}
	}
	if err == io.EOF {
		err = nil
	}
	return readed, err
}
func (f *FPerf) remove(name string) {
	os.Remove(filepath.Join(f.Path, name))
}

func (f *FPerf) Rw(name string, bs int64, count int) error {
	err := f.Write(name, bs, count)
	if err == nil {
		_, err = f.Read(name)
		return err
	} else {
		return err
	}

}

func (f *FPerf) Perf4MultiW(pref, logf string, fc, max int, bs int64, count int) (int64, error) {
	var terr error
	used, err := DoPerfV(fc, max, logf, func(v int) {
		tp := fmt.Sprintf("%v%v", pref, v)
		err := f.Write(tp, bs, count)
		if err != nil {
			terr = err
			log.E("TestMultiW error to (%v) on path(%v) error->%v", tp, f.Path, err.Error())
		}
		if f.Clear {
			f.remove(tp)
		}
	})
	if err == nil {
		err = terr
	}
	return used, err
}

func (f *FPerf) Perf4MultiRw(pref, logf string, fc, max int, bs int64, count int) (int64, error) {
	var terr error
	used, err := DoPerfV(fc, max, logf, func(v int) {
		tp := fmt.Sprintf("%v%v", pref, v)
		err := f.Rw(tp, bs, count)
		if err != nil {
			terr = err
			log.E("TestMultiRw error to (%v) on path(%v) error->%v", tp, f.Path, err.Error())
		}
		if f.Clear {
			f.remove(tp)
		}
	})
	if err == nil {
		err = terr
	}
	return used, err
}

func (f *FPerf) Perf4MultiR(pref, logf string, beg, end, max int) (int64, int64, error) {
	var terr error
	var readed int64 = 0
	used, err := DoPerfV(end-beg, max, logf, func(v int) {
		tp := fmt.Sprintf("%v%v", pref, beg+v)
		tr, err := f.Read(tp)
		atomic.AddInt64(&readed, tr)
		if err != nil {
			terr = err
			log.E("TestMultiR error to (%v) on path(%v) error->%v", tp, f.Path, err.Error())
		}
	})
	if err == nil {
		err = terr
	}
	return used, readed, err
}
