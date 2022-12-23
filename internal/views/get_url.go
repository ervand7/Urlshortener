package views

import (
	"context"
	"database/sql"
	"net/http"
	"time"

	"github.com/ervand7/urlshortener/internal/config"
	e "github.com/ervand7/urlshortener/internal/errors"
	"github.com/ervand7/urlshortener/internal/logger"
)

// GetURL GET ("/{id}")
func (server *Server) GetURL(w http.ResponseWriter, r *http.Request) {
	endpoint := r.URL.Path
	short := config.GetConfig().BaseURL + endpoint

	ctx, cancel := context.WithTimeout(r.Context(), ctxTime*time.Second)
	defer cancel()
	origin, err := server.Storage.Get(ctx, short)
	if err != nil {
		logger.Logger.Error(err.Error())
		if _, ok := err.(*e.URLNotActiveError); ok {
			http.Error(w, err.Error(), http.StatusGone)
			return
		}
		if err.Error() == sql.ErrNoRows.Error() {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if origin == "" {
		http.Error(w, "not exists", http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, origin, http.StatusTemporaryRedirect)
}
