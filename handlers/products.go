package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/mic3ael/product-api/data"
)

type Products struct {
	logger *log.Logger
}

func NewProducts(l *log.Logger) *Products {
	return &Products{l}
}

func (p *Products) GetProducts(w http.ResponseWriter, r *http.Request) {
	p.logger.Println("Handle GET Product")

	lp := data.GetProducts()

	err := lp.ToJSON(w)

	if err != nil {
		http.Error(w, "Unable to marshal json", http.StatusInternalServerError)
		return
	}
}

func (p *Products) AddProducts(w http.ResponseWriter, req *http.Request) {
	p.logger.Println("Handle POST Product")
	prod := req.Context().Value(KeyProduct{}).(data.Product)
	data.AddProduct(&prod)
}

func (p *Products) UpdateProducts(w http.ResponseWriter, req *http.Request) {

	vars := mux.Vars(req)
	rID := vars["id"]
	id, err := strconv.Atoi(rID)

	if err != nil {
		http.Error(w, fmt.Sprintf("product id - %s has an invalid format", rID), http.StatusBadRequest)
		return
	}

	p.logger.Println("Handle PUT Product")
	prod := req.Context().Value(KeyProduct{}).(data.Product)

	err = data.UpdateProduct(id, &prod)
	if err == data.ErrProductNotFound {
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "Product not found", http.StatusInternalServerError)
		return
	}
}

type KeyProduct struct{}

func (p Products) MiddlewareValidateProduct(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		prod := data.Product{}

		err := prod.FromJSON(r.Body)

		if err != nil {
			http.Error(rw, "Unable to unmarshal json", http.StatusBadRequest)
			return
		}

		ctx := context.WithValue(r.Context(), KeyProduct{}, prod)
		req := r.WithContext(ctx)

		next.ServeHTTP(rw, req)
	})
}
