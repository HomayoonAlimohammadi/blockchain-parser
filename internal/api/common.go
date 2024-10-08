package api

import (
	"encoding/json"
	"net/http"
)

type standardResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type standardError struct {
	Status string      `json:"status"`
	Error  string      `json:"error"`
	Data   interface{} `json:"data,omitempty"`
}

// JSONResponse sends a JSON response with the given status code, message, and data
func JSONResponse(w http.ResponseWriter, status int, message string, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	resp := standardResponse{
		Status:  http.StatusText(status),
		Message: message,
		Data:    data,
	}

	json.NewEncoder(w).Encode(resp)
}

// JSONError sends a JSON error response with the given status code, error, and data
func JSONError(w http.ResponseWriter, status int, err error, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	resp := standardError{
		Status: http.StatusText(status),
		Error:  err.Error(),
		Data:   data,
	}

	json.NewEncoder(w).Encode(resp)
}
