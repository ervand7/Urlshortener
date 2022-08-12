package views

import (
	"github.com/ervand7/urlshortener/internal/app/controllers"
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
		url := string(body)
		if url == "" {
			http.Error(w, "param url not filled", http.StatusBadRequest)
			return
		}
		w.Header().Add("Content-type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusCreated)

		shortenedURL := controllers.ShortenURL()
		server.Storage.Set(shortenedURL, url)
		w.Write([]byte(shortenedURL))
	}
}

// URLGet GET /{id}
func (server *Server) URLGet() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		endpoint := r.URL.Path
		shortenedURL := controllers.BaseURL + endpoint
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
