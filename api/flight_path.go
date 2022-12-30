package api

import (
	"net/http"

	"github.com/manigandand/adk/errors"
	"github.com/manigandand/adk/respond"
)

func flightPathCalculatorHandler(w http.ResponseWriter, r *http.Request) *errors.AppError {
	// ctx := r.Context()
	respond.OK(w, "todo..")
	return nil
}
