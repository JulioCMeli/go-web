package products

import (
	"encoding/json"
	"fmt"
	"os"
)

type Product struct {
	Id          int     `json:"id"`
	Name        string  `json:"name"`
	Quantity    int     `json:"quantity"`
	CodeValue   string  `json:"code_value"`
	IsPublished bool    `json:"is_published"`
	Expiration  string  `json:"expiration"`
	Price       float64 `json:"price"`
}

func ReadJson() ([]Product, error) {
	// Leer el contenido del archivo JSON
	fileContent, err := os.ReadFile("products.json")
	if err != nil {
		fmt.Println("Error al leer el archivo:", err)
		return nil, err
	}

	// Crear un slice para almacenar los datos parseados
	var products []Product

	// Parsear el contenido del archivo JSON en el slice de personas
	err = json.Unmarshal(fileContent, &products)
	if err != nil {
		fmt.Println("Error al parsear el archivo JSON:", err)
		return nil, err
	}

	// Imprimir el slice de personas
	fmt.Println("Datos parseados:")
	for _, p := range products {
		fmt.Printf("%v \n", p)
	}

	return products, nil
}
