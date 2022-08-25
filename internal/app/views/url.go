package views

import (
	"encoding/json"
	"github.com/ervand7/urlshortener/internal/app/config"
	"github.com/ervand7/urlshortener/internal/app/controllers/compress"
	"github.com/ervand7/urlshortener/internal/app/controllers/generatedata"
	"io"
	"net/http"
)

// URLShorten POST ("/")
func (server *Server) URLShorten() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "invalid body", http.StatusBadRequest)
			return
		}
		if r.Header.Get("Content-Encoding") == "gzip" {
			body, _ = compress.Decompress(body)
		}

		url := string(body)
		if url == "" {
			http.Error(w, "param url not filled", http.StatusBadRequest)
			return
		}
		shortenedURL := generatedata.ShortenURL()
		server.Storage.Set(shortenedURL, url)

		w.Header().Add("Content-type", "text/plain; charset=utf-8")
		if r.Header.Get("Accept-Encoding") == "gzip" {
			compressed, _ := compress.Compress([]byte(shortenedURL))
			w.WriteHeader(http.StatusCreated)
			w.Write(compressed)
		}
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(shortenedURL))
	}
}

// URLGet GET ("/{id}")
func (server *Server) URLGet() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		endpoint := r.URL.Path
		shortenedURL := config.GetConfig().BaseURL + endpoint
		originURL := server.Storage.Get(shortenedURL)
		if originURL == "" {
			http.Error(w, "not exists", http.StatusBadRequest)
			return
		}
		http.Redirect(
			w,
			r,
			originURL,
			http.StatusTemporaryRedirect,
		)
	}
}

// URLShortenJSON POST ("/api/shorten")
func (server *Server) URLShortenJSON() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		type (
			ReqBody struct {
				URL string `json:"url"`
			}
			RespBody struct {
				Result string `json:"result"`
			}
		)

		bodyBytes, err := io.ReadAll(r.Body)
		if r.Header.Get("Content-Encoding") == "gzip" {
			bodyBytes, _ = compress.Decompress(bodyBytes)
		}
		if err != nil {
			http.Error(w, "invalid body", http.StatusBadRequest)
			return
		}
		if len(bodyBytes) == 0 {
			http.Error(w, "param url not filled", http.StatusBadRequest)
			return
		}

		reqBody := ReqBody{}
		if err := json.Unmarshal(bodyBytes, &reqBody); err != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
		shortenedURL := generatedata.ShortenURL()
		respBody := RespBody{
			Result: shortenedURL,
		}
		marshaledBody, err := json.Marshal(respBody)
		if err != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		server.Storage.Set(shortenedURL, reqBody.URL)
		w.Header().Add("Content-type", "application/json")
		if r.Header.Get("Accept-Encoding") == "gzip" {
			compressed, _ := compress.Compress(marshaledBody)
			w.WriteHeader(http.StatusCreated)
			w.Write(compressed)
		}
		w.WriteHeader(http.StatusCreated)
		w.Write(marshaledBody)
	}
}
