package views

import (
	"bytes"
	s "github.com/ervand7/urlshortener/internal/controllers/storage"
	"github.com/ervand7/urlshortener/internal/database"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

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
					database.Downgrade()
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
