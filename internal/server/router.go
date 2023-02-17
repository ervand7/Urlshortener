package server

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/ervand7/urlshortener/internal/controllers/storage"
	"github.com/ervand7/urlshortener/internal/server/middlewares"
	"github.com/ervand7/urlshortener/internal/views"
)

func newRouter() chi.Router {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middlewares.GzipMiddleware)
	// r.Use(middlewares.NewTrustedNetwork().Handler)

	server := views.Server{
		Storage: storage.GetStorage(),
	}

	r.Route("/", func(r chi.Router) {
		r.Post("/", server.ShortenURL)
		r.Get("/{id:[a-zA-Z]+}", server.GetURL)
		r.Post("/api/shorten", server.APIShortenURL)
		r.Get("/api/user/urls", server.UserURLs)
		r.Post("/api/shorten/batch", server.APIShortenBatch)
		r.Delete("/api/user/urls", server.UserURLsDelete)
		r.Get("/ping", server.PingDB)
		// r.Get("/api/internal/stats", server.Stats)
	})

	StatsGroup := r.Group(nil)
	StatsGroup.Use(middlewares.NewTrustedNetwork().Handler)
	StatsGroup.Post("/api/internal/stats", server.Stats)

	return r
}
