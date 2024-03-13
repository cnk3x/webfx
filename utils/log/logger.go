package log

import (
	"io"
)

type Logger interface {
	Output(callDepth int, level, format string, v ...any)
	Printf(format string, v ...any)
	Debugf(format string, v ...any)
	Infof(format string, v ...any)
	Warnf(format string, v ...any)
	Fatalf(format string, v ...any)
	SetLevel(level string)
	SetOutput(w io.Writer)
}
