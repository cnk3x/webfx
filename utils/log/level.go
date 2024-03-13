package log

import "strings"

// logger level
const (
	FatalLevel = 1 << iota
	WarnLevel
	InfoLevel
	DebugLevel
	PrintLevel = 0

	PRINT = "PRINT"
	FATAL = "FATAL"
	WARN  = "WARN"
	INFO  = "INFO"
	DEBUG = "DEBUG"
)

func LevelInt(level string) int {
	return levelInt(strings.ToUpper(level))
}

func LevelString(level int) string {
	return levelString(level)
}

func levelInt(level string) int {
	switch level {
	case DEBUG:
		return DebugLevel
	case INFO:
		return InfoLevel
	case WARN:
		return WarnLevel
	case FATAL:
		return FatalLevel
	case PRINT:
		return PrintLevel
	}
	return -1
}

func levelString(level int) string {
	switch level {
	case DebugLevel:
		return DEBUG
	case InfoLevel:
		return INFO
	case WarnLevel:
		return WARN
	case FatalLevel:
		return FATAL
	case PrintLevel:
		return PRINT
	}
	return ""
}
