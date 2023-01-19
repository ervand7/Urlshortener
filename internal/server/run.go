package server

import (
	"log"
	"net/http"

	"golang.org/x/crypto/acme/autocert"

	"github.com/ervand7/urlshortener/internal/config"
)

func runSecure(router http.Handler) {
	manager := &autocert.Manager{
		Cache:      autocert.DirCache("cache-dir"),
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist(),
	}
	server := &http.Server{
		Addr:      config.GetServerAddress(),
		Handler:   router,
		TLSConfig: manager.TLSConfig(),
	}
	log.Fatal(
		server.ListenAndServeTLS("", ""),
	)
}

// Run launches app.
func Run() {
	router := newRouter()
	if config.GetEnableHTTPS() == "true" {
		runSecure(router)
	} else {
		log.Fatal(
			http.ListenAndServe(
				config.GetServerAddress(),
				router,
			),
		)
	}
}
