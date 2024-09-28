package logging

import (
	"fmt"
	"log/slog"
	"os"
	"strings"
)

func NewLogger(dev bool, level string) *slog.Logger {
	var log *slog.Logger

	if dev {
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level:       getLogLevel(level),
			AddSource:   true,
			ReplaceAttr: replaceAttr(),
		}))
	} else {
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level:       getLogLevel(level),
			AddSource:   true,
			ReplaceAttr: replaceAttr(),
		}))
	}

	return log
}

func NewLoggerFromEnv() *slog.Logger {
	dev := strings.TrimSpace(strings.ToLower(os.Getenv("LOG_MODE"))) == "dev"
	level := strings.TrimSpace(strings.ToLower(os.Getenv("LOG_LEVEL")))
	return NewLogger(dev, level)
}

type slogAttr func(groups []string, attr slog.Attr) slog.Attr

func replaceAttr() slogAttr {
	return func(groups []string, attr slog.Attr) slog.Attr {
		if attr.Key == slog.TimeKey {
			attr.Key = "time"
			attr.Value = slog.TimeValue(attr.Value.Time().UTC())
		}
		if attr.Key == slog.MessageKey {
			attr.Key = "msg"
		}
		if attr.Key == slog.SourceKey {
			source := attr.Value.Any().(*slog.Source)
			attr.Key = "caller"
			attr.Value = slog.StringValue(fmt.Sprintf("%s:%d", source.Function, source.Line))
		}
		return attr
	}
}

const (
	levelDebug = "DEBUG"
	levelInfo  = "INFO"
	levelWarn  = "WARN"
	levelError = "ERROR"
)

// getLogLevel converts parameter string to the appropriate slog level value.
func getLogLevel(level string) slog.Level {
	switch strings.ToUpper(strings.TrimSpace(level)) {
	case levelDebug:
		return slog.LevelDebug
	case levelInfo:
		return slog.LevelInfo
	case levelWarn:
		return slog.LevelWarn
	case levelError:
		return slog.LevelError
	}
	return slog.LevelInfo
}
