package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strconv"

	"github.com/bootcamp-go/go-web/internal/products"
	"github.com/go-chi/chi/v5"
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
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(code)
		json.NewEncoder(w).Encode(body)
	}
}

// Get returns a handler for the GET /products/{productId} route
func (h *MyHandler) GetById() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		idStr := chi.URLParam(r, "productId")

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
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(code)
			json.NewEncoder(w).Encode(body)
		} else {
			code := http.StatusNotFound
			body := MyResponse{Message: "Not Found", Data: nil}
			w.Header().Set("content-type", "application/json")
			w.WriteHeader(code)
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
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(code)
			json.NewEncoder(w).Encode(body)
		} else {
			code := http.StatusNotFound
			body := MyResponse{Message: "Not Found", Data: nil}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(code)
			json.NewEncoder(w).Encode(body)
		}

	}
}

type BodyRequestProductJSON struct {
	Name        string  `json:"name"`
	Quantity    int     `json:"quantity"`
	CodeValue   string  `json:"code_value"`
	IsPublished bool    `json:"is_published"`
	Expiration  string  `json:"expiration"`
	Price       float64 `json:"price"`
}

// Post returns a handler for the POST /products route
func (h *MyHandler) Post() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var body BodyRequestProductJSON
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("invalid body"))
			return
		}

		fmt.Println(body)

		// Get next or avaliable id
		var nextId = getNextId(&h.data)

		// Serealizing
		newProduct := products.Product{
			Id:          nextId,
			Name:        body.Name,
			Quantity:    body.Quantity,
			CodeValue:   body.CodeValue,
			IsPublished: body.IsPublished,
			Expiration:  body.Expiration,
			Price:       body.Price,
		}

		var err = validBody(&newProduct)
		var isCodeExist = codeExist(newProduct.CodeValue, &h.data)

		var code = http.StatusInternalServerError
		var bodyR any

		if err != nil {
			code = http.StatusBadRequest
			bodyR = MyResponse{Message: err.Error(), Data: nil}
		} else if isCodeExist {
			code = http.StatusConflict
			bodyR = MyResponse{Message: "El campo code_value debe ser único para cada producto.", Data: nil}
		} else {
			//Saving new product
			h.data = append(h.data, newProduct)
			// send response
			code = http.StatusCreated
			bodyR = MyResponse{Message: "OK", Data: newProduct}

		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(code)
		json.NewEncoder(w).Encode(bodyR)

	}
}

func validBody(p *products.Product) error {

	// fmt.Println("validate", *p)
	// fmt.Println("validate", p.Quantity)

	if p.Name == "" {
		return errors.New("name is required")
	}
	if p.Quantity <= 0 {
		return errors.New("quantity must be greater than 0")
	}
	if p.CodeValue == "" {
		return errors.New("codeValue is required")
	}
	// if p.IsPublished == "" {
	// 	return errors.New("IsPublishe is required")
	// }
	if !isAValidDate(p.Expiration) {
		return errors.New("expiration is required or is not valid format DD/MM/YYYY")
	}
	if p.Price <= 0.0 {
		return errors.New("price must be greater than 0.0")
	}

	return nil

}

func getNextId(lstProducts *[]products.Product) int {
	r := 1
	for _, p := range *lstProducts {
		if p.Id > r {
			break
		} else {
			r++
		}
	}

	return r
}

func codeExist(code string, lstProducts *[]products.Product) bool {
	r := false
	for _, p := range *lstProducts {
		if p.CodeValue == code {
			r = true
			break
		}
	}

	return r
}

func isAValidDate(dateStr string) bool {

	var patter = regexp.MustCompile(`^(0[1-9]|[12][0-9]|3[01])/(0[1-9]|1[0-2])/\d{4}$`)

	return patter.MatchString(dateStr)

}
