package models

// StandardResponse is the base structure for all API responses
type StandardResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
}

// SuccessResponse is used for successful API responses
type SuccessResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// ErrorResponse is used for error API responses
type ErrorResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Error   interface{} `json:"error,omitempty"`
}

// ValidationError represents a field-specific validation error
type ValidationError struct {
	Field string `json:"field"`
	Issue string `json:"issue"`
}

// NewSuccessResponse creates a new success response
func NewSuccessResponse(message string, data interface{}) SuccessResponse {
	return SuccessResponse{
		Success: true,
		Message: message,
		Data:    data,
	}
}

// NewErrorResponse creates a new error response
func NewErrorResponse(message string, err interface{}) ErrorResponse {
	return ErrorResponse{
		Success: false,
		Message: message,
		Error:   err,
	}
}

// NewValidationError creates a new validation error
func NewValidationError(field, issue string) ValidationError {
	return ValidationError{
		Field: field,
		Issue: issue,
	}
}
