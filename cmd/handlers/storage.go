package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type Product struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Quantity    int     `json:"quantity"`
	CodeValue   string  `json:"code_value"`
	IsPublished bool    `json:"is_published"`
	Expiration  string  `json:"expiration"`
	Price       float64 `json:"price"`
}

type Storage struct {
	Products []Product
}

func GetProductsFromJSON(filepath string) (st Storage, err error) {
	jsonFile, err := os.Open(filepath)
	if err != nil {
		err = fmt.Errorf("error opening file: %w", err)
		return
	}
	defer jsonFile.Close()

	data, err := io.ReadAll(jsonFile)
	if err != nil {
		err = fmt.Errorf("error reading file: %w", err)
		return
	}

	var products []Product
	err = json.Unmarshal(data, &products)
	if err != nil {
		err = fmt.Errorf("error unmarshaling file: %w", err)
		return
	}

	st = Storage{Products: products}
	return st, nil
}

func (st *Storage) GetProductById(id int) Product {
	for _, p := range st.Products {
		if p.ID == id {
			return p
		}
	}
	var p Product
	return p
}

func (st *Storage) GetProductsLessThan(price float64) []Product {

	var foundProducts []Product

	for _, p := range st.Products {
		if p.Price <= price {
			foundProducts = append(foundProducts, p)
		}
	}
	return foundProducts
}
