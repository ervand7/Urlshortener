// Package main is an entry point to the program.
// Run example:
// $ go build -ldflags="-X 'main.buildVersion=v1.0.0' -X 'main.buildDate=$(date +'%Y/%m/%d %H:%M:%S')'" -o shortener && ./shortener
package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"google.golang.org/grpc"

	"github.com/ervand7/urlshortener/internal/logger"
	g "github.com/ervand7/urlshortener/internal/server/grpc"
	h "github.com/ervand7/urlshortener/internal/server/http"
)

var (
	buildVersion = "N/A"
	buildDate    = "N/A"
	buildCommit  = "N/A"
)

func logBuildInfo() {
	logger.Logger.Info(fmt.Sprintf("Build version: %s", buildVersion))
	logger.Logger.Info(fmt.Sprintf("Build date: %s", buildDate))
	logger.Logger.Info(fmt.Sprintf("Build commit: %s", buildCommit))
}

func shutdownGraceful(httpRunner *http.Server, grpcRunner *grpc.Server) {
	logger.Logger.Info("Server shutdown gracefully")
	time.Sleep(time.Second)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := httpRunner.Shutdown(ctx); err != nil {
		logger.Logger.Fatal(err.Error())
	}
	grpcRunner.Stop()
}

func main() {
	logBuildInfo()
	termChan := make(chan os.Signal, 1)
	signal.Notify(termChan, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	httpServer := h.GetServer()
	go h.Run(httpServer)

	grpcServer := g.GetServer()
	go g.Run(grpcServer)

	<-termChan
	shutdownGraceful(httpServer, grpcServer)
}
