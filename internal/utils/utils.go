package utils

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// Envelope is a type alias for a map that allows storing dynamic key-value pairs.
// It is used to structure the response data in a flexible way (similar to a JSON object).
type Envelope map[string]interface{}

// WriteJSON writes a JSON response to the client.
// It takes the HTTP response writer, the status code, and the data to send as JSON.
// Returns an error if the encoding fails.
func WriteJSON(w http.ResponseWriter, status int, data Envelope) error {
	// MarshalIndent is used to convert the `data` (which is of type Envelope) to a formatted JSON string.
	// The second and third arguments help to pretty-print the JSON (with indentation).
	js, err := json.MarshalIndent(data, "", " ")
	if err != nil {
		return err // Return an error if JSON encoding fails
	}

	// Append a newline character to the JSON string for better readability.
	js = append(js, '\n')

	// Set the Content-Type header to "application/json" to tell the client that the response is in JSON format.
	w.Header().Set("Content-Type", "application/json")

	// Set the status code for the response (e.g., 200 OK, 400 Bad Request, etc.)
	w.WriteHeader(status)

	// Write the JSON data to the response body.
	w.Write(js)

	return nil
}

func ReadIDParam(r *http.Request) (uuid.UUID, error) {
	idParam := chi.URLParam(r, "id")

	if idParam == "" {
		return uuid.Nil, errors.New("invalid id parameter")
	}

	// Convert the "id" parameter from string to int64 (a 64-bit integer)
	id, err := uuid.Parse(idParam)
	if err != nil {
		return uuid.Nil, errors.New("invalid id parameter type")
	}

	return id, nil
}
