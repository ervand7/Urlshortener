package main

import (
	"log"
	"net/http"
	_ "net/http/pprof"

	"github.com/ervand7/urlshortener/internal/logger"
	"github.com/ervand7/urlshortener/internal/server"
)

func main() {
	go func() {
		log.Println(http.ListenAndServe(":6060", nil))
	}()

	logger.Logger.Info("server started")
	server.Run()
}
