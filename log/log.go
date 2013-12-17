package log

import (
	"fmt"
	"io"
	"log"
	"os"
)

type LogLevel int

const (
	DEBUG   LogLevel = 1
	INFO             = 2
	WARNING          = 3
	ERROR            = 4
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

func D(v ...interface{}) {
	dlog.D(v...)
}

func I(v ...interface{}) {
	dlog.I(v...)
}

func W(v ...interface{}) {
	dlog.W(v...)
}

func E(v ...interface{}) {
	dlog.E(v...)
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
func (t *Log) D(v ...interface{}) {
	t.log(DEBUG, "[D]", fmt.Sprint(v...))
}

func (t *Log) I(v ...interface{}) {
	t.log(INFO, "[I]", fmt.Sprint(v...))
}

func (t *Log) W(v ...interface{}) {
	t.log(WARNING, "[W]", fmt.Sprint(v...))
}

func (t *Log) E(v ...interface{}) {
	t.log(ERROR, "[E]", fmt.Sprint(v...))
}

func (t *Log) log(l LogLevel, v ...interface{}) {
	if l < t.level {
		return
	}
	t.Output(3, fmt.Sprintln(v...))
}
