package server

import (
	"log"
	"net/http"
)

func Run() {
	router := newRouter()
	log.Fatal(http.ListenAndServe(":8080", router))
}
