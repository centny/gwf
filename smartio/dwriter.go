package smartio

import (
	"fmt"
	"github.com/Centny/gwf/util"
	"io"
	"os"
	"path/filepath"
	"time"
)

// func SetFMode(fmod int) {
// 	fs_fmod = fmod
// }

type DateSwitchWriter struct {
	ws    string
	cfn   string
	NameF string
	F     *os.File
	io.Writer
	FMODE os.FileMode
}

func NewDateSwitchWriter(ws string, fm os.FileMode) *DateSwitchWriter {
	dsw := &DateSwitchWriter{}
	dsw.ws = ws
	dsw.cfn = ""
	dsw.F = nil
	dsw.FMODE = fm
	return dsw
}
func NewDateSwitchWriter2(ws string) *DateSwitchWriter {
	return NewDateSwitchWriter(ws, os.ModePerm)
}

func (d *DateSwitchWriter) Write(p []byte) (int, error) {
	var fname = time.Now().Format("2006-01-02")
	if len(d.NameF) > 0 {
		fname = fmt.Sprintf(d.NameF, fname)
	}
	if d.cfn != fname {
		if d.F != nil {
			d.F.Close()
		}
		d.F = nil
	}
	//create new log writer
	if d.F == nil {
		err := d.reopen(fname)
		if err != nil {
			return 0, err
		}
	}
	l, err := d.F.Write(p)
	if err == nil {
		return l, err
	} else { //if writing error,try again.
		d.F.Close()
		d.F = nil
		fmt.Fprintln(LOG, "writing data error:"+err.Error())
		time.Sleep(time.Second)
		return d.Write(p)
	}
}
func (d *DateSwitchWriter) reopen(fname string) error {
	fpath := filepath.Join(d.ws, fname)
	err := util.FTouch2(fpath, d.FMODE)
	if err != nil {
		return err
	}
	f, err := os.OpenFile(fpath, os.O_WRONLY|os.O_APPEND, d.FMODE)
	d.cfn = fname
	d.F = f
	fmt.Fprintln(LOG, "open file:", fpath)
	return err
}
func (d *DateSwitchWriter) Close() {
	if d.F != nil {
		fmt.Fprintln(LOG, "close file:", d.FilePath())
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
