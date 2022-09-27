package handlers

import (
	"net/http"

	"github.com/pmohanj/go-microservices/product-api/data"
)

// swagger:route PUT /products products updateProduct
// Update a product details
//
// reponses:
//	201: noContentResponse
// 404: errorResponse
// 422: errorValidation

// Update handles PUT requests to update products
func (p *Products) Update(rw http.ResponseWriter, r *http.Request) {
	id := getProductID(r)

	// fetch the product from the context
	prod := r.Context().Value(KeyProduct{}).(*data.Product)
	p.l.Println("[DEBUG] updating record id", id)

	err := data.UpdateProduct(id, prod)
	if err == data.ErrProductNotFound {
		p.l.Println("[ERROR] product not found", err)

		rw.WriteHeader(http.StatusNotFound)
		data.ToJSON(&GenericError{Message: "Product not found in database"}, rw)
		return
	}

	// write the no content success header
	rw.WriteHeader(http.StatusNoContent)
}
