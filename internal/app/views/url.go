package views

import (
	"github.com/ervand7/urlshortener/internal/app/controllers"
	"io"
	"net/http"
	"sync"
)

type URLStorage struct {
	HashTable map[string]string
	mutex     sync.Mutex
}

func (s *URLStorage) Get(short string) (origin string) {
	origin, exist := s.HashTable[short]
	if !exist {
		return ""
	}
	return origin
}

func (s *URLStorage) Set(short, origin string) {
	s.HashTable[short] = origin
}

// URLShorten POST ("/")
func (s *URLStorage) URLShorten() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		s.mutex.Lock()
		defer r.Body.Close()
		defer s.mutex.Unlock()

		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "invalid body", http.StatusBadRequest)
			return
		}
		url := string(body)
		if url == "" {
			http.Error(w, "param url not filled", http.StatusBadRequest)
			return
		}
		w.Header().Add("Content-type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusCreated)

		shortenedURL := controllers.ShortenURL()
		s.Set(shortenedURL, url)
		w.Write([]byte(shortenedURL))
	}
}

// URLGet GET /{id}
func (s *URLStorage) URLGet() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		s.mutex.Lock()
		defer s.mutex.Unlock()

		endpoint := r.URL.Path
		shortenedURL := controllers.BaseURL + endpoint
		originURL := s.Get(shortenedURL)
		if originURL == "" {
			http.Error(w, "not exists", http.StatusBadRequest)
			return
		}
		http.Redirect(
			w,
			r,
			originURL,
			http.StatusTemporaryRedirect,
		)
	}
}
