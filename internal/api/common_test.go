package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestJSONResponse(t *testing.T) {
	rr := httptest.NewRecorder()
	message := "Success"
	data := map[string]string{"key": "value"}

	JSONResponse(rr, http.StatusOK, message, data)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := standardResponse{
		Status:  http.StatusText(http.StatusOK),
		Message: message,
		Data:    data,
	}

	var actual standardResponse
	if err := json.NewDecoder(rr.Body).Decode(&actual); err != nil {
		t.Errorf("could not decode response: %v", err)
	}

	if actual != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", actual, expected)
	}
}

func TestJSONError(t *testing.T) {
	rr := httptest.NewRecorder()
	data := map[string]string{"key": "value"}

	JSONError(rr, http.StatusInternalServerError, http.ErrBodyNotAllowed, data)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusInternalServerError)
	}

	expected := standardError{
		Status: http.StatusText(http.StatusInternalServerError),
		Error:  http.ErrBodyNotAllowed.Error(),
		Data:   data,
	}

	var actual standardError
	if err := json.NewDecoder(rr.Body).Decode(&actual); err != nil {
		t.Errorf("could not decode response: %v", err)
	}

	if actual != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", actual, expected)
	}
}
