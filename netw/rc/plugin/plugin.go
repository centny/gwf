package plugin

import (
	"github.com/Centny/gwf/log"
)

const (
	TS_NONE      = 0
	TS_NOT_START = 1
	TS_RUNNING   = 2
)

//ShowLog is setting to log level
var ShowLog = 0

func slog(format string, args ...interface{}) {
	if ShowLog > 0 {
		log.D_(1, format, args...)
	}
}
