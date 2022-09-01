package views

import (
	"encoding/json"
	"github.com/ervand7/urlshortener/internal/app/config"
	"github.com/ervand7/urlshortener/internal/app/controllers/generatedata"
	"github.com/ervand7/urlshortener/internal/app/utils"
	"io"
	"net/http"
)

// URLShorten POST ("/")
func (server *Server) URLShorten(w http.ResponseWriter, r *http.Request) {
	defer func() {
		err := r.Body.Close()
		if err != nil {
			utils.Logger.Warn(err.Error())
		}
	}()

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
	w.Header().Add("Content-type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusCreated)

	shortenedURL := generatedata.ShortenURL()
	switch config.GetConfig().FileStoragePath {
	case "":
		server.MemoryStorage.Set(shortenedURL, url)
	default:
		if err := server.FileStorage.Set(shortenedURL, url); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	_, writeErr := w.Write([]byte(shortenedURL))
	if writeErr != nil {

		utils.Logger.Warn(writeErr.Error() + "")
	}
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
	defer func() {
		err := r.Body.Close()
		if err != nil {
			utils.Logger.Warn(err.Error())
		}
	}()

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

	shortenedURL := generatedata.ShortenURL()
	respBody := RespBody{
		Result: shortenedURL,
	}
	marshaledBody, err := json.Marshal(respBody)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Add("Content-type", "application/json")
	w.WriteHeader(http.StatusCreated)

	switch config.GetConfig().FileStoragePath {
	case "":
		server.MemoryStorage.Set(shortenedURL, reqBody.URL)
	default:
		if err := server.FileStorage.Set(shortenedURL, reqBody.URL); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	_, writeErr := w.Write(marshaledBody)
	if writeErr != nil {
		utils.Logger.Warn(writeErr.Error())
	}
}
