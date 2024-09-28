package main

import (
	"github.com/dlbarduzzi/bookshop/internal/logging"
)

func main() {
	log := logging.NewLoggerFromEnv()
	log.Info("Welcome to my bookshop!")
}
