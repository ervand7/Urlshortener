package views

import (
	"context"
	"encoding/json"
	g "github.com/ervand7/urlshortener/internal/app/controllers/generatedata"
	"github.com/ervand7/urlshortener/internal/app/logger"
	"github.com/ervand7/urlshortener/internal/app/models"
	"io"
	"net/http"
	"time"
)

// APIShortenBatch POST ("/api/shorten/batch")
func (server *Server) APIShortenBatch(w http.ResponseWriter, r *http.Request) {
	defer server.CloseBody(r)
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if len(body) == 0 {
		http.Error(w, "body is empty", http.StatusBadRequest)
		return
	}

	type (
		ReqPair struct {
			CorrelationID string `json:"correlation_id"`
			OriginURL     string `json:"original_url"`
		}
		RespPair struct {
			CorrelationID string `json:"correlation_id"`
			ShortURL      string `json:"short_url"`
		}
	)
	var (
		reqPairs  []ReqPair
		respPairs []RespPair
		dbEntries []models.Entry
	)
	if err := json.Unmarshal(body, &reqPairs); err != nil {
		logger.Logger.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	userID := server.GetOrCreateUserIDFromCookie(w, r)
	for _, val := range reqPairs {
		short := g.ShortenURL()
		respPair := RespPair{CorrelationID: val.CorrelationID, ShortURL: short}
		respPairs = append(respPairs, respPair)

		entry := models.Entry{UserID: userID, Short: short, Origin: val.OriginURL}
		dbEntries = append(dbEntries, entry)
	}
	marshaledBody, err := json.Marshal(respPairs)
	if err != nil {
		logger.Logger.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), ctxTime*time.Second)
	defer cancel()
	if err = server.Storage.SetMany(ctx, dbEntries); err != nil {
		logger.Logger.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-type", "application/json")
	w.WriteHeader(http.StatusCreated)
	server.Write(marshaledBody, w)
}
