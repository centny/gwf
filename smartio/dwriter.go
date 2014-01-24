package smartio

import (
	"fmt"
	"github.com/Centny/Cny4go/util"
	"io"
	"os"
	"path/filepath"
	"time"
)

var fs_fmod int = 0755

func SetFMode(fmod int) {
	fs_fmod = fmod
}

type DateSwitchWriter struct {
	ws  string
	cfn string
	F   *os.File
	io.Writer
}

func NewDateSwitchWriter(ws string) *DateSwitchWriter {
	dsw := &DateSwitchWriter{}
	dsw.ws = ws
	dsw.cfn = ""
	dsw.F = nil
	return dsw
}

func (d *DateSwitchWriter) Write(p []byte) (n int, err error) {
	fname := time.Now().Format("2006-1-2.log")
	if d.cfn != fname {
		if d.F != nil {
			d.F.Close()
		}
		d.F = nil
	}
	//create new log writer
	if d.F == nil {
		fpath := filepath.Join(d.ws, fname)
		err := util.FTouch(fpath)
		if err != nil {
			return 0, err
		}
		f, err := os.OpenFile(fpath, os.O_RDWR|os.O_APPEND, os.ModePerm)
		if err != nil {
			return 0, err
		}
		d.cfn = fname
		d.F = f
		os.Chmod(fpath, os.FileMode(fs_fmod))
		fmt.Println("open file:" + fpath)
	}
	return d.F.Write(p)
}
func (d *DateSwitchWriter) Close() {
	if d.F != nil {
		fmt.Println("close file:", d.FilePath())
		d.F.Close()
		d.F = nil
	}
}
func (d *DateSwitchWriter) FilePath() string {
	if d.cfn == "" {
		return ""
	}
	return filepath.Join(d.ws, d.cfn)
}
