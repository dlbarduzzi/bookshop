package logging

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"sync"
)

// contextKey is the logger string type used to avoid collisions.
type contextKey string

// loggerKey identifies the logger value stored in the context.
const loggerKey = contextKey("logger")

var (
	// defaultLogger is the default logger that should be initialized only once per package.
	defaultLogger     *slog.Logger
	defaultLoggerOnce sync.Once
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

func DefaultLogger() *slog.Logger {
	defaultLoggerOnce.Do(func() {
		defaultLogger = NewLoggerFromEnv()
	})
	return defaultLogger
}

func LoggerFromContext(ctx context.Context) *slog.Logger {
	if log, ok := ctx.Value(loggerKey).(*slog.Logger); ok {
		return log
	}
	return DefaultLogger()
}

func LoggerWithContext(ctx context.Context, logger *slog.Logger) context.Context {
	return context.WithValue(ctx, loggerKey, logger)
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
			attr.Value = slog.StringValue(fmt.Sprintf("%s:%d", source.File, source.Line))
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
