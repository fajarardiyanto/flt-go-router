package interfaces

import (
	"context"
	"encoding/json"
	"mime/multipart"
	"net/http"
)

// GetAllParams returns all route params stored in http.Request.
func GetAllParams(r *http.Request) ParamsMapType {
	if values, ok := r.Context().Value(ContextKey).(ParamsMapType); ok {
		return values
	}
	return nil
}

// GetParam returns parameter from url using a key
func GetParam(r *http.Request, key string) string {
	return GetAllParams(r)[key]
}

// GetQuery returns data from query params using a key
func GetQuery(r *http.Request, key string, fallback ...string) string {
	queries := r.URL.Query()
	query := queries.Get(key)
	if query == "" {
		if len(fallback) > 0 {
			return fallback[0]
		}
		return ""
	}
	return query
}

// GetFormData returns data from form using a key
func GetFormData(r *http.Request, key string, fallback ...string) string {
	value := r.PostFormValue(key)
	if value == "" {
		if len(fallback) > 0 {
			return fallback[0]
		}
		return ""
	}

	return value
}

// GetFormFile returns file from form using a key
func GetFormFile(r *http.Request, key string, maxFileSize ...int64) (*multipart.FileHeader, error) {
	var maxFileSizeParsed int64 = 10 << 20
	if len(maxFileSize) > 0 {
		maxFileSizeParsed = maxFileSize[0]
	}

	err := r.ParseMultipartForm(maxFileSizeParsed)
	if err != nil {
		return nil, err
	}

	file, header, err := r.FormFile(key)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return header, nil
}

// GetJSON get json data from requst body and store to variable reference
func GetJSON(r *http.Request, store interface{}) error {
	if err := json.NewDecoder(r.Body).Decode(&store); err != nil {
		return err
	}

	return nil
}

// SetContextValue returns http request object with new context that contains value
func SetContextValue(r *http.Request, key interface{}, value interface{}) *http.Request {
	ctx := r.Context()
	ctx = context.WithValue(ctx, key, value)
	r = r.WithContext(ctx)
	return r
}

// GetContextValue returns value from http request context
func GetContextValue(r *http.Request, key interface{}) interface{} {
	value := r.Context().Value(key)
	return value
}

// RequestError is used to pass an error during the request through the application
// with web specific context.
type RequestError struct {
	Status int   `json:"status"`
	Err    error `json:"err"`
}

// NewRequestError wraps a provided error with an HTTP status code.
// This function should be used when handlers encounter expected
// (trusted) errors.
func NewRequestError(err error, status int) *RequestError {
	return &RequestError{
		Status: status,
		Err:    err,
	}
}

// Error implements error interface.
func (re RequestError) Error() string {
	return re.Err.Error()
}
