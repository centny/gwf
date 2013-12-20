package smartio

import (
	"github.com/Centny/Cny4go/util"
	"io"
	"os"
	"path/filepath"
	"time"
)

type DateSwitchWriter struct {
	ws  string
	cfn string
	fw  *os.File
	io.Writer
}

func NewDateSwitchWriter(ws string) *DateSwitchWriter {
	dsw := &DateSwitchWriter{}
	dsw.ws = ws
	dsw.cfn = ""
	dsw.fw = nil
	return dsw
}

func (d *DateSwitchWriter) Write(p []byte) (n int, err error) {
	fname := time.Now().Format("2006-1-2.log")
	if d.cfn != fname {
		if d.fw != nil {
			d.fw.Close()
		}
		d.fw = nil
	}
	//create new log writer
	if d.fw == nil {
		fpath := filepath.Join(d.ws, fname)
		err := util.FTouch(fpath)
		if err != nil {
			return 0, err
		}
		f, err := os.OpenFile(fpath, os.O_RDWR|os.O_APPEND, os.ModePerm)
		if err != nil {
			return 0, err
		}
		d.fw = f
	}
	return d.fw.Write(p)
}
func (d *DateSwitchWriter) Close() {
	if d.fw != nil {
		d.fw.Close()
		d.fw = nil
	}
}
