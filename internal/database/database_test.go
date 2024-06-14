package database

import (
	"testing"
	"time"
)

func TestParseConfig(t *testing.T) {
	t.Parallel()

	cfg := &Config{
		ConnectionURL: "postgres://user:pass@127.0.0.1:5432/db?sslmode=disable",
	}

	cfg, err := cfg.parseConfig()
	if err != nil {
		t.Fatalf("expected error to be nil; got %v", err)
	}

	wantMaxOpenConns := 10
	wantMaxIdleConns := 10
	wantConnMaxIdleTime := time.Minute * 5

	if cfg.MaxOpenConns != wantMaxOpenConns {
		t.Errorf("expected max-open-conns to be %d; got %d", wantMaxOpenConns, cfg.MaxOpenConns)
	}

	if cfg.MaxIdleConns != wantMaxIdleConns {
		t.Errorf("expected max-idle-conns to be %d; got %d", wantMaxIdleConns, cfg.MaxIdleConns)
	}

	if cfg.ConnMaxIdleTime != wantConnMaxIdleTime {
		t.Errorf("expected conn-max-idle-time to be %v; got %v", wantConnMaxIdleTime, cfg.ConnMaxIdleTime)
	}

	cfg.ConnectionURL = ""

	_, err = cfg.parseConfig()
	wantErr := "invalid database connection url"

	if err == nil || err.Error() != wantErr {
		t.Fatalf("expected error to be %v; got nil", wantErr)
	}
}
