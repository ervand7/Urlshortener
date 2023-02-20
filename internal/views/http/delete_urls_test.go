package http

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"

	"github.com/ervand7/urlshortener/internal/config"
	"github.com/ervand7/urlshortener/internal/controllers/storage/dbstorage"
)

func TestUserURLsDelete(t *testing.T) {
	if os.Getenv("DATABASE_DSN") != config.TestDBAddr {
		return
	}
	defer dbstorage.Downgrade()
	db := dbstorage.Database{}
	db.Run()
	server := Server{
		Storage: dbstorage.NewDBStorage(db),
	}

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
	time.Sleep(dbstorage.Timeout * time.Second)

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
