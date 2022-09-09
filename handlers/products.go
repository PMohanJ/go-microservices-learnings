package handlers

import (
	"log"
	"net/http"
	"regexp"
	"strconv"

	"github.com/pmohanj/go-microservices/data"
)

// This type is for handle Products
type Products struct {
	l *log.Logger
}

// This is kind of like returning a obj of Products struct
func NewProduct(l *log.Logger) *Products {
	return &Products{l}
}

// This ServeHTTP is the actual handlers that will be invoked when a http req is
// made to the URI path, as servehttp implements handler interface
func (p *Products) ServeHTTP(rw http.ResponseWriter, r *http.Request) {

	// Retrieve products
	if r.Method == http.MethodGet {
		p.getProducts(rw, r)
		return
	}

	// To create (add) poducts
	if r.Method == http.MethodPost {
		p.addProduct(rw, r)
		return
	}

	// For PUT method we need id of the peoduct to update
	if r.Method == http.MethodPut {
		// expect the id in the URI
		reg := regexp.MustCompile(`/([0-9]+)`)
		g := reg.FindAllStringSubmatch(r.URL.Path, -1)

		if len(g) != 1 {
			p.l.Println("Invalid URL more than one id")
			http.Error(rw, "Invalid URI", http.StatusBadRequest)
			return
		}

		if len(g[0]) != 2 {
			p.l.Println("Invalid URL more than one capture group")
			http.Error(rw, "Invalid URI", http.StatusBadRequest)
			return
		}

		// converting aquired id to int format
		idString := g[0][1]
		id, err := strconv.Atoi(idString)
		if err != nil {
			p.l.Println("Invalid URL unable to convert to number", idString)
			http.Error(rw, "Invalid URI", http.StatusBadRequest)
			return
		}

		p.updateProducts(id, rw, r)
		return

	}

	rw.WriteHeader(http.StatusMethodNotAllowed)
}

func (p *Products) getProducts(rw http.ResponseWriter, r *http.Request) {
	lp := data.GetProducts()

	err := lp.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Unable to marshal json", http.StatusInternalServerError)
	}
}

func (p *Products) addProduct(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle POST Products")

	prod := &data.Product{}
	err := prod.FromJSON(r.Body)
	if err != nil {
		http.Error(rw, "Unable to unmarshal json", http.StatusBadRequest)
		return
	}

	p.l.Printf("Prod: %v", prod)
	data.AddProduct(prod)
}

func (p *Products) updateProducts(id int, rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle PUT Products")

	prod := &data.Product{}
	err := prod.FromJSON(r.Body)
	if err != nil {
		http.Error(rw, "Unable to unmarshal json", http.StatusBadRequest)
		return
	}

	// calling method to update product
	err = data.UpdateProduct(id, prod)
	if err == data.ErrProcutNotFound {
		http.Error(rw, "Product not found", http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(rw, "Product some error", http.StatusInternalServerError)
		return
	}
}
