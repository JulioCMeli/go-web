package main

import (
	"net/http"

	"github.com/bootcamp-go/go-web/internal/handler"

	"github.com/go-chi/chi/v5"
)

func main() {

	r := chi.NewRouter()

	h := handler.NewHandler()

	r.Get("/ping", h.Ping())
	r.Get("/", h.Get())
	r.Get("/products/{productId}", h.GetById())
	r.Get("/products", h.GetByQuery())

	http.ListenAndServe(":8080", r)
}
