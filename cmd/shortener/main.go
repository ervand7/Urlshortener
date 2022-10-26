package main

import (
	"github.com/ervand7/urlshortener/internal/logger"
	"github.com/ervand7/urlshortener/internal/server"
)

func main() {
	logger.Logger.Info("server started")
	server.Run()
}
