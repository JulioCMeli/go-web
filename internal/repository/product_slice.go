package repository

import (
	"errors"
	"fmt"

	"github.com/bootcamp-go/go-web/internal/products"
	"github.com/bootcamp-go/go-web/internal/storage"
)

// ProductRepository is a struct with a slice of products
type ProductRepository struct {
	data []products.Product
}

// NewProductRepository returns a new ProductRepository
func NewProductRepository() *ProductRepository {

	var repository ProductRepository

	repository.data, _ = storage.ReadJson()

	return &repository
}

func (r *ProductRepository) GetAll() []products.Product {
	return r.data
}

func (r *ProductRepository) GetById(id int) (products.Product, error) {

	fmt.Println("GetById ", id)
	var pRes products.Product
	wasFound := false
	for _, p := range r.data {
		if p.Id == id {
			pRes = p
			wasFound = true
			break
		}
	}
	fmt.Println("pRes", pRes)

	if wasFound {
		return pRes, nil
	}
	return pRes, errors.New("product with id not found")

}

func (h *ProductRepository) GetByQuery(id int, priceGt float64) (products.Product, error) {

	var pRes products.Product
	wasFound := false
	for _, p := range h.data {
		if p.Id == id && p.Price >= priceGt {
			pRes = p
			wasFound = true
			break
		}
	}
	if !wasFound {
		return pRes, errors.New("product with this query not found")
	}

	return pRes, nil

}

// Post returns a handler for the POST /products route
func (h *ProductRepository) Save(product products.Product) error {

	wasFound := false
	index := 0
	for i, p := range h.data {
		if p.Id == product.Id {
			index = i
			wasFound = true
			break
		}
	}

	if wasFound {
		h.data[index] = product
	} else {
		h.data = append(h.data, product)
	}
	storage.SaveJson(h.data)
	return nil
}

func (h *ProductRepository) GetNextId() int {
	return h.data[len(h.data)-1].Id + 1
}

func (h *ProductRepository) DeleteById(id int) {

	for index, p := range h.data {
		if p.Id == id {
			h.data = append(h.data[:index], h.data[index+1:]...)

		}
	}
	storage.SaveJson(h.data)

}

func (h *ProductRepository) IsCodeExist(code string) bool {
	r := false
	for _, p := range h.data {
		if p.CodeValue == code {
			r = true
			break
		}
	}

	return r
}
