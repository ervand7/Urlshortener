package main

import (
	"github.com/ervand7/urlshortener/internal/app/logger"
	"github.com/ervand7/urlshortener/internal/app/server"
)

func main() {
	logger.Logger.Info("server started")
	server.Run()
}
