package server

import (
	"github.com/ervand7/urlshortener/internal/app/views"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func Run() {
	rtr := mux.NewRouter()
	rtr.HandleFunc("/", views.UrlShorten).Methods("POST")
	rtr.HandleFunc("/{id:[0-9a-zA-Z]+}", views.UrlGet).Methods("GET")

	http.Handle("/", rtr)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
