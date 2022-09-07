package views

import (
	"encoding/hex"
	"encoding/json"
	"github.com/ervand7/urlshortener/internal/app/config"
	"github.com/ervand7/urlshortener/internal/app/controllers/generatedata"
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

	userID := server.GetUserIDFromCookie(w, r)
	w.Header().Add("Content-type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusCreated)

	shortenURL := generatedata.ShortenURL()
	server.SaveURL(userID, shortenURL, url, w)
	server.Write([]byte(shortenURL), w)
}

// URLGet GET ("/{id}")
func (server *Server) URLGet(w http.ResponseWriter, r *http.Request) {
	endpoint := r.URL.Path
	shortenedURL := config.GetConfig().BaseURL + endpoint
	var originURL string

	switch config.GetConfig().FileStoragePath {
	case "":
		originURL = server.MemoryStorage.Get(shortenedURL)
	default:
		result, err := server.FileStorage.Get(shortenedURL)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		originURL = result
	}

	if originURL == "" {
		http.Error(w, "not exists", http.StatusBadRequest)
		return
	}
	http.Redirect(w, r, originURL, http.StatusTemporaryRedirect)
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

	shortenURL := generatedata.ShortenURL()
	respBody := RespBody{
		Result: shortenURL,
	}
	marshaledBody, err := json.Marshal(respBody)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	userID := server.GetUserIDFromCookie(w, r)
	w.Header().Add("Content-type", "application/json")
	w.WriteHeader(http.StatusCreated)

	server.SaveURL(userID, shortenURL, reqBody.URL, w)
	server.Write(marshaledBody, w)
}

// URLUserAll GET ("/api/user/urls")
func (server *Server) URLUserAll(w http.ResponseWriter, r *http.Request) {
	userID, err := r.Cookie("userID")
	if err != nil {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	decoded, _ := hex.DecodeString(userID.Value)

	if err != nil {
		utils.Logger.Warn(err.Error())
	}

	userURLs := server.GetUserURLs(string(decoded))
	msg, err := json.Marshal(userURLs)
	if err != nil {
		utils.Logger.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-type", "application/json")
	server.Write(msg, w)
}
