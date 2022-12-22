package views

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ervand7/urlshortener/internal/config"
	"github.com/ervand7/urlshortener/internal/controllers/algorithms"
	s "github.com/ervand7/urlshortener/internal/controllers/storage"
	"github.com/ervand7/urlshortener/internal/controllers/storage/dbstorage"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAPIShortenURL(t *testing.T) {
	lenResultURL := len(config.GetBaseURL()) +
		len("/") +
		algorithms.ShortenEndpointLen
	const Endpoint string = "/api/shorten"

	type want struct {
		contentType  string
		statusCode   int
		lenResultURL int
	}
	tests := []struct {
		name      string
		endpoint  string
		body      string
		checkBody bool
		method    string
		want      want
	}{
		{
			name:      "success 201",
			endpoint:  Endpoint,
			body:      `{"url":"https://practicum.yandex.ru/learn/go-advanced/courses/14d6ff29-c8b6-43bf-9c55-12e8fe25b1b0/sprints/39172/topics/add19e4a-79bf-416e-9d13-0df2005ec81e/lessons/4280c1ec-45c9-4b73-bdf1-8ea438b18212/"}`,
			checkBody: true,
			method:    http.MethodPost,
			want: want{
				contentType:  "application/json",
				statusCode:   http.StatusCreated,
				lenResultURL: lenResultURL,
			},
		},
		{
			name:      "fail 405",
			endpoint:  Endpoint,
			body:      "https://practicum.yandex.ru/learn/go-advanced/courses/14d6ff29-c8b6-43bf-9c55-12e8fe25b1b0/sprints/39172/topics/add19e4a-79bf-416e-9d13-0df2005ec81e/lessons/4280c1ec-45c9-4b73-bdf1-8ea438b18212/",
			checkBody: false,
			method:    http.MethodGet,
			want: want{
				contentType: "",
				statusCode:  http.StatusMethodNotAllowed,
			},
		},
		{
			name:      "fail 400 no body",
			endpoint:  Endpoint,
			body:      "",
			checkBody: false,
			method:    http.MethodPost,
			want: want{
				contentType: "text/plain; charset=utf-8",
				statusCode:  http.StatusBadRequest,
			},
		},
		{
			name:      "fail 404",
			endpoint:  Endpoint + "hello",
			body:      "https://practicum.yandex.ru/learn/go-advanced/courses/14d6ff29-c8b6-43bf-9c55-12e8fe25b1b0/sprints/39172/topics/add19e4a-79bf-416e-9d13-0df2005ec81e/lessons/4280c1ec-45c9-4b73-bdf1-8ea438b18212/",
			checkBody: false,
			method:    http.MethodPost,
			want: want{
				contentType: "text/plain; charset=utf-8",
				statusCode:  http.StatusNotFound,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var body = []byte(tt.body)
			request := httptest.NewRequest(
				tt.method,
				tt.endpoint,
				bytes.NewBuffer(body),
			)

			server := Server{
				Storage: s.GetStorage(),
			}
			defer func() {
				switch server.Storage.(type) {
				case *dbstorage.DBStorage:
					dbstorage.Downgrade()
				}
			}()

			router := mux.NewRouter()
			router.HandleFunc("/api/shorten", server.APIShortenURL).Methods("POST")
			writer := httptest.NewRecorder()
			router.ServeHTTP(writer, request)

			response := writer.Result()
			assert.Equal(t, tt.want.statusCode, response.StatusCode)
			assert.Equal(t, tt.want.contentType, response.Header.Get("Content-Type"))
			if !tt.checkBody {
				return
			}

			marshalled := make(map[string]string, 0)
			bytesBody, err := io.ReadAll(response.Body)
			require.NoError(t, err)

			err = json.Unmarshal(bytesBody, &marshalled)
			require.NoError(t, err)
			assert.Equal(t, tt.want.lenResultURL, len(marshalled["result"]))

			err = response.Body.Close()
			require.NoError(t, err)
		})
	}
}
