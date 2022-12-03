package views

import (
	"bytes"
	"encoding/json"
	"github.com/ervand7/urlshortener/internal/config"
	s "github.com/ervand7/urlshortener/internal/controllers/storage"
	"github.com/ervand7/urlshortener/internal/database"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestAPIShortenBatch(t *testing.T) {
	if os.Getenv("DATABASE_DSN") != config.TestDBAddr {
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

			defer database.Downgrade()
			db := database.Database{}
			db.Run()
			server := Server{
				Storage: s.NewDBStorage(db),
			}

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