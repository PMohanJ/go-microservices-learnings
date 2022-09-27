package handlers

import (
	"context"
	"net/http"

	"github.com/pmohanj/go-microservices/product-api/data"
)

// A middleware is also a handler maily intented to use for data validation
// or authentication before the request is passed to any actual main handlers
func (p Products) MiddlewareValidateProduct(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		prod := &data.Product{}

		err := data.FromJSON(prod, r.Body)
		if err != nil {
			p.l.Println("[ERROR] deserializing product", err)

			rw.WriteHeader(http.StatusBadRequest)
			data.ToJSON(&GenericError{Message: err.Error()}, rw)
			return
		}

		// validate the product
		errs := p.v.Validate(prod)
		if err != nil {
			p.l.Println("[ERROR] validating product", err)

			// return the validation messages as an array
			rw.WriteHeader(http.StatusUnprocessableEntity)
			data.ToJSON(&ValidationError{Messages: errs.Errors()}, rw)
			return
		}

		// add product to the context
		ctx := context.WithValue(r.Context(), KeyProduct{}, prod)
		req := r.WithContext(ctx)

		// calling next handler, which can be another middlerware or final handler
		next.ServeHTTP(rw, req)
	})
}
