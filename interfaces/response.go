package interfaces

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// String write string data to response body
func String(w http.ResponseWriter, data string) error {
	w.Header().Add("content-type", "text/plain")
	_, err := fmt.Fprintln(w, data)
	return err
}

// StringWithStatus write string data to response body with status code
func StringWithStatus(w http.ResponseWriter, status int, data string) error {
	_ = String(w, data)
	w.WriteHeader(status)
	return nil
}

// JSON write JSON data to response
func JSON(w http.ResponseWriter, status int, data interface{}) error {
	w.Header().Add("content-type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		return err
	}

	return nil
}

// APIResponseError is a response error format for REST API.
type APIResponseError struct {
	Code    int    `json:"code"`
	Fields  error  `json:"fields,omitempty"`
	Message string `json:"message,omitempty"`
}

// NewAPIResponseError creates a new response error.
func NewAPIResponseError(message string, code int) *APIResponseError {
	return &APIResponseError{
		Code:    code,
		Message: message,
	}
}

// NewAPIValidationError creates a new validation error.
func NewAPIValidationError(fields error, code int) *APIResponseError {
	return &APIResponseError{
		Code:    code,
		Fields:  fields,
		Message: "Validation(s) error",
	}
}
