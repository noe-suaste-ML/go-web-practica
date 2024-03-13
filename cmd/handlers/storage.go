package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"
)

type Product struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Quantity    int       `json:"quantity"`
	CodeValue   string    `json:"code_value"`
	IsPublished bool      `json:"is_published"`
	Expiration  time.Time `json:"expiration"`
	Price       float64   `json:"price"`
}

type ProductJSON struct {
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

	var productsJSON []ProductJSON
	err = json.NewDecoder(jsonFile).Decode(&productsJSON)
	if err != nil {
		err = fmt.Errorf("error decoding file: %w", err)
		return
	}

	var products []Product
	for _, p := range productsJSON {
		expDate, errDate := time.Parse("02/01/2006", p.Expiration)
		if errDate != nil {
			err = fmt.Errorf("error parsing date: %w", errDate)
			return
		}
		products = append(products, Product{
			ID:          p.ID,
			Name:        p.Name,
			Quantity:    p.Quantity,
			CodeValue:   p.CodeValue,
			IsPublished: p.IsPublished,
			Expiration:  expDate,
			Price:       p.Price,
		})
	}
	st = Storage{Products: products}
	return st, nil
}

var (
	ErrorRequiredField = errors.New("error required field")
	ErrorInvalidField  = errors.New("error invalid field")
)

func (st *Storage) ValidateProduct(p Product) error {
	if p.Name == "" {
		return ErrorRequiredField
	}
	if p.Quantity < 0 {
		return ErrorInvalidField
	}
	if p.CodeValue == "" {
		return ErrorRequiredField
	}
	for _, product := range st.Products {
		if product.CodeValue == p.CodeValue {
			return ErrorInvalidField
		}
	}
	if p.Expiration.IsZero() {
		return ErrorRequiredField
	}
	if p.Price <= 0 {
		return ErrorInvalidField
	}
	return nil
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
