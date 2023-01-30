package views

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ervand7/urlshortener/internal/config"
	"github.com/ervand7/urlshortener/internal/controllers/storage/dbstorage"
)

func TestPingDB(t *testing.T) {
	if os.Getenv("DATABASE_DSN") != config.TestDBAddr {
		return
	}
	defer dbstorage.Downgrade()

	db := dbstorage.Database{}
	db.Run()
	server := Server{
		Storage: dbstorage.NewDBStorage(db),
	}
	request := httptest.NewRequest(http.MethodGet, "/ping", nil)
	router := mux.NewRouter()
	router.HandleFunc("/ping", server.PingDB).Methods(http.MethodGet)
	writer := httptest.NewRecorder()
	router.ServeHTTP(writer, request)

	response := writer.Result()
	assert.Equal(t, response.StatusCode, http.StatusOK)
	err := response.Body.Close()
	require.NoError(t, err)
}
