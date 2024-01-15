package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/bootcamp-go/go-web/internal/handler"
	"github.com/bootcamp-go/go-web/internal/repository"
	"github.com/bootcamp-go/go-web/internal/storage"

	"github.com/go-chi/chi/v5"
)

func main() {

	var token string = "GO_API_TOKEN_1"
	os.Setenv(token, "julio")

	// Obtener el valor de una variable de entorno
	tokenValue := os.Getenv(token)
	fmt.Println("Valor de la token/variable de entorno:", tokenValue)

	auth := handler.NewAuthenticator(tokenValue)

	r := chi.NewRouter()

	r.Use(auth.Auth)

	data, _ := storage.ReadJson()

	rp := repository.NewProductRepository(&data)

	h := handler.NewProductHandler(rp)

	r.Get("/ping", h.Ping())
	r.Get("/", h.Get())
	r.Get("/products/{productId}", h.GetById())
	r.Get("/products", h.GetByQuery())
	r.Post("/products", h.Post())
	r.Delete("/products/{productId}", h.DeleteById())
	r.Put("/products", h.Put())
	r.Patch("/products", h.Patch())

	http.ListenAndServe(":8080", r)
}
