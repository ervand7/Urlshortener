package views

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"github.com/ervand7/urlshortener/internal/app/config"
	"github.com/ervand7/urlshortener/internal/app/controllers/generatedata"
	s "github.com/ervand7/urlshortener/internal/app/controllers/storage"
	"github.com/ervand7/urlshortener/internal/app/database"
	q "github.com/ervand7/urlshortener/internal/app/database/rawqueries"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

func TestShortenURL(t *testing.T) {
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
				Storage: s.GetStorage(),
			}
			defer func() {
				switch server.Storage.(type) {
				case *s.DBStorage:
					db := database.Database{}
					db.Run()
					_, err := db.Conn.Exec(q.DropAll)
					assert.NoError(t, err)
				}
			}()

			router := mux.NewRouter()
			router.HandleFunc("/", server.ShortenURL).Methods("POST")
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

func TestShortenURL409(t *testing.T) {
	if os.Getenv("DATABASE_DSN") != "user=ervand password=ervand dbname=urlshortener_test host=localhost port=5432 sslmode=disable" {
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

			db := database.Database{}
			db.Run()
			server := Server{
				Storage: s.NewDBStorage(db),
			}
			defer func() {
				_, err := db.Conn.Exec(q.DropAll)
				assert.NoError(t, err)
			}()

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

func TestGetURL(t *testing.T) {
	short := generatedata.ShortenURL()

	type want struct {
		statusCode int
		origin     string
	}
	tests := []struct {
		name     string
		endpoint string
		short    string
		method   string
		want     want
	}{
		{
			name:     "success 307",
			endpoint: short,
			short:    short,
			method:   http.MethodGet,
			want: want{
				statusCode: 307,
				origin:     "https://practicum.yandex.ru/learn/go-advanced/courses/14d6ff29-c8b6-43bf-9c55-12e8fe25b1b0/sprints/39172/topics/add19e4a-79bf-416e-9d13-0df2005ec81e/lessons/4280c1ec-45c9-4b73-bdf1-8ea438b18212/",
			},
		},
		{
			name:     "fail 405",
			endpoint: short,
			short:    short,
			method:   http.MethodPost,
			want: want{
				statusCode: http.StatusMethodNotAllowed,
				origin:     "",
			},
		},
		{
			name:     "fail 400",
			endpoint: short + "Hello",
			short:    short,
			method:   http.MethodGet,
			want: want{
				statusCode: http.StatusBadRequest,
				origin:     "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			server := Server{
				Storage: s.GetStorage(),
			}
			defer func() {
				switch server.Storage.(type) {
				case *s.DBStorage:
					db := database.Database{}
					db.Run()
					_, err := db.Conn.Exec(q.DropAll)
					assert.NoError(t, err)
				}
			}()

			ctx := context.Background()
			err := server.Storage.Set(ctx, uuid.New().String(), tt.short, tt.want.origin)
			assert.NoError(t, err)

			request := httptest.NewRequest(tt.method, tt.endpoint, nil)
			router := mux.NewRouter()
			router.HandleFunc("/{id:[a-zA-Z]+}", server.GetURL).Methods("GET")
			writer := httptest.NewRecorder()
			router.ServeHTTP(writer, request)

			response := writer.Result()
			err = response.Body.Close()
			assert.NoError(t, err)

			assert.Equal(t, tt.want.statusCode, response.StatusCode)
			assert.Equal(t, tt.want.origin, response.Header.Get("Location"))
		})
	}
}

func TestGetURL410(t *testing.T) {
	if os.Getenv("DATABASE_DSN") != "user=ervand password=ervand dbname=urlshortener_test host=localhost port=5432 sslmode=disable" {
		return
	}

	db := database.Database{}
	db.Run()
	server := Server{
		Storage: s.NewDBStorage(db),
	}
	defer func() {
		_, err := db.Conn.Exec(q.DropAll)
		assert.NoError(t, err)
	}()

	userID := uuid.New().String()
	short := generatedata.ShortenURL()
	origin := "world"
	_, err := db.Conn.Exec(
		`insert into url ("user_id", "short", "origin", "active") values ($1, $2, $3, $4)`,
		userID, short, origin, false,
	)
	assert.NoError(t, err)

	request := httptest.NewRequest(http.MethodGet, short, nil)
	router := mux.NewRouter()
	router.HandleFunc("/{id:[a-zA-Z]+}", server.GetURL).Methods("GET")
	writer := httptest.NewRecorder()
	router.ServeHTTP(writer, request)

	response := writer.Result()
	err = response.Body.Close()
	assert.NoError(t, err)
	assert.Equal(t, response.StatusCode, http.StatusGone)
}

func TestAPIShortenURL(t *testing.T) {
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
				Storage: s.GetStorage(),
			}
			defer func() {
				switch server.Storage.(type) {
				case *s.DBStorage:
					db := database.Database{}
					db.Run()
					_, err := db.Conn.Exec(q.DropAll)
					assert.NoError(t, err)
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
			bytesBody, err := ioutil.ReadAll(response.Body)
			require.NoError(t, err)

			err = json.Unmarshal(bytesBody, &marshalled)
			require.NoError(t, err)
			assert.Equal(t, tt.want.lenResultURL, len(marshalled["result"]))

			err = response.Body.Close()
			require.NoError(t, err)
		})
	}
}

