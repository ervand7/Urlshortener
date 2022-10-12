package views

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/ervand7/urlshortener/internal/app/config"
	"github.com/ervand7/urlshortener/internal/app/controllers/algorithms"
	g "github.com/ervand7/urlshortener/internal/app/controllers/generatedata"
	e "github.com/ervand7/urlshortener/internal/app/errors"
	"github.com/ervand7/urlshortener/internal/app/utils"
)

// URLShorten POST ("/")
func (server *Server) URLShorten(w http.ResponseWriter, r *http.Request) {
	defer server.CloseBody(r)
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	url := string(body)
	if url == "" {
		http.Error(w, "param url not filled", http.StatusBadRequest)
		return
	}

	userID := server.GetOrCreateUserIDFromCookie(w, r)
	short := g.ShortenURL()
	httpStatus := http.StatusCreated
	if err := server.SaveURL(userID, short, url, r); err != nil {
		if errData, ok := err.(*e.ShortAlreadyExistsError); ok {
			short = errData.Error()
			httpStatus = http.StatusConflict
		} else {
			utils.Logger.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	w.Header().Add("Content-type", "text/plain; charset=utf-8")
	w.WriteHeader(httpStatus)
	server.Write([]byte(short), w)
}

// URLGet GET ("/{id}")
func (server *Server) URLGet(w http.ResponseWriter, r *http.Request) {
	endpoint := r.URL.Path
	short := config.GetConfig().BaseURL + endpoint
	origin, err := server.GetOriginByShort(short, r)
	if err != nil {
		if _, ok := err.(*e.URLNotActiveError); ok {
			http.Error(w, err.Error(), http.StatusGone)
			return
		}
		utils.Logger.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if origin == "" {
		http.Error(w, "not exists", http.StatusBadRequest)
		return
	}
	http.Redirect(w, r, origin, http.StatusTemporaryRedirect)
}

// URLShortenJSON POST ("/api/shorten")
func (server *Server) URLShortenJSON(w http.ResponseWriter, r *http.Request) {
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
	if err = server.SaveURL(userID, short, reqBody.URL, r); err != nil {
		if errData, ok := err.(*e.ShortAlreadyExistsError); ok {
			short = errData.Error()
			httpStatus = http.StatusConflict
		} else {
			utils.Logger.Error(err.Error())
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

// URLUserAll GET ("/api/user/urls")
func (server *Server) URLUserAll(w http.ResponseWriter, r *http.Request) {
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

	userURLs, err := server.GetUserURLs(string(decodedUserID), r)
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

// URLShortenBatch POST ("/api/shorten/batch")
func (server *Server) URLShortenBatch(w http.ResponseWriter, r *http.Request) {
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
		dbEntries []utils.DBEntry
	)

	if err := json.Unmarshal(body, &reqPairs); err != nil {
		utils.Logger.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	userID := server.GetOrCreateUserIDFromCookie(w, r)
	for _, val := range reqPairs {
		short := g.ShortenURL()
		respPair := RespPair{CorrelationID: val.CorrelationID, ShortURL: short}
		respPairs = append(respPairs, respPair)

		entry := utils.DBEntry{UserID: userID, Short: short, Origin: val.OriginURL}
		dbEntries = append(dbEntries, entry)
	}

	marshaledBody, err := json.Marshal(respPairs)
	if err != nil {
		utils.Logger.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err = server.SaveURLs(dbEntries, r); err != nil {
		utils.Logger.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-type", "application/json")
	w.WriteHeader(http.StatusCreated)
	server.Write(marshaledBody, w)
}

// URLDeleteMany delete ("/api/user/urls")
func (server *Server) URLDeleteMany(w http.ResponseWriter, r *http.Request) {
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

	var urlsFromRequest []string
	if err = json.NewDecoder(r.Body).Decode(&urlsFromRequest); err != nil {
		utils.Logger.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userURLs, err := server.GetUserURLs(string(decodedUserID), r)
	var userUrlsFromDB []string
	for _, val := range userURLs {
		userUrlsFromDB = append(userUrlsFromDB, val["short_url"])
	}
	if !algorithms.IsIntersection(userUrlsFromDB, urlsFromRequest) {
		utils.Logger.Warn("user can delete only his own urls")
	}

	l := fmt.Sprintf("mass delete: url from request %v, db urls %v,  user %q", urlsFromRequest, userUrlsFromDB, string(decodedUserID))
	utils.Logger.Info(l)

	go func() {
		server.DeleteURLs(urlsFromRequest)
	}()
	w.WriteHeader(http.StatusAccepted)
}
