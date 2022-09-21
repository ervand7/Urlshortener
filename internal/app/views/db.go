package views

import (
	d "github.com/ervand7/urlshortener/internal/app/database"
	"github.com/ervand7/urlshortener/internal/app/utils"
	"net/http"
	"os"
)

// DBPing GET ("/ping")
func (server *Server) DBPing(w http.ResponseWriter, _ *http.Request) {
	err := os.Setenv("DATABASE_DSN", "user=ervand password=ervand dbname=urlshortener host=localhost port=5432 sslmode=disable")
	if err != nil {
		utils.Logger.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	db := d.Database{}
	defer func() {
		if err := os.Unsetenv("DATABASE_DSN"); err != nil {
			utils.Logger.Error(err.Error())
		}
		if err := db.ConnClose(); err != nil {
			utils.Logger.Error(err.Error())
		}
	}()
	err = db.ConnStart()
	if err != nil {
		utils.Logger.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = db.Conn.Ping()
	if err != nil {
		utils.Logger.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		server.Write([]byte(err.Error()), w)
	} else {
		w.WriteHeader(http.StatusOK)
	}
}
