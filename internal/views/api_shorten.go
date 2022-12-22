package views

import (
	"context"
	"encoding/json"
	"github.com/ervand7/urlshortener/internal/controllers/algorithms"
	e "github.com/ervand7/urlshortener/internal/errors"
	"github.com/ervand7/urlshortener/internal/logger"
	"io"
	"net/http"
	"time"
)

// APIShortenURL POST ("/api/shorten")
func (server *Server) APIShortenURL(w http.ResponseWriter, r *http.Request) {
	defer server.CloseBody(r)
	type (
		ReqBody struct {
			URL string `json:"url"`
		}
		RespBody struct {
			Result string `json:"result"`
		}
	)

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if len(body) == 0 {
		http.Error(w, "param url not filled", http.StatusBadRequest)
		return
	}
	reqBody := ReqBody{}
	if err = json.Unmarshal(body, &reqBody); err != nil {
		logger.Logger.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	userID := server.GetOrCreateUserIDFromCookie(w, r)
	short := algorithms.GenerateShortURL()
	httpStatus := http.StatusCreated

	ctx, cancel := context.WithTimeout(r.Context(), ctxTime*time.Second)
	defer cancel()
	if err = server.Storage.Set(ctx, userID, short, reqBody.URL); err != nil {
		logger.Logger.Error(err.Error())
		if errData, ok := err.(*e.ShortAlreadyExistsError); ok {
			short = errData.Error()
			httpStatus = http.StatusConflict
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	respBody := RespBody{
		Result: short,
	}
	marshaledBody, err := json.Marshal(respBody)
	if err != nil {
		logger.Logger.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-type", "application/json")
	w.WriteHeader(httpStatus)
	server.Write(marshaledBody, w)
}
