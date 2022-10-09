package views

import (
	"context"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"github.com/ervand7/urlshortener/internal/app/config"
	"github.com/ervand7/urlshortener/internal/app/controllers/algorithms"
	g "github.com/ervand7/urlshortener/internal/app/controllers/generatedata"
	e "github.com/ervand7/urlshortener/internal/app/errors"
	"github.com/ervand7/urlshortener/internal/app/utils"
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
	short := g.ShortenURL()
	httpStatus := http.StatusCreated

	ctx, cancel := context.WithTimeout(r.Context(), ctxTime*time.Second)
	defer cancel()
	if err := server.Storage.Set(ctx, userID, short, origin); err != nil {
		utils.Logger.Error(err.Error())
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

// GetURL GET ("/{id}")
func (server *Server) GetURL(w http.ResponseWriter, r *http.Request) {
	endpoint := r.URL.Path
	short := config.GetConfig().BaseURL + endpoint

	ctx, cancel := context.WithTimeout(r.Context(), ctxTime*time.Second)
	defer cancel()
	origin, err := server.Storage.Get(ctx, short)
	if err != nil {
		utils.Logger.Error(err.Error())
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
	if err := json.Unmarshal(body, &reqBody); err != nil {
		utils.Logger.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	userID := server.GetOrCreateUserIDFromCookie(w, r)
	short := g.ShortenURL()
	httpStatus := http.StatusCreated

	ctx, cancel := context.WithTimeout(r.Context(), ctxTime*time.Second)
	defer cancel()
	if err = server.Storage.Set(ctx, userID, short, reqBody.URL); err != nil {
		utils.Logger.Error(err.Error())
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
		utils.Logger.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-type", "application/json")
	w.WriteHeader(httpStatus)
	server.Write(marshaledBody, w)
}

// UserURLs GET ("/api/user/urls")
func (server *Server) UserURLs(w http.ResponseWriter, r *http.Request) {
	userID, err := r.Cookie("userID")
	if err != nil {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	decodedUserID, err := hex.DecodeString(userID.Value)
	if err != nil {
		utils.Logger.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), ctxTime*time.Second)
	defer cancel()
	userURLs, err := server.Storage.GetUserURLs(ctx, string(decodedUserID))
	if err != nil {
		utils.Logger.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	msg, err := json.Marshal(userURLs)
	if err != nil {
		utils.Logger.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	server.Write(msg, w)
}

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
	var (
		reqPairs  []utils.ReqPair
		respPairs []utils.RespPair
		dbEntries []utils.Entry
	)
	if err := json.Unmarshal(body, &reqPairs); err != nil {
		utils.Logger.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	userID := server.GetOrCreateUserIDFromCookie(w, r)
	for _, val := range reqPairs {
		short := g.ShortenURL()
		respPair := utils.RespPair{CorrelationID: val.CorrelationID, ShortURL: short}
		respPairs = append(respPairs, respPair)

		entry := utils.Entry{UserID: userID, Short: short, Origin: val.OriginURL}
		dbEntries = append(dbEntries, entry)
	}
	marshaledBody, err := json.Marshal(respPairs)
	if err != nil {
		utils.Logger.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), ctxTime*time.Second)
	defer cancel()
	if err = server.Storage.SetMany(ctx, dbEntries); err != nil {
		utils.Logger.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-type", "application/json")
	w.WriteHeader(http.StatusCreated)
	server.Write(marshaledBody, w)
}

// UserURLsDelete delete ("/api/user/urls")
func (server *Server) UserURLsDelete(w http.ResponseWriter, r *http.Request) {
	defer server.CloseBody(r)
	userID, err := r.Cookie("userID")
	if err != nil {
		utils.Logger.Error(err.Error())
		w.WriteHeader(http.StatusNoContent)
		return
	}
	decodedUserID, err := hex.DecodeString(userID.Value)
	if err != nil {
		utils.Logger.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var (
		urlsFromRequest []string
		userUrlsFromDB  []string
	)
	if err = json.NewDecoder(r.Body).Decode(&urlsFromRequest); err != nil {
		utils.Logger.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), ctxTime*time.Second)
	defer cancel()
	userURLs, err := server.Storage.GetUserURLs(ctx, string(decodedUserID))
	if err != nil {
		utils.Logger.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	for _, val := range userURLs {
		userUrlsFromDB = append(userUrlsFromDB, val["short_url"])
	}
	if !algorithms.Issubset(userUrlsFromDB, urlsFromRequest) {
		utils.Logger.Warn("user can delete only his own urls")
	}

	algorithms.PrepareShortened(urlsFromRequest)
	go func() {
		server.Storage.DeleteUserURLs(urlsFromRequest)
	}()
	w.WriteHeader(http.StatusAccepted)
}
