package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strconv"

	"github.com/bootcamp-go/go-web/internal/products"
	"github.com/bootcamp-go/go-web/internal/repository"
	"github.com/go-chi/chi/v5"
)

// ProductHandler is a handler with a map of users as data
type ProductHandler struct {
	rp repository.ProductRepository
}

// NewHandler returns a new ProductHandler
func NewProductHandler(rp *repository.ProductRepository) *ProductHandler {

	var myHandler ProductHandler
	myHandler.rp = *rp
	return &myHandler
}

// Get returns a handler for the GET /ping route
func (h *ProductHandler) Ping() http.HandlerFunc {
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
func (h *ProductHandler) Get() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		code := http.StatusOK
		body := MyResponse{Message: "OK", Data: h.rp.GetAll()}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(code)
		json.NewEncoder(w).Encode(body)
	}
}

// Get returns a handler for the GET /products/{productId} route
func (h *ProductHandler) GetById() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		idStr := chi.URLParam(r, "productId")

		id, err := strconv.Atoi(idStr)
		if err != nil {
			fmt.Println("Error al convertir a entero:", err)
			// Manejar el error según sea necesario
			return
		}

		pRes, err := h.rp.GetById(id)

		if err == nil {
			code := http.StatusOK
			body := MyResponse{Message: "OK", Data: pRes}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(code)
			json.NewEncoder(w).Encode(body)
		} else {
			code := http.StatusNotFound
			body := MyResponse{Message: "Not Found", Data: err.Error()}
			w.Header().Set("content-type", "application/json")
			w.WriteHeader(code)
			json.NewEncoder(w).Encode(body)
		}

	}
}

// Delte returns a handler for the DELETE /products/{productId} route
func (h *ProductHandler) DeleteById() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		idStr := chi.URLParam(r, "productId")

		id, err := strconv.Atoi(idStr)
		if err != nil {
			fmt.Println("Error al convertir a entero:", err)
			// Manejar el error según sea necesario
			return
		}

		h.rp.DeleteById(id)
		code := http.StatusOK
		body := MyResponse{Message: "OK", Data: nil}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(code)
		json.NewEncoder(w).Encode(body)

	}
}

// Put returns a handler for the Put /products route
func (h *ProductHandler) Put() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var body products.Product
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("invalid body"))
			return
		}

		fmt.Println(body)

		pRes, err2 := h.rp.GetById(body.Id)

		if err2 == nil {

			// Serealizing
			newProduct := products.Product{
				Id:          body.Id,
				Name:        body.Name,
				Quantity:    body.Quantity,
				CodeValue:   body.CodeValue,
				IsPublished: body.IsPublished,
				Expiration:  body.Expiration,
				Price:       body.Price,
			}

			var err = validBody(&newProduct)
			var isCodeExist = h.rp.IsCodeExist(newProduct.CodeValue)

			var code = http.StatusInternalServerError
			var bodyR any

			if err != nil {
				code = http.StatusBadRequest
				bodyR = MyResponse{Message: err.Error(), Data: nil}
			} else if isCodeExist && newProduct.Id != pRes.Id {
				code = http.StatusConflict
				bodyR = MyResponse{Message: "El campo code_value debe ser único para cada producto.", Data: nil}
			} else {
				//Saving new product
				h.rp.DeleteById(pRes.Id)
				h.rp.Save(newProduct)
				// send response
				code = http.StatusCreated
				bodyR = MyResponse{Message: "OK", Data: newProduct}

			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(code)
			json.NewEncoder(w).Encode(bodyR)

		} else {
			code := http.StatusNotFound
			body := MyResponse{Message: "Not Found", Data: err2.Error()}
			w.Header().Set("content-type", "application/json")
			w.WriteHeader(code)
			json.NewEncoder(w).Encode(body)
		}
	}
}

// Patch returns a handler for the Patch /products/ route
func (h *ProductHandler) Patch() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var body products.Product
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("invalid body"))
			return
		}

		fmt.Println(body)

		pRes, err2 := h.rp.GetById(body.Id)

		if err2 == nil {

			// Serealizing
			newProduct := products.Product{
				Id:          body.Id,
				Name:        body.Name,
				Quantity:    body.Quantity,
				CodeValue:   body.CodeValue,
				IsPublished: body.IsPublished,
				Expiration:  body.Expiration,
				Price:       body.Price,
			}

			var err = validBody(&newProduct)
			var isCodeExist = h.rp.IsCodeExist(newProduct.CodeValue)

			var code = http.StatusInternalServerError
			var bodyR any

			if err != nil {
				code = http.StatusBadRequest
				bodyR = MyResponse{Message: err.Error(), Data: nil}
			} else if isCodeExist && newProduct.Id != pRes.Id {
				code = http.StatusConflict
				bodyR = MyResponse{Message: "El campo code_value debe ser único para cada producto.", Data: nil}
			} else {
				//Saving new product
				h.rp.Save(newProduct)
				// send response
				code = http.StatusCreated
				bodyR = MyResponse{Message: "OK", Data: newProduct}

			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(code)
			json.NewEncoder(w).Encode(bodyR)

		} else {
			code := http.StatusNotFound
			body := MyResponse{Message: "Not Found", Data: err2.Error()}
			w.Header().Set("content-type", "application/json")
			w.WriteHeader(code)
			json.NewEncoder(w).Encode(body)
		}

	}
}

// GetByQuery returns a handler for the GET /products?productId={id} route
func (h *ProductHandler) GetByQuery() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		idStr := r.URL.Query().Get("productId")
		priceGtStr := r.URL.Query().Get("priceGt")

		id, _ := strconv.Atoi(idStr)
		priceGt, _ := strconv.Atoi(priceGtStr)

		pRes, err := h.rp.GetByQuery(id, float64(priceGt))

		if err == nil {
			code := http.StatusOK
			body := MyResponse{Message: "OK", Data: pRes}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(code)
			json.NewEncoder(w).Encode(body)
		} else {
			code := http.StatusNotFound
			body := MyResponse{Message: "Not Found", Data: err.Error()}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(code)
			json.NewEncoder(w).Encode(body)
		}

	}
}

type BodyRequestNewProductJSON struct {
	Name        string  `json:"name"`
	Quantity    int     `json:"quantity"`
	CodeValue   string  `json:"code_value"`
	IsPublished bool    `json:"is_published"`
	Expiration  string  `json:"expiration"`
	Price       float64 `json:"price"`
}

// Post returns a handler for the POST /products route
func (h *ProductHandler) Post() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var body BodyRequestNewProductJSON
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("invalid body"))
			return
		}

		fmt.Println(body)

		// Get next or avaliable id
		var nextId = h.rp.GetNextId()

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
		var isCodeExist = h.rp.IsCodeExist(newProduct.CodeValue)

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
			h.rp.Save(newProduct)
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

func isAValidDate(dateStr string) bool {

	var patter = regexp.MustCompile(`^(0[1-9]|[12][0-9]|3[01])/(0[1-9]|1[0-2])/\d{4}$`)

	return patter.MatchString(dateStr)

}
