// Package main is an entry point to the program.
// Run example:
// $ go build -ldflags="-X 'main.buildVersion=v1.0.0' -X 'main.buildDate=$(date +'%Y/%m/%d %H:%M:%S')'" -o shortener && ./shortener
package main

import (
	"fmt"

	"github.com/ervand7/urlshortener/internal/logger"
	"github.com/ervand7/urlshortener/internal/server"
)

var (
	buildVersion = "N/A"
	buildDate    = "N/A"
	buildCommit  = "N/A"
)

func main() {
	logger.Logger.Info(fmt.Sprintf("Build version: %s", buildVersion))
	logger.Logger.Info(fmt.Sprintf("Build date: %s", buildDate))
	logger.Logger.Info(fmt.Sprintf("Build commit: %s", buildCommit))
	logger.Logger.Info("server started")
	server.Run()
}
