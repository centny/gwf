package log

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/Centny/gwf/smartio"
)

type LogLevel int

const (
	DEBUG   LogLevel = 1
	INFO             = 2
	WARNING          = 3
	ERROR            = 4
	NONE             = 5
)

//default log flags.
const D_LOG_FLAGS int = log.LstdFlags | log.Lshortfile

type Log struct {
	level LogLevel
	*log.Logger
}

var dlog Log = NewDefaultLog(os.Stdout)
var dlev LogLevel = DEBUG

func SetWriter(w io.Writer) {
	SetWriterFlag(w, "", D_LOG_FLAGS)
}

func SetWriterFlag(w io.Writer, prefix string, flag int) {
	dlog = NewLog(w, prefix, flag)
	SetLevel(dlev)
}

func SetLevel(level LogLevel) {
	dlev = level
	dlog.SetLevel(level)
}

func SetSLevel(level string) {
	var dlev LogLevel
	switch level {
	case "DEBUG":
		dlev = DEBUG
		break
	case "INFO":
		dlev = INFO
		break
	case "WARNING":
		dlev = WARNING
		break
	case "ERROR":
		dlev = ERROR
		break
	default:
		dlev = INFO
		fmt.Errorf("unknow log level")
		break
	}
	SetLevel(dlev)
}

func D(format string, v ...interface{}) {
	dlog.log(DEBUG, "[D] %s", fmt.Sprintf(format, v...))
}
func D_(d int, format string, v ...interface{}) {
	dlog.log_(DEBUG, d, "[D] %s", fmt.Sprintf(format, v...))
}

func I(format string, v ...interface{}) {
	dlog.log(INFO, "[I] %s", fmt.Sprintf(format, v...))
}

func W(format string, v ...interface{}) {
	dlog.log(WARNING, "[W] %s", fmt.Sprintf(format, v...))
}

func E(format string, v ...interface{}) {
	dlog.log(ERROR, "[E] %s", fmt.Sprintf(format, v...))
}
func NewDefaultLog(w io.Writer) Log {
	l := Log{}
	l.level = DEBUG
	l.Logger = log.New(w, "", D_LOG_FLAGS)
	return l
}

func NewLog(w io.Writer, prefix string, flag int) Log {
	l := Log{}
	l.level = DEBUG
	l.Logger = log.New(w, prefix, flag)
	return l
}
func (t *Log) SetLevel(l LogLevel) {
	t.level = l
}
func (t *Log) D(format string, v ...interface{}) {
	t.log(DEBUG, "[D] %s", fmt.Sprintf(format, v...))
}
func (t *Log) D_(d int, format string, v ...interface{}) {
	t.log_(DEBUG, d, "[D] %s", fmt.Sprintf(format, v...))
}

func (t *Log) I(format string, v ...interface{}) {
	t.log(INFO, "[I] %s", fmt.Sprintf(format, v...))
}

func (t *Log) W(format string, v ...interface{}) {
	t.log(WARNING, "[W] %s", fmt.Sprintf(format, v...))
}

func (t *Log) E(format string, v ...interface{}) {
	t.log(ERROR, "[E] %s", fmt.Sprintf(format, v...))
}

func (t *Log) log(l LogLevel, format string, v ...interface{}) {
	if l < t.level {
		return
	}
	t.Output(3, fmt.Sprintf(format, v...)+"\n")
}

func (t *Log) log_(l LogLevel, d int, format string, v ...interface{}) {
	if l < t.level {
		return
	}
	t.Output(d+3, fmt.Sprintf(format, v...)+"\n")
}

func Redirect(out_l, err_l string) {
	RedirectV(out_l, err_l, true)
}

func RedirectV(out_l, err_l string, sys bool) {
	// var out_l = fcfg.Val2("out_l", "")
	// var err_l = fcfg.Val2("err_l", "")
	if len(out_l) < 1 && len(err_l) < 1 {
		return
	}
	fmt.Printf("redirect stdout to file(%v) and stderr to file(%v)\n", out_l, err_l)
	if len(out_l) > 0 {
		smartio.RedirectStdout4(out_l, sys)
	}
	if len(err_l) > 0 {
		smartio.RedirectStderr4(err_l, sys)
	}
	SetWriter(os.Stdout)
}

func ResetStd() {
	smartio.ResetStd()
}

func Flush() {
	if smartio.Stdout != nil {
		smartio.Stdout.Flush()
	}
	if smartio.Stderr != nil {
		smartio.Stderr.Flush()
	}
}
