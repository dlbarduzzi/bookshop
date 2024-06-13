package logging

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"sync"
	"time"
)

// contextKey is the logger string type used to avoid collisions.
type contextKey string

// loggerKey identifies the logger value stored in the context.
const loggerKey = contextKey("logger")

var (
	// _defaultLogger is the default logger that should be initialized only once per package.
	_defaultLogger     *slog.Logger
	_defaultLoggerOnce sync.Once
)

// NewLogger creates and returns a new instance of logger with custom configurations.
func NewLogger(dev bool, level string) *slog.Logger {
	var log *slog.Logger

	if dev {
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level:       getLogLevel(level),
			AddSource:   false,
			ReplaceAttr: getReplaceAttr(false),
		}))
	} else {
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level:       getLogLevel(level),
			AddSource:   true,
			ReplaceAttr: getReplaceAttr(true),
		}))
	}

	return log
}

type slogAttr func(groups []string, attr slog.Attr) slog.Attr

// getReplaceAttr returns the slog attributes function with custom key and/or values pairs defined.
func getReplaceAttr(nano bool) slogAttr {
	return func(groups []string, attr slog.Attr) slog.Attr {
		if attr.Key == slog.TimeKey {
			attr.Key = "time"
			attr.Value = slog.StringValue(getTimeFormat(attr.Value.Time(), nano))
		}
		if attr.Key == slog.MessageKey {
			attr.Key = "message"
		}
		if attr.Key == slog.SourceKey {
			source := attr.Value.Any().(*slog.Source)
			attr.Key = "caller"
			attr.Value = slog.StringValue(fmt.Sprintf("%s:%d", source.File, source.Line))
		}
		return attr
	}
}

// NewLoggerFromEnv creates a new logger by evaluating a few environment variables to
// determine some of the custom configurations available to be set.
func NewLoggerFromEnv() *slog.Logger {
	dev := strings.ToLower(strings.TrimSpace(os.Getenv("LOG_MODE"))) == "dev"
	level := strings.ToLower(strings.TrimSpace(os.Getenv("LOG_LEVEL")))
	return NewLogger(dev, level)
}

// DefaultLogger returns the default logger for the package.
func DefaultLogger() *slog.Logger {
	_defaultLoggerOnce.Do(func() {
		_defaultLogger = NewLoggerFromEnv()
	})
	return _defaultLogger
}

// LoggerFromContext returns the logger stored in the context. If no logger exists,
// a default logger is returned.
func LoggerFromContext(ctx context.Context) *slog.Logger {
	if log, ok := ctx.Value(loggerKey).(*slog.Logger); ok {
		return log
	}
	return DefaultLogger()
}

// LoggerWithContext create a new context with the provided logger attached.
func LoggerWithContext(ctx context.Context, logger *slog.Logger) context.Context {
	return context.WithValue(ctx, loggerKey, logger)
}

const (
	_levelDebug = "DEBUG"
	_levelInfo  = "INFO"
	_levelWarn  = "WARN"
	_levelError = "ERROR"
)

// getLogLevel converts parameter string to the appropriate slog level value.
func getLogLevel(level string) slog.Level {
	switch strings.ToUpper(strings.TrimSpace(level)) {
	case _levelDebug:
		return slog.LevelDebug
	case _levelInfo:
		return slog.LevelInfo
	case _levelWarn:
		return slog.LevelWarn
	case _levelError:
		return slog.LevelError
	}
	return slog.LevelInfo
}

// getTimeFormate returns time with or without nanoseconds based on given parameter.
func getTimeFormat(t time.Time, nano bool) string {
	if nano {
		return t.Format(time.RFC3339Nano)
	} else {
		return t.Format(time.RFC3339)
	}
}
