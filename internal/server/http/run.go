package http

import (
	"log"
	"net/http"

	"golang.org/x/crypto/acme/autocert"

	"github.com/ervand7/urlshortener/internal/config"
	"github.com/ervand7/urlshortener/internal/logger"
)

// GetServer creates new http.Server
func GetServer() *http.Server {
	router := NewRouter()
	runner := http.Server{Addr: config.GetServerAddress(), Handler: router}
	return &runner
}

// Run launches http.Server
func Run(runner *http.Server) {
	logger.Logger.Info("===== HTTP server started =====")
	if config.GetEnableHTTPS() == "true" {
		runSecure(runner)
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
