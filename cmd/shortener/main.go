package main

import (
	"github.com/ervand7/urlshortener/internal/logger"
	"github.com/ervand7/urlshortener/internal/server"
	"log"
	"net/http"
	//_ "net/http/pprof"
)

func main() {
	go func() {
		log.Println(http.ListenAndServe(":6060", nil))
	}()

	logger.Logger.Info("server started")
	server.Run()
}
