package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/pmohanj/go-microservices/data"
)

type KeyProduct struct{}

// This type is for handle Products
type Products struct {
	l *log.Logger
	v *data.Validation
}

// This is kind of like returning a obj of Products struct
func NewProduct(l *log.Logger, v *data.Validation) *Products {
	return &Products{l, v}
}

var ErrInvalidProductPath = fmt.Errorf("Invalid path, path should be /products/[id]")

type GenericError struct {
	Message string `json:"message"`
}

type ValidationError struct {
	Messages []string `json:"messages"`
}

// getProductID returns the product ID from the URL
// Panics if cannot convert the id into an integer
// this should never happen as the router ensures that
// this is a valid number
func getProductID(r *http.Request) int {
	// parse the product id from the url
	vars := mux.Vars(r)

	// convert the id into an integer and return
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		// should never happen
		panic(err)
	}

	return id
}
