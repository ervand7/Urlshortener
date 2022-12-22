package views

import (
	"github.com/ervand7/urlshortener/internal/controllers/storage/db_storage"
	"github.com/ervand7/urlshortener/internal/logger"
	"net/http"
	"os"
)

// PingDB GET ("/ping")
func (server *Server) PingDB(w http.ResponseWriter, _ *http.Request) {
	dsn := "user=ervand password=ervand dbname=urlshortener host=localhost port=5432 sslmode=disable"
	err := os.Setenv("DATABASE_DSN", dsn)
	if err != nil {
		logger.Logger.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	db := db_storage.Database{}
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
