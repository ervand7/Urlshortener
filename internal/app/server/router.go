package server

import (
	"github.com/ervand7/urlshortener/internal/app/models"
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

	storage := make(models.Storage, 0)
	r.Route("/", func(r chi.Router) {
		r.Post("/", views.URLShorten(storage))
		r.Get("/{id:[a-zA-Z]+}", views.URLGet(storage))
	})

	return r
}
