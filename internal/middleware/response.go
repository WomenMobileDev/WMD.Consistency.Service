package middleware

import (
	"net/http"

	"github.com/WomenMobileDev/WMD.Consistency.Service/internal/models"
	"github.com/gin-gonic/gin"
)

// ResponseFormatter standardizes API responses
func ResponseFormatter() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Process the request first
		c.Next()

		// Skip if response is already committed
		if c.Writer.Written() {
			return
		}

		// Get the status code that would be sent
		statusCode := c.Writer.Status()

		// Get the data from the context
		data, exists := c.Get("response_data")
		if !exists {
			return
		}

		// Check if the response is already standardized
		if _, ok := data.(models.SuccessResponse); ok {
			c.JSON(statusCode, data)
			return
		}
		if _, ok := data.(models.ErrorResponse); ok {
			c.JSON(statusCode, data)
			return
		}

		// Format the response based on status code
		if statusCode >= 200 && statusCode < 300 {
			// Success response
			message := "Request successful"
			if statusCode == http.StatusCreated {
				message = "Resource created successfully"
			} else if statusCode == http.StatusOK {
				message = "Resource fetched successfully"
			} else if statusCode == http.StatusNoContent {
				message = "Resource deleted successfully"
			}

			c.JSON(statusCode, models.NewSuccessResponse(message, data))
		} else {
			// Error response
			message := "Request failed"
			if statusCode == http.StatusBadRequest {
				message = "Invalid request"
			} else if statusCode == http.StatusUnauthorized {
				message = "Unauthorized"
			} else if statusCode == http.StatusForbidden {
				message = "Forbidden"
			} else if statusCode == http.StatusNotFound {
				message = "Resource not found"
			} else if statusCode == http.StatusConflict {
				message = "Resource conflict"
			} else if statusCode >= 500 {
				message = "Server error"
			}

			var errorObj interface{}
			if errMap, ok := data.(gin.H); ok && errMap["error"] != nil {
				errorObj = errMap["error"]
			} else {
				errorObj = data
			}

			c.JSON(statusCode, models.NewErrorResponse(message, errorObj))
		}
	}
}

// ErrorHandler is a middleware that handles errors and returns standardized responses
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// Check if there were any errors
		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err
			
			// Handle different types of errors
			switch err := err.(type) {
			case *gin.Error:
				c.JSON(http.StatusBadRequest, models.NewErrorResponse("Validation failed", err.Error()))
			default:
				c.JSON(http.StatusInternalServerError, models.NewErrorResponse("Internal server error", err.Error()))
			}
		}
	}
}

// Helper functions for handlers to use standardized responses

// RespondWithSuccess sends a standardized success response
func RespondWithSuccess(c *gin.Context, statusCode int, message string, data interface{}) {
	c.JSON(statusCode, models.NewSuccessResponse(message, data))
}

// RespondWithError sends a standardized error response
func RespondWithError(c *gin.Context, statusCode int, message string, err interface{}) {
	c.JSON(statusCode, models.NewErrorResponse(message, err))
}

// RespondWithValidationError sends a standardized validation error response
func RespondWithValidationError(c *gin.Context, field, issue string) {
	c.JSON(http.StatusBadRequest, models.NewErrorResponse(
		"Validation failed",
		models.NewValidationError(field, issue),
	))
}

// RespondWithCreated sends a standardized created response
func RespondWithCreated(c *gin.Context, data interface{}) {
	RespondWithSuccess(c, http.StatusCreated, "Resource created successfully", data)
}

// RespondWithOK sends a standardized OK response
func RespondWithOK(c *gin.Context, data interface{}) {
	RespondWithSuccess(c, http.StatusOK, "Resource fetched successfully", data)
}

// RespondWithNoContent sends a standardized no content response
func RespondWithNoContent(c *gin.Context) {
	c.Status(http.StatusNoContent)
}

// RespondWithBadRequest sends a standardized bad request response
func RespondWithBadRequest(c *gin.Context, err interface{}) {
	RespondWithError(c, http.StatusBadRequest, "Invalid request", err)
}

// RespondWithUnauthorized sends a standardized unauthorized response
func RespondWithUnauthorized(c *gin.Context) {
	RespondWithError(c, http.StatusUnauthorized, "Unauthorized", nil)
}

// RespondWithForbidden sends a standardized forbidden response
func RespondWithForbidden(c *gin.Context) {
	RespondWithError(c, http.StatusForbidden, "Forbidden", nil)
}

// RespondWithNotFound sends a standardized not found response
func RespondWithNotFound(c *gin.Context, resource string) {
	RespondWithError(c, http.StatusNotFound, resource+" not found", nil)
}

// RespondWithConflict sends a standardized conflict response
func RespondWithConflict(c *gin.Context, message string) {
	RespondWithError(c, http.StatusConflict, message, nil)
}

// RespondWithInternalError sends a standardized internal server error response
func RespondWithInternalError(c *gin.Context, err interface{}) {
	RespondWithError(c, http.StatusInternalServerError, "Internal server error", err)
}
