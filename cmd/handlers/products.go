package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type ProductsHandler struct {
	St Storage
}

func NewProductsHandler(st Storage) *ProductsHandler {
	return &ProductsHandler{St: st}
}

func (ph *ProductsHandler) Ping() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("pong"))
	}
}

func (ph *ProductsHandler) GetProducts() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := json.NewEncoder(w).Encode(ph.St.Products)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func (ph *ProductsHandler) GetProductByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, _ := strconv.Atoi(chi.URLParam(r, "id"))
		err := json.NewEncoder(w).Encode(ph.St.GetProductById(id))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func (ph *ProductsHandler) SearchProducts() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		price, _ := strconv.ParseFloat(r.URL.Query().Get("price"), 64)
		err := json.NewEncoder(w).Encode(ph.St.GetProductsLessThan(float64(price)))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
