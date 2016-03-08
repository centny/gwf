package smartio

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

var SysOut *os.File = os.Stdout
var SysErr *os.File = os.Stderr
var LOG io.Writer = os.Stdout

var Stdout *TimeFlushWriter = nil
var Stderr *TimeFlushWriter = nil

var ShowLog bool = false

func slog(format string, args ...interface{}) {
	if ShowLog {
		fmt.Fprintf(LOG, format, args...)
	}
}

func NewRedirect(ws, name_f string, bsize int, cdelay int64, sys *os.File) (*os.File, *TimeFlushWriter, error) {
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
		tw = NewTimeWriter(sw, bsize, time.Duration(cdelay))
		if sys == nil {
			go io.Copy(tw, r)
		} else {
			go io.Copy(io.MultiWriter(tw, sys), r)
		}
	}
	return w, tw, nil
}
func RedirectStdoutV(ws, name_f string, bsize int, cdelay int64, sys bool) error {
	var sys_out *os.File = nil
	if sys {
		sys_out = SysOut
	}
	var w, tw, err = NewRedirect(ws, name_f, bsize, cdelay, sys_out)
	if err == nil {
		os.Stdout, Stdout = w, tw
	}
	return err
}
func RedirectStdout(ws, name_f string, bsize int, cdelay int64) error {
	return RedirectStdoutV(ws, name_f, bsize, cdelay, true)
}
func RedirectStdout2(path_f string, bsize int, cdelay int64) error {
	var ws, name_f = filepath.Split(path_f)
	return RedirectStdout(ws, name_f, bsize, cdelay)
}

func RedirectStdout3(path_f string) error {
	return RedirectStdout2(path_f, 1024, 3000)
}

func RedirectStderrV(ws, name_f string, bsize int, cdelay int64, sys bool) error {
	var sys_err *os.File = nil
	if sys {
		sys_err = SysErr
	}
	var w, tw, err = NewRedirect(ws, name_f, bsize, cdelay, sys_err)
	if err == nil {
		os.Stderr, Stderr = w, tw
	}
	return err
}
func RedirectStderr(ws, name_f string, bsize int, cdelay int64) error {
	return RedirectStderrV(ws, name_f, bsize, cdelay, true)
}
func RedirectStderr2(path_f string, bsize int, cdelay int64) error {
	var ws, name_f = filepath.Split(path_f)
	return RedirectStderr(ws, name_f, bsize, cdelay)
}

func RedirectStderr3(path_f string) error {
	return RedirectStderr2(path_f, 1024, 3000)
}

func ResetStd() {
	if Stdout != nil {
		Stdout.Stop()
		if cs, ok := Stdout.sw.(io.Closer); ok {
			cs.Close()
		}
		Stdout = nil
	}
	if Stderr != nil {
		Stderr.Stop()
		if cs, ok := Stderr.sw.(io.Closer); ok {
			cs.Close()
		}
		Stderr = nil
	}
	os.Stdout, os.Stderr = SysOut, SysErr
}

func NewNamedWriter(ws, name_f string, bsize int, cdelay int64) (*DateSwitchWriter, *TimeFlushWriter) {
	var sw = NewDateSwitchWriter2(ws)
	sw.NameF = name_f
	return sw, NewTimeWriter(sw, bsize, time.Duration(cdelay))
}
