package handlers

import (
	"log"
	"net/http"
	"regexp"
	"strconv"

	"github.com/mic3ael/product-api/data"
)

type Products struct {
	logger *log.Logger
}

func NewProducts(l *log.Logger) *Products {
	return &Products{l}
}

func (p *Products) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {
		p.getProducts(w, r)
		return
	}

	if r.Method == http.MethodPost {
		p.addProducts(w, r)
		return
	}

	if r.Method == http.MethodPut {
		reg := regexp.MustCompile(`/([0-9]+)`)
		g := reg.FindAllStringSubmatch(r.URL.Path, -1)

		if len(g) != 1 {
			p.logger.Println("Invalid URI more than one id")
			http.Error(w, "Invalid URI", http.StatusBadRequest)
			return
		}

		if len(g[0]) != 2 {
			p.logger.Println("Invalid URI more than one capture group")
			http.Error(w, "Invalid URI", http.StatusBadRequest)
			return
		}

		idString := g[0][1]
		id, err := strconv.Atoi(idString)

		if err != nil {
			p.logger.Println("Invalid URI unable to convert to number - ", idString)
			http.Error(w, "Invalid URI", http.StatusBadRequest)
			return
		}

		p.logger.Println("got it", id)

		p.updateProducts(id, w, r)
		return
	}

	w.WriteHeader(http.StatusMethodNotAllowed)
}

func (p *Products) getProducts(w http.ResponseWriter, r *http.Request) {
	p.logger.Println("Handle GET Product")

	lp := data.GetProducts()

	err := lp.ToJSON(w)

	if err != nil {
		http.Error(w, "Unable to marshal json", http.StatusInternalServerError)
		return
	}
}

func (p *Products) addProducts(w http.ResponseWriter, req *http.Request) {
	p.logger.Println("Handle POST Product")

	prod := &data.Product{}
	err := prod.FromJSON(req.Body)

	if err != nil {
		http.Error(w, "Unable to unmarshal json", http.StatusBadRequest)
		return
	}

	p.logger.Printf("Prod: %#v\n", prod)
	data.AddProduct(prod)
}

func (p *Products) updateProducts(id int, w http.ResponseWriter, req *http.Request) {
	p.logger.Println("Handle PUT Product")

	prod := &data.Product{}
	err := prod.FromJSON(req.Body)

	if err != nil {
		http.Error(w, "Unable to unmarshal json", http.StatusBadRequest)
		return
	}

	err = data.UpdateProduct(id, prod)
	if err == data.ErrProductNotFound {
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "Product not found", http.StatusInternalServerError)
		return
	}
}
