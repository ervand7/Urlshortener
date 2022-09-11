package views

import (
	"encoding/hex"
	"encoding/json"
	"github.com/ervand7/urlshortener/internal/app/config"
	g "github.com/ervand7/urlshortener/internal/app/controllers/generatedata"
	"github.com/ervand7/urlshortener/internal/app/utils"
	"io"
	"net/http"
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
	if err := server.SaveURL(userID, short, url, r); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusCreated)
	server.Write([]byte(short), w)
}

// URLGet GET ("/{id}")
func (server *Server) URLGet(w http.ResponseWriter, r *http.Request) {
	endpoint := r.URL.Path
	short := config.GetConfig().BaseURL + endpoint
	origin, err := server.GetOriginByShort(short, r)
	if err != nil {
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
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	short := g.ShortenURL()
	respBody := RespBody{
		Result: short,
	}
	marshaledBody, err := json.Marshal(respBody)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	userID := server.GetOrCreateUserIDFromCookie(w, r)
	if err = server.SaveURL(userID, short, reqBody.URL, r); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-type", "application/json")
	w.WriteHeader(http.StatusCreated)
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
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	userURLs, err := server.GetUserURLs(string(decodedUserID), r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	msg, err := json.Marshal(userURLs)
	if err != nil {
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
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err = server.SaveURLs(dbEntries, r); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-type", "application/json")
	w.WriteHeader(http.StatusCreated)
	server.Write(marshaledBody, w)
}
