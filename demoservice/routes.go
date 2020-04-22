package demoservice

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/rs/cors"
)

// Router returns the API router
func Router() http.Handler {
	mux := httprouter.New()

	// Service endpoints
	mux.POST("/v1/checkout", checkoutHandler)
	mux.POST("/v1/preauth", addPreauthHandler)
	mux.POST("/v2/resolve/:pa", resolverHandler)
	mux.POST("/v1/callback", callbackHandler)
	mux.POST("/v1/success", successHandler)

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedHeaders: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "DELETE"},
	})

	handler := c.Handler(mux)

	return handler
}
