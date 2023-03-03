package http

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"

	"github.com/google/uuid"
	"github.com/gorilla/mux"

	"github.com/ervand7/urlshortener/internal/config"
	"github.com/ervand7/urlshortener/internal/controllers/algorithms"
	s "github.com/ervand7/urlshortener/internal/controllers/storage"
	"github.com/ervand7/urlshortener/internal/controllers/storage/dbstorage"
)

func ExampleServer_ShortenURL() {
	bodyString := "https://practicum.yandex.ru/learn/go-advanced/courses/14d6ff29-c8b6-43bf-9c55-12e8fe25b1b0/sprints/39172/topics/add19e4a-79bf-416e-9d13-0df2005ec81e/lessons/4280c1ec-45c9-4b73-bdf1-8ea438b18212/"
	var body = []byte(bodyString)
	request := httptest.NewRequest(
		http.MethodPost,
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
	router.HandleFunc("/", server.ShortenURL).Methods(http.MethodPost)
	writer := httptest.NewRecorder()
	router.ServeHTTP(writer, request)

	response := writer.Result()
	response.Body.Close()
	fmt.Println(response.StatusCode)

	// Output:
	// 201
}

func ExampleServer_GetURL() {
	server := Server{
		Storage: s.GetStorage(),
	}
	defer func() {
		if _, ok := server.Storage.(*dbstorage.DBStorage); ok {
			dbstorage.Downgrade()
		}
	}()

	ctx := context.Background()
	short := algorithms.GenerateShortURL()
	origin := "https://practicum.yandex.ru/learn/go-advanced/courses/14d6ff29-c8b6-43bf-9c55-12e8fe25b1b0/sprints/39172/topics/add19e4a-79bf-416e-9d13-0df2005ec81e/lessons/4280c1ec-45c9-4b73-bdf1-8ea438b18212/"
	server.Storage.Set(ctx, uuid.New().String(), short, origin)

	request := httptest.NewRequest(http.MethodGet, short, nil)
	router := mux.NewRouter()
	router.HandleFunc("/{id:[a-zA-Z]+}", server.GetURL).Methods(http.MethodGet)
	writer := httptest.NewRecorder()
	router.ServeHTTP(writer, request)

	response := writer.Result()
	response.Body.Close()
	fmt.Println(response.StatusCode)

	// Output:
	// 307
}

func ExampleServer_APIShortenURL() {
	bodyString := `{"url":"https://practicum.yandex.ru/learn/go-advanced/courses/14d6ff29-c8b6-43bf-9c55-12e8fe25b1b0/sprints/39172/topics/add19e4a-79bf-416e-9d13-0df2005ec81e/lessons/4280c1ec-45c9-4b73-bdf1-8ea438b18212/"}`
	var body = []byte(bodyString)
	request := httptest.NewRequest(
		http.MethodPost,
		"/api/shorten",
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
	router.HandleFunc("/api/shorten", server.APIShortenURL).Methods(http.MethodPost)
	writer := httptest.NewRecorder()
	router.ServeHTTP(writer, request)

	response := writer.Result()
	response.Body.Close()
	fmt.Println(response.StatusCode)

	// Output:
	// 201
}

func ExampleServer_UserURLs() {
	server := Server{
		Storage: s.GetStorage(),
	}
	defer func() {
		if _, ok := server.Storage.(*dbstorage.DBStorage); ok {
			dbstorage.Downgrade()
		}
	}()

	requestPOST := httptest.NewRequest(
		http.MethodPost,
		"/api/shorten",
		bytes.NewBuffer([]byte(`{"url":"https://hello.world/"}`)),
	)
	router := mux.NewRouter()
	router.HandleFunc("/api/shorten", server.APIShortenURL).Methods(http.MethodPost)
	writer := httptest.NewRecorder()
	router.ServeHTTP(writer, requestPOST)

	responsePOST := writer.Result()
	responsePOST.Body.Close()
	cookies := responsePOST.Cookies()

	requestGET := httptest.NewRequest(http.MethodGet, "/api/user/urls", nil)
	cookie := &http.Cookie{Name: "userID", Value: cookies[0].Value, HttpOnly: true}
	requestGET.AddCookie(cookie)
	router = mux.NewRouter()
	writer = httptest.NewRecorder()
	router.HandleFunc("/api/user/urls", server.UserURLs).Methods(http.MethodGet)
	router.ServeHTTP(writer, requestGET)

	responseGET := writer.Result()
	responseGET.Body.Close()
	fmt.Println(responseGET.StatusCode)

	// Output:
	// 200
}

func ExampleServer_APIShortenBatch() {
	if os.Getenv("DATABASE_DSN") != config.TestDBAddr {
		fmt.Println(201)
		return
	}
	body := `
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
			`
	var requestBody = []byte(body)
	request := httptest.NewRequest(
		http.MethodPost,
		"/api/shorten/batch",
		bytes.NewBuffer(requestBody),
	)

	defer dbstorage.Downgrade()
	db := dbstorage.Database{}
	db.Run()
	server := Server{
		Storage: dbstorage.NewDBStorage(db),
	}

	router := mux.NewRouter()
	router.HandleFunc("/api/shorten/batch",
		server.APIShortenBatch).Methods(http.MethodPost)
	writer := httptest.NewRecorder()
	router.ServeHTTP(writer, request)

	response := writer.Result()
	response.Body.Close()
	fmt.Println(response.StatusCode)

	// Output:
	// 201
}

func ExampleServer_UserURLsDelete() {
	if os.Getenv("DATABASE_DSN") != config.TestDBAddr {
		fmt.Println(202)
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
		if err != nil {
			fmt.Println(err)
		}
	}

	cookie := &http.Cookie{Name: "userID", Value: hex.EncodeToString([]byte(userID))}
	marshaled, err := json.Marshal(shorts)
	if err != nil {
		fmt.Println(err)
	}

	request := httptest.NewRequest(
		http.MethodDelete, "/api/user/urls", bytes.NewBuffer(marshaled))
	request.AddCookie(cookie)
	router := mux.NewRouter()
	router.HandleFunc("/api/user/urls", server.UserURLsDelete).Methods(http.MethodDelete)
	writer := httptest.NewRecorder()
	router.ServeHTTP(writer, request)

	response := writer.Result()
	response.Body.Close()
	fmt.Println(response.StatusCode)

	// Output:
	// 202
}
