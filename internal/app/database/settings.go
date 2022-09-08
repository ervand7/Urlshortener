package database

import (
	"github.com/ervand7/urlshortener/internal/app/config"
	"github.com/ervand7/urlshortener/internal/app/utils"
	"os"
	"os/signal"
	"syscall"
)

func ManageDB() {
	if config.GetConfig().DatabaseDSN == "" {
		return
	}

	err := DB.ConnStart()
	if err != nil {
		utils.Logger.Fatal(err.Error())
	}

	DB.SetConnPool()

	err = DB.CreateAll()
	if err != nil {
		utils.Logger.Fatal(err.Error())
	}

	ch := make(chan os.Signal, 3)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		<-ch
		signal.Stop(ch)
		err = DB.ConnClose()
		if err != nil {
			utils.Logger.Error(err.Error())
		} else {
			utils.Logger.Info("Connection to DB was closed")
		}
		os.Exit(0)
	}()
}