func TestUserURLs(t *testing.T) {
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
				Storage: s.GetStorage(),
			}
			defer func() {
				switch server.Storage.(type) {
				case *s.DBStorage:
					db := database.Database{}
					db.Run()
					_, err := db.Conn.Exec(q.DropAll)
					assert.NoError(t, err)
				}
			}()

			requestPOST := httptest.NewRequest(
				http.MethodPost,
				"/api/shorten",
				bytes.NewBuffer([]byte(`{"url":"https://hello.world/"}`)),
			)
			router1 := mux.NewRouter()
			router1.HandleFunc("/api/shorten", server.APIShortenURL).Methods("POST")
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
			router2.HandleFunc("/api/user/urls", server.UserURLs).Methods("GET")
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

func TestAPIShortenBatch(t *testing.T) {
	if os.Getenv("DATABASE_DSN") != "user=ervand password=ervand dbname=urlshortener_test host=localhost port=5432 sslmode=disable" {
		return
	}

	type want struct {
		statusCode int
	}
	tests := []struct {
		name      string
		method    string
		body      string
		checkBody bool
		want      want
	}{
		{
			name:   "success 201",
			method: http.MethodPost,
			body: `
				[
				  {
					"correlation_id": "799eb140-c1f8-45c4-8268-656c5b20595b",
					"original_url": "https://practicum.yandex.ru/learn/go-advanced/courses/14d6ff29-c8b6-43bf-9c55-12e8fe25b1b0/sprints/39172/topics/add19e4a-79bf-416e-9d13-0df2005ec81e/lessons/74ebe51b-de58-4142-83b7-3d21ba695a9f/"
				  },
				  {
					"correlation_id": "6dde456c-1777-40f5-9f8c-8f6ea6454715",
					"original_url": "https://practicum/799eb140-c1f8-45c4-8268-656c5b20595bb6-43bf-9c55-12e8fe25b1b0/sprints/39172/656c5b20595bb6-43bf-9c55-12e8fe81e/lessons/74ebe51b-de58-4142-83b7-3d2-656c5b20595bb6-43bf-1ba695a9f/"
				  }
				]
			`,
			checkBody: true,
			want: want{
				statusCode: 201,
			},
		},
		{
			name:      "fail 405",
			method:    http.MethodGet,
			body:      "",
			checkBody: false,
			want: want{
				statusCode: 405,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var reqBody = []byte(tt.body)
			request := httptest.NewRequest(
				tt.method,
				"/api/shorten/batch",
				bytes.NewBuffer(reqBody),
			)

			db := database.Database{}
			db.Run()
			server := Server{
				Storage: s.NewDBStorage(db),
			}
			defer func() {
				_, err := db.Conn.Exec(q.DropAll)
				assert.NoError(t, err)
			}()

			router := mux.NewRouter()
			router.HandleFunc("/api/shorten/batch",
				server.APIShortenBatch).Methods("POST")
			writer := httptest.NewRecorder()
			router.ServeHTTP(writer, request)

			response := writer.Result()
			assert.Equal(t, tt.want.statusCode, response.StatusCode)
			if !tt.checkBody {
				return
			}

			type RespPair struct {
				CorrelationID string `json:"correlation_id"`
				ShortURL      string `json:"short_url"`
			}
			var respPairs []RespPair
			respBody, err := ioutil.ReadAll(response.Body)
			require.NoError(t, err)

			err = json.Unmarshal(respBody, &respPairs)
			require.NoError(t, err)
			assert.Equal(t, respPairs[0].CorrelationID, "799eb140-c1f8-45c4-8268-656c5b20595b")
			assert.Equal(t, respPairs[1].CorrelationID, "6dde456c-1777-40f5-9f8c-8f6ea6454715")

			err = response.Body.Close()
			require.NoError(t, err)
		})
	}
}

func TestUserURLsDelete(t *testing.T) {
	if os.Getenv("DATABASE_DSN") != "user=ervand password=ervand dbname=urlshortener_test host=localhost port=5432 sslmode=disable" {
		return
	}

	db := database.Database{}
	db.Run()
	server := Server{
		Storage: s.NewDBStorage(db),
	}
	defer func() {
		_, err := db.Conn.Exec(q.DropAll)
		assert.NoError(t, err)
	}()

	shorts := []string{
		"hello1",
		"hello2",
		"hello3",
		"hello4",
		"hello5",
		"hello6",
	}
	origins := []string{
		"world1",
		"world2",
		"world3",
		"world4",
		"world5",
		"world6",
	}
	userID := uuid.New().String()
	for idx := range shorts {
		_, err := db.Conn.Exec(
			`insert into url ("user_id", "short", "origin") values ($1, $2, $3)`,
			userID, shorts[idx], origins[idx],
		)
		assert.NoError(t, err)
	}

	cookie := &http.Cookie{Name: "userID", Value: hex.EncodeToString([]byte(userID))}
	marshaled, err := json.Marshal(shorts)
	assert.NoError(t, err)

	request := httptest.NewRequest(http.MethodDelete, "/api/user/urls", bytes.NewBuffer(marshaled))
	request.AddCookie(cookie)
	router := mux.NewRouter()
	router.HandleFunc("/api/user/urls", server.UserURLsDelete).Methods("DELETE")
	writer := httptest.NewRecorder()
	router.ServeHTTP(writer, request)

	response := writer.Result()
	assert.Equal(t, response.StatusCode, http.StatusAccepted)
	err = response.Body.Close()
	assert.NoError(t, err)
	time.Sleep(s.Timeout * time.Second)

	for _, short := range shorts {
		ctx := context.TODO()
		rows := db.Conn.QueryRowContext(
			ctx, `select "active" from url where "short" = $1;`, short,
		)
		var active bool

		err = rows.Scan(&active)
		assert.NoError(t, err)
		assert.Equal(t, false, active)
		err = rows.Err()
		assert.NoError(t, err)
	}
}
