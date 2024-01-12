package main

import (
	"net/http"

	"github.com/bootcamp-go/go-web/internal/handler"
	"github.com/bootcamp-go/go-web/internal/repository"

	"github.com/go-chi/chi/v5"
)

func main() {

	r := chi.NewRouter()

	rp := repository.NewProductRepository()

	h := handler.NewProductHandler(rp)

	r.Get("/ping", h.Ping())
	r.Get("/", h.Get())
	r.Get("/products/{productId}", h.GetById())
	r.Get("/products", h.GetByQuery())
	r.Post("/products", h.Post())
	r.Delete("/products/{productId}", h.DeleteById())

	http.ListenAndServe(":8080", r)
}
