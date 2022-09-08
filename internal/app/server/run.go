package server

import (
	"github.com/ervand7/urlshortener/internal/app/config"
	"github.com/ervand7/urlshortener/internal/app/database"
	"log"
	"net/http"
)

func Run() {
	database.ManageDB()
	router := newRouter()
	log.Fatal(
		http.ListenAndServe(
			config.GetConfig().ServerAddress,
			router,
		),
	)
}
