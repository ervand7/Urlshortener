package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang.org/x/crypto/acme/autocert"

	"github.com/ervand7/urlshortener/internal/config"
	"github.com/ervand7/urlshortener/internal/logger"
)

// Run launches app.
func Run() {
	router := newRouter()
	runner := http.Server{Addr: config.GetServerAddress(), Handler: router}
	go shutdownGraceful(&runner)

	if config.GetEnableHTTPS() == "true" {
		runSecure(&runner)
	} else {
		log.Fatal(
			runner.ListenAndServe(),
		)
	}
}

func runSecure(runner *http.Server) {
	manager := &autocert.Manager{
		Cache:      autocert.DirCache("cache-dir"),
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist(),
	}
	runner.TLSConfig = manager.TLSConfig()

	log.Fatal(
		runner.ListenAndServeTLS("", ""),
	)
}

func shutdownGraceful(runner *http.Server) {
	termChan := make(chan os.Signal, 1)
	signal.Notify(termChan, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	<-termChan
	logger.Logger.Info("Server shutdown gracefully")
	time.Sleep(time.Second)
	if err := runner.Shutdown(ctx); err != nil {
		logger.Logger.Fatal(err.Error())
	}
}
