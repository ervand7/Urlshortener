package views

import (
	"github.com/ervand7/urlshortener/internal/app/controllers"
	"github.com/ervand7/urlshortener/internal/app/database"
	"io"
	"net/http"
)

// UrlShorten POST ("/")
func UrlShorten(writer http.ResponseWriter, request *http.Request) {
	defer request.Body.Close()
	body, err := io.ReadAll(request.Body)
	if err != nil {
		http.Error(
			writer,
			"invalid body",
			http.StatusBadRequest)
		return
	}
	url := string(body)
	if url == "" {
		http.Error(
			writer,
			"param url not filled",
			http.StatusBadRequest,
		)
		return
	}
	writer.WriteHeader(http.StatusCreated)
	writer.Header().Add("Content-type", "text/plain; charset=utf-8")

	shortenedUrl := controllers.ShortenUrl()
	database.OriginUrlStorage[shortenedUrl] = url
	writer.Write([]byte(shortenedUrl))
	return
}

// UrlGet GET /{id}
func UrlGet(writer http.ResponseWriter, request *http.Request) {
	key := controllers.BaseUrl + request.URL.String()
	OriginUrl := database.OriginUrlStorage[key]
	if OriginUrl == "" {
		http.Error(
			writer,
			"not exists",
			http.StatusBadRequest)
		return
	}
	http.Redirect(
		writer,
		request,
		OriginUrl,
		http.StatusTemporaryRedirect)
}
