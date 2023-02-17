package views

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/ervand7/urlshortener/internal/logger"
)

// Stats POST ("/api/internal/stats").
func (server *Server) Stats(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	stats, err := server.GetStats(ctx)
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
