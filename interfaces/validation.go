package interfaces

import "encoding/json"

// FieldError contains error and field for validation error.
type FieldError struct {
	Field string `json:"field"`
	Error string `json:"error"`
}

// ValidationError represents validation errors.
type ValidationError []FieldError

// NewValidationError creates a new validation error.
func NewValidationError(fields ...FieldError) ValidationError {
	fieldErrors := make([]FieldError, 0)

	fieldErrors = append(fieldErrors, fields...)

	return fieldErrors
}

// Error implements the error interface.
func (ve ValidationError) Error() string {
	d, err := json.Marshal(ve)
	if err != nil {
		return err.Error()
	}

	return string(d)
}
