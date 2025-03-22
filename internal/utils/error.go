package utils

import "errors"

// Predefined custom errors
var (
	ErrTaskNotFound = errors.New("task not found")
	ErrDBConnection = errors.New("failed to connect to database")
)

// Custom error response struct
type ErrorResponse struct {
	Message string `json:"message"`
}

// Utility function to format error response
func NewErrorResponse(err error) ErrorResponse {
	return ErrorResponse{Message: err.Error()}
}
