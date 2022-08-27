package main

import (
	"github.com/ervand7/urlshortener/internal/app/server"
	"github.com/ervand7/urlshortener/internal/app/utils"
)

func main() {
	utils.Logger.Info("============= server started =============")
	server.Run()
}
