package log

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

func Standard() Logger {
	return defaultLogger
}

type standard struct {
	prefix string
	level  int
	w      io.Writer
}

func (s *standard) SetLevel(level string) {
	if l := LevelInt(level); l != -1 {
		s.level = l
	}
}

func (s *standard) SetOutput(w io.Writer) {
	if w != nil {
		s.w = w
	}
}

func (s *standard) output(callDepth int, level int, format string, v ...any) {
	if level > s.level {
		return
	}
	var fh string
	if _, path, line, ok := runtime.Caller(callDepth + 1); ok {
		fh = fmt.Sprintf("%16s:%-3d", filepath.Base(path), line)
	}

	msg := fmt.Sprintf(
		"%s%s %5s %s %s\n",
		s.prefix,
		time.Now().Format("01/02 15:04:05"),
		levelString(level),
		fh,
		strings.TrimSpace(format),
	)

	fmt.Fprintf(s.w, msg, v...)
}

func (s *standard) Output(callDepth int, level, format string, v ...any) {
	s.output(callDepth+1, LevelInt(level), format, v...)
}

func (s *standard) Printf(format string, v ...any) {
	s.output(1, PrintLevel, format, v...)
}

func (s *standard) Debugf(format string, v ...any) {
	s.output(1, DebugLevel, format, v...)
}

func (s *standard) Infof(format string, v ...any) {
	s.output(1, InfoLevel, format, v...)
}

func (s *standard) Warnf(format string, v ...any) {
	s.output(1, WarnLevel, format, v...)
}

func (s *standard) Fatalf(format string, v ...any) {
	s.output(1, FatalLevel, format, v...)
	os.Exit(1)
}
