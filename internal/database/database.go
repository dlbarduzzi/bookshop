package database

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"
	"time"

	_ "github.com/lib/pq"
)

type Config struct {
	ConnectionURL   string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxIdleTime time.Duration
}

func NewDatabase(config *Config) (*sql.DB, error) {
	cfg, err := config.parse()
	if err != nil {
		return nil, err
	}

	db, err := sql.Open("postgres", cfg.ConnectionURL)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetConnMaxIdleTime(cfg.ConnMaxIdleTime)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return nil, fmt.Errorf("database ping connection failed: %v", err)
	}

	return db, nil
}

func (c *Config) parse() (*Config, error) {
	_, err := url.ParseRequestURI(c.ConnectionURL)
	if err != nil {
		// We are not returning the error here on purpose. If the connection url
		// is invalid, and the connection url has database username and password
		// credentials, the error message could potentially leak these credentials.
		return nil, fmt.Errorf("invalid database connection url")
	}

	if c.MaxOpenConns < 1 {
		c.MaxOpenConns = 10
	}

	if c.MaxIdleConns < 1 {
		c.MaxIdleConns = 10
	}

	if c.ConnMaxIdleTime < time.Minute {
		c.ConnMaxIdleTime = time.Minute * 5
	}

	return c, nil
}
