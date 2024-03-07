package main

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/go-chi/chi/v5"
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

func getProductsFromJSON(products *[]Product) {
	filename := "products.json"

	jsonFile, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer jsonFile.Close()

	data, err := io.ReadAll(jsonFile)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(data, &products)
	if err != nil {
		panic(err)
	}
}

func getProductById(products []Product, id int) Product {
	for _, p := range products {
		if p.ID == id {
			return p
		}
	}
	var p Product
	return p
}

func getProductsLessThan(products []Product, price float64) []Product {

	var foundProducts []Product

	for _, p := range products {
		if p.Price <= price {
			foundProducts = append(foundProducts, p)
		}
	}
	return foundProducts
}

func main() {

	// Make slice of products from JSON file
	var products []Product
	getProductsFromJSON(&products)

	// Setting up the server
	rt := chi.NewRouter()

	rt.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("pong"))
	})

	rt.Get("/products", func(w http.ResponseWriter, r *http.Request) {
		err := json.NewEncoder(w).Encode(products)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	rt.Get("/products/{id}", func(w http.ResponseWriter, r *http.Request) {
		id, _ := strconv.Atoi(chi.URLParam(r, "id"))
		err := json.NewEncoder(w).Encode(getProductById(products, id))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	rt.Get("/search", func(w http.ResponseWriter, r *http.Request) {
		price, _ := strconv.ParseFloat(r.URL.Query().Get("price"), 32)
		err := json.NewEncoder(w).Encode(getProductsLessThan(products, float64(price)))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	if err := http.ListenAndServe(":8080", rt); err != nil {
		panic(err)
	}

}
