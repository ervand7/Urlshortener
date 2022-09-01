package views

import (
	"bytes"
	"encoding/json"
	"github.com/ervand7/urlshortener/internal/app/config"
	"github.com/ervand7/urlshortener/internal/app/controllers/generatedata"
	"github.com/ervand7/urlshortener/internal/app/models/url"
	"github.com/ervand7/urlshortener/internal/app/utils"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUrlShorten(t *testing.T) {
	lenRespBody := len(config.GetConfig().BaseURL) +
		len("/") +
		generatedata.ShortenEndpointLen

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
				MemoryStorage: &url.MemoryStorage{
					HashTable: make(map[string]url.ShortenURLStruct, 0),
				},
			}
			router := mux.NewRouter()
			router.HandleFunc("/", server.URLShorten).Methods("POST")
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
	shortened := generatedata.ShortenURL()

	type want struct {
		statusCode int
		originURL  string
	}
	tests := []struct {
		name         string
		endpoint     string
		shortenedURL string
		method       string
		want         want
	}{
		{
			name:         "success 307",
			endpoint:     shortened,
			shortenedURL: shortened,
			method:       http.MethodGet,
			want: want{
				statusCode: 307,
				originURL:  "https://practicum.yandex.ru/learn/go-advanced/courses/14d6ff29-c8b6-43bf-9c55-12e8fe25b1b0/sprints/39172/topics/add19e4a-79bf-416e-9d13-0df2005ec81e/lessons/4280c1ec-45c9-4b73-bdf1-8ea438b18212/",
			},
		},
		{
			name:         "fail 405",
			endpoint:     shortened,
			shortenedURL: shortened,
			method:       http.MethodPost,
			want: want{
				statusCode: http.StatusMethodNotAllowed,
				originURL:  "",
			},
		},
		{
			name:         "fail 400",
			endpoint:     shortened + "Hello",
			shortenedURL: shortened,
			method:       http.MethodGet,
			want: want{
				statusCode: http.StatusBadRequest,
				originURL:  "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			server := Server{
				MemoryStorage: &url.MemoryStorage{
					HashTable: make(map[string]url.ShortenURLStruct, 0),
				},
			}
			server.MemoryStorage.Set("", tt.shortenedURL, tt.want.originURL)

			request := httptest.NewRequest(tt.method, tt.endpoint, nil)
			router := mux.NewRouter()
			router.HandleFunc("/{id:[a-zA-Z]+}", server.URLGet).Methods("GET")
			writer := httptest.NewRecorder()
			router.ServeHTTP(writer, request)

			response := writer.Result()
			if err := response.Body.Close(); err != nil {
				utils.Logger.Warn(err.Error())
			}
			assert.Equal(t, tt.want.statusCode, response.StatusCode)
			assert.Equal(t, tt.want.originURL, response.Header.Get("Location"))
		})
	}
}

func TestUrlShortenJSON(t *testing.T) {
	lenResultURL := len(config.GetConfig().BaseURL) +
		len("/") +
		generatedata.ShortenEndpointLen
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
				MemoryStorage: &url.MemoryStorage{
					HashTable: make(map[string]url.ShortenURLStruct, 0),
				},
			}
			router := mux.NewRouter()
			router.HandleFunc("/api/shorten", server.URLShortenJSON).Methods("POST")
			writer := httptest.NewRecorder()
			router.ServeHTTP(writer, request)

			response := writer.Result()
			assert.Equal(t, tt.want.statusCode, response.StatusCode)
			assert.Equal(t, tt.want.contentType, response.Header.Get("Content-Type"))
			if !tt.checkBody {
				return
			}

			marshalled := make(map[string]string, 0)
			bytesBody, err := ioutil.ReadAll(response.Body)
			require.NoError(t, err)

			unmarshalErr := json.Unmarshal(bytesBody, &marshalled)
			require.NoError(t, unmarshalErr)
			assert.Equal(t, tt.want.lenResultURL, len(marshalled["result"]))

			err = response.Body.Close()
			require.NoError(t, err)
		})
	}
}

func TestURLUserAll(t *testing.T) {
	type want struct {
		statusCode int
		bodyChunk  string
	}
	tests := []struct {
		name      string
		method    string
		userID    string
		checkBody bool
		setCookie bool
		want      want
	}{
		{
			name:      "success 200",
			method:    http.MethodGet,
			checkBody: true,
			setCookie: true,
			want: want{
				statusCode: 200,
				bodyChunk:  `[{"original_url":"https://hello.world/","short_url":"http://localhost:8080/`,
			},
		},
		{
			name:      "success 204",
			method:    http.MethodGet,
			checkBody: false,
			setCookie: false,
			want: want{
				statusCode: 204,
			},
		},
		{
			name:      "fail 405",
			method:    http.MethodPost,
			checkBody: false,
			setCookie: false,
			want: want{
				statusCode: 405,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := Server{
				MemoryStorage: &url.MemoryStorage{
					HashTable: make(map[string]url.ShortenURLStruct, 0),
				},
			}
			requestPOST := httptest.NewRequest(
				http.MethodPost,
				"/api/shorten",
				bytes.NewBuffer([]byte(`{"url":"https://hello.world/"}`)),
			)
			router1 := mux.NewRouter()
			router1.HandleFunc("/api/shorten", server.URLShortenJSON).Methods("POST")
			writer1 := httptest.NewRecorder()
			router1.ServeHTTP(writer1, requestPOST)

			responsePOST := writer1.Result()
			cookies := responsePOST.Cookies()
			assert.Equal(t, cookies[0].Name, "userID")
			assert.NotEmpty(t, cookies[0].Value)
			err := responsePOST.Body.Close()
			require.NoError(t, err)

			requestGET := httptest.NewRequest(tt.method, "/api/user/urls", nil)
			if tt.setCookie {
				cookie := &http.Cookie{Name: "userID", Value: cookies[0].Value, HttpOnly: true}
				requestGET.AddCookie(cookie)
			}
			router2 := mux.NewRouter()
			writer2 := httptest.NewRecorder()
			router2.HandleFunc("/api/user/urls", server.URLUserAll).Methods("GET")
			router2.ServeHTTP(writer2, requestGET)

			responseGET := writer2.Result()
			assert.Equal(t, responseGET.StatusCode, tt.want.statusCode)

			if tt.checkBody == false {
				return
			}
			body, err := ioutil.ReadAll(responseGET.Body)
			require.NoError(t, err)
			assert.Contains(t, string(body), tt.want.bodyChunk)
			err = responseGET.Body.Close()
			require.NoError(t, err)
		})
	}
}
