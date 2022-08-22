package server

import (
	"github.com/ervand7/urlshortener/internal/app/models/url"
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

	server := views.Server{
		MemoryStorage: &url.MemoryStorage{
			HashTable: make(map[string]string, 0),
		},
		FileStorage: &url.FileStorage{},
	}
	r.Route("/", func(r chi.Router) {
		r.Post("/", server.URLShorten())
		r.Post("/api/shorten", server.URLShortenJSON())
		r.Get("/{id:[a-zA-Z]+}", server.URLGet())
	})

	return r
}
