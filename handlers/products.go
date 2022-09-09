package handlers

import (
	"log"
	"net/http"

	"github.com/pmohanj/go-microservices/data"
)

// This type is for handle Products
type Products struct {
	l *log.Logger
}

// This is kind of like returning a obj of Products handler
func NewProduct(l *log.Logger) *Products {
	return &Products{l}
}

// This ServeHTTP is the actual handlers that will be invoked when a http req is
// made to the URI path as servehttp implements handler interface
func (p *Products) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		p.getProducts(rw, r)
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
