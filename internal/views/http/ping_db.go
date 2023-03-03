package http

import (
	"net/http"
	"os"

	"github.com/ervand7/urlshortener/internal/controllers/storage/dbstorage"
	"github.com/ervand7/urlshortener/internal/logger"
)

// PingDB GET ("/ping").
func (server *Server) PingDB(w http.ResponseWriter, _ *http.Request) {
	dsn := "user=ervand password=ervand dbname=urlshortener host=localhost port=5432 sslmode=disable"
	err := os.Setenv("DATABASE_DSN", dsn)
	if err != nil {
		logger.Logger.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	db := dbstorage.Database{}
	defer func() {
		if err = os.Unsetenv("DATABASE_DSN"); err != nil {
			logger.Logger.Error(err.Error())
		}
		if err = db.ConnClose(); err != nil {
			logger.Logger.Error(err.Error())
		}
	}()
	err = db.ConnStart()
	if err != nil {
		logger.Logger.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = db.Conn.Ping()
	if err != nil {
		logger.Logger.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		server.Write([]byte(err.Error()), w)
	} else {
		w.WriteHeader(http.StatusOK)
	}
}
