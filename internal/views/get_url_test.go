package views

import (
	"context"
	"github.com/ervand7/urlshortener/internal/config"
	"github.com/ervand7/urlshortener/internal/controllers/generatedata"
	s "github.com/ervand7/urlshortener/internal/controllers/storage"
	"github.com/ervand7/urlshortener/internal/database"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

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
					database.Downgrade()
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
	if os.Getenv("DATABASE_DSN") != config.TestDBAddr {
		return
	}
	defer database.Downgrade()
	db := database.Database{}
	db.Run()
	server := Server{
		Storage: s.NewDBStorage(db),
	}

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