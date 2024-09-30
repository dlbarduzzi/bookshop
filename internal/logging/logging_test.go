package logging

import (
	"context"
	"log/slog"
	"testing"
	"time"
)

func TestNewLogger(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		level string
	}{
		{
			name:  "dev",
			level: "debug",
		},
		{
			name:  "prod",
			level: "info",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			log := NewLogger(tt.name == "dev", tt.level)
			if log == nil {
				t.Fatal("expected logger not to be nil")
			}
		})
	}
}

func TestDefaultLogger(t *testing.T) {
	t.Parallel()

	log1 := DefaultLogger()
	if log1 == nil {
		t.Fatal("expected logger not to be nil")
	}

	log2 := DefaultLogger()
	if log2 == nil {
		t.Fatal("expected logger not to be nil")
	}

	if log1 != log2 {
		t.Errorf("expected logger %#v to be equal %#v", log1, log2)
	}
}

func TestLoggerContext(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	log1 := LoggerFromContext(ctx)
	if log1 == nil {
		t.Fatal("expected logger not to be nil")
	}

	ctx = LoggerWithContext(ctx, log1)

	log2 := LoggerFromContext(ctx)
	if log1 != log2 {
		t.Errorf("expected logger %#v to be equal %#v", log1, log2)
	}
}

func TestGetLogLevel(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		level     string
		wantLevel slog.Level
	}{
		{
			name:      "empty",
			level:     "",
			wantLevel: slog.LevelInfo,
		},
		{
			name:      "invalid",
			level:     "invalid",
			wantLevel: slog.LevelInfo,
		},
		{
			name:      "debug",
			level:     "debug",
			wantLevel: slog.LevelDebug,
		},
		{
			name:      "info",
			level:     "info",
			wantLevel: slog.LevelInfo,
		},
		{
			name:      "warn",
			level:     "warn",
			wantLevel: slog.LevelWarn,
		},
		{
			name:      "error",
			level:     "error",
			wantLevel: slog.LevelError,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			level := getLogLevel(tt.level)
			if level != tt.wantLevel {
				t.Errorf("expected logger level to be %v; got %v", tt.wantLevel, level)
			}
		})
	}
}

func TestGetReplaceAttr(t *testing.T) {
	t.Parallel()

	tm := time.Now()
	sr := &slog.Source{Function: "main.main", File: "/path/to/file", Line: 12}

	tests := []struct {
		name      string
		nano      bool
		logKey    string
		logValue  slog.Value
		wantKey   string
		wantValue string
	}{
		{
			name:      "time",
			logKey:    slog.TimeKey,
			logValue:  slog.TimeValue(tm),
			wantValue: tm.UTC().String(),
		},
		{
			name:    "message",
			logKey:  slog.MessageKey,
			wantKey: "msg",
		},
		{
			name:      "source",
			logKey:    slog.SourceKey,
			logValue:  slog.AnyValue(sr),
			wantKey:   "caller",
			wantValue: "/path/to/file:12",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			fn := replaceAttr()

			attr := slog.Attr{
				Key:   tt.logKey,
				Value: tt.logValue,
			}

			resp := fn(nil, attr)

			if tt.logKey == slog.TimeKey {
				got := resp.Value.String()
				if got != tt.wantValue {
					t.Errorf("expected time to be %s; got %s", tt.wantValue, got)
				}
			}

			if tt.logKey == slog.MessageKey {
				got := resp.Key
				if got != tt.wantKey {
					t.Errorf("expected message key to be %s; got %s", tt.wantKey, got)
				}
			}

			if tt.logKey == slog.SourceKey {
				gotKey := resp.Key
				if gotKey != tt.wantKey {
					t.Errorf("expected source key to be %s; got %s", tt.wantKey, gotKey)
				}
				gotValue := resp.Value
				if gotValue.String() != tt.wantValue {
					t.Errorf("expected source value to be %s; got %s", tt.wantValue, gotValue)
				}
			}
		})
	}
}
