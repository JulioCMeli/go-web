package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	//"github.com/bootcamp-go/go-web/internal/products"
	"github.com/bootcamp-go/go-web/internal/products"
	//. "github.com/bootcamp-go/go-web/internal/products"
)

// MyHandler is a handler with a map of users as data
type MyHandler struct {
	data []products.Product
}

// NewHandler returns a new MyHandler
func NewHandler() *MyHandler {

	var myHandler MyHandler

	myHandler.data, _ = products.ReadJson()

	return &myHandler
}

// Get returns a handler for the GET /ping route
func (h *MyHandler) Ping() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	}
}

// MyResponse is an struct for the response
type MyResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// Get returns a handler for the GET /products route
func (h *MyHandler) Get() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		code := http.StatusOK
		body := MyResponse{Message: "OK", Data: h.data}
		w.WriteHeader(code)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(body)
	}
}

// Get returns a handler for the GET /products/{productId} route
func (h *MyHandler) GetById() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		//idStr := chi.URLParam(r, "productId")
		idStr := "123"

		id, err := strconv.Atoi(idStr)
		if err != nil {
			fmt.Println("Error al convertir a entero:", err)
			// Manejar el error según sea necesario
			return
		}

		fmt.Println(id)
		var pRes products.Product
		wasFound := false
		for _, p := range h.data {
			if p.Id == id {
				pRes = p
				wasFound = true
				break
			}
		}

		if wasFound {
			code := http.StatusOK
			body := MyResponse{Message: "OK", Data: pRes}
			w.WriteHeader(code)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(body)
		} else {
			code := http.StatusNotFound
			body := MyResponse{Message: "Not Found", Data: nil}
			w.WriteHeader(code)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(body)
		}

	}
}

// GetByQuery returns a handler for the GET /products?productId={id} route
func (h *MyHandler) GetByQuery() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		idStr := r.URL.Query().Get("productId")
		priceGtStr := r.URL.Query().Get("priceGt")

		id, _ := strconv.Atoi(idStr)
		priceGt, _ := strconv.Atoi(priceGtStr)

		fmt.Println(id)
		fmt.Println(priceGt)
		var pRes products.Product
		wasFound := false
		for _, p := range h.data {
			if p.Id == id && p.Price >= float64(priceGt) {
				pRes = p
				wasFound = true
				break
			}
		}
		if wasFound {
			code := http.StatusOK
			body := MyResponse{Message: "OK", Data: pRes}
			w.WriteHeader(code)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(body)
		} else {
			code := http.StatusNotFound
			body := MyResponse{Message: "Not Found", Data: nil}
			w.WriteHeader(code)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(body)
		}

	}
}
