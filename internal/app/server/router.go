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
	r.Use(GzipMiddleware)

	server := views.Server{
		Storage: GetServerStorage(),
	}

	r.Route("/", func(r chi.Router) {
		r.Post("/", server.URLShorten)
		r.Post("/api/shorten", server.URLShortenJSON)
		r.Post("/api/shorten/batch", server.URLShortenBatch)
		r.Get("/{id:[a-zA-Z]+}", server.URLGet)
		r.Get("/api/user/urls", server.URLUserAll)
		r.Get("/ping", server.DBPing)
	})

	return r
}
