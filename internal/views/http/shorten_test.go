package http

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ervand7/urlshortener/internal/config"
	"github.com/ervand7/urlshortener/internal/controllers/algorithms"
	s "github.com/ervand7/urlshortener/internal/controllers/storage"
	"github.com/ervand7/urlshortener/internal/controllers/storage/dbstorage"
)

func TestShortenURL(t *testing.T) {
	lenRespBody := len(config.GetBaseURL()) +
		len("/") +
		algorithms.ShortenEndpointLen

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

			server := Server{
				Storage: s.GetStorage(),
			}
			defer func() {
				if _, ok := server.Storage.(*dbstorage.DBStorage); ok {
					dbstorage.Downgrade()
				}
			}()

			router := mux.NewRouter()
			router.HandleFunc("/", server.ShortenURL).Methods("POST")
			writer := httptest.NewRecorder()
			router.ServeHTTP(writer, request)

			response := writer.Result()
			assert.Equal(t, tt.want.statusCode, response.StatusCode)
			assert.Equal(t, tt.want.contentType, response.Header.Get("Content-Type"))

			responseBody, err := io.ReadAll(response.Body)
			assert.Equal(t, tt.want.lenRespBody, len(responseBody))
			require.NoError(t, err)
			err = response.Body.Close()
			require.NoError(t, err)
		})
	}
}

func TestShortenURL409(t *testing.T) {
	if os.Getenv("DATABASE_DSN") != config.TestDBAddr {
		return
	}
	type want struct {
		contentType string
	}
	tests := []struct {
		name     string
		handler  string
		endpoint string
		body     string
		method   string
		want     want
	}{
		{
			name:     "success 409 ShortenURL",
			handler:  "ShortenURL",
			endpoint: "/",
			body:     "https://practicum.yandex.ru/learn/go-advanced/courses/14d6ff29-c8b6-43bf-9c55-12e8fe25b1b0/sprints/39172/topics/add19e4a-79bf-416e-9d13-0df2005ec81e/lessons/4280c1ec-45c9-4b73-bdf1-8ea438b18212/",
			method:   http.MethodPost,
			want: want{
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name:     "success 409 APIShortenURL",
			handler:  "APIShortenURL",
			endpoint: "/api/shorten",
			body:     `{"url":"https://practicum.yandex.ru/learn/go-advanced/courses/14d6ff29-c8b6-43bf-9c55-12e8fe25b1b0/sprints/39172/topics/add19e4a-79bf-416e-9d13-0df2005ec81e/lessons/4280c1ec-45c9-4b73-bdf1-8ea438b18212/"}`,
			method:   http.MethodPost,
			want: want{
				contentType: "application/json",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var body = []byte(tt.body)
			request1 := httptest.NewRequest(
				tt.method,
				tt.endpoint,
				bytes.NewBuffer(body),
			)

			defer dbstorage.Downgrade()
			db := dbstorage.Database{}
			db.Run()
			server := Server{
				Storage: dbstorage.NewDBStorage(db),
			}

			var handler func(w http.ResponseWriter, r *http.Request)
			if tt.handler == "ShortenURL" {
				handler = server.ShortenURL
			} else {
				handler = server.APIShortenURL
			}

			router1 := mux.NewRouter()
			router1.HandleFunc(tt.endpoint,
				handler).Methods("POST")
			writer1 := httptest.NewRecorder()
			router1.ServeHTTP(writer1, request1)

			response1 := writer1.Result()
			assert.Equal(t, response1.StatusCode, http.StatusCreated)

			request2 := httptest.NewRequest(
				tt.method,
				tt.endpoint,
				bytes.NewBuffer(body),
			)
			router2 := mux.NewRouter()
			router2.HandleFunc(tt.endpoint,
				handler).Methods("POST")
			writer2 := httptest.NewRecorder()
			router2.ServeHTTP(writer2, request2)

			response2 := writer2.Result()
			assert.Equal(t, response2.StatusCode, http.StatusConflict)
			assert.Equal(t, tt.want.contentType, response2.Header.Get("Content-Type"))

			err := response1.Body.Close()
			require.NoError(t, err)
			err = response2.Body.Close()
			require.NoError(t, err)
		})
	}
}
