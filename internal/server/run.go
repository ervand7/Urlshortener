package server

import (
	"log"
	"net/http"

	"github.com/ervand7/urlshortener/internal/config"
)

func Run() {
	router := newRouter()
	log.Fatal(
		http.ListenAndServe(
			config.GetServerAddress(),
			router,
		),
	)
}
