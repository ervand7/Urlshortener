package server

import (
	"github.com/ervand7/urlshortener/internal/app/config"
	"log"
	"net/http"
)

func Run() {
	router := newRouter()
	log.Fatal(
		http.ListenAndServe(
			config.GetConfig().ServerAddress,
			router,
		),
	)
}
