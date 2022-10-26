package views

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"github.com/ervand7/urlshortener/internal/app/controllers/algorithms"
	"github.com/ervand7/urlshortener/internal/app/logger"
	"net/http"
	"time"
)

// UserURLsDelete delete ("/api/user/urls")
func (server *Server) UserURLsDelete(w http.ResponseWriter, r *http.Request) {
	defer server.CloseBody(r)
	userID, err := r.Cookie("userID")
	if err != nil {
		logger.Logger.Error(err.Error())
		w.WriteHeader(http.StatusNoContent)
		return
	}
	decodedUserID, err := hex.DecodeString(userID.Value)
	if err != nil {
		logger.Logger.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var (
		urlsFromRequest []string
		userUrlsFromDB  []string
	)
	if err = json.NewDecoder(r.Body).Decode(&urlsFromRequest); err != nil {
		logger.Logger.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), ctxTime*time.Second)
	defer cancel()
	userURLs, err := server.Storage.GetUserURLs(ctx, string(decodedUserID))
	if err != nil {
		logger.Logger.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	for _, val := range userURLs {
		userUrlsFromDB = append(userUrlsFromDB, val["short_url"])
	}
	if !algorithms.Issubset(userUrlsFromDB, urlsFromRequest) {
		logger.Logger.Warn("user can delete only his own urls")
	}

	algorithms.PrepareShortened(urlsFromRequest)
	go func() {
		server.Storage.DeleteUserURLs(urlsFromRequest)
	}()
	w.WriteHeader(http.StatusAccepted)
}
