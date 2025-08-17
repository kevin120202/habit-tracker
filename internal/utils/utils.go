package utils

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func ReadIDParam(r *http.Request) (int64, error) {
	idParam := chi.URLParam(r, "id")

	if idParam == "" {
		return 0, errors.New("invalid id parameter")
	}

	// Convert the "id" parameter from string to int64 (a 64-bit integer)
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		return 0, errors.New("invalid id parameter type")
	}

	return id, nil
}
