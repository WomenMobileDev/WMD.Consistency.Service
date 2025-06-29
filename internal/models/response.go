package models

type StandardResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type SuccessResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type ErrorObject struct {
	Code    string      `json:"code"`              // e.g., "USER_NOT_FOUND"
	Message string      `json:"message"`           // user-friendly error message
	Details interface{} `json:"details,omitempty"` // optional: field errors, invalid values, etc.
}

type ErrorResponse struct {
	Success bool        `json:"success"`
	Error   ErrorObject `json:"error"`
}

type ValidationError struct {
	Field string `json:"field"`
	Issue string `json:"issue"`
}

func NewSuccessResponse(message string, data interface{}) SuccessResponse {
	return SuccessResponse{
		Success: true,
		Message: message,
		Data:    data,
	}
}

func NewErrorResponse(code, message string, details interface{}) ErrorResponse {
	return ErrorResponse{
		Success: false,
		Error: ErrorObject{
			Code:    code,
			Message: message,
			Details: details,
		},
	}
}

func NewValidationError(field, issue string) ValidationError {
	return ValidationError{
		Field: field,
		Issue: issue,
	}
}

type AppError struct {
	Code    string
	Message string
	Details interface{}
}

func (e *AppError) Error() string {
	return e.Message
}
