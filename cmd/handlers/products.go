package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

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

type RequestProductCreate struct {
	Name        string  `json:"name"`
	Quantity    int     `json:"quantity"`
	CodeValue   string  `json:"code_value"`
	IsPublished bool    `json:"is_published"`
	Expiration  string  `json:"expiration"`
	Price       float64 `json:"price"`
}

type ResponseProductCreate struct {
	Message string   `json:"message"`
	Data    *Product `json:"data"`
}

func (ph *ProductsHandler) CreateProduct() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var p RequestProductCreate
		err := json.NewDecoder(r.Body).Decode(&p)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(ResponseProductCreate{Message: "invalid request body", Data: nil})
			return
		}

		expDate, err := time.Parse("02/01/2006", p.Expiration)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(ResponseProductCreate{Message: "invalid date format", Data: nil})
			return
		}

		product := &Product{
			ID:          len(ph.St.Products) + 1,
			Name:        p.Name,
			Quantity:    p.Quantity,
			CodeValue:   p.CodeValue,
			IsPublished: p.IsPublished,
			Expiration:  expDate,
			Price:       p.Price,
		}

		if err := ph.St.ValidateProduct(*product); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(ResponseProductCreate{Message: err.Error(), Data: nil})
			return
		}

		ph.St.Products = append(ph.St.Products, *product)
		w.WriteHeader(http.StatusCreated)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(ResponseProductCreate{Message: "product created", Data: product})
	}
}
