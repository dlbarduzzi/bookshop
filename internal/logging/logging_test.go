package logging

import (
	"context"
	"log/slog"
	"testing"
	"time"
)

func TestNewLogget(t *testing.T) {
	testCases := []struct {
		mode  string
		level string
	}{
		{mode: "dev", level: "debug"},
		{mode: "prod", level: "debug"},
	}

	for _, tc := range testCases {
		t.Run(tc.mode, func(t *testing.T) {
			t.Parallel()
			logger := NewLogger(tc.mode, tc.level)
			if logger == nil {
				t.Fatal("expected logger not to be nil")
			}
		})
	}
}

func TestDefaultLogger(t *testing.T) {
	loggerA := DefaultLogger()
	if loggerA == nil {
		t.Fatal("expected logger not to be nil")
	}

	loggerB := DefaultLogger()
	if loggerB == nil {
		t.Fatal("expected logger not to be nil")
	}

	if loggerA != loggerB {
		t.Errorf("expected logger %#v to be equal %#v", loggerA, loggerB)
	}
}

func TestLoggerContext(t *testing.T) {
	ctx := context.Background()

	loggerA := LoggerFromContext(ctx)
	if loggerA == nil {
		t.Fatal("expected logger not to be nil")
	}

	ctx = LoggerWithContext(ctx, loggerA)

	loggerB := LoggerFromContext(ctx)
	if loggerA != loggerB {
		t.Errorf("expected logger %#v to be equal %#v", loggerA, loggerB)
	}
}

func TestGetLogLevel(t *testing.T) {
	testCases := []struct {
		mode      string
		level     string
		wantLevel slog.Level
	}{
		{mode: "empty", level: "", wantLevel: slog.LevelInfo},
		{mode: "invalid", level: "invalid", wantLevel: slog.LevelInfo},
		{mode: "debug", level: "DEBUG", wantLevel: slog.LevelDebug},
		{mode: "info", level: "info", wantLevel: slog.LevelInfo},
		{mode: "warn", level: "warn", wantLevel: slog.LevelWarn},
		{mode: "error", level: "error", wantLevel: slog.LevelError},
	}

	for _, tc := range testCases {
		t.Run(tc.mode, func(t *testing.T) {
			t.Parallel()
			level := getLogLevel(tc.level)
			if level != tc.wantLevel {
				t.Errorf("expected logger level to be %v; got %v", tc.wantLevel, level)
			}
		})
	}
}

func TestReplaceAttr(t *testing.T) {
	now := time.Now()

	attr := slog.Attr{
		Key:   slog.TimeKey,
		Value: slog.TimeValue(now),
	}

	fn := replaceAttr("dev")
	res := fn(nil, attr)

	gotTime := res.Value.String()
	wantTime := now.UTC().String()

	if gotTime != wantTime {
		t.Errorf("expected time value to be %s; got %s", wantTime, gotTime)
	}

	attr = slog.Attr{
		Key:   slog.MessageKey,
		Value: slog.StringValue("test-message"),
	}

	fn = replaceAttr("dev")
	res = fn(nil, attr)

	gotKey := res.Key
	wantKey := "message"

	if gotKey != wantKey {
		t.Errorf("expected message key to be %s; got %s", wantKey, gotKey)
	}

	src := &slog.Source{Function: "app.main", File: "/path/to/file", Line: 12}

	attr = slog.Attr{
		Key:   slog.SourceKey,
		Value: slog.AnyValue(src),
	}

	fn = replaceAttr("dev")
	res = fn(nil, attr)

	gotKey = res.Key
	wantKey = "caller"

	if gotKey != wantKey {
		t.Errorf("expected source key to be %s; got %s", wantKey, gotKey)
	}

	gotSourcePath := res.Value
	wantSourcePath := "app.main:12"

	if gotSourcePath.String() != wantSourcePath {
		t.Errorf("expected source value to be %s; got %s", wantSourcePath, gotSourcePath)
	}

	fn = replaceAttr("prod")
	res = fn(nil, attr)

	gotSourcePath = res.Value
	wantSourcePath = "/path/to/file:12"

	if gotSourcePath.String() != wantSourcePath {
		t.Errorf("expected source value to be %s; got %s", wantSourcePath, gotSourcePath)
	}
}
