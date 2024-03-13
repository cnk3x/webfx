package log

import "os"

var defaultLogger Logger = &standard{prefix: "", w: os.Stderr, level: InfoLevel}

var (
	Output    = defaultLogger.Output
	Printf    = defaultLogger.Printf
	Debugf    = defaultLogger.Debugf
	Infof     = defaultLogger.Infof
	Warnf     = defaultLogger.Warnf
	Fatalf    = defaultLogger.Fatalf
	SetLevel  = defaultLogger.SetLevel
	SetOutput = defaultLogger.SetOutput
)

func SetDefault(l Logger) {
	defaultLogger = l
}

func Default() Logger {
	return defaultLogger
}
