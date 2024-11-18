package main

import (
	"github.com/dlbarduzzi/guestbook/internal/logging"
)

func main() {
	logger := logging.NewLoggerFromEnv()
	logger.Info("Hello")
}
