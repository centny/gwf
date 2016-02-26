package smartio

import (
	"io"
	"os"
	"path/filepath"
	"time"
)

var LOG io.Writer = os.Stdout

var Stdout io.Writer = nil
var Stderr io.Writer = nil

func NewRedirect(ws, name_f string, bsize int, cdelay int64) (*os.File, *TimeFlushWriter, error) {
	if len(ws) > 0 {
		err := os.MkdirAll(ws, os.ModePerm)
		if err != nil {
			return nil, nil, err
		}
	}
	var tw *TimeFlushWriter
	r, w, err := os.Pipe()
	if err == nil {
		var sw = NewDateSwitchWriter2(ws)
		sw.NameF = name_f
		var tw = NewTimeWriter(sw, bsize, time.Duration(cdelay))
		go io.Copy(tw, r)
	}
	return w, tw, nil
}

func RedirectStdout(ws, name_f string, bsize int, cdelay int64) error {
	var w, tw, err = NewRedirect(ws, name_f, bsize, cdelay)
	if err == nil {
		os.Stdout, Stdout = w, tw
	}
	return err
}

func RedirectStdout2(path_f string, bsize int, cdelay int64) error {
	var ws, name_f = filepath.Split(path_f)
	return RedirectStdout(ws, name_f, bsize, cdelay)
}

func RedirectStdout3(path_f string) error {
	return RedirectStdout2(path_f, 1024, 3000)
}

func RedirectStderr(ws, name_f string, bsize int, cdelay int64) error {
	var w, tw, err = NewRedirect(ws, name_f, bsize, cdelay)
	if err == nil {
		os.Stderr, Stderr = w, tw
	}
	return err
}

func RedirectStderr2(path_f string, bsize int, cdelay int64) error {
	var ws, name_f = filepath.Split(path_f)
	return RedirectStderr(ws, name_f, bsize, cdelay)
}

func RedirectStderr3(path_f string) error {
	return RedirectStderr2(path_f, 1024, 3000)
}
