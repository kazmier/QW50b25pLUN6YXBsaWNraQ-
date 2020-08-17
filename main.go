package main

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"

	"github.com/QW50b25pLUN6YXBsaWNraQ-/url"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Post("/api/fetcher", url.SaveUrl)
	r.Get("/api/fetcher/{id}", url.GetUrl)
	r.Delete("/api/fetcher/{id}", url.DeleteUrl)
	r.Get("/api/fetcher", url.GetAllUrls)
	r.Get("/api/fetcher/{id}/history", url.GetUrlHistory)

	http.ListenAndServe(":8080", r)
}
