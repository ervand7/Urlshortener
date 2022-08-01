package server

import (
	"github.com/ervand7/urlshortener/internal/app/storage"
	"github.com/ervand7/urlshortener/internal/app/views"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func Run() {
	st := make(storage.Storage, 0)

	router := mux.NewRouter()
	router.HandleFunc("/", views.UrlShorten(st)).Methods("POST")
	router.HandleFunc("/{id:[a-zA-Z]+}", views.UrlGet(st)).Methods("GET")

	http.Handle("/", router)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
