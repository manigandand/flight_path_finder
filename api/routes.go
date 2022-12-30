package api

import (
	"net/http"

	"github.com/go-chi/chi"
	adkapi "github.com/manigandand/adk/api"
)

// Routes - all the registered routes
func Routes(router chi.Router) {
	router.Get("/", adkapi.IndexHandeler)
	router.Get("/top", adkapi.HealthHandeler)
	router.Method(http.MethodPost, "/calculate", adkapi.Handler(flightPathCalculatorHandler))
}
