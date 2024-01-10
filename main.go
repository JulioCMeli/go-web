package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/bootcamp-go/go-web/internal/products"
	"github.com/go-chi/chi/v5"
)

func main() {

	lstProducts, err := products.ReadJson()

	if err != nil {
		os.Exit(1)
	}

	// router
	router := chi.NewRouter()
	// - register endpoints
	// - ping
	router.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		// request
		fmt.Println("GET /ping")
		fmt.Println("method:", r.Method)
		fmt.Println("url:", r.URL)
		fmt.Println("header:", r.Header)

		w.Write([]byte(`{"message": "pong"}`))
	})

	http.ListenAndServe(":8080", router)
}
