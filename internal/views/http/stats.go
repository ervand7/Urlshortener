package http

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/ervand7/urlshortener/internal/logger"
	"github.com/ervand7/urlshortener/internal/views"
)

// Stats POST ("/api/internal/stats").
func (server *Server) Stats(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	stats, err := views.GetStats(ctx, server.Storage)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(stats)
	if err != nil {
		logger.Logger.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	server.Write(data, w)
}
