package views

import (
	"github.com/ervand7/urlshortener/internal/app/controllers"
	"github.com/ervand7/urlshortener/internal/app/models"
	"io"
	"net/http"
)

// URLShorten POST ("/")
func URLShorten(storage models.URLRepository) func(writer http.ResponseWriter, request *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		defer request.Body.Close()
		body, err := io.ReadAll(request.Body)
		if err != nil {
			http.Error(writer, "invalid body", http.StatusBadRequest)
			return
		}
		url := string(body)
		if url == "" {
			http.Error(writer, "param url not filled", http.StatusBadRequest)
			return
		}
		writer.Header().Add("Content-type", "text/plain; charset=utf-8")
		writer.WriteHeader(http.StatusCreated)

		shortenedURL := controllers.ShortenURL()
		storage.Set(shortenedURL, url)
		writer.Write([]byte(shortenedURL))
	}
}

// URLGet GET /{id}
func URLGet(storage models.URLRepository) func(writer http.ResponseWriter, request *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		endpoint := request.URL.Path
		shortenedURL := controllers.BaseURL + endpoint
		originURL := storage.Get(shortenedURL)
		if originURL == "" {
			http.Error(writer, "not exists", http.StatusBadRequest)
			return
		}
		http.Redirect(
			writer,
			request,
			originURL,
			http.StatusTemporaryRedirect,
		)
	}
}
