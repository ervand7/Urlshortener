package views

import (
	"github.com/ervand7/urlshortener/internal/app/database"
	"github.com/ervand7/urlshortener/internal/app/utils"
	"net/http"
)

// DBPing GET ("/ping")
func (server *Server) DBPing(w http.ResponseWriter, r *http.Request) {
	db := database.Database{}
	defer db.Close()

	err := db.StartConn()
	if err != nil {
		utils.Logger.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = db.Conn.Ping()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		server.Write([]byte(err.Error()), w)
	} else {
		w.WriteHeader(http.StatusOK)
	}
}
