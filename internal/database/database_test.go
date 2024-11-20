package database

import (
	"testing"
	"time"
)

func TestConfigParse(t *testing.T) {
	t.Parallel()

	config := &Config{}

	_, err := config.parse()
	wantErr := "invalid database connection url"

	if err == nil || err.Error() != wantErr {
		t.Fatalf("expected error to be %v; got %v", wantErr, err)
	}

	config = &Config{
		ConnectionURL: "postgres://user:pass@127.0.0.1:5432/db?sslmode=disable",
	}

	cfg, err := config.parse()
	if err != nil {
		t.Fatalf("expected error to be nil; got %v", err)
	}

	if cfg.MaxOpenConns != 10 {
		t.Errorf("expected max-open-conns to be %d; got %d", 10, cfg.MaxOpenConns)
	}

	if cfg.MaxIdleConns != 10 {
		t.Errorf("expected max-idle-conns to be %d; got %d", 10, cfg.MaxIdleConns)
	}

	idleTime := time.Minute * 5
	if cfg.ConnMaxIdleTime != idleTime {
		t.Errorf("expected conn-max-idle-time to be %v; got %v", idleTime, cfg.ConnMaxIdleTime)
	}
}
