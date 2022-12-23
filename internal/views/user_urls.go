package views

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"time"

	"github.com/ervand7/urlshortener/internal/logger"
)

// UserURLs GET ("/api/user/urls")
func (server *Server) UserURLs(w http.ResponseWriter, r *http.Request) {
	userID, err := r.Cookie("userID")
	if err != nil {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	decodedUserID, err := hex.DecodeString(userID.Value)
	if err != nil {
		logger.Logger.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
	msg, err := json.Marshal(userURLs)
	if err != nil {
		logger.Logger.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	server.Write(msg, w)
}
