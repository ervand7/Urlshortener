package views

import (
	"bytes"
	"github.com/ervand7/urlshortener/internal/app/controllers"
	"github.com/ervand7/urlshortener/internal/app/models"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUrlShorten(t *testing.T) {
	lenRespBody := len(controllers.BaseUrl) + len("/") + controllers.ShortenEndpointLen

	type want struct {
		contentType string
		statusCode  int
		lenRespBody int
	}
	tests := []struct {
		name   string
		body   string
		method string
		want   want
	}{
		{
			name:   "success 201",
			body:   "https://practicum.yandex.ru/learn/go-advanced/courses/14d6ff29-c8b6-43bf-9c55-12e8fe25b1b0/sprints/39172/topics/add19e4a-79bf-416e-9d13-0df2005ec81e/lessons/4280c1ec-45c9-4b73-bdf1-8ea438b18212/",
			method: http.MethodPost,
			want: want{
				contentType: "text/plain; charset=utf-8",
				statusCode:  http.StatusCreated,
				lenRespBody: lenRespBody,
			},
		},
		{
			name:   "fail 405",
			body:   "https://practicum.yandex.ru/learn/go-advanced/courses/14d6ff29-c8b6-43bf-9c55-12e8fe25b1b0/sprints/39172/topics/add19e4a-79bf-416e-9d13-0df2005ec81e/lessons/4280c1ec-45c9-4b73-bdf1-8ea438b18212/",
			method: http.MethodGet,
			want: want{
				contentType: "",
				statusCode:  http.StatusMethodNotAllowed,
				lenRespBody: 0,
			},
		},
		{
			name:   "fail 400 no body",
			body:   "",
			method: http.MethodPost,
			want: want{
				contentType: "text/plain; charset=utf-8",
				statusCode:  http.StatusBadRequest,
				lenRespBody: len("param url not filled "),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			var body = []byte(tt.body)
			request := httptest.NewRequest(
				tt.method,
				"/",
				bytes.NewBuffer(body),
			)

			storage := make(models.Storage, 0)
			router := mux.NewRouter()
			router.HandleFunc("/", UrlShorten(storage)).Methods("POST")
			writer := httptest.NewRecorder()
			router.ServeHTTP(writer, request)

			response := writer.Result()
			assert.Equal(t, tt.want.statusCode, response.StatusCode)
			assert.Equal(t, tt.want.contentType, response.Header.Get("Content-Type"))

			responseBody, err := ioutil.ReadAll(response.Body)
			assert.Equal(t, tt.want.lenRespBody, len(responseBody))
			require.NoError(t, err)
			err = response.Body.Close()
			require.NoError(t, err)
		})
	}
}

func TestUrlGet(t *testing.T) {
	shortened := controllers.ShortenUrl()

	type want struct {
		statusCode int
		originUrl  string
	}
	tests := []struct {
		name         string
		endpoint     string
		shortenedUrl string
		method       string
		want         want
	}{
		{
			name:         "success 307",
			endpoint:     shortened,
			shortenedUrl: shortened,
			method:       http.MethodGet,
			want: want{
				statusCode: 307,
				originUrl:  "https://practicum.yandex.ru/learn/go-advanced/courses/14d6ff29-c8b6-43bf-9c55-12e8fe25b1b0/sprints/39172/topics/add19e4a-79bf-416e-9d13-0df2005ec81e/lessons/4280c1ec-45c9-4b73-bdf1-8ea438b18212/",
			},
		},
		{
			name:         "fail 405",
			endpoint:     shortened,
			shortenedUrl: shortened,
			method:       http.MethodPost,
			want: want{
				statusCode: http.StatusMethodNotAllowed,
				originUrl:  "",
			},
		},
		{
			name:         "fail 400",
			endpoint:     shortened + "Hello",
			shortenedUrl: shortened,
			method:       http.MethodGet,
			want: want{
				statusCode: http.StatusBadRequest,
				originUrl:  "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			storage := make(models.Storage, 0)
			storage.Set(tt.shortenedUrl, tt.want.originUrl)

			request := httptest.NewRequest(tt.method, tt.endpoint, nil)
			router := mux.NewRouter()
			router.HandleFunc("/{id:[a-zA-Z]+}", UrlGet(storage)).Methods("GET")
			writer := httptest.NewRecorder()
			router.ServeHTTP(writer, request)

			response := writer.Result()
			assert.Equal(t, tt.want.statusCode, response.StatusCode)
			assert.Equal(t, tt.want.originUrl, response.Header.Get("Location"))
		})
	}
}
