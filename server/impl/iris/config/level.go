package config

import (
	"fmt"
	"strings"
)

const (
	LogLevelDisable LogLevel = "disable"
	LogLevelFatal   LogLevel = "fatal"
	LogLevelError   LogLevel = "error"
	LogLevelWarn    LogLevel = "warn"
	LogLevelInfo    LogLevel = "info"
	LogLevelDebug   LogLevel = "debug"
)

type LogLevel string

func (l *LogLevel) UnmarshalText(text []byte) error {
	if len(text) == 0 {
		*l = LogLevelInfo
		return nil
	}

	switch string(text) {
	case string(LogLevelDisable), strings.ToUpper(string(LogLevelDisable)):
		*l = LogLevelDisable
	case string(LogLevelFatal), strings.ToUpper(string(LogLevelFatal)):
		*l = LogLevelFatal
	case string(LogLevelError), strings.ToUpper(string(LogLevelError)):
		*l = LogLevelError
	case string(LogLevelWarn), strings.ToUpper(string(LogLevelWarn)):
		*l = LogLevelWarn
	case string(LogLevelInfo), strings.ToUpper(string(LogLevelInfo)):
		*l = LogLevelInfo
	case string(LogLevelDebug), strings.ToUpper(string(LogLevelDebug)):
		*l = LogLevelDebug
	default:
		return fmt.Errorf("invalid log level: %s", text)
	}
	return nil
}

func (l LogLevel) MarshalText() ([]byte, error) {
	return []byte(l), nil
}
