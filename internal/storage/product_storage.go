package storage

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/bootcamp-go/go-web/internal/products"
)

var fileName = "/Users/jdiazaltamir/Desktop/go-web/products2.json"

func ReadJson() ([]products.Product, error) {
	// Leer el contenido del archivo JSON
	fileContent, err := os.ReadFile(fileName)
	if err != nil {
		fmt.Println("Error al leer el archivo:", err)
		return nil, err
	}

	// Crear un slice para almacenar los datos parseados
	var products []products.Product

	// Parsear el contenido del archivo JSON en el slice de personas
	err = json.Unmarshal(fileContent, &products)
	if err != nil {
		fmt.Println("Error al parsear el archivo JSON:", err)
		return nil, err
	}

	// Imprimir el slice de personas
	// fmt.Println("Datos parseados:")
	// for _, p := range products {
	// 	fmt.Printf("%v \n", p)
	// }

	return products, nil
}

func SaveJson(products []products.Product) error {
	// Leer el contenido del archivo JSON
	content, _ := json.MarshalIndent(products, "", " ")

	err := os.WriteFile(fileName, content, 0777)
	if err != nil {
		fmt.Println("Error al salvar el archivo:", err)
		return err
	}

	return nil
}
