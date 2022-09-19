package handlers

import (
	"net/http"

	"github.com/pmohanj/go-microservices/data"
)

// swagger:route GET /products products listProducts
// Returns a list of products
// responses:
// 	200: productsResponse

// ListAll returns the products from the data store
func (p *Products) ListAll(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("[DEBUG] get all records")

	prods := data.GetProducts()

	err := data.ToJSON(prods, rw)
	if err != nil {
		// we should never be here but log the error just incase
		p.l.Println("[ERROR] serializing product", err)
	}
}

// swagger:route GET /products/{id} products listSingleProduct
// Return a single product form database
// responses:
//  200: productResponse
//  404: errorResponse

// ListSingle handles Get Requests
func (p *Products) ListSingle(rw http.ResponseWriter, r *http.Request) {
	id := getProductID(r)

	p.l.Println("[Debug] get record id", id)

	prod, err := data.GetProductByID(id)

	switch err {
	case nil:
	case data.ErrProductNotFound:
		p.l.Println("[ERROR] fetching product", err)

		rw.WriteHeader(http.StatusNotFound)
		data.ToJSON(&GenericError{Message: err.Error()}, rw)
		return
	default:
		p.l.Println("[ERROR] fetching product", err)

		rw.WriteHeader(http.StatusInternalServerError)
		data.ToJSON(&GenericError{Message: err.Error()}, rw)
		return
	}

	err = data.ToJSON(prod, rw)
	if err != nil {
		// we should never be here but log the error just incase
		p.l.Println("[ERROR] serializing product", err)
	}
}
