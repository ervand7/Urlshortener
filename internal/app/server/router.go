package server

import (
	"github.com/ervand7/urlshortener/internal/app/views"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func newRouter() chi.Router {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	storage := views.URLStorage{HashTable: make(map[string]string, 0)}
	r.Route("/", func(r chi.Router) {
		r.Post("/", storage.URLShorten())
		r.Get("/{id:[a-zA-Z]+}", storage.URLGet())
	})

	return r
}
