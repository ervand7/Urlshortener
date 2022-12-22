package views

import (
	"context"
	"github.com/ervand7/urlshortener/internal/controllers/algorithms"
	e "github.com/ervand7/urlshortener/internal/errors"
	"github.com/ervand7/urlshortener/internal/logger"
	"io"
	"net/http"
	"time"
)

// ShortenURL POST ("/")
func (server *Server) ShortenURL(w http.ResponseWriter, r *http.Request) {
	defer server.CloseBody(r)
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	origin := string(body)
	if origin == "" {
		http.Error(w, "param url not filled", http.StatusBadRequest)
		return
	}

	userID := server.GetOrCreateUserIDFromCookie(w, r)
	short := algorithms.GenerateShortURL()
	httpStatus := http.StatusCreated

	ctx, cancel := context.WithTimeout(r.Context(), ctxTime*time.Second)
	defer cancel()
	if err = server.Storage.Set(ctx, userID, short, origin); err != nil {
		logger.Logger.Error(err.Error())
		if errData, ok := err.(*e.ShortAlreadyExistsError); ok {
			short = errData.Error()
			httpStatus = http.StatusConflict
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	w.Header().Add("Content-type", "text/plain; charset=utf-8")
	w.WriteHeader(httpStatus)
	server.Write([]byte(short), w)
}
