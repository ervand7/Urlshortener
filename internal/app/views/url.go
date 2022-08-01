package views

import (
	"github.com/ervand7/urlshortener/internal/app/controllers"
	"github.com/ervand7/urlshortener/internal/app/storage"
	"io"
	"net/http"
)

// UrlShorten POST ("/")
func UrlShorten(_storage storage.UrlRepository) func(writer http.ResponseWriter, request *http.Request) {
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

		shortenedUrl := controllers.ShortenUrl()
		_storage.Set(shortenedUrl, url)
		writer.Write([]byte(shortenedUrl))
		return
	}
}

// UrlGet GET /{id}
func UrlGet(_storage storage.UrlRepository) func(writer http.ResponseWriter, request *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		endpoint := request.URL.Path
		shortenedUrl := controllers.BaseUrl + endpoint
		originUrl := _storage.Get(shortenedUrl)
		if originUrl == "" {
			http.Error(writer, "not exists", http.StatusBadRequest)
			return
		}
		http.Redirect(
			writer,
			request,
			originUrl,
			http.StatusTemporaryRedirect,
		)
	}
}
