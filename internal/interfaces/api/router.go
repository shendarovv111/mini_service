package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewRouter(handler *Handler) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.Logger, middleware.Recoverer)

	r.Route("/quotes", func(r chi.Router) {
		r.Post("/", handler.CreateQuote)
		r.Get("/", handler.GetAllQuotes)
		r.Get("/random", handler.GetRandomQuote)
		r.Delete("/{id}", handler.DeleteQuote)
	})

	return r
}
