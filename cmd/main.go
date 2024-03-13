package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/noe-suaste-ML/go-web-practica/cmd/handlers"
)

const (
	productsFilepath = "docs/db/products.json"
)

func main() {

	// Make slice of products from JSON file
	st, err := handlers.GetProductsFromJSON(productsFilepath)
	if err != nil {
		log.Println(err)
		return
	}

	// Create a new products handler
	ph := handlers.NewProductsHandler(st)

	// Setting up the server
	rt := chi.NewRouter()

	rt.Route("/products", func(rt chi.Router) {
		// GET all products
		rt.Get("/", ph.GetProducts())
		// GET product by ID
		rt.Get("/{id}", ph.GetProductByID())
		// GET products less than a price
		rt.Get("/search", ph.SearchProducts())
		// GET ping
		rt.Get("/ping", ph.Ping())
		// POST create product
		rt.Post("/", ph.CreateProduct())
	})

	if err := http.ListenAndServe(":8080", rt); err != nil {
		panic(err)
	}

}
