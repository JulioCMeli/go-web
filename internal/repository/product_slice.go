package repository

import (
	"errors"
	"fmt"

	"github.com/bootcamp-go/go-web/internal/products"
)

// ProductRepository is a struct with a slice of products
type ProductRepository struct {
	data []products.Product
}

// NewProductRepository returns a new ProductRepository
func NewProductRepository() *ProductRepository {

	var repository ProductRepository

	repository.data, _ = products.ReadJson()

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
func (h *ProductRepository) Save(p products.Product) error {
	h.data = append(h.data, p)
	return nil
}

func (h *ProductRepository) GetNextId() int {
	r := 1
	for _, p := range h.data {
		if p.Id > r {
			break
		} else {
			r++
		}
	}

	return r
}

func (h *ProductRepository) DeleteById(id int) {

	for index, p := range h.data {
		if p.Id == id {
			h.data = append(h.data[:index], h.data[index+1:]...)

		}
	}

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
